package main

import (
	"fmt"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/shadercross"
)

// InitGPUDevice initializes the SDL_GPU device, claims the window, and sets VSync
func InitGPUDevice(win *sdl.Window) (*sdl.GPUDevice, error) {
	// 1. Load shadercross
	if err := shadercross.LoadLibrary(shadercross.Path()); err != nil {
		return nil, fmt.Errorf("failed to load shadercross library: %w", err)
	}

	// 2. Create GPU Device
	gpuDev, err := sdl.CreateGPUDevice(
		sdl.GPU_SHADERFORMAT_SPIRV|sdl.GPU_SHADERFORMAT_DXIL|sdl.GPU_SHADERFORMAT_MSL,
		false,
		"",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create GPU device: %w", err)
	}

	// 3. Claim Window & Set VSync
	if err := gpuDev.ClaimWindow(win); err != nil {
		gpuDev.Destroy()
		return nil, fmt.Errorf("failed to claim window: %w", err)
	}

	gpuDev.SetSwapchainParameters(win, sdl.GPU_SWAPCHAINCOMPOSITION_SDR, sdl.GPU_PRESENTMODE_VSYNC)

	return gpuDev, nil
}

// CompileShader compiles HLSL into a GPUShader (100% GC-safe)
func CompileShader(gpuDev *sdl.GPUDevice, source, entrypoint string, stage shadercross.ShaderStage, numSamplers, numUniformBuffers uint32) (*sdl.GPUShader, error) {
	spirvBytecode, err := shadercross.CompileSPIRVFromHLSL(&shadercross.HLSLInfo{
		Source:      source,
		Entrypoint:  entrypoint,
		ShaderStage: stage,
	})
	if err != nil {
		return nil, fmt.Errorf("HLSL->SPIRV error (%s): %w", entrypoint, err)
	}

	resourceInfo := &shadercross.GraphicsShaderResourceInfo{
		NumSamplers:        numSamplers,
		NumUniformBuffers:  numUniformBuffers,
		NumStorageBuffers:  0,
		NumStorageTextures: 0,
	}

	shader, err := shadercross.CompileGraphicsShaderFromSPIRV(
		gpuDev,
		&shadercross.SPIRVInfo{
			Bytecode:    spirvBytecode,
			Entrypoint:  entrypoint,
			ShaderStage: stage,
		},
		resourceInfo,
		0,
	)
	if err != nil {
		return nil, fmt.Errorf("SPIRV->GPUShader error (%s): %w", entrypoint, err)
	}

	return shader, nil
}
