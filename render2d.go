package main

import (
	"fmt"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/shadercross"
	"github.com/Zyko0/go-sdl3/ttf"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex2D struct {
	X, Y, Z    float32
	R, G, B, A float32
}

type TextVertex struct {
	X, Y       float32
	U, V       float32
	R, G, B, A float32
}

type SpriteVertex struct {
	X, Y       float32
	U, V       float32
	R, G, B, A float32
	Fx         float32 // Effect ID
	FxParam    float32 // Effect progress / intensity parameter
}
type UniformBlock2D struct {
	Ortho mgl32.Mat4
	Time  float32
	Pad   [3]float32 // 16-byte alignment for GPU uniform buffers
}
type DrawCmdType uint8

const (
	CmdNone    DrawCmdType = 0
	CmdShapes  DrawCmdType = 1
	CmdText    DrawCmdType = 2
	CmdTexture DrawCmdType = 3
)

type DrawCommand struct {
	Type        DrawCmdType
	Texture     *sdl.GPUTexture
	VertexStart uint32
	VertexCount uint32
}

const shader2DHLSL = `
cbuffer UniformBlock : register(b0, space1)
{
    float4x4 OrthoMatrix;
    float    Time;
    float3   _pad;
};

// 1. SOLID COLOR SHADER
struct VSInput
{
    float3 Position : TEXCOORD0;
    float4 Color    : TEXCOORD1;
};

struct PSInput
{
    float4 Position : SV_Position;
    float4 Color    : TEXCOORD0;
};

PSInput VSMain(VSInput input)
{
    PSInput output;
    output.Position = mul(OrthoMatrix, float4(input.Position, 1.0));
    output.Color    = input.Color;
    return output;
}

float4 PSMain(PSInput input) : SV_Target
{
    return input.Color;
}

// 2. TEXT SHADER
Texture2D    AtlasTexture : register(t0, space2);
SamplerState TextSampler  : register(s0, space2);

struct VSTextInput
{
    float2 Position : TEXCOORD0;
    float2 UV       : TEXCOORD1;
    float4 Color    : TEXCOORD2;
};

struct PSTextInput
{
    float4 Position : SV_Position;
    float2 UV       : TEXCOORD0;
    float4 Color    : TEXCOORD1;
};

PSTextInput VSText(VSTextInput input)
{
    PSTextInput output;
    output.Position = mul(OrthoMatrix, float4(input.Position, 0.0, 1.0));
    output.UV       = input.UV;
    output.Color    = input.Color;
    return output;
}

float4 PSText(PSTextInput input) : SV_Target
{
    float4 sampled = AtlasTexture.Sample(TextSampler, input.UV);
    return float4(input.Color.rgb, input.Color.a * sampled.a);
}

// 3. SPRITE TEXTURE SHADER
Texture2D    SpriteTexture : register(t0, space2);
SamplerState SpriteSampler : register(s0, space2);

struct VSSpriteInput
{
    float2 Position : TEXCOORD0;
    float2 UV       : TEXCOORD1;
    float4 Color    : TEXCOORD2;
    float2 FxData   : TEXCOORD3; // x: Fx ID, y: FxParam
};

struct PSSpriteInput
{
    float4 Position : SV_Position;
    float2 UV       : TEXCOORD0;
    float4 Color    : TEXCOORD1;
    float2 FxData   : TEXCOORD2;
    float  Time     : TEXCOORD3;
};

PSSpriteInput VSSprite(VSSpriteInput input)
{
    PSSpriteInput output;
    output.Position = mul(OrthoMatrix, float4(input.Position, 0.0, 1.0));
    output.UV       = input.UV;
    output.Color    = input.Color;
    output.FxData   = input.FxData;
    output.Time     = Time;
    return output;
}

float4 PSSprite(PSSpriteInput input) : SV_Target
{
    float4 texColor = SpriteTexture.Sample(SpriteSampler, input.UV);
    float4 finalColor = texColor * input.Color;

    // Fx 1: Grayscale
    if (input.FxData.x >= 0.5 && input.FxData.x < 1.5)
    {
        float gray = dot(finalColor.rgb, float3(0.299, 0.587, 0.114));
        finalColor.rgb = float3(gray, gray, gray);
    }
    // Fx 2: Sepia
    else if (input.FxData.x >= 1.5 && input.FxData.x < 2.5)
    {
        float3 sepia;
        sepia.r = dot(finalColor.rgb, float3(0.393, 0.769, 0.189));
        sepia.g = dot(finalColor.rgb, float3(0.349, 0.686, 0.168));
        sepia.b = dot(finalColor.rgb, float3(0.272, 0.534, 0.131));
        finalColor.rgb = sepia;
    }
    // Fx 3: Posterize
    else if (input.FxData.x >= 2.5 && input.FxData.x < 3.5)
    {
        float bands = (input.FxData.y <= 1.0) ? 4.0 : input.FxData.y;
        finalColor.rgb = floor(finalColor.rgb * bands) / bands;
    }
    // Fx 4: Scanlines (Hologram)
    else if (input.FxData.x >= 3.5)
    {
        float scanline = sin((input.UV.y * 120.0) + (input.Time * 8.0)) * 0.5 + 0.5;
        finalColor.rgb *= (0.7 + 0.3 * scanline);
        finalColor.rgb = lerp(finalColor.rgb, float3(0.2, 0.8, 1.0) * finalColor.rgb, 0.3);
    }

    return finalColor;
}
`

type Render2D struct {
	gpuDev       *sdl.GPUDevice
	textEngine   *ttf.TextEngine
	pipeFill     *sdl.GPUGraphicsPipeline
	pipeText     *sdl.GPUGraphicsPipeline
	pipeSprite   *sdl.GPUGraphicsPipeline
	samplerText  *sdl.GPUSampler
	samplerPixel *sdl.GPUSampler
	vBufFill     *sdl.GPUBuffer
	vBufText     *sdl.GPUBuffer
	vBufSprite   *sdl.GPUBuffer
	fillVerts    []Vertex2D
	textVerts    []TextVertex
	spriteVerts  []SpriteVertex
	commands     []DrawCommand
	textAtlas    *sdl.GPUTexture
	uniforms     UniformBlock2D
}

var R2D *Render2D

func InitRender2D(gpuDev *sdl.GPUDevice, win *sdl.Window, scrW, scrH int) (*Render2D, error) {
	vShader, err := CompileShader(gpuDev, shader2DHLSL, "VSMain", shadercross.SHADERSTAGE_VERTEX, 0, 1)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(vShader)

	pShader, err := CompileShader(gpuDev, shader2DHLSL, "PSMain", shadercross.SHADERSTAGE_FRAGMENT, 0, 0)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(pShader)

	vShaderText, err := CompileShader(gpuDev, shader2DHLSL, "VSText", shadercross.SHADERSTAGE_VERTEX, 0, 1)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(vShaderText)

	pShaderText, err := CompileShader(gpuDev, shader2DHLSL, "PSText", shadercross.SHADERSTAGE_FRAGMENT, 1, 0)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(pShaderText)

	// Compile VSSprite
	vShaderSprite, err := CompileShader(gpuDev, shader2DHLSL, "VSSprite", shadercross.SHADERSTAGE_VERTEX, 0, 1)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(vShaderSprite)

	pShaderSprite, err := CompileShader(gpuDev, shader2DHLSL, "PSSprite", shadercross.SHADERSTAGE_FRAGMENT, 1, 0)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(pShaderSprite)

	textEngine, err := ttf.CreateGPUTextEngine(gpuDev)
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU text engine: %w", err)
	}

	swapchainFormat := gpuDev.SwapchainTextureFormat(win)

	vertexAttrs := []sdl.GPUVertexAttribute{
		{Location: 0, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT3, Offset: 0},
		{Location: 1, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT4, Offset: 12},
	}
	bufferDescs := []sdl.GPUVertexBufferDescription{
		{Slot: 0, Pitch: uint32(unsafe.Sizeof(Vertex2D{})), InputRate: sdl.GPU_VERTEXINPUTRATE_VERTEX},
	}

	targetInfo := sdl.GPUGraphicsPipelineTargetInfo{
		ColorTargetDescriptions: []sdl.GPUColorTargetDescription{
			{
				Format: swapchainFormat,
				BlendState: sdl.GPUColorTargetBlendState{
					EnableBlend:         true,
					SrcColorBlendfactor: sdl.GPU_BLENDFACTOR_SRC_ALPHA,
					DstColorBlendfactor: sdl.GPU_BLENDFACTOR_ONE_MINUS_SRC_ALPHA,
					ColorBlendOp:        sdl.GPU_BLENDOP_ADD,
					SrcAlphaBlendfactor: sdl.GPU_BLENDFACTOR_ONE,
					DstAlphaBlendfactor: sdl.GPU_BLENDFACTOR_ZERO,
					AlphaBlendOp:        sdl.GPU_BLENDOP_ADD,
				},
			},
		},
		HasDepthStencilTarget: true,
		DepthStencilFormat:    sdl.GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT,
	}

	depthState2D := sdl.GPUDepthStencilState{
		EnableDepthTest:  false,
		EnableDepthWrite: false,
	}

	pipeFill, err := gpuDev.CreateGraphicsPipeline(&sdl.GPUGraphicsPipelineCreateInfo{
		VertexShader:      vShader,
		FragmentShader:    pShader,
		VertexInputState:  sdl.GPUVertexInputState{VertexAttributes: vertexAttrs, VertexBufferDescriptions: bufferDescs},
		PrimitiveType:     sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
		DepthStencilState: depthState2D,
		TargetInfo:        targetInfo,
	})
	if err != nil {
		return nil, err
	}

	attrsTex := []sdl.GPUVertexAttribute{
		{Location: 0, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: 0},
		{Location: 1, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: uint32(unsafe.Offsetof(TextVertex{}.U))},
		{Location: 2, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT4, Offset: uint32(unsafe.Offsetof(TextVertex{}.R))},
	}
	descTex := []sdl.GPUVertexBufferDescription{
		{Slot: 0, Pitch: uint32(unsafe.Sizeof(TextVertex{})), InputRate: sdl.GPU_VERTEXINPUTRATE_VERTEX},
	}

	pipeText, err := gpuDev.CreateGraphicsPipeline(&sdl.GPUGraphicsPipelineCreateInfo{
		VertexShader:      vShaderText,
		FragmentShader:    pShaderText,
		VertexInputState:  sdl.GPUVertexInputState{VertexAttributes: attrsTex, VertexBufferDescriptions: descTex},
		PrimitiveType:     sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
		DepthStencilState: depthState2D,
		TargetInfo:        targetInfo,
	})
	if err != nil {
		return nil, err
	}

	// Sprite Vertex Layout: Location 3 takes Fx & FxParam as a FLOAT2
	attrsSprite := []sdl.GPUVertexAttribute{
		{Location: 0, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: 0},
		{Location: 1, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: uint32(unsafe.Offsetof(SpriteVertex{}.U))},
		{Location: 2, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT4, Offset: uint32(unsafe.Offsetof(SpriteVertex{}.R))},
		{Location: 3, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: uint32(unsafe.Offsetof(SpriteVertex{}.Fx))},
	}
	descSprite := []sdl.GPUVertexBufferDescription{
		{Slot: 0, Pitch: uint32(unsafe.Sizeof(SpriteVertex{})), InputRate: sdl.GPU_VERTEXINPUTRATE_VERTEX},
	}

	pipeSprite, err := gpuDev.CreateGraphicsPipeline(&sdl.GPUGraphicsPipelineCreateInfo{
		VertexShader:      vShaderSprite,
		FragmentShader:    pShaderSprite,
		VertexInputState:  sdl.GPUVertexInputState{VertexAttributes: attrsSprite, VertexBufferDescriptions: descSprite},
		PrimitiveType:     sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
		DepthStencilState: depthState2D,
		TargetInfo:        targetInfo,
	})
	if err != nil {
		return nil, err
	}

	samplerText, _ := gpuDev.CreateSampler(&sdl.GPUSamplerCreateInfo{
		MinFilter:    sdl.GPU_FILTER_LINEAR,
		MagFilter:    sdl.GPU_FILTER_LINEAR,
		MipmapMode:   sdl.GPU_SAMPLERMIPMAPMODE_LINEAR,
		AddressModeU: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeV: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeW: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
	})

	samplerPixel, _ := gpuDev.CreateSampler(&sdl.GPUSamplerCreateInfo{
		MinFilter:    sdl.GPU_FILTER_NEAREST,
		MagFilter:    sdl.GPU_FILTER_NEAREST,
		MipmapMode:   sdl.GPU_SAMPLERMIPMAPMODE_NEAREST,
		AddressModeU: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeV: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
		AddressModeW: sdl.GPU_SAMPLERADDRESSMODE_CLAMP_TO_EDGE,
	})

	initCap := 32768
	vBufFill, _ := gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_VERTEX, Size: uint32(initCap * int(unsafe.Sizeof(Vertex2D{})))})
	vBufText, _ := gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_VERTEX, Size: uint32(initCap * int(unsafe.Sizeof(TextVertex{})))})
	vBufSprite, _ := gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_VERTEX, Size: uint32(initCap * int(unsafe.Sizeof(SpriteVertex{})))})

	orthoMatrix := mgl32.Ortho(0, float32(scrW), float32(scrH), 0, 0, -1)

	r2d := &Render2D{
		gpuDev:       gpuDev,
		textEngine:   textEngine,
		pipeFill:     pipeFill,
		pipeText:     pipeText,
		pipeSprite:   pipeSprite,
		samplerText:  samplerText,
		samplerPixel: samplerPixel,
		vBufFill:     vBufFill,
		vBufText:     vBufText,
		vBufSprite:   vBufSprite,
		fillVerts:    make([]Vertex2D, 0, 2048),
		textVerts:    make([]TextVertex, 0, 2048),
		spriteVerts:  make([]SpriteVertex, 0, 2048),
		commands:     make([]DrawCommand, 0, 128),
		uniforms: UniformBlock2D{
			Ortho: orthoMatrix,
			Time:  0,
		},
	}

	R2D = r2d
	return r2d, nil
}

func (r2d *Render2D) Destroy() {
	if r2d == nil {
		return
	}
	r2d.textEngine.DestroyGPU()
	r2d.gpuDev.ReleaseGraphicsPipeline(r2d.pipeFill)
	r2d.gpuDev.ReleaseGraphicsPipeline(r2d.pipeText)
	r2d.gpuDev.ReleaseGraphicsPipeline(r2d.pipeSprite)
	r2d.gpuDev.ReleaseSampler(r2d.samplerText)
	r2d.gpuDev.ReleaseSampler(r2d.samplerPixel)
	r2d.gpuDev.ReleaseBuffer(r2d.vBufFill)
	r2d.gpuDev.ReleaseBuffer(r2d.vBufText)
	r2d.gpuDev.ReleaseBuffer(r2d.vBufSprite)
}

func (r2d *Render2D) Begin() {
	r2d.fillVerts = r2d.fillVerts[:0]
	r2d.textVerts = r2d.textVerts[:0]
	r2d.spriteVerts = r2d.spriteVerts[:0]
	r2d.commands = r2d.commands[:0]
	r2d.textAtlas = nil
}

func (r2d *Render2D) addShapeVertices(verts ...Vertex2D) {
	start := uint32(len(r2d.fillVerts))
	r2d.fillVerts = append(r2d.fillVerts, verts...)
	count := uint32(len(verts))

	if len(r2d.commands) > 0 && r2d.commands[len(r2d.commands)-1].Type == CmdShapes {
		r2d.commands[len(r2d.commands)-1].VertexCount += count
	} else {
		r2d.commands = append(r2d.commands, DrawCommand{Type: CmdShapes, VertexStart: start, VertexCount: count})
	}
}

func (r2d *Render2D) addTextVertices(verts []TextVertex) {
	if len(verts) == 0 {
		return
	}
	start := uint32(len(r2d.textVerts))
	r2d.textVerts = append(r2d.textVerts, verts...)
	count := uint32(len(verts))

	if len(r2d.commands) > 0 && r2d.commands[len(r2d.commands)-1].Type == CmdText {
		r2d.commands[len(r2d.commands)-1].VertexCount += count
	} else {
		r2d.commands = append(r2d.commands, DrawCommand{Type: CmdText, VertexStart: start, VertexCount: count})
	}
}

func (r2d *Render2D) addSpriteVertices(tex *sdl.GPUTexture, verts []SpriteVertex) {
	if len(verts) == 0 || tex == nil {
		return
	}
	start := uint32(len(r2d.spriteVerts))
	r2d.spriteVerts = append(r2d.spriteVerts, verts...)
	count := uint32(len(verts))

	if len(r2d.commands) > 0 && r2d.commands[len(r2d.commands)-1].Type == CmdTexture && r2d.commands[len(r2d.commands)-1].Texture == tex {
		r2d.commands[len(r2d.commands)-1].VertexCount += count
	} else {
		r2d.commands = append(r2d.commands, DrawCommand{Type: CmdTexture, Texture: tex, VertexStart: start, VertexCount: count})
	}
}

func (r2d *Render2D) Upload(cmd *sdl.GPUCommandBuffer) {
	fillCount := len(r2d.fillVerts)
	textCount := len(r2d.textVerts)
	spriteCount := len(r2d.spriteVerts)
	if fillCount == 0 && textCount == 0 && spriteCount == 0 {
		return
	}

	fillBytes := uint32(fillCount * int(unsafe.Sizeof(Vertex2D{})))
	textBytes := uint32(textCount * int(unsafe.Sizeof(TextVertex{})))
	spriteBytes := uint32(spriteCount * int(unsafe.Sizeof(SpriteVertex{})))
	totalBytes := fillBytes + textBytes + spriteBytes

	tBuf, err := r2d.gpuDev.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD, Size: totalBytes})
	if err != nil {
		return
	}

	mapped, err := r2d.gpuDev.MapTransferBuffer(tBuf, false)
	if err != nil {
		r2d.gpuDev.ReleaseTransferBuffer(tBuf)
		return
	}

	memPtr := (*byte)(unsafe.Pointer(mapped))

	if fillBytes > 0 {
		copy(unsafe.Slice(memPtr, fillBytes), unsafe.Slice((*byte)(unsafe.Pointer(&r2d.fillVerts[0])), fillBytes))
	}
	if textBytes > 0 {
		textDest := (*byte)(unsafe.Add(unsafe.Pointer(memPtr), fillBytes))
		copy(unsafe.Slice(textDest, textBytes), unsafe.Slice((*byte)(unsafe.Pointer(&r2d.textVerts[0])), textBytes))
	}
	if spriteBytes > 0 {
		spriteDest := (*byte)(unsafe.Add(unsafe.Pointer(memPtr), fillBytes+textBytes))
		copy(unsafe.Slice(spriteDest, spriteBytes), unsafe.Slice((*byte)(unsafe.Pointer(&r2d.spriteVerts[0])), spriteBytes))
	}

	r2d.gpuDev.UnmapTransferBuffer(tBuf)

	copyPass := cmd.BeginCopyPass()
	if fillBytes > 0 {
		copyPass.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: 0}, &sdl.GPUBufferRegion{Buffer: r2d.vBufFill, Offset: 0, Size: fillBytes}, false)
	}
	if textBytes > 0 {
		copyPass.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: fillBytes}, &sdl.GPUBufferRegion{Buffer: r2d.vBufText, Offset: 0, Size: textBytes}, false)
	}
	if spriteBytes > 0 {
		copyPass.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: fillBytes + textBytes}, &sdl.GPUBufferRegion{Buffer: r2d.vBufSprite, Offset: 0, Size: spriteBytes}, false)
	}
	copyPass.End()

	r2d.gpuDev.ReleaseTransferBuffer(tBuf)
}

func (r2d *Render2D) Draw(cmd *sdl.GPUCommandBuffer, renderPass *sdl.GPURenderPass) {
	if len(r2d.commands) == 0 {
		return
	}

	// 🚀 Push Uniforms with updated running time
	r2d.uniforms.Time = float32(sdl.TicksNS()) / 1e9
	uBytes := unsafe.Slice((*byte)(unsafe.Pointer(&r2d.uniforms)), unsafe.Sizeof(r2d.uniforms))
	cmd.PushVertexUniformData(0, uBytes)

	lastType := CmdNone
	var lastTex *sdl.GPUTexture

	for _, c := range r2d.commands {
		switch c.Type {
		case CmdShapes:
			if lastType != CmdShapes {
				renderPass.BindGraphicsPipeline(r2d.pipeFill)
				renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: r2d.vBufFill, Offset: 0}})
				lastType = CmdShapes
			}
			renderPass.DrawPrimitives(c.VertexCount, 1, c.VertexStart, 0)

		case CmdText:
			if r2d.textAtlas != nil {
				if lastType != CmdText {
					renderPass.BindGraphicsPipeline(r2d.pipeText)
					renderPass.BindFragmentSamplers([]sdl.GPUTextureSamplerBinding{
						{Texture: r2d.textAtlas, Sampler: r2d.samplerText},
					})
					renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: r2d.vBufText, Offset: 0}})
					lastType = CmdText
				}
				renderPass.DrawPrimitives(c.VertexCount, 1, c.VertexStart, 0)
			}

		case CmdTexture:
			if c.Texture != nil {
				if lastType != CmdTexture || lastTex != c.Texture {
					renderPass.BindGraphicsPipeline(r2d.pipeSprite)
					renderPass.BindFragmentSamplers([]sdl.GPUTextureSamplerBinding{
						{Texture: c.Texture, Sampler: r2d.samplerPixel},
					})
					renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: r2d.vBufSprite, Offset: 0}})
					lastType = CmdTexture
					lastTex = c.Texture
				}
				renderPass.DrawPrimitives(c.VertexCount, 1, c.VertexStart, 0)
			}
		}
	}
}

func V2D(x, y float32, c sdl.Color) Vertex2D {
	rf, gf, bf, af := color2float(c)
	return Vertex2D{X: x, Y: y, Z: 0, R: rf, G: gf, B: bf, A: af}
}
