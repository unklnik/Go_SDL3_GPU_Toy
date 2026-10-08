package main

import (
	"math"

	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	SceneLights, LightsOff []*Light3D
	selectedLightID        int = -1
	lightsPosVis           bool
	sceneLightsActiv       = true
)

type LightType int32

const (
	LightDirectional LightType = 0 // Global Sun/Moon
	LightPoint       LightType = 1 // Torches, lamps, fire
	LightSpot        LightType = 2 // Flashlights, cones
)

// GPULight matches HLSL memory layout (exactly 64 bytes, 16-byte aligned)
type GPULight struct {
	Position  [3]float32
	Type      int32
	Direction [3]float32
	Intensity float32
	Color     [4]float32
	Range     float32
	SpotAngle float32
	Padding   [2]float32
}

type Light3D struct {
	Name      string
	Type      LightType
	Color     sdl.Color
	Intensity float32
	Range     float32 // Distance light travels
	SpotAngle float32 // Cone cutoff for spotlights (cosine)

	// 3D Transform
	WorldX, WorldY, WorldZ float32
	DirX, DirY, DirZ       float32

	// Torch / Fire Flicker
	Flicker       bool
	FlickerSpeed  float32
	FlickerAmount float32
	flickerTimer  float32
	CurIntensity  float32 // Evaluated per frame

	// Attached Object ID (-1 if standalone)
	AttachedObjID int
}

// UpdateLights updates torch flicker and moves attached lights
func UpdateLights(dt float32) {
	for _, l := range SceneLights {
		// 1. Update Torch / Fire Flicker
		l.CurIntensity = l.Intensity
		if l.Flicker {
			l.flickerTimer += dt * l.FlickerSpeed
			// Multi-octave organic turbulence
			noise := float32(math.Sin(float64(l.flickerTimer)*7.3))*0.5 +
				float32(math.Sin(float64(l.flickerTimer)*14.1))*0.3 +
				float32(math.Sin(float64(l.flickerTimer)*28.7))*0.2
			factor := 1.0 + (noise * l.FlickerAmount)
			if factor < 0 {
				factor = 0
			}
			l.CurIntensity = l.Intensity * factor
		}

		// 2. Sync position with attached object (if attached to an OBJ3D)
		if l.AttachedObjID >= 0 && l.AttachedObjID < len(lev3d.o3d) {
			obj := &lev3d.o3d[l.AttachedObjID]
			l.WorldX = obj.WorldX
			l.WorldY = obj.WorldY + (obj.Scale * 1.2) // Float above object
			l.WorldZ = obj.WorldZ
		}
	}
}

// ConvertToGPULights packs active lights for the shader uniform buffer (up to 8 lights)
func ConvertToGPULights() ([8]GPULight, int32) {
	var gpuLights [8]GPULight
	count := int32(0)

	for _, l := range SceneLights {
		if count >= 8 {
			break
		}

		rf, gf, bf, _ := color2float(l.Color)
		gpuLights[count] = GPULight{
			Position:  [3]float32{l.WorldX, l.WorldY, l.WorldZ},
			Type:      int32(l.Type),
			Direction: [3]float32{l.DirX, l.DirY, l.DirZ},
			Intensity: l.CurIntensity,
			Color:     [4]float32{rf, gf, bf, 1.0},
			Range:     l.Range,
			SpotAngle: l.SpotAngle,
		}
		count++
	}

	return gpuLights, count
}

// Default Scene Lighting (1 Sun + 1 Warm Torch)
func InitDefaultLights() {
	SceneLights = []*Light3D{
		// 1. Sun (Directional)
		{
			Name:      "Sunlight",
			Type:      LightDirectional,
			Color:     Col.White,
			Intensity: 1.0,
			DirX:      0.577, DirY: 0.707, DirZ: -0.577,
			AttachedObjID: -1,
		},
		// 2. Dungeon Torch (Point Light with Flicker)
		{
			Name:      "Dungeon Torch",
			Type:      LightPoint,
			Color:     Col.GoldWarm,
			Intensity: 3.5,
			Range:     8.0,
			WorldX:    -2.0, WorldY: 1.5, WorldZ: 0.0,
			Flicker:       true,
			FlickerSpeed:  4.0,
			FlickerAmount: 0.25,
			AttachedObjID: -1,
		},
	}
}
