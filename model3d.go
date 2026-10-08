package main

import (
	"bytes"
	"fmt"
	"image"
	_ "image/jpeg" // Support embedded JPEG textures
	_ "image/png"  // Support embedded PNG textures
	"math"

	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
)

type Vertex3D struct {
	X, Y, Z    float32 // Position (Offset 0)
	NX, NY, NZ float32 // Normal   (Offset 12)
	U, V       float32 // UV Coord (Offset 24)
	R, G, B, A float32 // Color    (Offset 32)
}

type TextureData struct {
	Width, Height int
	Pixels        []byte // RGBA8888 pixel bytes
}

type Mesh3D struct {
	Name        string
	Vertices    []Vertex3D
	TriIndices  []uint32
	LineIndices []uint32
	Texture     *TextureData // Embedded texture (nil if untextured)
}

type Model3D struct {
	Name   string
	Meshes []Mesh3D
}

// LoadGLB loads any .glb file, extracting geometry AND embedded texture maps
func LoadGLB(path string) (*Model3D, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open GLB: %w", err)
	}

	model := &Model3D{Name: path}

	for _, mesh := range doc.Meshes {
		for _, prim := range mesh.Primitives {
			posIdx, ok := prim.Attributes[gltf.POSITION]
			if !ok || int(posIdx) >= len(doc.Accessors) {
				continue
			}
			positions, err := modeler.ReadPosition(doc, doc.Accessors[posIdx], nil)
			if err != nil || len(positions) == 0 {
				continue
			}

			var normals [][3]float32
			if normIdx, ok := prim.Attributes[gltf.NORMAL]; ok && int(normIdx) < len(doc.Accessors) {
				normals, _ = modeler.ReadNormal(doc, doc.Accessors[normIdx], nil)
			}

			var uvs [][2]float32
			if uvIdx, ok := prim.Attributes[gltf.TEXCOORD_0]; ok && int(uvIdx) < len(doc.Accessors) {
				uvs, _ = modeler.ReadTextureCoord(doc, doc.Accessors[uvIdx], nil)
			}

			// Material & Texture Extraction
			baseColor := [4]float32{1.0, 1.0, 1.0, 1.0}
			var texData *TextureData

			if prim.Material != nil && int(*prim.Material) < len(doc.Materials) {
				mat := doc.Materials[*prim.Material]
				if mat.PBRMetallicRoughness != nil {
					bcf := mat.PBRMetallicRoughness.BaseColorFactor
					baseColor = [4]float32{float32(bcf[0]), float32(bcf[1]), float32(bcf[2]), float32(bcf[3])}

					// Extract Embedded Base Color Texture
					if mat.PBRMetallicRoughness.BaseColorTexture != nil {
						texIdx := mat.PBRMetallicRoughness.BaseColorTexture.Index
						if texIdx < len(doc.Textures) && doc.Textures[texIdx].Source != nil {
							imgIdx := *doc.Textures[texIdx].Source
							if imgIdx < len(doc.Images) {
								texData = extractGLBImage(doc, doc.Images[imgIdx])
							}
						}
					}
				}
			}

			var triIndices []uint32
			if prim.Indices != nil && int(*prim.Indices) < len(doc.Accessors) {
				triIndices, _ = modeler.ReadIndices(doc, doc.Accessors[*prim.Indices], nil)
			}
			if len(triIndices) == 0 {
				triIndices = make([]uint32, len(positions))
				for i := range positions {
					triIndices[i] = uint32(i)
				}
			}

			var lineIndices []uint32
			for i := 0; i+2 < len(triIndices); i += 3 {
				i0, i1, i2 := triIndices[i], triIndices[i+1], triIndices[i+2]
				lineIndices = append(lineIndices, i0, i1, i1, i2, i2, i0)
			}

			meshData := Mesh3D{
				Name:        mesh.Name,
				TriIndices:  triIndices,
				LineIndices: lineIndices,
				Texture:     texData,
				Vertices:    make([]Vertex3D, len(positions)),
			}

			for i := range positions {
				meshData.Vertices[i].X = positions[i][0]
				meshData.Vertices[i].Y = positions[i][1]
				meshData.Vertices[i].Z = positions[i][2]

				if len(normals) > i {
					meshData.Vertices[i].NX = normals[i][0]
					meshData.Vertices[i].NY = normals[i][1]
					meshData.Vertices[i].NZ = normals[i][2]
				}
				if len(uvs) > i {
					meshData.Vertices[i].U = uvs[i][0]
					meshData.Vertices[i].V = uvs[i][1]
				}

				meshData.Vertices[i].R = baseColor[0]
				meshData.Vertices[i].G = baseColor[1]
				meshData.Vertices[i].B = baseColor[2]
				meshData.Vertices[i].A = baseColor[3]
			}

			model.Meshes = append(model.Meshes, meshData)
		}
	}

	if len(model.Meshes) == 0 {
		return nil, fmt.Errorf("no valid geometry in %s", path)
	}

	normalizeModel(model)
	return model, nil
}

// Decodes embedded PNG/JPEG image from the GLB binary buffer
func extractGLBImage(doc *gltf.Document, imgMeta *gltf.Image) *TextureData {
	if imgMeta.BufferView == nil || int(*imgMeta.BufferView) >= len(doc.BufferViews) {
		return nil
	}

	bv := doc.BufferViews[*imgMeta.BufferView]
	if bv.Buffer >= len(doc.Buffers) {
		return nil
	}

	rawBuf := doc.Buffers[bv.Buffer].Data
	imgBytes := rawBuf[bv.ByteOffset : bv.ByteOffset+bv.ByteLength]

	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	rgbaBytes := make([]byte, w*h*4)

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

	return &TextureData{
		Width:  w,
		Height: h,
		Pixels: rgbaBytes,
	}
}

func normalizeModel(m *Model3D) {
	minX, minY, minZ := float32(math.MaxFloat32), float32(math.MaxFloat32), float32(math.MaxFloat32)
	maxX, maxY, maxZ := float32(-math.MaxFloat32), float32(-math.MaxFloat32), float32(-math.MaxFloat32)

	for _, mesh := range m.Meshes {
		for _, v := range mesh.Vertices {
			if v.X < minX {
				minX = v.X
			}
			if v.X > maxX {
				maxX = v.X
			}
			if v.Y < minY {
				minY = v.Y
			}
			if v.Y > maxY {
				maxY = v.Y
			}
			if v.Z < minZ {
				minZ = v.Z
			}
			if v.Z > maxZ {
				maxZ = v.Z
			}
		}
	}

	centerX := (minX + maxX) / 2
	centerY := (minY + maxY) / 2
	centerZ := (minZ + maxZ) / 2

	maxSize := max(maxX-minX, max(maxY-minY, maxZ-minZ))
	if maxSize == 0 {
		maxSize = 1
	}
	scaleFactor := 2.0 / maxSize

	for mi := range m.Meshes {
		for vi := range m.Meshes[mi].Vertices {
			m.Meshes[mi].Vertices[vi].X = (m.Meshes[mi].Vertices[vi].X - centerX) * scaleFactor
			m.Meshes[mi].Vertices[vi].Y = (m.Meshes[mi].Vertices[vi].Y - centerY) * scaleFactor
			m.Meshes[mi].Vertices[vi].Z = (m.Meshes[mi].Vertices[vi].Z - centerZ) * scaleFactor
		}
	}
}
