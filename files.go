package main

import (
	"log"
	"path/filepath"

	"github.com/Zyko0/go-sdl3/sdl"
)

type ModelFile struct {
	mod      *GPUModel
	thumb    *Texture
	nm, path string
}

func handleDropEvent(event sdl.Event) {
	switch event.Type {
	case sdl.EVENT_DROP_POSITION:
		drop := event.DropEvent()
		_ = drop.X
		_ = drop.Y
	case sdl.EVENT_DROP_BEGIN:
	case sdl.EVENT_DROP_COMPLETE:
	}
}

// openNativeFileDialog: Open file OS window
func openNativeFileDialog(nm string, siz int) {
	var callback sdl.DialogFileCallback
	var filters []sdl.DialogFileFilter

	switch nm {
	case "texture":
		filters = []sdl.DialogFileFilter{
			{Name: "Image Files (*.png, *.jpg)", Pattern: "png;jpg;jpeg"},
			{Name: "All Files (*.*)", Pattern: "*"},
		}
		callback = sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
			if len(fileList) == 0 {
				return
			}
			for _, filePath := range fileList {
				if HasTex(filePath) {
					mInfoTXT("Texture already loaded: " + filepath.Base(filePath))
					continue
				}
				tex, err := LoadTexture(gpuDev, filePath)
				if err != nil {
					log.Printf("Failed to load 2D texture: %v", err)
					return
				}
				high, wide := false, false
				if tex.W > un3 {
					wide = true
				}
				if tex.H > un3 {
					high = true
				}
				usr.tex2d = append(usr.tex2d, TEX2D{
					path:      filePath,
					nm:        filepath.Base(filePath),
					recBorder: R(0, 0, tex.W, tex.H),
					cnt:       reccntr(R(0, 0, tex.W, tex.H)),
					W:         tex.W,
					H:         tex.H,
					tex:       tex,
					tags:      []string{"examples", "user"},
					c:         Col.White,
					cShadow:   CA(Col.BlackCharcoal, 100),
					shadowX:   -unh,
					shadowY:   unh,
					wide:      wide,
					high:      high,
					scale:     1,
					curAnim:   -1,
				})
			}
		})

	case "model":
		filters = []sdl.DialogFileFilter{
			{Name: "3D Model Files (*.glb)", Pattern: "glb"},
			{Name: "All Files (*.*)", Pattern: "*"},
		}
		callback = sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
			if len(fileList) == 0 {
				return
			}
			for _, filePath := range fileList {
				if HasModel(filePath) {
					mInfoTXT("Model already loaded: " + filepath.Base(filePath))
					continue
				}
				rawModel, err := LoadGLB(filePath)
				if err != nil {
					log.Printf("Failed loading GLB model: %v", err)
					continue
				}

				gpuModel, err := r3d.UploadModel(rawModel)
				if err != nil {
					log.Printf("Failed uploading GLB model to GPU: %v", err)
					continue
				}

				icon, err := GenerateModelThumbnail(gpuDev, win, r3d, gpuModel, siz)
				if err != nil {
					log.Printf("Failed generating thumbnail: %v", err)
				}

				usr.o3d = append(usr.o3d, OBJ3D{
					path:    filePath,
					nm:      filepath.Base(filePath),
					thumb:   icon,
					mod:     gpuModel,
					tags:    []string{"examples", "user"},
					Color:   Col.White,
					UVScale: 1.0,
					Scale:   0.5,
				})
			}
		})

	// 🚀 NEW: File dialog to pick a 2D texture for the currently selected 3D model!
	case "object_texture":
		filters = []sdl.DialogFileFilter{
			{Name: "Image Files (*.png, *.jpg)", Pattern: "png;jpg;jpeg"},
			{Name: "All Files (*.*)", Pattern: "*"},
		}
		callback = sdl.NewDialogFileCallback(func(fileList []string, filter int32) {
			if len(fileList) == 0 || selOBJ3Dnum < 0 || selOBJ3Dnum >= len(lev3d.o3d) {
				return
			}
			filePath := fileList[0]
			tex, err := LoadTexture(gpuDev, filePath)
			if err != nil {
				log.Printf("Failed to load texture for object: %v", err)
				return
			}
			// Assign texture directly to the selected object!
			lev3d.o3d[selOBJ3Dnum].CustomTex = tex
		})
	}

	sdl.ShowOpenFileDialog(callback, win, filters, "", true)
}

// HasModel checks if a model with the given file path is already loaded in usr.o3d
func HasModel(path string) bool {
	if usr == nil {
		return false
	}
	cleanTarget := filepath.Clean(path)

	// 🚀 Fast zero-copy loop (doesn't copy OBJ3D structs)
	for i := range usr.o3d {
		if filepath.Clean(usr.o3d[i].path) == cleanTarget {
			return true // Duplicate found!
		}
	}
	return false
}
func HasTex(path string) bool {
	if usr == nil {
		return false
	}
	cleanTarget := filepath.Clean(path)

	// 🚀 Fast zero-copy loop (doesn't copy OBJ3D structs)
	for i := range usr.tex2d {
		if filepath.Clean(usr.tex2d[i].path) == cleanTarget {
			return true // Duplicate found!
		}
	}
	return false
}
