package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	Keys              = make(map[sdl.Keycode]bool)
	ms                MouseState
	msp, mslcP, msrcP sdl.FPoint
	activeCursor      string      = iconCursor // Default cursor icon
	LastKey           sdl.Keycode              // 🚀 Stores the key pressed on this frame

	LastKeyRepeat, isCursorActiv, msinUI bool

	msScrollSpeed float32 = unh

	mYrel, mXrel float32
)

type MouseState struct {
	X, Y           float32
	ld, rd, rc, lc bool
}

func uINP() {
	ms.lc = false
	ms.rc = false
	LastKey = 0
	//MOUSE
	if cms(sbL.r) || cms(sbR.r) || cms(qiconsBGR) || PaletteModeAlpha || PaletteMode {
		msinUI = true
	} else {
		msinUI = false
	}
}

// Process SDL Events
func HandleEvent(event sdl.Event) bool {
	switch event.Type {
	case sdl.EVENT_QUIT:
		return true // Request quit

	case sdl.EVENT_KEY_DOWN:
		kb := event.KeyboardEvent()
		Keys[kb.Key] = true
		LastKey = kb.Key // 🚀 Capture for text input
		LastKeyRepeat = kb.Repeat

		if kb.Key == sdl.K_ESCAPE {
			return true // Request quit
		}

		// Single-press toggles
		if !kb.Repeat {

			if kb.Key == sdl.K_F1 {
				settingsVis = !settingsVis
			}
			if kb.Key == sdl.K_F2 {
				hideUI = !hideUI
				if hideUI {
					mWarnTXT("F2 key to restore UI")
				}
			}
			if kb.Key == sdl.K_F3 {

			}
			if kb.Key == sdl.K_F9 {

			}
			if kb.Key == sdl.K_F10 {
				DebugMode = !DebugMode
			}
			if kb.Key == sdl.K_F11 {
				PaletteMode = !PaletteMode
			}
			if kb.Key == sdl.K_F12 {
				PaletteModeAlpha = !PaletteModeAlpha
			}
		}
		//3D MODELS IN EDITOR MODES
		k := event.KeyboardEvent().Key
		if k == sdl.K_1 {
			R3D.Mode = ModeShaded
		}
		if k == sdl.K_2 {
			R3D.Mode = ModeSolid
		}
		if k == sdl.K_3 {
			R3D.Mode = ModeWireframe
		}
		if k == sdl.K_ESCAPE {
			selOBJ3Dnum = -1 // Deselect
		}

	case sdl.EVENT_KEY_UP:
		Keys[event.KeyboardEvent().Key] = false

	case sdl.EVENT_MOUSE_MOTION:
		m := event.MouseMotionEvent()
		ms.X = m.X
		ms.Y = m.Y
		mXrel = m.Xrel
		mYrel = m.Yrel
		msp = sdl.FPoint{m.X, m.Y}

		//MOVE 2D LINES IN EDITOR
		if lineChangeActiv == -2 {
			lev2d.li[selLINE2Dnum].Move(m.Xrel, m.Yrel)
			mInfoTXT("Left click to drop in new position, right click to cancel")
			if ms.lc {
				lineChangeActiv = -1
				ms.lc = false
			}
		}

		//MOVE 2D RECTANGLES IN EDITOR
		if recMoveActiv != -1 {
			lev2d.rec[recMoveActiv].r.X += m.Xrel
			lev2d.rec[recMoveActiv].r.Y += m.Yrel
		}

		//2D TEXTURES IN EDITOR
		if isMovingTEX && selTEXnum >= 0 && selTEXnum < len(lev2d.tex) {
			lev2d.tex[selTEXnum].cnt.X += m.Xrel
			lev2d.tex[selTEXnum].cnt.Y += m.Yrel
			lev2d.tex[selTEXnum].recBorder = recmovedcntr(lev2d.tex[selTEXnum].recBorder, lev2d.tex[selTEXnum].cnt)
		}

		//3D MODELS IN EDITOR
		// 1. If dragging the Rotation Gizmo (Rotate)
		if isDraggingGizmo && selOBJ3Dnum >= 0 && selOBJ3Dnum < len(lev3d.o3d) {
			lev3d.o3d[selOBJ3Dnum].RotY += m.Xrel * 0.02
			lev3d.o3d[selOBJ3Dnum].RotX += m.Yrel * 0.02
		}

		// 2. 🚀 If dragging the Object Body (Translate / Move in 3D)
		if isMovingObj && selOBJ3Dnum >= 0 && selOBJ3Dnum < len(lev3d.o3d) {
			aspect := SET.SCRW / SET.SCRH
			visibleHalfH := R3D.CamDist * 0.41421356 // tan(45deg / 2)
			visibleHalfW := visibleHalfH * aspect

			// Convert screen pixel delta to 3D world space delta
			worldDeltaX := (m.Xrel / SET.SCRW) * (visibleHalfW * 2)
			worldDeltaY := -(m.Yrel / SET.SCRH) * (visibleHalfH * 2)

			lev3d.o3d[selOBJ3Dnum].WorldX += worldDeltaX
			lev3d.o3d[selOBJ3Dnum].WorldY += worldDeltaY
		}

	case sdl.EVENT_MOUSE_BUTTON_DOWN:

		b := event.MouseButtonEvent()
		if b.Button == 1 {
			ms.ld = true
			ms.lc = true
			mslcP = msp
		}
		if b.Button == 3 {
			ms.rc = true
			ms.rd = true
			msrcP = msp
		}

	case sdl.EVENT_MOUSE_BUTTON_UP:
		b := event.MouseButtonEvent()
		if b.Button == 1 {
			ms.ld = false
		}
		if b.Button == 3 {
			ms.rd = false
		}
		//3D MODELS IN EDITOR
		if event.MouseButtonEvent().Button == 1 {
			isDraggingGizmo = false
			isMovingObj = false // 🚀 Release object dragging
			isMovingTEX = false
		}

	case sdl.EVENT_MOUSE_WHEEL:
		w := event.MouseWheelEvent()

		//SIDEBARS
		if inSBL {
			if sbL.ys == 0 {
				if w.Y < 0 {
					sbL.ys += w.Y * msScrollSpeed
				}
			} else if sbL.ys < 0 {
				sbL.ys += w.Y * msScrollSpeed
			}
			if sbL.ys > 0 {
				sbL.ys = 0
			}
		}
		debugnum = int(w.Y)
		//3D
		R3D.CamDist -= w.Y * 0.3
		if R3D.CamDist < 1.5 {
			R3D.CamDist = 1.5
		}
		if R3D.CamDist > 20.0 {
			R3D.CamDist = 20.0
		}

	}

	return false
}
