package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
)

var (
	prim3DwinOBJ []OBJ3D
)

// mPrim3D: Make default primitive 3D shapes & thumbnails
func mPrim3D() {
	m := NewCubeModel()
	siz := thumbWdef - int(un8th)
	gpuM, err := r3d.UploadModel(m)
	if err != nil {
		log.Printf("[WARNING] Could not upload default 3D primitive model: %v", err)
	}
	thumb, err := GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Cube", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "cube"}})

	m6 := NewPlaneModel()
	gpuM, err = r3d.UploadModel(m6)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Plane", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "plane"}})

	m2 := NewSphereModel(24, 32)
	gpuM, err = r3d.UploadModel(m2)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Sphere", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "sphere"}})

	m3 := NewPyramidModel()
	gpuM, err = r3d.UploadModel(m3)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Pyramid", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "pyramid"}})

	m4 := NewConeModel(24)
	gpuM, err = r3d.UploadModel(m4)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Cone", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "cone"}})

	m5 := NewCylinderModel(24)
	gpuM, err = r3d.UploadModel(m5)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Cylinder", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "cylinder"}})

	m7 := NewTorusModel(16, 24)
	gpuM, err = r3d.UploadModel(m7)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Torus", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "torus"}})

	m8 := NewOctahedronModel()
	gpuM, err = r3d.UploadModel(m8)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Octahedron", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "octahedron"}})

	m9 := NewDodecahedronModel()
	gpuM, err = r3d.UploadModel(m9)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Dodecahedron", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "dodecahedron"}})

	m10 := NewPentagonalPyramidModel()
	gpuM, err = r3d.UploadModel(m10)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Pentagonal Pyramid", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "pentagonalpyramid"}})

	m11 := NewTetrahedronModel()
	gpuM, err = r3d.UploadModel(m11)
	if err != nil {
		log.Printf("[WARNING] Could not load default 3D primitive model: %v", err)
	}
	thumb, err = GenerateModelThumbnail(gpuDev, win, r3d, gpuM, siz)
	if err != nil {
		log.Printf("[WARNING] Could not generate default 3D primitive thumbnail: %v", err)
	}
	prim3DwinOBJ = append(prim3DwinOBJ, OBJ3D{mod: gpuM, nm: "Tetrahedron", thumb: thumb, Scale: 0.5, tags: []string{"primitive", "tetrahedron"}})
}

// SetTexture loads a 2D image (PNG/JPG) and applies it to all meshes in the model
func (m *Model3D) SetTexture(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open texture %s: %w", path, err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode texture %s: %w", path, err)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgba := make([]byte, w*h*4)

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			idx := (y*w + x) * 4
			rgba[idx+0] = byte(r >> 8)
			rgba[idx+1] = byte(g >> 8)
			rgba[idx+2] = byte(b >> 8)
			rgba[idx+3] = byte(a >> 8)
		}
	}

	texData := &TextureData{
		Width:  w,
		Height: h,
		Pixels: rgba,
	}

	for i := range m.Meshes {
		m.Meshes[i].Texture = texData
	}

	return nil
}

// =========================================================================
// 1. CUBE (Unit Cube [-1, 1], 6 faces with sharp normals and clean UVs)
// =========================================================================

func NewCubeModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	// 24 Vertices (4 per face for sharp per-face normals)
	verts := []Vertex3D{
		// Front (+Z)
		{X: -1, Y: -1, Z: 1, NX: 0, NY: 0, NZ: 1, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: 1, NX: 0, NY: 0, NZ: 1, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: 1, NX: 0, NY: 0, NZ: 1, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 1, Z: 1, NX: 0, NY: 0, NZ: 1, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Back (-Z)
		{X: 1, Y: -1, Z: -1, NX: 0, NY: 0, NZ: -1, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: -1, NX: 0, NY: 0, NZ: -1, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 1, Z: -1, NX: 0, NY: 0, NZ: -1, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: -1, NX: 0, NY: 0, NZ: -1, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Top (+Y)
		{X: -1, Y: 1, Z: 1, NX: 0, NY: 1, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: 1, NX: 0, NY: 1, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: -1, NX: 0, NY: 1, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 1, Z: -1, NX: 0, NY: 1, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Bottom (-Y)
		{X: -1, Y: -1, Z: -1, NX: 0, NY: -1, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: -1, NX: 0, NY: -1, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: 1, NX: 0, NY: -1, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: 1, NX: 0, NY: -1, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Right (+X)
		{X: 1, Y: -1, Z: 1, NX: 1, NY: 0, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: -1, NX: 1, NY: 0, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: -1, NX: 1, NY: 0, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 1, Z: 1, NX: 1, NY: 0, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Left (-X)
		{X: -1, Y: -1, Z: -1, NX: -1, NY: 0, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: 1, NX: -1, NY: 0, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 1, Z: 1, NX: -1, NY: 0, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 1, Z: -1, NX: -1, NY: 0, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
	}

	triIndices := make([]uint32, 0, 36)
	for i := uint32(0); i < 24; i += 4 {
		triIndices = append(triIndices, i, i+1, i+2, i, i+2, i+3)
	}

	// Clean 12 wireframe box edges
	lineIndices := []uint32{
		0, 1, 1, 2, 2, 3, 3, 0, // Front
		4, 5, 5, 6, 6, 7, 7, 4, // Back
		0, 5, 1, 4, 2, 7, 3, 6, // Pillars
	}

	return &Model3D{
		Name: "Cube",
		Meshes: []Mesh3D{
			{
				Name:        "CubeMesh",
				Vertices:    verts,
				TriIndices:  triIndices,
				LineIndices: lineIndices,
			},
		},
	}
}

// =========================================================================
// 2. SPHERE (UV Sphere with smooth normals and spherical UV wrap)
// =========================================================================

func NewSphereModel(rings, sectors int) *Model3D {
	if rings < 4 {
		rings = 16
	}
	if sectors < 4 {
		sectors = 24
	}

	white := [4]float32{1, 1, 1, 1}
	verts := make([]Vertex3D, 0, (rings+1)*(sectors+1))

	R := float32(1.0)
	for r := 0; r <= rings; r++ {
		theta := float32(r) * float32(math.Pi) / float32(rings)
		sinTheta := float32(math.Sin(float64(theta)))
		cosTheta := float32(math.Cos(float64(theta)))

		for s := 0; s <= sectors; s++ {
			phi := float32(s) * 2.0 * float32(math.Pi) / float32(sectors)
			sinPhi := float32(math.Sin(float64(phi)))
			cosPhi := float32(math.Cos(float64(phi)))

			x := cosPhi * sinTheta
			y := cosTheta
			z := sinPhi * sinTheta

			u := float32(s) / float32(sectors)
			v := float32(r) / float32(rings)

			verts = append(verts, Vertex3D{
				X: x * R, Y: y * R, Z: z * R,
				NX: x, NY: y, NZ: z,
				U: u, V: v,
				R: white[0], G: white[1], B: white[2], A: white[3],
			})
		}
	}

	triIndices := make([]uint32, 0, rings*sectors*6)
	lineIndices := make([]uint32, 0, rings*sectors*6)

	for r := 0; r < rings; r++ {
		for s := 0; s < sectors; s++ {
			first := uint32(r*(sectors+1) + s)
			second := first + uint32(sectors+1)

			triIndices = append(triIndices, first, second, first+1)
			triIndices = append(triIndices, second, second+1, first+1)

			// Wireframe latitude and longitude lines
			lineIndices = append(lineIndices, first, first+1, first, second)
		}
	}

	return &Model3D{
		Name: "Sphere",
		Meshes: []Mesh3D{
			{
				Name:        "SphereMesh",
				Vertices:    verts,
				TriIndices:  triIndices,
				LineIndices: lineIndices,
			},
		},
	}
}

// =========================================================================
// 3. CYLINDER (Top cap, bottom cap, and smooth circular wall)
// =========================================================================

func NewCylinderModel(segments int) *Model3D {
	if segments < 6 {
		segments = 24
	}
	white := [4]float32{1, 1, 1, 1}

	var verts []Vertex3D
	var triIndices []uint32
	var lineIndices []uint32

	radius := float32(1.0)
	halfH := float32(1.0)

	// --- Side Wall ---
	for i := 0; i <= segments; i++ {
		theta := float32(i) * 2.0 * float32(math.Pi) / float32(segments)
		cosT := float32(math.Cos(float64(theta)))
		sinT := float32(math.Sin(float64(theta)))
		u := float32(i) / float32(segments)

		// Bottom vertex
		verts = append(verts, Vertex3D{
			X: cosT * radius, Y: -halfH, Z: sinT * radius,
			NX: cosT, NY: 0, NZ: sinT,
			U: u, V: 1, R: white[0], G: white[1], B: white[2], A: white[3],
		})

		// Top vertex
		verts = append(verts, Vertex3D{
			X: cosT * radius, Y: halfH, Z: sinT * radius,
			NX: cosT, NY: 0, NZ: sinT,
			U: u, V: 0, R: white[0], G: white[1], B: white[2], A: white[3],
		})
	}

	for i := 0; i < segments; i++ {
		b1 := uint32(i * 2)
		t1 := b1 + 1
		b2 := b1 + 2
		t2 := b1 + 3

		triIndices = append(triIndices, b1, t1, b2)
		triIndices = append(triIndices, t1, t2, b2)

		lineIndices = append(lineIndices, b1, b2, t1, t2, b1, t1)
	}

	// --- Top Cap (+Y) ---
	topCenterIdx := uint32(len(verts))
	verts = append(verts, Vertex3D{X: 0, Y: halfH, Z: 0, NX: 0, NY: 1, NZ: 0, U: 0.5, V: 0.5, R: white[0], G: white[1], B: white[2], A: white[3]})

	topStartIdx := uint32(len(verts))
	for i := 0; i <= segments; i++ {
		theta := float32(i) * 2.0 * float32(math.Pi) / float32(segments)
		cosT := float32(math.Cos(float64(theta)))
		sinT := float32(math.Sin(float64(theta)))
		verts = append(verts, Vertex3D{
			X: cosT * radius, Y: halfH, Z: sinT * radius,
			NX: 0, NY: 1, NZ: 0,
			U: 0.5 + cosT*0.5, V: 0.5 + sinT*0.5, R: white[0], G: white[1], B: white[2], A: white[3],
		})
	}
	for i := 0; i < segments; i++ {
		triIndices = append(triIndices, topCenterIdx, topStartIdx+uint32(i), topStartIdx+uint32(i+1))
		lineIndices = append(lineIndices, topStartIdx+uint32(i), topStartIdx+uint32(i+1))
	}

	// --- Bottom Cap (-Y) ---
	botCenterIdx := uint32(len(verts))
	verts = append(verts, Vertex3D{X: 0, Y: -halfH, Z: 0, NX: 0, NY: -1, NZ: 0, U: 0.5, V: 0.5, R: white[0], G: white[1], B: white[2], A: white[3]})

	botStartIdx := uint32(len(verts))
	for i := 0; i <= segments; i++ {
		theta := float32(i) * 2.0 * float32(math.Pi) / float32(segments)
		cosT := float32(math.Cos(float64(theta)))
		sinT := float32(math.Sin(float64(theta)))
		verts = append(verts, Vertex3D{
			X: cosT * radius, Y: -halfH, Z: sinT * radius,
			NX: 0, NY: -1, NZ: 0,
			U: 0.5 + cosT*0.5, V: 0.5 + sinT*0.5, R: white[0], G: white[1], B: white[2], A: white[3],
		})
	}
	for i := 0; i < segments; i++ {
		triIndices = append(triIndices, botCenterIdx, botStartIdx+uint32(i+1), botStartIdx+uint32(i))
		lineIndices = append(lineIndices, botStartIdx+uint32(i), botStartIdx+uint32(i+1))
	}

	return &Model3D{
		Name: "Cylinder",
		Meshes: []Mesh3D{
			{Name: "CylinderMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 4. CONE (Circular base with a pointed apex)
// =========================================================================

func NewConeModel(segments int) *Model3D {
	if segments < 6 {
		segments = 24
	}
	white := [4]float32{1, 1, 1, 1}

	var verts []Vertex3D
	var triIndices []uint32
	var lineIndices []uint32

	radius := float32(1.0)
	halfH := float32(1.0)

	// Slant normal factor
	slope := radius / (2.0 * halfH)
	normLen := float32(math.Sqrt(float64(1.0 + slope*slope)))

	// --- Cone Sides ---
	for i := 0; i < segments; i++ {
		theta1 := float32(i) * 2.0 * float32(math.Pi) / float32(segments)
		theta2 := float32(i+1) * 2.0 * float32(math.Pi) / float32(segments)

		c1, s1 := float32(math.Cos(float64(theta1))), float32(math.Sin(float64(theta1)))
		c2, s2 := float32(math.Cos(float64(theta2))), float32(math.Sin(float64(theta2)))
		cm, sm := (c1+c2)/2, (s1+s2)/2

		idx := uint32(len(verts))

		// Apex
		verts = append(verts, Vertex3D{
			X: 0, Y: halfH, Z: 0,
			NX: cm / normLen, NY: slope / normLen, NZ: sm / normLen,
			U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3],
		})
		// Base Left
		verts = append(verts, Vertex3D{
			X: c1 * radius, Y: -halfH, Z: s1 * radius,
			NX: c1 / normLen, NY: slope / normLen, NZ: s1 / normLen,
			U: float32(i) / float32(segments), V: 1, R: white[0], G: white[1], B: white[2], A: white[3],
		})
		// Base Right
		verts = append(verts, Vertex3D{
			X: c2 * radius, Y: -halfH, Z: s2 * radius,
			NX: c2 / normLen, NY: slope / normLen, NZ: s2 / normLen,
			U: float32(i+1) / float32(segments), V: 1, R: white[0], G: white[1], B: white[2], A: white[3],
		})

		triIndices = append(triIndices, idx, idx+1, idx+2)
		lineIndices = append(lineIndices, idx, idx+1, idx+1, idx+2)
	}

	// --- Base Cap ---
	baseCenterIdx := uint32(len(verts))
	verts = append(verts, Vertex3D{X: 0, Y: -halfH, Z: 0, NX: 0, NY: -1, NZ: 0, U: 0.5, V: 0.5, R: white[0], G: white[1], B: white[2], A: white[3]})

	baseStart := uint32(len(verts))
	for i := 0; i <= segments; i++ {
		theta := float32(i) * 2.0 * float32(math.Pi) / float32(segments)
		c := float32(math.Cos(float64(theta)))
		s := float32(math.Sin(float64(theta)))
		verts = append(verts, Vertex3D{
			X: c * radius, Y: -halfH, Z: s * radius,
			NX: 0, NY: -1, NZ: 0,
			U: 0.5 + c*0.5, V: 0.5 + s*0.5, R: white[0], G: white[1], B: white[2], A: white[3],
		})
	}
	for i := 0; i < segments; i++ {
		triIndices = append(triIndices, baseCenterIdx, baseStart+uint32(i+1), baseStart+uint32(i))
	}

	return &Model3D{
		Name: "Cone",
		Meshes: []Mesh3D{
			{Name: "ConeMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 5. PYRAMID (Square base with 4 triangular faces)
// =========================================================================

func NewPyramidModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	// Slant normal calculation for 45 deg slope: 1/sqrt(2) = 0.707
	n := float32(0.7071)

	verts := []Vertex3D{
		// Front Face (+Z)
		{X: 0, Y: 1, Z: 0, NX: 0, NY: n, NZ: n, U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: 1, NX: 0, NY: n, NZ: n, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: 1, NX: 0, NY: n, NZ: n, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Right Face (+X)
		{X: 0, Y: 1, Z: 0, NX: n, NY: n, NZ: 0, U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: 1, NX: n, NY: n, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: -1, NX: n, NY: n, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Back Face (-Z)
		{X: 0, Y: 1, Z: 0, NX: 0, NY: n, NZ: -n, U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: -1, NX: 0, NY: n, NZ: -n, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: -1, NX: 0, NY: n, NZ: -n, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Left Face (-X)
		{X: 0, Y: 1, Z: 0, NX: -n, NY: n, NZ: 0, U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: -1, NX: -n, NY: n, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: 1, NX: -n, NY: n, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},

		// Bottom Base (-Y)
		{X: -1, Y: -1, Z: -1, NX: 0, NY: -1, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: -1, NX: 0, NY: -1, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: -1, Z: 1, NX: 0, NY: -1, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: -1, Z: 1, NX: 0, NY: -1, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
	}

	triIndices := []uint32{
		0, 1, 2, // Front
		3, 4, 5, // Right
		6, 7, 8, // Back
		9, 10, 11, // Left
		12, 13, 14, 12, 14, 15, // Base
	}

	lineIndices := []uint32{
		1, 2, 4, 5, 7, 8, 10, 11, // Square base outline
		0, 1, 3, 4, 6, 7, 9, 10, // Slant edges to apex
	}

	return &Model3D{
		Name: "Pyramid",
		Meshes: []Mesh3D{
			{Name: "PyramidMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 6. PLANE / FLOOR QUAD (Horizontal XZ plane for ground/floors)
// =========================================================================

func NewPlaneModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	verts := []Vertex3D{
		{X: -1, Y: 0, Z: 1, NX: 0, NY: 1, NZ: 0, U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 0, Z: 1, NX: 0, NY: 1, NZ: 0, U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: 1, Y: 0, Z: -1, NX: 0, NY: 1, NZ: 0, U: 1, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
		{X: -1, Y: 0, Z: -1, NX: 0, NY: 1, NZ: 0, U: 0, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
	}

	triIndices := []uint32{0, 1, 2, 0, 2, 3}
	lineIndices := []uint32{0, 1, 1, 2, 2, 3, 3, 0}

	return &Model3D{
		Name: "Plane",
		Meshes: []Mesh3D{
			{Name: "PlaneMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 7. TORUS (Donut shape with customizable ring and tube density)
// =========================================================================

func NewTorusModel(radialSegments, tubularSegments int) *Model3D {
	if radialSegments < 6 {
		radialSegments = 16
	}
	if tubularSegments < 6 {
		tubularSegments = 24
	}

	white := [4]float32{1, 1, 1, 1}
	mainRadius := float32(0.7) // Radius from center to middle of tube
	tubeRadius := float32(0.3) // Thickness of the tube (0.7 + 0.3 = 1.0 fits in unit box)

	var verts []Vertex3D
	var triIndices []uint32
	var lineIndices []uint32

	for i := 0; i <= tubularSegments; i++ {
		u := float32(i) * 2.0 * float32(math.Pi) / float32(tubularSegments)
		cosU := float32(math.Cos(float64(u)))
		sinU := float32(math.Sin(float64(u)))

		for j := 0; j <= radialSegments; j++ {
			v := float32(j) * 2.0 * float32(math.Pi) / float32(radialSegments)
			cosV := float32(math.Cos(float64(v)))
			sinV := float32(math.Sin(float64(v)))

			// Surface Position
			x := (mainRadius + tubeRadius*cosV) * cosU
			y := tubeRadius * sinV
			z := (mainRadius + tubeRadius*cosV) * sinU

			// Surface Normal
			nx := cosV * cosU
			ny := sinV
			nz := cosV * sinU

			texU := float32(i) / float32(tubularSegments)
			texV := float32(j) / float32(radialSegments)

			verts = append(verts, Vertex3D{
				X: x, Y: y, Z: z,
				NX: nx, NY: ny, NZ: nz,
				U: texU, V: texV,
				R: white[0], G: white[1], B: white[2], A: white[3],
			})
		}
	}

	for i := 0; i < tubularSegments; i++ {
		for j := 0; j < radialSegments; j++ {
			p1 := uint32(i*(radialSegments+1) + j)
			p2 := uint32((i+1)*(radialSegments+1) + j)
			p3 := p1 + 1
			p4 := p2 + 1

			triIndices = append(triIndices, p1, p2, p3)
			triIndices = append(triIndices, p2, p4, p3)

			lineIndices = append(lineIndices, p1, p3, p1, p2)
		}
	}

	return &Model3D{
		Name: "Torus",
		Meshes: []Mesh3D{
			{Name: "TorusMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 8. TETRAHEDRON (4 triangular faces, sharp faceted normals)
// =========================================================================

func NewTetrahedronModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	// 4 base vertices of a regular tetrahedron normalized to radius 1
	invSqrt3 := float32(1.0 / math.Sqrt(3.0))
	p := [4][3]float32{
		{invSqrt3, invSqrt3, invSqrt3},
		{-invSqrt3, -invSqrt3, invSqrt3},
		{-invSqrt3, invSqrt3, -invSqrt3},
		{invSqrt3, -invSqrt3, -invSqrt3},
	}

	// 4 triangular faces
	faceIndices := [4][3]int{
		{0, 2, 1},
		{0, 1, 3},
		{0, 3, 2},
		{1, 2, 3},
	}

	var verts []Vertex3D
	var triIndices []uint32

	for i, f := range faceIndices {
		p0, p1, p2 := p[f[0]], p[f[1]], p[f[2]]
		norm := calcFaceNormal(p0, p1, p2)

		idx := uint32(i * 3)
		verts = append(verts,
			Vertex3D{X: p0[0], Y: p0[1], Z: p0[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0.5, V: 0.0, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p1[0], Y: p1[1], Z: p1[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0.0, V: 1.0, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p2[0], Y: p2[1], Z: p2[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 1.0, V: 1.0, R: white[0], G: white[1], B: white[2], A: white[3]},
		)

		triIndices = append(triIndices, idx, idx+1, idx+2)
	}

	// 6 unique wireframe edges
	lineIndices := []uint32{
		0, 1, 1, 2, 2, 0, // Face 0 edges
		3, 5, 4, 8, 7, 10, // Remaining silhouette edges
	}

	return &Model3D{
		Name: "Tetrahedron",
		Meshes: []Mesh3D{
			{Name: "TetrahedronMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 9. OCTAHEDRON (8 triangular faces, double pyramid)
// =========================================================================

func NewOctahedronModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	// 6 vertices representing the poles (+X, -X, +Y, -Y, +Z, -Z)
	poles := [6][3]float32{
		{0, 1, 0},  // 0: +Y
		{0, -1, 0}, // 1: -Y
		{1, 0, 0},  // 2: +X
		{-1, 0, 0}, // 3: -X
		{0, 0, 1},  // 4: +Z
		{0, 0, -1}, // 5: -Z
	}

	// 8 triangular faces
	faces := [8][3]int{
		{0, 4, 2}, {0, 2, 5}, {0, 5, 3}, {0, 3, 4}, // Top 4
		{1, 2, 4}, {1, 5, 2}, {1, 3, 5}, {1, 4, 3}, // Bottom 4
	}

	var verts []Vertex3D
	var triIndices []uint32

	for i, f := range faces {
		p0, p1, p2 := poles[f[0]], poles[f[1]], poles[f[2]]
		norm := calcFaceNormal(p0, p1, p2)

		idx := uint32(i * 3)
		verts = append(verts,
			Vertex3D{X: p0[0], Y: p0[1], Z: p0[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0.5, V: 0.0, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p1[0], Y: p1[1], Z: p1[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0.0, V: 1.0, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p2[0], Y: p2[1], Z: p2[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 1.0, V: 1.0, R: white[0], G: white[1], B: white[2], A: white[3]},
		)

		triIndices = append(triIndices, idx, idx+1, idx+2)
	}

	// 12 unique wireframe edges (no diagonals)
	lineIndices := []uint32{
		1, 2, 4, 5, 7, 8, 10, 11, // Equator belt
		0, 1, 3, 4, 6, 7, 9, 10, // Top apex edges
		12, 13, 15, 16, 18, 19, 21, 22, // Bottom apex edges
	}

	return &Model3D{
		Name: "Octahedron",
		Meshes: []Mesh3D{
			{Name: "OctahedronMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 10. PENTAGONAL PYRAMID (5-sided base + 5 triangular sides)
// =========================================================================

func NewPentagonalPyramidModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	var verts []Vertex3D
	var triIndices []uint32
	var lineIndices []uint32

	radius := float32(1.0)
	halfH := float32(1.0)

	// Pre-calculate 5 base perimeter points
	var basePoints [5][3]float32
	for i := 0; i < 5; i++ {
		theta := float32(i) * 2.0 * float32(math.Pi) / 5.0
		basePoints[i] = [3]float32{
			float32(math.Cos(float64(theta))) * radius,
			-halfH,
			float32(math.Sin(float64(theta))) * radius,
		}
	}

	apex := [3]float32{0, halfH, 0}

	// 1. Five Triangular Sides
	for i := 0; i < 5; i++ {
		p1 := basePoints[i]
		p2 := basePoints[(i+1)%5]
		norm := calcFaceNormal(apex, p1, p2)

		idx := uint32(len(verts))
		verts = append(verts,
			Vertex3D{X: apex[0], Y: apex[1], Z: apex[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0.5, V: 0, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p1[0], Y: p1[1], Z: p1[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 0, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
			Vertex3D{X: p2[0], Y: p2[1], Z: p2[2], NX: norm[0], NY: norm[1], NZ: norm[2], U: 1, V: 1, R: white[0], G: white[1], B: white[2], A: white[3]},
		)

		triIndices = append(triIndices, idx, idx+1, idx+2)
		lineIndices = append(lineIndices, idx+1, idx+2, idx, idx+1) // Base edge and slant edge
	}

	// 2. Flat Pentagonal Base Cap
	baseCenterIdx := uint32(len(verts))
	verts = append(verts, Vertex3D{X: 0, Y: -halfH, Z: 0, NX: 0, NY: -1, NZ: 0, U: 0.5, V: 0.5, R: white[0], G: white[1], B: white[2], A: white[3]})

	baseStart := uint32(len(verts))
	for i := 0; i < 5; i++ {
		p := basePoints[i]
		u := 0.5 + (p[0]/radius)*0.5
		v := 0.5 + (p[2]/radius)*0.5
		verts = append(verts, Vertex3D{X: p[0], Y: -halfH, Z: p[2], NX: 0, NY: -1, NZ: 0, U: u, V: v, R: white[0], G: white[1], B: white[2], A: white[3]})
	}

	for i := 0; i < 5; i++ {
		triIndices = append(triIndices, baseCenterIdx, baseStart+uint32((i+1)%5), baseStart+uint32(i))
	}

	return &Model3D{
		Name: "PentagonalPyramid",
		Meshes: []Mesh3D{
			{Name: "PentagonalPyramidMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// 11. DODECAHEDRON (12 regular pentagonal faces, 36 triangles)
// =========================================================================

func NewDodecahedronModel() *Model3D {
	white := [4]float32{1, 1, 1, 1}

	// Golden ratio phi = (1 + sqrt(5)) / 2
	phi := float32((1.0 + math.Sqrt(5.0)) / 2.0)
	invPhi := 1.0 / phi

	// 20 vertices of a regular dodecahedron normalized to radius 1.0
	scale := float32(1.0 / math.Sqrt(3.0))

	rawP := [20][3]float32{
		// (+-1, +-1, +-1)
		{-1, -1, -1}, {1, -1, -1}, {1, 1, -1}, {-1, 1, -1},
		{-1, -1, 1}, {1, -1, 1}, {1, 1, 1}, {-1, 1, 1},
		// (0, +-1/phi, +-phi)
		{0, -invPhi, -phi}, {0, invPhi, -phi}, {0, -invPhi, phi}, {0, invPhi, phi},
		// (+-1/phi, +-phi, 0)
		{-invPhi, -phi, 0}, {invPhi, -phi, 0}, {invPhi, phi, 0}, {-invPhi, phi, 0},
		// (+-phi, 0, +-1/phi)
		{-phi, 0, -invPhi}, {phi, 0, -invPhi}, {phi, 0, invPhi}, {-phi, 0, invPhi},
	}

	var p [20][3]float32
	for i := range rawP {
		p[i] = [3]float32{rawP[i][0] * scale, rawP[i][1] * scale, rawP[i][2] * scale}
	}

	// 12 pentagonal faces (defined by 5 indices each)
	dodecFaceIndices := [12][5]int{
		{0, 8, 9, 3, 16},
		{0, 16, 19, 4, 12},
		{0, 12, 13, 1, 8},
		{1, 17, 2, 9, 8},
		{1, 13, 5, 18, 17},
		{2, 14, 15, 3, 9},
		{2, 17, 18, 6, 14},
		{3, 15, 7, 19, 16},
		{4, 10, 11, 7, 19},
		{4, 12, 13, 5, 10},
		{5, 18, 6, 11, 10},
		{6, 14, 15, 7, 11},
	}

	var verts []Vertex3D
	var triIndices []uint32
	var lineIndices []uint32

	for _, face := range dodecFaceIndices {
		p0, p1, p2 := p[face[0]], p[face[1]], p[face[2]]
		norm := calcFaceNormal(p0, p1, p2)

		startIdx := uint32(len(verts))

		// 5 vertices for this pentagonal face
		for j, ptIdx := range face {
			ang := float32(j) * 2.0 * float32(math.Pi) / 5.0
			u := 0.5 + float32(math.Cos(float64(ang)))*0.5
			v := 0.5 + float32(math.Sin(float64(ang)))*0.5

			pt := p[ptIdx]
			verts = append(verts, Vertex3D{
				X: pt[0], Y: pt[1], Z: pt[2],
				NX: norm[0], NY: norm[1], NZ: norm[2],
				U: u, V: v,
				R: white[0], G: white[1], B: white[2], A: white[3],
			})
		}

		// Triangulate pentagon into 3 triangles: (0, 1, 2), (0, 2, 3), (0, 3, 4)
		triIndices = append(triIndices, startIdx, startIdx+1, startIdx+2)
		triIndices = append(triIndices, startIdx, startIdx+2, startIdx+3)
		triIndices = append(triIndices, startIdx, startIdx+3, startIdx+4)

		// 5 clean boundary edges for wireframe
		lineIndices = append(lineIndices,
			startIdx, startIdx+1,
			startIdx+1, startIdx+2,
			startIdx+2, startIdx+3,
			startIdx+3, startIdx+4,
			startIdx+4, startIdx,
		)
	}

	return &Model3D{
		Name: "Dodecahedron",
		Meshes: []Mesh3D{
			{Name: "DodecahedronMesh", Vertices: verts, TriIndices: triIndices, LineIndices: lineIndices},
		},
	}
}

// =========================================================================
// HELPER: Face Normal Calculation
// =========================================================================

func calcFaceNormal(p0, p1, p2 [3]float32) [3]float32 {
	e1 := [3]float32{p1[0] - p0[0], p1[1] - p0[1], p1[2] - p0[2]}
	e2 := [3]float32{p2[0] - p0[0], p2[1] - p0[1], p2[2] - p0[2]}

	nx := e1[1]*e2[2] - e1[2]*e2[1]
	ny := e1[2]*e2[0] - e1[0]*e2[2]
	nz := e1[0]*e2[1] - e1[1]*e2[0]

	lenN := float32(math.Sqrt(float64(nx*nx + ny*ny + nz*nz)))
	if lenN > 0 {
		return [3]float32{nx / lenN, ny / lenN, nz / lenN}
	}
	return [3]float32{0, 1, 0}
}
