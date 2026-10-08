package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Zyko0/go-sdl3/bin/binsdl"
	"github.com/Zyko0/go-sdl3/bin/binttf"
	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/shadercross"
	"github.com/Zyko0/go-sdl3/ttf"
)

var (

	//CORE
	CNT    sdl.FPoint
	scrREC sdl.FRect
	win    *sdl.Window
	err    error
	r3d    *Render3D
	r2d    *Render2D
	gpuDev *sdl.GPUDevice
	timer  *TimeManager

	//2D
	isMovingTEX bool

	//3D MODEL EDITOR GIZMO
	culled3Dcount, drawn3Dcount int
	isDraggingGizmo             bool = false
	isMovingObj                 bool
)

func main() {
	logFile, err := os.OpenFile("engine.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err == nil {
		defer logFile.Close()
		// Send all log.Print, log.Printf, and log.Fatal output to BOTH terminal and engine.log:
		log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	}
	//	LoadSettings()
	//	defer SaveSettings()
	SET = DefaultSettings()

	defer binsdl.Load().Unload()
	defer binttf.Load().Unload()

	// 1. Initialize SDL & TTF
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	if err = ttf.Init(); err != nil {
		log.Fatal(err)
	}
	defer ttf.Quit()

	scrW, scrH := int(SET.SCRW), int(SET.SCRH)
	CNT = P(scrW/2, scrH/2)
	scrREC = R(0, 0, SET.SCRW, SET.SCRH)

	// 2. Create Window
	win, err = sdl.CreateWindow("ENGINE 08/2026", scrW, scrH, sdl.WINDOW_RESIZABLE)
	if err != nil {
		log.Fatal(err)
	}
	defer win.Destroy()

	// 3. Initialize GPU Device via gpu.go
	gpuDev, err = InitGPUDevice(win)
	if err != nil {
		log.Fatal(err)
	}
	defer gpuDev.Destroy()
	defer shadercross.CloseLibrary()

	// 4. Initialize 2D Batch Renderer (Shapes + Text)
	r2d, err = InitRender2D(gpuDev, win, scrW, scrH)
	if err != nil {
		log.Fatal(err)
	}
	defer r2d.Destroy()
	defer func() {
		// Clean up all user imported 3D models and icons on program exit
		for _, tex := range usr.tex2d {
			if tex.tex != nil {
				tex.tex.Destroy(gpuDev)
			}
		}
		for _, tex := range usr.tex2dSave {
			if tex.tex != nil {
				tex.tex.Destroy(gpuDev)
			}
		}
		for _, tex := range usr.tex2dTrash {
			if tex.tex != nil {
				tex.tex.Destroy(gpuDev)
			}
		}
		for _, tex := range lev2d.tex {
			if tex.tex != nil {
				tex.tex.Destroy(gpuDev)
			}
		}
	}()

	// 5. Initialize 3D Renderer (Shaders + Z-Buffer)
	r3d, err = InitRender3D(gpuDev, win, scrW, scrH)
	if err != nil {
		log.Fatal(err)
	}
	defer r3d.Destroy()
	defer func() {
		// Clean up all user imported 3D models and icons on program exit
		for _, obj := range usr.o3d {
			if obj.mod != nil {
				r3d.ReleaseGPUModel(obj.mod)
			}
			if obj.thumb != nil {
				obj.thumb.Destroy(gpuDev)
			}
		}
		for _, obj := range usr.o3dSave {
			if obj.mod != nil {
				r3d.ReleaseGPUModel(obj.mod)
			}
			if obj.thumb != nil {
				obj.thumb.Destroy(gpuDev)
			}
		}
		for _, obj := range usr.o3dTrash {
			if obj.mod != nil {
				r3d.ReleaseGPUModel(obj.mod)
			}
			if obj.thumb != nil {
				obj.thumb.Destroy(gpuDev)
			}
		}
		for _, obj := range lev3d.o3d {
			if obj.mod != nil {
				r3d.ReleaseGPUModel(obj.mod)
			}
			if obj.thumb != nil {
				obj.thumb.Destroy(gpuDev)
			}
		}
	}()
	// 6. Initialize Fonts
	if err = InitFonts(); err != nil {
		log.Fatal(err)
	}
	defer CloseFonts()

	//MARK: INITIAL MAKE
	sdl.HideCursor()
	uiAlphaActiv = true
	if uiAlphaActiv {
		SET.UIC = CA(SET.UIC, uiAlpha)
	}
	if len(users) == 0 {
		mUsers()
	}
	mUI()
	InitDefaultLights()
	mPrim3D()
	mLevDefault()

	// 8. FPS Controller
	timer = NewTimeManager(SET.TargetFPS)

	// Main Engine Loop
	sdl.RunLoop(func() error {
		timer.BeginFrame()
		UP()

		// Event Handling
		var event sdl.Event
		for sdl.PollEvent(&event) {
			if HandleEvent(event) {
				return sdl.EndLoop
			}

			// Mode Switch Hotkeys:
			switch event.Type {

			// 3D Camera Orbit Controls
			/*
				case sdl.EVENT_MOUSE_BUTTON_DOWN:
					if event.MouseButtonEvent().Button == 1 && !MouseInUI() {
						isOrbiting = true
					}
				case sdl.EVENT_MOUSE_BUTTON_UP:
					if event.MouseButtonEvent().Button == 1 {
						isOrbiting = false
					}
				case sdl.EVENT_MOUSE_MOTION:
					m := event.MouseMotionEvent()
					if isOrbiting {
						r3d.RotY += m.Xrel * 0.01
						r3d.RotX += m.Yrel * 0.01
					}
			*/

			// Zoom in / out
			case sdl.EVENT_MOUSE_WHEEL:
				w := event.MouseWheelEvent()
				r3d.CamDist -= w.Y * 0.3
				if r3d.CamDist < 1.5 {
					r3d.CamDist = 1.5
				}
				if r3d.CamDist > 20.0 {
					r3d.CamDist = 20.0
				}
			}
		}
		if !msinUI {
			if drop3D { //MARK: DROP3D
				if cms(scrREC) && ms.lc {
					// Convert drop screen coordinates to 3D world space
					aspect := SET.SCRW / SET.SCRH
					visibleHalfH := r3d.CamDist * 0.4142 // tan(45deg / 2)
					visibleHalfW := visibleHalfH * aspect

					dropO3D.WorldX = ((ms.X/SET.SCRW)*2 - 1) * visibleHalfW
					dropO3D.WorldY = -((ms.Y/SET.SCRH)*2 - 1) * visibleHalfH
					dropO3D.WorldZ = 0
					dropO3D.Scale = 0.5
					dropO3D.RotX = 0
					dropO3D.RotY = 0
					dropO3D.cnt = V3{dropO3D.WorldX, dropO3D.WorldY, dropO3D.WorldZ}

					addOBJ3D([]OBJ3D{dropO3D})
					selOBJ3Dnum = len(lev3d.o3d) - 1 // Auto-select newly dropped object
					sbLmenuNum = 1
				}
				mInfoTXT("Left click drop, right click to exit")
				if ms.rc {
					drop3D = false
				}
			}
			if drop2D { //MARK: DROP2D TEX
				if cms(scrREC) && ms.lc {
					addTEX(droptex2D)
					selTEXnum = len(lev2d.tex) - 1
					sbLmenuNum = 2
				}
			}

			if polylineNewActiv || lineNewActiv { //MARK: LINE2D
				if ms.lc {
					lineVal.p = addPointUnique(lineVal.p, mslcP)
					if lineNewActiv && len(lineVal.p) == 2 {
						lev2d.li = append(lev2d.li, lineVal)
						lineVal = lineDef()
						selLINE2Dnum = len(lev2d.li) - 1
						lev2d.li[selLINE2Dnum].cnt = lev2d.li[selLINE2Dnum].Mid()
						lev2d.li[selLINE2Dnum].l = lev2d.li[selLINE2Dnum].Length()
						lineNewActiv = false
					}
					ms.lc = false
				}
				if ms.rc {
					if len(lineVal.p) < 2 {
						lineVal = lineDef()
						polylineNewActiv = false
						lineNewActiv = false
					} else {
						lev2d.li = append(lev2d.li, lineVal)
						selLINE2Dnum = len(lev2d.li) - 1
						lineVal = lineDef()
						lev2d.li[selLINE2Dnum].cnt = lev2d.li[selLINE2Dnum].Mid()
						lev2d.li[selLINE2Dnum].l = lev2d.li[selLINE2Dnum].Length()
						polylineNewActiv = false
					}
					ms.rc = false
				}
			}
			if polylineNewActiv {
				mInfoTXT("Left click to add points, right click to end line")
			} else if lineNewActiv {
				mInfoTXT("Left click to add points, right click to cancel")
			}
		}

		//MARK: DRAW 2D
		r2d.Begin()

		if len(lev2d.gr) > 0 { //MARK: DRAW 2D GRID
			for i := len(lev2d.gr) - 1; i >= 0; i-- {
				dGrid(&lev2d.gr[i], 1)
				if cms(lev2d.gr[i].recBorder) {
					drrlw(lev2d.gr[i].recBorder, 2, RandomColor())
				}
			}
		}
		if len(lev2d.li) > 0 { //MARK: DRAW 2D LINES
			inLineRec := false
			if selLINE2Dnum != -1 {
				sbLmenuNum = 3
			}
			for i := len(lev2d.li) - 1; i >= 0; i-- {
				lev2d.li[i].Draw()
				if SET.DrawHidden {
					if len(lev2d.li[i].p) > 1 {
						for j := range len(lev2d.li[i].p) { //MOVE LINE POINT
							r := recfromcntr(lev2d.li[i].p[j], unq, unq)
							if cms(r) {
								inLineRec = true
								mInfoTXT("Point " + fmt.Sprint(j) + ": Left click move point, right click select line")
								drr(r, SET.UIC3)
								if lineChangeActiv == -1 && lineChangeActiv != -2 && ms.lc {
									lineVal = lineDef()
									lineVal.pPrevious = lev2d.li[i].p[j]
									lineChangeActiv = j
									selLINE2Dnum = i
									ms.lc = false
								}
								if ms.rc {
									selLINE2Dnum = i
								}
							} else {
								drrl(r, SET.UIC3)
							}
						}
					}
					if len(lev2d.li[i].cnt) > 0 { //MOVE ENTIRE LINE
						for j := range len(lev2d.li[i].cnt) {
							r := recfromcntr(lev2d.li[i].cnt[j], unq, unq)
							if cms(r) {
								inLineRec = true
								if lineChangeActiv == -2 {
									mInfoTXT("Left click to drop in new position, right click to cancel")
								} else {
									mInfoTXT("Midpoint " + fmt.Sprint(j) + ": Left click move entire line, right click select line")
								}
								drr(r, SET.UIC3)
								if lineChangeActiv == -2 && ms.lc {
									lineChangeActiv = -1
									ms.lc = false
								}
								if lineChangeActiv == -1 && ms.lc {
									lineChangeActiv = -2
									selLINE2Dnum = i
									lineVal = lev2d.li[i].Clone()
									ms.lc = false
								}
								if ms.rc {
									selLINE2Dnum = i
								}
							} else {
								drrl(r, SET.UIC4)
							}
						}
					}
				}
			}

			//CLEAR LINE SELECT IF MOUSE CLICK & NOT COLLIDING WITH LINE HIT BOXES
			if !msinUI && selLINE2Dnum != -1 && !inLineRec {
				if ms.lc || ms.rc {
					selLINE2Dnum = -1
					sbLmenuNum = -1
				}
			}

			//LINE POINT MOVE
			if lineChangeActiv != -1 && lineChangeActiv != -2 {
				mInfoTXT("Left click to set new point, right click to cancel")
				lev2d.li[selLINE2Dnum].p[lineChangeActiv] = msp
				lev2d.li[selLINE2Dnum].l = lev2d.li[selLINE2Dnum].Length()
				lev2d.li[selLINE2Dnum].cnt = lev2d.li[selLINE2Dnum].Mid()
				if ms.lc {
					lev2d.li[selLINE2Dnum].l = lev2d.li[selLINE2Dnum].Length()
					lineChangeActiv = -1
					lineVal = lineDef()
					ms.lc = false
				}
				if ms.rc {
					lev2d.li[selLINE2Dnum].p[lineChangeActiv] = lineVal.pPrevious
					lev2d.li[selLINE2Dnum].cnt = lev2d.li[selLINE2Dnum].Mid()
					lineChangeActiv = -1
					lineVal = lineDef()
					ms.rc = false
				}

			} else if lineChangeActiv == -2 {
				if ms.rc {
					lev2d.li[selLINE2Dnum] = lineVal
					lineVal = lineDef()
					lineChangeActiv = -1
					ms.rc = false
				}
			}

		}
		if len(lev2d.rec) > 0 { //MARK: DRAW 2D REC
			inRec := false
			for i := len(lev2d.rec) - 1; i >= 0; i-- {
				if lev2d.rec[i].xyPrev != P(lev2d.rec[i].r.X, lev2d.rec[i].r.Y) {
					lev2d.rec[i].p = recRotatedPoints(lev2d.rec[i].r, lev2d.rec[i].ro) //MOVED
					lev2d.rec[i].xyPrev = P(lev2d.rec[i].r.X, lev2d.rec[i].r.Y)
				}
				lev2d.rec[i].Draw() // DRAW
				if selREC2Dnum == i {
					lev2d.rec[i].DrawSelected()
					if !msinUI && lev2d.rec[i].HasMouse() {
						if recMoveActiv == -1 && !recNewActiv {
							mInfoTXT("Left click to move")
							if ms.lc {
								prevREC2Dpos = P(lev2d.rec[i].r.X, lev2d.rec[i].r.Y)
								recMoveActiv = i
								ms.lc = false
							}
						} else if recMoveActiv != -1 && !recNewActiv {
							mInfoTXT("Left click to drop in new position, right click to cancel")
							if ms.lc {
								recMoveActiv = -1
								ms.lc = false
							}
							if ms.rc {
								lev2d.rec[i].r.X = prevREC2Dpos.X
								lev2d.rec[i].r.Y = prevREC2Dpos.Y
								prevREC2Dpos = sdl.FPoint{}
								recMoveActiv = -1
								ms.rc = false
							}
						}
					}
				}
				if !msinUI && lev2d.rec[i].HasMouse() { //SELECT RECTANGLES
					inRec = true
					if ms.lc {
						selREC2Dnum = i
						ms.lc = false
					}
				}

			}
			if !inRec && !msinUI {
				if ms.lc || ms.rc {
					selREC2Dnum = -1
				}
			}
		}
		if len(lev2d.tri) > 0 { //MARK: DRAW 2D TRI
			for i := len(lev2d.tri) - 1; i >= 0; i-- {
				lev2d.tri[i].Draw()
				if selTRI2Dnum == i {
					lev2d.tri[i].DrawSelected()
				}
				if ctrip(msp, lev2d.tri[i].p) {
					if ms.lc {
						selTRI2Dnum = i
					}
				}
			}
		}
		if len(lev2d.tex) > 0 { //MARK: DRAW 2D TEXTURES
			for i := len(lev2d.tex) - 1; i >= 0; i-- {
				dtex2Dcntr(&lev2d.tex[i], lev2d.tex[i].cnt)
			}
		}
		if recNewActiv { //MARK:NEW REC2D
			if !msinUI {
				if ms.lc && len(recVal.p) == 0 {
					recVal.p = append(recVal.p, mslcP)
					ms.lc = false
				} else if ms.lc && len(recVal.p) == 1 {
					recVal.p = recPointsFrom2Points(recVal.p[0], mslcP)
					recVal.r = recFromPoints(recVal.p[0], recVal.p[2])
					recVal.xyPrev = P(recVal.r.X, recVal.r.Y)
					lev2d.rec = append(lev2d.rec, recVal)
					recVal = recDef()
					selREC2Dnum = len(lev2d.rec) - 1
					recNewActiv = false
				}
			}
			if len(recVal.p) == 1 {
				drcntrWH(recVal.p[0], un8th, un8th, SET.UICgrn)
				drppl(recVal.p[0], msp, SET.Rec2DLineC)
			}
		}
		if triNewActiv { //MARK: NEW TRI2D
			dtrilcntr(msp, SET.Tri2DSideW, SET.Tri2DLineC)
			if !msinUI && ms.lc {
				triVal.cnt = mslcP
				triVal.p = triPointsCntr(triVal.cnt, SET.Tri2DSideW)
				lev2d.tri = append(lev2d.tri, triVal)
				triVal = triDef()
				triNewActiv = false
			}
		}

		if polylineNewActiv || lineNewActiv { //MARK: NEW POLYLINE2D LINE2D
			if len(lineVal.p) >= 1 {
				drcntrWH(lineVal.p[0], un8th, un8th, SET.UICgrn)
				if len(lineVal.p) > 1 {
					lineVal.Draw()
				}
				dliwp(lineVal.p[len(lineVal.p)-1], msp, lineVal.w, lineVal.c)
			}
		}

		//DROP3D
		// Reset cursor icon for the frame
		activeCursor = iconCursor
		// Evaluate Bounding Boxes for all placed 3D objects
		if !drop3D && len(lev3d.o3d) > 0 {
			hoveredIdx := -1
			for i := len(lev3d.o3d) - 1; i >= 0; i-- {
				obj := &lev3d.o3d[i]
				bbox := r3d.GetObjectBoundingBox(obj)
				if !bbox.Valid {
					continue
				}

				// Check if mouse is inside the 2D bounding box
				isHovered := ms.X >= bbox.MinX && ms.X <= bbox.MaxX &&
					ms.Y >= bbox.MinY && ms.Y <= bbox.MaxY

				if isHovered {
					hoveredIdx = i
					// 🚀 Show hand icon when hovering over a 3D model!
					activeCursor = iconHand
				}

				// Draw Bounding Boxes & Gizmo:
				if selOBJ3Dnum == i {
					// Selected: Bright Gold Box
					Draw3DBoundingBox(bbox, SET.UIC4, 2.0)

					// Gizmo Handle Position (above model)
					gizmoX := (bbox.MinX + bbox.MaxX) / 2
					gizmoY := bbox.MinY - 30
					gizmoRadius := float32(18)
					gizmoRect := R(gizmoX-gizmoRadius, gizmoY-gizmoRadius, gizmoRadius*2, gizmoRadius*2)
					gizmoHovered := cms(gizmoRect)

					// Gizmo takes priority over body dragging
					if gizmoHovered {
						activeCursor = iconMove
						if ms.lc {
							isDraggingGizmo = true
						}
					}

					gizmoCol := Col.CyanBright
					if isDraggingGizmo {
						gizmoCol = Col.PinkHot
					} else if gizmoHovered {
						gizmoCol = Col.YellowCyber
					}

					drround(gizmoRect.X, gizmoRect.Y, gizmoRect.W, gizmoRect.H, gizmoRadius, CA(Col.BlackPure, 200))
					drlround(gizmoRect.X, gizmoRect.Y, gizmoRect.W, gizmoRect.H, gizmoRadius, 2.0, gizmoCol)
					dtxtreccnt(Fon.IconSml, iconMove, gizmoRect, gizmoCol)
					dliw(gizmoX, gizmoY+gizmoRadius, gizmoX, bbox.MinY, 1.5, gizmoCol)

				} else if isHovered {
					// Hovered: Cyan Box
					Draw3DBoundingBox(bbox, Col.CyanBright, 1.5)
				}
			}

			// Left click selection and move initiation:
			if !msinUI {
				if ms.lc && !isDraggingGizmo {
					selOBJ3Dnum = hoveredIdx
					if selOBJ3Dnum != -1 {
						isMovingObj = true
					}
				}
			}
		}

		// Keep hand cursor active while actively moving an object
		if !msinUI && drop3D {
			dtexmscntrscale(dropO3D.thumb, 1, CA(Col.White, 150))
		}
		if drop2D {
			if !msinUI {
				dtexmscntrscale(droptex2D.tex, droptex2D.scale, CA(Col.White, 150))
				mInfoTXT("Left click drop, right click to exit")
			}
			if ms.rc {
				droptex2D = TEX2D{}
				drop2D = false
			}
		} else {
			if len(lev2d.tex) > 0 { //SELECT MOVE TEXTURE 2D
				hoveredIdx := -1
				for i := len(lev2d.tex) - 1; i >= 0; i-- {
					if lev2d.tex[i].ro != 0 { //ROTATION
						if cmsrro4slice([4]sdl.FPoint(lev2d.tex[i].recPoints)) {
							activeCursor = iconHand
							hoveredIdx = i
							drlro(lev2d.tex[i].recPoints, 2, RandomColor())
							mInfoTXT("Left click edit, left click hold move, right click to copy")
						}
						if selTEXnum == i {
							drlro(lev2d.tex[i].recPoints, 2, SET.UIC4)
						}
					} else { //NO ROTATION
						if cms(lev2d.tex[i].recBorder) {
							hoveredIdx = i
							drrlw(lev2d.tex[i].recBorder, 2, RandomColor())
							if selectFrame {
								mInfoTXT("Left click to add/remove frames from animation sequence, right click exit")
							} else {
								activeCursor = iconHand
								mInfoTXT("Left click edit, left click hold move, right click to copy")
							}
						}
						if selTEXnum == i {
							drrlw(lev2d.tex[i].recBorder, 2, SET.UIC4)
						}
					}
				}
				if !msinUI {
					if ms.lc {
						selTEXnum = hoveredIdx
						if selTEXnum == -1 {
							if sbLmenuNum == 2 {
								sbLmenuNum = -1
							}
						} else {
							isMovingTEX = true
						}

					}
				}
			}
		}

		if !hideUI { //UI
			dUI()
		}
		if settingsVis {
			dSettings()
		}
		if PaletteMode { //PALETTE
			dPalette(unh)
		}
		if PaletteModeAlpha { //PALETTE ALPHA
			dPaletteAlpha(unh)
		}
		if DebugMode { //DEBUG
			dDebug()
		}
		// WARN TXT
		if warntxtTimer > 0 {
			dtxtTopBottomWin(fl, warntxt, 0, false, SET.UICred)
			warntxtTimer--
		}
		//INFO TXT
		if infotxtTimer == 0 {
			infotxt = ""
		} else {
			dtxtTopBottomWin(fd, infotxt, unh+un8th, true, SET.UIC4)
			infotxtTimer--
		}
		if infotxt2Timer == 0 {
			infotxt2 = ""
		} else {
			dtxtTopBottomWin(fd, infotxt2, unh+un8th+fonh(fd), true, SET.UIC4)
			infotxt2Timer--
		}
		//CURSOR
		cursorC := SET.CursorColor
		if isMovingObj || isMovingTEX {
			activeCursor = iconHand
			cursorC = SET.UIC3
		} else if isCursorActiv {
			cursorC = SET.UIC3
		}
		dtxt(id, activeCursor, ms.X-9, ms.Y-7, Col.BlackCharcoal)
		dtxt(id, activeCursor, ms.X-8, ms.Y-8, cursorC)
		// ==========================================
		// 2. GPU EXECUTION & RENDER PASS
		// ==========================================
		cmd, err := gpuDev.AcquireCommandBuffer()
		if err != nil {
			return nil
		}

		// Upload dynamic 2D shapes & text to VRAM
		r2d.Upload(cmd)

		bgr, bgg, bgb, bga := color2float(SET.BGC)
		swapchainTex, err := cmd.WaitAndAcquireGPUSwapchainTexture(win)
		if err == nil && swapchainTex != nil {
			colorTarget := sdl.GPUColorTargetInfo{
				Texture:    swapchainTex.Texture,
				LoadOp:     sdl.GPU_LOADOP_CLEAR,
				StoreOp:    sdl.GPU_STOREOP_STORE,
				ClearColor: sdl.FColor{R: bgr, G: bgg, B: bgb, A: bga}, // Background color
			}

			// Hardware Depth Buffer (Z-Buffer) Target
			depthTarget := sdl.GPUDepthStencilTargetInfo{
				Texture:    r3d.depthTex,
				LoadOp:     sdl.GPU_LOADOP_CLEAR,
				StoreOp:    sdl.GPU_STOREOP_DONT_CARE,
				ClearDepth: 1.0,
				Cycle:      true,
			}

			//MARK: DRAW 3D
			// Begin Unified 3D + 2D Render Pass
			renderPass := cmd.BeginRenderPass([]sdl.GPUColorTargetInfo{colorTarget}, &depthTarget)
			cullcount, drawncount := 0, 0
			if len(lev3d.o3d) > 0 {
				for i := range lev3d.o3d {
					obj := &lev3d.o3d[i]
					if !r3d.IsObjectVisible(obj) {
						cullcount++
						continue
					}
					if obj.mod != nil {
						drawncount++
						// 🚀 Draws position, rotation, scale, color tint, UV tiling, and custom textures!
						r3d.DrawObj(cmd, renderPass, obj)
					}
				}
			}

			culled3Dcount, drawn3Dcount = cullcount, drawncount

			// Step B: Draw 2D UI & Text Overlay on Top
			r2d.Draw(cmd, renderPass)

			renderPass.End()
		}

		cmd.Submit()
		AFTER()
		timer.EndFrame()
		return nil
	})
}
