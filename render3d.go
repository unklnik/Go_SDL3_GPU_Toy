package main

import (
	"fmt"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/shadercross"
	"github.com/go-gl/mathgl/mgl32"
)

type RenderMode3D int32

const (
	ModeShaded    RenderMode3D = 0
	ModeSolid     RenderMode3D = 1
	ModeWireframe RenderMode3D = 2
)

// 3D Uniform Buffer with Dynamic Multi-Light Array (704 bytes, 16-byte aligned)
type Uniform3D struct {
	MVP          [16]float32
	Model        [16]float32
	CameraPos    [3]float32
	RenderMode   int32
	AmbientColor [4]float32
	ColorTint    [4]float32
	UVScale      [2]float32
	NumLights    int32
	Padding      float32
	Lights       [8]GPULight
}

type GPUMesh struct {
	VertexBuffer *sdl.GPUBuffer
	TriBuffer    *sdl.GPUBuffer
	LineBuffer   *sdl.GPUBuffer
	Texture      *sdl.GPUTexture
	NumTriangles uint32
	NumLines     uint32
}

type GPUModel struct {
	Meshes []GPUMesh
}

// 3D HLSL Shader with Dynamic Point, Directional, and Spot Lights
const shader3DHLSL = `
struct GPULight
{
    float3 Position;
    int    Type; // 0=Directional, 1=Point, 2=Spot
    float3 Direction;
    float  Intensity;
    float4 Color;
    float  Range;
    float  SpotAngle;
    float2 Padding;
};

cbuffer UniformData3D : register(b0, space1)
{
    float4x4 MVP;
    float4x4 Model;
    float3   CameraPos;
    int      RenderMode;
    float4   AmbientColor;
    float4   ColorTint;
    float2   UVScale;
    int      NumLights;
    float    Padding;
    GPULight Lights[8];
};

struct VS3DInput
{
    float3 Position : TEXCOORD0; // Location 0
    float3 Normal   : TEXCOORD1; // Location 1
    float2 UV       : TEXCOORD2; // Location 2
    float4 Color    : TEXCOORD3; // Location 3
};

struct PS3DInput
{
    float4 Position : SV_Position;
    float2 UV       : TEXCOORD0;
    float4 Color    : TEXCOORD1;
};

PS3DInput VS3D(VS3DInput input)
{
    PS3DInput output;
    output.Position = mul(MVP, float4(input.Position, 1.0));
    output.UV       = input.UV * UVScale;

    // Mode 2: Wireframe
    if (RenderMode == 2)
    {
        output.Color = float4(0.2, 1.0, 0.5, -1.0);
        return output;
    }

    float4 baseColor = input.Color * ColorTint;

    // Mode 1: Solid Unlit Geometry
    if (RenderMode == 1)
    {
        output.Color = baseColor;
        return output;
    }

    // Mode 0: Dynamic Multi-Light Shading
    float3 worldPos = mul(Model, float4(input.Position, 1.0)).xyz;
    float3 n = normalize(mul((float3x3)Model, input.Normal));
    float3 v = normalize(CameraPos - worldPos);

    // Start with Ambient light
    float3 totalLight = AmbientColor.rgb * AmbientColor.a;

    // Accumulate all active lights (Sun, Torches, Spotlights)
    for (int i = 0; i < NumLights; i++)
    {
        float3 lColor = Lights[i].Color.rgb;
        float lInt = Lights[i].Intensity;

        // 0: DIRECTIONAL SUNLIGHT
        if (Lights[i].Type == 0)
        {
            float3 l = normalize(-Lights[i].Direction);
            float diff = max(dot(n, l), 0.0);

            float3 h = normalize(l + v);
            float spec = pow(max(dot(n, h), 0.0), 16.0);

            totalLight += (diff * lColor + spec * 0.3 * lColor) * lInt;
        }
        // 1: POINT LIGHT (Torches, Lanterns, Fire)
        else if (Lights[i].Type == 1)
        {
            float3 toLight = Lights[i].Position - worldPos;
            float dist = length(toLight);
            if (dist < Lights[i].Range)
            {
                float3 l = toLight / dist;
                float diff = max(dot(n, l), 0.0);

                // Smooth quadratic attenuation falloff
                float atten = clamp(1.0 - (dist / Lights[i].Range), 0.0, 1.0);
                atten = atten * atten;

                float3 h = normalize(l + v);
                float spec = pow(max(dot(n, h), 0.0), 16.0);

                totalLight += (diff * lColor + spec * 0.35 * lColor) * (lInt * atten);
            }
        }
        // 2: SPOT LIGHT
        else if (Lights[i].Type == 2)
        {
            float3 toLight = Lights[i].Position - worldPos;
            float dist = length(toLight);
            if (dist < Lights[i].Range)
            {
                float3 l = toLight / dist;
                float cosAngle = dot(-l, normalize(Lights[i].Direction));
                if (cosAngle > Lights[i].SpotAngle)
                {
                    float spotFactor = clamp((cosAngle - Lights[i].SpotAngle) / (1.0 - Lights[i].SpotAngle), 0.0, 1.0);
                    float atten = clamp(1.0 - (dist / Lights[i].Range), 0.0, 1.0);
                    atten = atten * atten * spotFactor;

                    float diff = max(dot(n, l), 0.0);
                    float3 h = normalize(l + v);
                    float spec = pow(max(dot(n, h), 0.0), 16.0);

                    totalLight += (diff * lColor + spec * 0.35 * lColor) * (lInt * atten);
                }
            }
        }
    }

    output.Color = float4(baseColor.rgb * totalLight, baseColor.a);
    return output;
}

Texture2D    AlbedoTexture : register(t0, space2);
SamplerState AlbedoSampler : register(s0, space2);

float4 PS3D(PS3DInput input) : SV_Target
{
    if (input.Color.a < 0.0)
    {
        return float4(0.2, 1.0, 0.5, 1.0); // Wireframe
    }

    float4 texColor = AlbedoTexture.Sample(AlbedoSampler, input.UV);
    return float4(texColor.rgb * input.Color.rgb, texColor.a * input.Color.a);
}
`

type Render3D struct {
	gpuDev        *sdl.GPUDevice
	pipeSolid     *sdl.GPUGraphicsPipeline
	pipeWireframe *sdl.GPUGraphicsPipeline
	depthTex      *sdl.GPUTexture
	defaultWhite  *sdl.GPUTexture
	sampler3D     *sdl.GPUSampler
	scrW, scrH    int
	RotX, RotY    float32
	CamDist       float32
	Mode          RenderMode3D
	AmbientColor  sdl.Color
}

var R3D *Render3D

type BBox2D struct {
	MinX, MaxX float32
	MinY, MaxY float32
	Points     [8]sdl.FPoint
	Valid      bool
}

var unitBoxCorners = [8]mgl32.Vec3{
	{-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1},
	{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1},
}

var boxEdges = [12][2]int{
	{0, 1}, {1, 2}, {2, 3}, {3, 0},
	{4, 5}, {5, 6}, {6, 7}, {7, 4},
	{0, 4}, {1, 5}, {2, 6}, {3, 7},
}

func InitRender3D(gpuDev *sdl.GPUDevice, win *sdl.Window, scrW, scrH int) (*Render3D, error) {
	vShader, err := CompileShader(gpuDev, shader3DHLSL, "VS3D", shadercross.SHADERSTAGE_VERTEX, 0, 1)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(vShader)

	pShader, err := CompileShader(gpuDev, shader3DHLSL, "PS3D", shadercross.SHADERSTAGE_FRAGMENT, 1, 0)
	if err != nil {
		return nil, err
	}
	defer gpuDev.ReleaseShader(pShader)

	depthTex, err := gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
		Type:              sdl.GPU_TEXTURETYPE_2D,
		Format:            sdl.GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT,
		Width:             uint32(scrW),
		Height:            uint32(scrH),
		LayerCountOrDepth: 1,
		NumLevels:         1,
		Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET,
	})
	if err != nil {
		return nil, err
	}

	defaultWhite, _ := gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
		Type:              sdl.GPU_TEXTURETYPE_2D,
		Format:            sdl.GPU_TEXTUREFORMAT_R8G8B8A8_UNORM,
		Width:             1,
		Height:            1,
		LayerCountOrDepth: 1,
		NumLevels:         1,
		Usage:             sdl.GPU_TEXTUREUSAGE_SAMPLER,
	})
	{
		tBuf, _ := gpuDev.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD, Size: 4})
		mapped, _ := gpuDev.MapTransferBuffer(tBuf, false)
		copy(unsafe.Slice((*byte)(unsafe.Pointer(mapped)), 4), []byte{255, 255, 255, 255})
		gpuDev.UnmapTransferBuffer(tBuf)
		cmd, _ := gpuDev.AcquireCommandBuffer()
		cp := cmd.BeginCopyPass()
		cp.UploadToGPUTexture(&sdl.GPUTextureTransferInfo{TransferBuffer: tBuf}, &sdl.GPUTextureRegion{Texture: defaultWhite, W: 1, H: 1, D: 1}, false)
		cp.End()
		cmd.Submit()
		gpuDev.ReleaseTransferBuffer(tBuf)
	}

	sampler3D, err := gpuDev.CreateSampler(&sdl.GPUSamplerCreateInfo{
		MinFilter:    sdl.GPU_FILTER_LINEAR,
		MagFilter:    sdl.GPU_FILTER_LINEAR,
		MipmapMode:   sdl.GPU_SAMPLERMIPMAPMODE_LINEAR,
		AddressModeU: sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
		AddressModeV: sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
		AddressModeW: sdl.GPU_SAMPLERADDRESSMODE_REPEAT,
	})
	if err != nil {
		return nil, err
	}

	attrs3D := []sdl.GPUVertexAttribute{
		{Location: 0, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT3, Offset: 0},
		{Location: 1, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT3, Offset: uint32(unsafe.Offsetof(Vertex3D{}.NX))},
		{Location: 2, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT2, Offset: uint32(unsafe.Offsetof(Vertex3D{}.U))},
		{Location: 3, BufferSlot: 0, Format: sdl.GPU_VERTEXELEMENTFORMAT_FLOAT4, Offset: uint32(unsafe.Offsetof(Vertex3D{}.R))},
	}
	desc3D := []sdl.GPUVertexBufferDescription{
		{Slot: 0, Pitch: uint32(unsafe.Sizeof(Vertex3D{})), InputRate: sdl.GPU_VERTEXINPUTRATE_VERTEX},
	}

	depthState := sdl.GPUDepthStencilState{
		EnableDepthTest:  true,
		EnableDepthWrite: true,
		CompareOp:        sdl.GPU_COMPAREOP_LESS_OR_EQUAL,
	}

	targetInfo := sdl.GPUGraphicsPipelineTargetInfo{
		ColorTargetDescriptions: []sdl.GPUColorTargetDescription{
			{Format: gpuDev.SwapchainTextureFormat(win)},
		},
		HasDepthStencilTarget: true,
		DepthStencilFormat:    sdl.GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT,
	}

	pipeSolid, err := gpuDev.CreateGraphicsPipeline(&sdl.GPUGraphicsPipelineCreateInfo{
		VertexShader:      vShader,
		FragmentShader:    pShader,
		VertexInputState:  sdl.GPUVertexInputState{VertexAttributes: attrs3D, VertexBufferDescriptions: desc3D},
		PrimitiveType:     sdl.GPU_PRIMITIVETYPE_TRIANGLELIST,
		DepthStencilState: depthState,
		TargetInfo:        targetInfo,
		RasterizerState:   sdl.GPURasterizerState{CullMode: sdl.GPU_CULLMODE_BACK, FrontFace: sdl.GPU_FRONTFACE_COUNTER_CLOCKWISE},
	})
	if err != nil {
		return nil, err
	}

	pipeWireframe, err := gpuDev.CreateGraphicsPipeline(&sdl.GPUGraphicsPipelineCreateInfo{
		VertexShader:      vShader,
		FragmentShader:    pShader,
		VertexInputState:  sdl.GPUVertexInputState{VertexAttributes: attrs3D, VertexBufferDescriptions: desc3D},
		PrimitiveType:     sdl.GPU_PRIMITIVETYPE_LINELIST,
		DepthStencilState: depthState,
		TargetInfo:        targetInfo,
		RasterizerState:   sdl.GPURasterizerState{CullMode: sdl.GPU_CULLMODE_NONE},
	})
	if err != nil {
		return nil, err
	}

	r3d := &Render3D{
		gpuDev:        gpuDev,
		pipeSolid:     pipeSolid,
		pipeWireframe: pipeWireframe,
		depthTex:      depthTex,
		defaultWhite:  defaultWhite,
		sampler3D:     sampler3D,
		scrW:          scrW,
		scrH:          scrH,
		RotX:          0.2,
		RotY:          0.0,
		CamDist:       4.0,
		Mode:          ModeShaded,
		AmbientColor:  Hex(0x353540), // Subtle ambient shadow color
	}

	R3D = r3d
	return r3d, nil
}

func (r3d *Render3D) Destroy() {
	if r3d == nil {
		return
	}
	r3d.gpuDev.ReleaseGraphicsPipeline(r3d.pipeSolid)
	r3d.gpuDev.ReleaseGraphicsPipeline(r3d.pipeWireframe)
	r3d.gpuDev.ReleaseTexture(r3d.depthTex)
	r3d.gpuDev.ReleaseTexture(r3d.defaultWhite)
	r3d.gpuDev.ReleaseSampler(r3d.sampler3D)
}

func (r3d *Render3D) UploadModel(m *Model3D) (*GPUModel, error) {
	if m == nil {
		return nil, nil
	}

	gpuModel := &GPUModel{}

	for _, mesh := range m.Meshes {
		if len(mesh.Vertices) == 0 {
			continue
		}

		vSize := uint32(len(mesh.Vertices) * int(unsafe.Sizeof(Vertex3D{})))
		triSize := uint32(len(mesh.TriIndices) * int(unsafe.Sizeof(uint32(0))))
		lineSize := uint32(len(mesh.LineIndices) * int(unsafe.Sizeof(uint32(0))))
		totalBufferSize := vSize + triSize + lineSize

		vBuf, _ := r3d.gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_VERTEX, Size: vSize})
		triBuf, _ := r3d.gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_INDEX, Size: triSize})
		lineBuf, _ := r3d.gpuDev.CreateBuffer(&sdl.GPUBufferCreateInfo{Usage: sdl.GPU_BUFFERUSAGE_INDEX, Size: lineSize})

		tBuf, _ := r3d.gpuDev.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD, Size: totalBufferSize})
		mapped, _ := r3d.gpuDev.MapTransferBuffer(tBuf, false)
		memPtr := (*byte)(unsafe.Pointer(mapped))

		copy(unsafe.Slice(memPtr, vSize), unsafe.Slice((*byte)(unsafe.Pointer(&mesh.Vertices[0])), vSize))
		triDest := (*byte)(unsafe.Add(unsafe.Pointer(memPtr), vSize))
		copy(unsafe.Slice(triDest, triSize), unsafe.Slice((*byte)(unsafe.Pointer(&mesh.TriIndices[0])), triSize))
		lineDest := (*byte)(unsafe.Add(unsafe.Pointer(memPtr), vSize+triSize))
		copy(unsafe.Slice(lineDest, lineSize), unsafe.Slice((*byte)(unsafe.Pointer(&mesh.LineIndices[0])), lineSize))

		r3d.gpuDev.UnmapTransferBuffer(tBuf)

		cmd, _ := r3d.gpuDev.AcquireCommandBuffer()
		cp := cmd.BeginCopyPass()
		cp.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: 0}, &sdl.GPUBufferRegion{Buffer: vBuf, Offset: 0, Size: vSize}, false)
		cp.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: vSize}, &sdl.GPUBufferRegion{Buffer: triBuf, Offset: 0, Size: triSize}, false)
		cp.UploadToGPUBuffer(&sdl.GPUTransferBufferLocation{TransferBuffer: tBuf, Offset: vSize + triSize}, &sdl.GPUBufferRegion{Buffer: lineBuf, Offset: 0, Size: lineSize}, false)

		var meshGPUTexture *sdl.GPUTexture = r3d.defaultWhite
		if mesh.Texture != nil {
			texW := uint32(mesh.Texture.Width)
			texH := uint32(mesh.Texture.Height)
			texBytes := uint32(len(mesh.Texture.Pixels))

			tex, err := r3d.gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
				Type:              sdl.GPU_TEXTURETYPE_2D,
				Format:            sdl.GPU_TEXTUREFORMAT_R8G8B8A8_UNORM,
				Width:             texW,
				Height:            texH,
				LayerCountOrDepth: 1,
				NumLevels:         1,
				Usage:             sdl.GPU_TEXTUREUSAGE_SAMPLER,
			})
			if err == nil {
				meshGPUTexture = tex
				texTransferBuf, _ := r3d.gpuDev.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD, Size: texBytes})
				mappedTex, _ := r3d.gpuDev.MapTransferBuffer(texTransferBuf, false)
				copy(unsafe.Slice((*byte)(unsafe.Pointer(mappedTex)), texBytes), mesh.Texture.Pixels)
				r3d.gpuDev.UnmapTransferBuffer(texTransferBuf)

				cp.UploadToGPUTexture(&sdl.GPUTextureTransferInfo{TransferBuffer: texTransferBuf}, &sdl.GPUTextureRegion{Texture: tex, W: texW, H: texH, D: 1}, false)
				defer r3d.gpuDev.ReleaseTransferBuffer(texTransferBuf)
			}
		}

		cp.End()
		cmd.Submit()
		r3d.gpuDev.ReleaseTransferBuffer(tBuf)

		gpuModel.Meshes = append(gpuModel.Meshes, GPUMesh{
			VertexBuffer: vBuf,
			TriBuffer:    triBuf,
			LineBuffer:   lineBuf,
			Texture:      meshGPUTexture,
			NumTriangles: uint32(len(mesh.TriIndices)),
			NumLines:     uint32(len(mesh.LineIndices)),
		})
	}

	return gpuModel, nil
}

func (r3d *Render3D) ReleaseGPUModel(model *GPUModel) {
	if model == nil {
		return
	}
	for _, m := range model.Meshes {
		r3d.gpuDev.ReleaseBuffer(m.VertexBuffer)
		r3d.gpuDev.ReleaseBuffer(m.TriBuffer)
		r3d.gpuDev.ReleaseBuffer(m.LineBuffer)
		if m.Texture != nil && m.Texture != r3d.defaultWhite {
			r3d.gpuDev.ReleaseTexture(m.Texture)
		}
	}
}

func (r3d *Render3D) DrawObj(cmd *sdl.GPUCommandBuffer, renderPass *sdl.GPURenderPass, obj *OBJ3D) {
	if obj == nil || obj.mod == nil {
		return
	}

	var customTex *sdl.GPUTexture
	if obj.CustomTex != nil {
		customTex = obj.CustomTex.GPU
	}

	uvScale := obj.UVScale
	if uvScale <= 0 {
		uvScale = 1.0
	}

	r3d.DrawAt(cmd, renderPass, obj.mod, obj.WorldX, obj.WorldY, obj.WorldZ, obj.RotX, obj.RotY, obj.RotZ, obj.Scale, obj.Color, uvScale, customTex)
}

func (r3d *Render3D) DrawAt(cmd *sdl.GPUCommandBuffer, renderPass *sdl.GPURenderPass, model *GPUModel, worldX, worldY, worldZ, rotX, rotY, rotZ, scale float32, tint sdl.Color, uvScale float32, customTex *sdl.GPUTexture) {
	if model == nil || len(model.Meshes) == 0 {
		return
	}

	aspect := float32(r3d.scrW) / float32(r3d.scrH)
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	camPos := mgl32.Vec3{0, 0, r3d.CamDist}
	view := mgl32.LookAtV(camPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	translation := mgl32.Translate3D(worldX, worldY, worldZ)
	scaling := mgl32.Scale3D(scale, scale, scale)
	rotation := mgl32.HomogRotate3DX(rotX).Mul4(mgl32.HomogRotate3DY(rotY)).Mul4(mgl32.HomogRotate3DZ(rotZ))
	modelMatrix := translation.Mul4(rotation).Mul4(scaling)
	mvp := proj.Mul4(view).Mul4(modelMatrix)

	if tint.A == 0 {
		tint = Col.White
	}
	rf, gf, bf, af := color2float(tint)

	if uvScale <= 0 {
		uvScale = 1.0
	}

	// Pack Dynamic Scene Lights
	gpuLights, numLights := ConvertToGPULights()

	ambR, ambG, ambB, _ := color2float(r3d.AmbientColor)

	uniforms := Uniform3D{
		MVP:          mvp,
		Model:        modelMatrix,
		CameraPos:    [3]float32{camPos.X(), camPos.Y(), camPos.Z()},
		RenderMode:   int32(r3d.Mode),
		AmbientColor: [4]float32{ambR, ambG, ambB, 0.3},
		ColorTint:    [4]float32{rf, gf, bf, af},
		UVScale:      [2]float32{uvScale, uvScale},
		NumLights:    numLights,
		Lights:       gpuLights,
	}

	uBytes := unsafe.Slice((*byte)(unsafe.Pointer(&uniforms)), unsafe.Sizeof(uniforms))
	cmd.PushVertexUniformData(0, uBytes)

	if r3d.Mode == ModeWireframe {
		renderPass.BindGraphicsPipeline(r3d.pipeWireframe)
		for _, m := range model.Meshes {
			renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: m.VertexBuffer, Offset: 0}})
			renderPass.BindIndexBuffer(&sdl.GPUBufferBinding{Buffer: m.LineBuffer, Offset: 0}, sdl.GPU_INDEXELEMENTSIZE_32BIT)
			renderPass.DrawIndexedPrimitives(m.NumLines, 1, 0, 0, 0)
		}
	} else {
		renderPass.BindGraphicsPipeline(r3d.pipeSolid)
		for _, m := range model.Meshes {
			texToBind := m.Texture
			if customTex != nil {
				texToBind = customTex
			}

			renderPass.BindFragmentSamplers([]sdl.GPUTextureSamplerBinding{
				{Texture: texToBind, Sampler: r3d.sampler3D},
			})
			renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: m.VertexBuffer, Offset: 0}})
			renderPass.BindIndexBuffer(&sdl.GPUBufferBinding{Buffer: m.TriBuffer, Offset: 0}, sdl.GPU_INDEXELEMENTSIZE_32BIT)
			renderPass.DrawIndexedPrimitives(m.NumTriangles, 1, 0, 0, 0)
		}
	}
}

func (r3d *Render3D) Draw(cmd *sdl.GPUCommandBuffer, renderPass *sdl.GPURenderPass, model *GPUModel) {
	r3d.DrawAt(cmd, renderPass, model, 0, 0, 0, r3d.RotX, r3d.RotY, 0, 1.0, Col.White, 1.0, nil)
}

func (r3d *Render3D) DrawCustom(cmd *sdl.GPUCommandBuffer, renderPass *sdl.GPURenderPass, model *GPUModel, rotX, rotY, camDist, aspect float32) {
	if model == nil || len(model.Meshes) == 0 {
		return
	}

	proj := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	camPos := mgl32.Vec3{0, 0, camDist}
	view := mgl32.LookAtV(camPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	modelRot := mgl32.HomogRotate3DX(rotX).Mul4(mgl32.HomogRotate3DY(rotY))
	mvp := proj.Mul4(view).Mul4(modelRot)

	rf, gf, bf, af := color2float(Col.White)
	uniforms := Uniform3D{
		MVP:          mvp,
		Model:        modelRot,
		CameraPos:    [3]float32{camPos.X(), camPos.Y(), camPos.Z()},
		RenderMode:   int32(ModeShaded),
		AmbientColor: [4]float32{0.4, 0.4, 0.45, 1.0},
		ColorTint:    [4]float32{rf, gf, bf, af},
		UVScale:      [2]float32{1.0, 1.0},
		NumLights:    1,
		Lights: [8]GPULight{
			{
				Type:      0, // Directional Sun
				Direction: [3]float32{0.577, 0.707, -0.577},
				Intensity: 1.0,
				Color:     [4]float32{1, 1, 1, 1},
			},
		},
	}

	uBytes := unsafe.Slice((*byte)(unsafe.Pointer(&uniforms)), unsafe.Sizeof(uniforms))
	cmd.PushVertexUniformData(0, uBytes)

	renderPass.BindGraphicsPipeline(r3d.pipeSolid)
	for _, m := range model.Meshes {
		renderPass.BindFragmentSamplers([]sdl.GPUTextureSamplerBinding{
			{Texture: m.Texture, Sampler: r3d.sampler3D},
		})
		renderPass.BindVertexBuffers([]sdl.GPUBufferBinding{{Buffer: m.VertexBuffer, Offset: 0}})
		renderPass.BindIndexBuffer(&sdl.GPUBufferBinding{Buffer: m.TriBuffer, Offset: 0}, sdl.GPU_INDEXELEMENTSIZE_32BIT)
		renderPass.DrawIndexedPrimitives(m.NumTriangles, 1, 0, 0, 0)
	}
}

func (r3d *Render3D) GetObjectBoundingBox(obj *OBJ3D) BBox2D {
	aspect := float32(r3d.scrW) / float32(r3d.scrH)
	proj := mgl32.Perspective(mgl32.DegToRad(45.0), aspect, 0.1, 100.0)
	camPos := mgl32.Vec3{0, 0, r3d.CamDist}
	view := mgl32.LookAtV(camPos, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})

	translation := mgl32.Translate3D(obj.WorldX, obj.WorldY, obj.WorldZ)
	scaling := mgl32.Scale3D(obj.Scale, obj.Scale, obj.Scale)
	rotation := mgl32.HomogRotate3DX(obj.RotX).Mul4(mgl32.HomogRotate3DY(obj.RotY)).Mul4(mgl32.HomogRotate3DZ(obj.RotZ))

	modelMatrix := translation.Mul4(rotation).Mul4(scaling)
	mvp := proj.Mul4(view).Mul4(modelMatrix)

	var bbox BBox2D

	for i, corner := range unitBoxCorners {
		clip := mvp.Mul4x1(mgl32.Vec4{corner.X(), corner.Y(), corner.Z(), 1.0})
		if clip.W() <= 0.1 {
			return BBox2D{Valid: false}
		}

		ndcX := clip.X() / clip.W()
		ndcY := clip.Y() / clip.W()

		sx := (ndcX + 1.0) * 0.5 * float32(r3d.scrW)
		sy := (1.0 - ndcY) * 0.5 * float32(r3d.scrH)

		bbox.Points[i] = sdl.FPoint{X: sx, Y: sy}
	}

	bbox.MinX, bbox.MaxX = bbox.Points[0].X, bbox.Points[0].X
	bbox.MinY, bbox.MaxY = bbox.Points[0].Y, bbox.Points[0].Y

	for i := 1; i < 8; i++ {
		px := bbox.Points[i].X
		py := bbox.Points[i].Y

		if px < bbox.MinX {
			bbox.MinX = px
		}
		if px > bbox.MaxX {
			bbox.MaxX = px
		}
		if py < bbox.MinY {
			bbox.MinY = py
		}
		if py > bbox.MaxY {
			bbox.MaxY = py
		}
	}

	bbox.Valid = true
	return bbox
}

func Draw3DBoundingBox(bbox BBox2D, col sdl.Color, thickness float32) {
	if !bbox.Valid {
		return
	}
	for _, edge := range boxEdges {
		p1 := bbox.Points[edge[0]]
		p2 := bbox.Points[edge[1]]
		dliw(p1.X, p1.Y, p2.X, p2.Y, thickness, col)
	}
}

func GenerateModelThumbnail(gpuDev *sdl.GPUDevice, win *sdl.Window, r3d *Render3D, model *GPUModel, size int) (*Texture, error) {
	if model == nil || len(model.Meshes) == 0 {
		return nil, fmt.Errorf("cannot generate thumbnail for empty model")
	}

	swapchainFormat := gpuDev.SwapchainTextureFormat(win)
	uSize := uint32(size)

	colorTex, err := gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
		Type:              sdl.GPU_TEXTURETYPE_2D,
		Format:            swapchainFormat,
		Width:             uSize,
		Height:            uSize,
		LayerCountOrDepth: 1,
		NumLevels:         1,
		Usage:             sdl.GPU_TEXTUREUSAGE_COLOR_TARGET | sdl.GPU_TEXTUREUSAGE_SAMPLER,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating thumbnail color texture: %w", err)
	}

	depthTex, err := gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
		Type:              sdl.GPU_TEXTURETYPE_2D,
		Format:            sdl.GPU_TEXTUREFORMAT_D24_UNORM_S8_UINT,
		Width:             uSize,
		Height:            uSize,
		LayerCountOrDepth: 1,
		NumLevels:         1,
		Usage:             sdl.GPU_TEXTUREUSAGE_DEPTH_STENCIL_TARGET,
	})
	if err != nil {
		gpuDev.ReleaseTexture(colorTex)
		return nil, fmt.Errorf("failed creating thumbnail depth texture: %w", err)
	}
	defer gpuDev.ReleaseTexture(depthTex)

	cmd, err := gpuDev.AcquireCommandBuffer()
	if err != nil {
		gpuDev.ReleaseTexture(colorTex)
		return nil, err
	}

	colorTarget := sdl.GPUColorTargetInfo{
		Texture:    colorTex,
		LoadOp:     sdl.GPU_LOADOP_CLEAR,
		StoreOp:    sdl.GPU_STOREOP_STORE,
		ClearColor: sdl.FColor{R: 0, G: 0, B: 0, A: 0},
	}

	depthTarget := sdl.GPUDepthStencilTargetInfo{
		Texture:    depthTex,
		LoadOp:     sdl.GPU_LOADOP_CLEAR,
		StoreOp:    sdl.GPU_STOREOP_DONT_CARE,
		ClearDepth: 1.0,
		Cycle:      true,
	}

	renderPass := cmd.BeginRenderPass([]sdl.GPUColorTargetInfo{colorTarget}, &depthTarget)
	r3d.DrawCustom(cmd, renderPass, model, 0.35, 0.65, 4.2, 1.0)
	renderPass.End()
	cmd.Submit()

	return &Texture{
		GPU:  colorTex,
		W:    float32(size),
		H:    float32(size),
		Path: "thumbnail",
	}, nil
}

func (r3d *Render3D) IsObjectVisible(obj *OBJ3D) bool {
	if obj == nil || obj.mod == nil {
		return false
	}

	bbox := r3d.GetObjectBoundingBox(obj)
	if !bbox.Valid {
		return false
	}

	scrW := float32(SET.SCRW)
	scrH := float32(SET.SCRH)

	if bbox.MaxX < 0 || bbox.MinX > scrW || bbox.MaxY < 0 || bbox.MinY > scrH {
		return false
	}

	return true
}
