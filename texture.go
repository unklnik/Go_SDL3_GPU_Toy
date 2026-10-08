package main

import (
	"fmt"
	"image"
	_ "image/png"
	"os"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
)

type Texture struct {
	GPU  *sdl.GPUTexture
	W, H float32
	Path string
}

// MARK: UTILS

// LoadTexture loads any PNG file directly into GPU VRAM
func LoadTexture(gpuDev *sdl.GPUDevice, path string) (*Texture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image %s: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image %s: %w", path, err)
	}

	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()
	totalBytes := uint32(w * h * 4)

	// 1. Convert Go image to raw RGBA8888 byte slice
	rgbaBytes := make([]byte, totalBytes)
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			idx := (y*w + x) * 4
			rgbaBytes[idx+0] = byte(r >> 8)
			rgbaBytes[idx+1] = byte(g >> 8)
			rgbaBytes[idx+2] = byte(b >> 8)
			rgbaBytes[idx+3] = byte(a >> 8)
		}
	}

	// 2. Create GPU VRAM Texture
	gpuTex, err := gpuDev.CreateTexture(&sdl.GPUTextureCreateInfo{
		Type:              sdl.GPU_TEXTURETYPE_2D,
		Format:            sdl.GPU_TEXTUREFORMAT_R8G8B8A8_UNORM,
		Width:             uint32(w),
		Height:            uint32(h),
		LayerCountOrDepth: 1,
		NumLevels:         1,
		Usage:             sdl.GPU_TEXTUREUSAGE_SAMPLER,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU texture: %w", err)
	}

	// 3. Upload to GPU via TransferBuffer + CopyPass
	tBuf, err := gpuDev.CreateTransferBuffer(&sdl.GPUTransferBufferCreateInfo{
		Usage: sdl.GPU_TRANSFERBUFFERUSAGE_UPLOAD,
		Size:  totalBytes,
	})
	if err != nil {
		gpuDev.ReleaseTexture(gpuTex)
		return nil, err
	}

	mapped, _ := gpuDev.MapTransferBuffer(tBuf, false)
	copy(unsafe.Slice((*byte)(unsafe.Pointer(mapped)), totalBytes), rgbaBytes)
	gpuDev.UnmapTransferBuffer(tBuf)

	cmd, _ := gpuDev.AcquireCommandBuffer()
	cp := cmd.BeginCopyPass()
	cp.UploadToGPUTexture(
		&sdl.GPUTextureTransferInfo{TransferBuffer: tBuf},
		&sdl.GPUTextureRegion{Texture: gpuTex, W: uint32(w), H: uint32(h), D: 1},
		false,
	)
	cp.End()
	cmd.Submit()
	gpuDev.ReleaseTransferBuffer(tBuf)

	return &Texture{
		GPU:  gpuTex,
		W:    float32(w),
		H:    float32(h),
		Path: path,
	}, nil
}

func (t *Texture) Destroy(gpuDev *sdl.GPUDevice) {
	if t != nil && t.GPU != nil {
		gpuDev.ReleaseTexture(t.GPU)
		t.GPU = nil
	}
}
