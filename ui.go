package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

var (

	//UI DEFAULTS
	uiAlpha              uint8 = 150
	uiIconsSiz, uispc          = unh, un16th * 3
	thumbWdef                  = int(un + unh)
	hideUI, uiAlphaActiv bool
	viewportW, viewportH float32

	//SIDEBARS
	inSBL, inSBR bool
	sbR, sbL     SIDEBAR
	sbRGeomIcons ICONLIST
	sbLmenuNum   = -1
	scrollBarSiz = unq + un8th

	nameEdit, widthEdit, heightEdit bool

	scaleStep, roStep, lineWstep float32 = 0.5, 5, 0.5

	selOBJ3Dnum, selTRI2Dnum, selTEXnum, selLINE2Dnum, selREC2Dnum int = -1, -1, -1, -1, -1

	//SIDEBAR ICONS
	sbQicons ICONLIST

	//TEXTURES ANIMATIONS
	animTemp                                     ANIM2D
	collisFramesVis, selectFrame, animPreviewVis bool
	frameListTemp                                []int

	//RECYCLE
	recycleVis bool

	//DROP
	drop3D, drop2D, dropGrid bool
	dropO3D                  OBJ3D
	droptex2D                TEX2D
	dropGridFrom, dropTimer  int

	//FOOTER ICONS
	qIconsDEF, qIcons3D ICONLIST
	qIcons              []ICONLIST
	qiconsBGR           sdl.FRect
)

type ICONLIST struct {
	nm, ic, hovertxt []string
	onoff            []bool
	siz              float32
}

type UIITEM struct {
	nm          string
	cat, winNum int
}

type SIDEBAR struct {
	r, rEditName              sdl.FRect
	ic                        ICONLIST
	tabVertNum, tabHorizNum   int
	thumbSiz, spc, ys, innerW float32
	tabHr                     []sdl.FRect
}

func dUI() { // MARK: DRAW UI
	if !settingsVis {
		dSidebarL()    //LEFT SIDEBAR MENU
		dSidebarR()    //RIGHT SIDEBAR MENU
		dFooterIcons() //FOOTER ICON MENU
	}

}
func dSBqicons(x, y float32) { // MARK: DRAW SIDEBAR QUICK ICONS
	r := R(x, y, sbQicons.siz, sbQicons.siz)
	for i := range sbQicons.ic {
		if cms(r) {
			dtxtreccnt(id2, sbQicons.ic[i], r, SET.UICgrn)
			mInfoTXT(sbQicons.hovertxt[i])
			if ms.lc {
				switch sbQicons.nm[i] {
				case "Save":
					switch sbLmenuNum {
					case 2: //TEX
						var added bool
						usr.tex2dSave, added = addUniqueTEX2D(usr.tex2dSave, lev2d.tex[selTEXnum])
						if added {
							mInfoTXT2("Saved")
						} else {
							mInfoTXT2("Not Saved, identical saved texture exists")
						}
					}
				case "Delete":
					switch sbLmenuNum {
					case 2: //TEX

						usr.tex2dTrash, _ = addUniqueTEX2D(usr.tex2dTrash, lev2d.tex[selTEXnum])
						lev2d.tex, selTEXnum = DeleteTEX2Dat(lev2d.tex, selTEXnum)
						mInfoTXT2("Deleted texture")
					}
				}
			}
		} else {
			dtxtreccnt(id2, sbQicons.ic[i], r, SET.UICtxt)
		}
		r.X += r.W
	}
}
func dSidebarR() { // MARK: DRAW SIDEBAR RIGHT
	//BG REC
	drr(sbR.r, SET.UIC)
	//TAB ICON LIST
	r := R(sbR.r.X, sbR.r.Y, sbR.ic.siz, sbR.ic.siz)
	for i, ic := range sbR.ic.ic {
		if sbR.tabVertNum == i {
			dtxtreccnt(id2, ic, r, SET.UIC3)
		} else {
			if cms(r) {
				dtxtreccnt(id2, ic, r, SET.UICgrn)
				if ms.lc {
					sbR.tabVertNum = i
				}
			} else {
				dtxtreccnt(id2, ic, r, SET.UICtxt)
			}
		}

		if cms(r) {
			mInfoTXT(sbR.ic.nm[i])
		}
		r.Y += sbR.ic.siz
	}
	//ACTIV SIDEBAR WINDOW CONTENTS
	var x, y, ox float32
	x = sbR.r.X + sbR.ic.siz + sbR.spc
	ox = x
	y = sbR.ys

	//MARK: RIGHT SIDEBAR CONTENTS

	switch sbR.tabVertNum {
	case 4: //2D GEOMETRY
		r := R(x, y, sbRGeomIcons.siz, sbRGeomIcons.siz)
		for i := range sbRGeomIcons.ic {
			if cms(r) {
				mInfoTXT("Create 2D " + sbRGeomIcons.nm[i])
				drr(r, SET.UIC3)
				if ms.lc {
					switch i {
					case 4: //TRIANGLE
						triVal = triDef()
						triNewActiv = true
						sbLmenuNum = 5
					case 3: //RECTANGLE
						recVal = recDef()
						recNewActiv = true
						sbLmenuNum = 4
					case 2: //POLYLINE
						lineVal = lineDef()
						polylineNewActiv = true
						sbLmenuNum = 3
					case 1: //LINE
						lineVal = lineDef()
						lineNewActiv = true
						sbLmenuNum = 3
					case 0: //GRID
						gridVal = gridDef()
						sbLmenuNum = 0
					}
					ms.lc = false
				}
			} else {
				drr(r, CA(SET.BGC, uiAlpha))
			}
			dtxtreccnt(il2, sbRGeomIcons.ic[i], r, SET.UICtxt)
			r.X += r.W + sbR.spc
			if r.X+r.W+sbR.spc > sbR.r.X+sbR.r.W {
				r.X = ox
				r.Y += r.H + sbR.spc
			}
		}

	case 3: //TEXTURES
		r := R(x, y, sbR.thumbSiz, sbR.thumbSiz)
		if len(usr.tex2dSave) > 0 {
			r = R(x, sbR.tabHr[0].Y+sbR.tabHr[0].H+uispc, sbR.thumbSiz, sbR.thumbSiz)
		}
		var rs []sdl.FRect
		switch sbR.tabHorizNum {
		case 0: //ADD NEW TEX
			if len(usr.tex2d) > 0 {
				for i := range len(usr.tex2d) + 1 {
					if i == len(usr.tex2d) {
						dAddFile(r, "Add texture", 1)
					} else {
						if usr.tex2d[i].wide {
							r.W = (sbR.thumbSiz * 2) + sbR.spc
							if r.X != ox {
								r.X = ox
								r.Y += r.H + sbR.spc
							}
						}
						if usr.tex2d[i].high {
							r.H = (sbR.thumbSiz * 2) + sbR.spc
						}
						if crrslice(r, rs) {
							if r.X == ox {
								r.Y += sbR.thumbSiz + sbR.spc
							}
						}
						rs = append(rs, r)
						drr(r, CA(SET.BGC, uiAlpha))
						dtexthumb(&usr.tex2d[i], r)

						if cms(r) {
							mInfoTXT("Left click texture select & left click again drop")
							if SET.TexPreviewHover {
								dTexPreview(&usr.tex2d[i])
							}
							if ms.lc && !drop2D {
								droptex2D = usr.tex2d[i]
								drop2D = true
							}
						}
						drrlw(r, 1, SET.UIC)
						if usr.tex2d[i].high && usr.tex2d[i].wide {
							r.Y += r.H + sbR.spc
							r.W = sbR.thumbSiz
							r.H = sbR.thumbSiz
						} else if usr.tex2d[i].high && !usr.tex2d[i].wide {
							r.H = sbR.thumbSiz
							r.X += r.W + sbR.spc
						} else if !usr.tex2d[i].high && usr.tex2d[i].wide {
							r.X = ox
							r.Y += r.H + sbR.spc
							r.W = sbR.thumbSiz
						} else {
							r.X += r.W + sbR.spc
							if r.X+r.W+sbR.spc > sbR.r.X+sbR.r.W {
								r.X = ox
								r.Y += r.H + sbR.spc
							}
						}
						if crrslice(r, rs) {
							if r.X == ox {
								r.Y += sbR.thumbSiz + sbR.spc
							}
						}

					}
				}
			} else {
				dAddFile(r, "Add texture", 1)
			}
		case 1: //SAVED TEX
			if len(usr.tex2dSave) > 0 {
				for i := range len(usr.tex2dSave) {
					if usr.tex2dSave[i].wide {
						r.W = (sbR.thumbSiz * 2) + sbR.spc
						if r.X != ox {
							r.X = ox
							r.Y += r.H + sbR.spc
						}
					}
					if usr.tex2dSave[i].high {
						r.H = (sbR.thumbSiz * 2) + sbR.spc
					}
					if crrslice(r, rs) {
						if r.X == ox {
							r.Y += sbR.thumbSiz + sbR.spc
						}
					}
					rs = append(rs, r)
					drr(r, CA(SET.BGC, uiAlpha))
					dtexthumb(&usr.tex2dSave[i], r)

					if cms(r) {
						mInfoTXT("Left click texture select & left click again drop")
						if SET.TexPreviewHover {
							dTexPreview(&usr.tex2dSave[i])
						}
						if ms.lc && !drop2D {
							droptex2D = usr.tex2dSave[i]
							drop2D = true
						}
					}
					drrlw(r, 1, SET.UIC)
					if usr.tex2dSave[i].high && usr.tex2dSave[i].wide {
						r.Y += r.H + sbR.spc
						r.W = sbR.thumbSiz
						r.H = sbR.thumbSiz
					} else if usr.tex2dSave[i].high && !usr.tex2dSave[i].wide {
						r.H = sbR.thumbSiz
						r.X += r.W + sbR.spc
					} else if !usr.tex2dSave[i].high && usr.tex2dSave[i].wide {
						r.X = ox
						r.Y += r.H + sbR.spc
						r.W = sbR.thumbSiz
					} else {
						r.X += r.W + sbR.spc
						if r.X+r.W+sbR.spc > sbR.r.X+sbR.r.W {
							r.X = ox
							r.Y += r.H + sbR.spc
						}
					}
					if crrslice(r, rs) {
						if r.X == ox {
							r.Y += sbR.thumbSiz + sbR.spc
						}
					}

				}
			}
		}
		//MARK: RIGHT SIDEBAR TEX HORIZ TABS
		if len(usr.tex2dSave) > 0 {
			switch sbR.tabHorizNum {
			case 0:
				if cms(sbR.tabHr[1]) {
					drr(sbR.tabHr[1], SET.UIC3)
					if ms.lc {
						sbR.tabHorizNum = 1
					}
				} else {
					drr(sbR.tabHr[1], SET.BGC)
				}
				drliwSides(sbR.tabHr[0], []int{1, 2}, 2, SET.UIC2)
				drliwSides(sbR.tabHr[1], []int{3}, 2, SET.UIC2)
			case 1:
				if cms(sbR.tabHr[0]) {
					drr(sbR.tabHr[0], SET.UIC3)
					if ms.lc {
						sbR.tabHorizNum = 0
					}
				} else {
					drr(sbR.tabHr[0], SET.BGC)
				}
				drliwSides(sbR.tabHr[1], []int{1, 4}, 2, SET.UIC2)
				drliwSides(sbR.tabHr[0], []int{3}, 2, SET.UIC2)
			}
			dtxtreccnt(fs, "+ Add New", sbR.tabHr[0], SET.UICtxt)
			dtxtreccnt(fs, "Saved", sbR.tabHr[1], SET.UICtxt)

		}
	case 2: //LIGHTS
		xtx := x
		ytx := y
		lineH := fonh(fs) + un8th
		contentW := (sbL.r.W - sbR.ic.siz) - sbL.spc*2

		// Add Buttons
		addTorchBtn := R(xtx, ytx, contentW, fonh(fs)+4)
		if cms(addTorchBtn) {
			drr(addTorchBtn, SET.UIC3)
			mInfoTXT("Add a flickering dungeon torch point light")
			if ms.lc {
				SceneLights = append(SceneLights, &Light3D{
					Name:      fmt.Sprintf("Torch %d", len(SceneLights)),
					Type:      LightPoint,
					Color:     Col.GoldWarm,
					Intensity: 3.0,
					Range:     10.0,
					WorldX:    0, WorldY: 1.0, WorldZ: 0,
					Flicker:       true,
					FlickerSpeed:  4.0,
					FlickerAmount: 0.25,
					AttachedObjID: -1,
				})
				ms.lc = false
			}
		} else {
			drr(addTorchBtn, SET.UIC2)
			drrl(addTorchBtn, SET.UIC3)
		}
		dtxtreccnt(fs, "+ Add Dungeon Torch", addTorchBtn, SET.UICtxt)
		ytx += lineH + 4

		// List active lights
		for i, l := range SceneLights {
			lightCard := R(xtx, ytx, contentW, fonh(fs)+6)
			isSel := selectedLightID == i

			if cms(lightCard) {
				if ms.lc {
					selectedLightID = i
					ms.lc = false
				}
			}

			cardCol := SET.UIC2
			if isSel {
				cardCol = SET.UIC3
			}
			drr(lightCard, cardCol)
			drrl(lightCard, Col.GrayDark)
			dtxt(fs, l.Name, xtx+8, ytx+3, Col.White)

			// Quick color indicator dot
			dr(xtx+contentW-24, ytx+4, 16, 16, l.Color)
			ytx += lineH + 2

			// If selected, show controls
			if isSel {
				paletteOpenBoxFOR(&l.Color, xtx+contentW-60, ytx, fonh(fs), l.Color, "Change light color")
				dtxt(fs, "Color: ", xtx, ytx, SET.UICtxt)
				ytx += lineH

				numStepperFloat(&l.Intensity, 0.2, 0.0, 20.0, xtx, ytx, contentW, fonh(fs)+4, false, "Intensity:", "Brightness of the light")
				ytx += lineH + 2

				numStepperFloat(&l.Range, 0.5, 0.5, 50.0, xtx, ytx, contentW, fonh(fs)+4, false, "Range:", "Distance the light travels")
				ytx += lineH + 2

				// Flicker Toggle
				dtxt(fs, "Flicker: ", xtx, ytx, SET.UICtxt)
				onoff(&l.Flicker, xtx+80, ytx, un8thx3, "Turn light flicker on/off")
				ytx += lineH + 4
			}
		}
	case 1: //3D MODELS
		r := R(x, y, sbR.thumbSiz, sbR.thumbSiz)
		if len(usr.o3d) > 0 {
			for i := range len(usr.o3d) + 1 {
				if i == len(usr.o3d) {
					dAddFile(r, "Add 3D models", 0)
				} else {
					drr(r, CA(SET.BGC, uiAlpha))
					dtexscale1(usr.o3d[i].thumb, r.X, r.Y)
					if cms(r) {
						mInfoTXT("Left click model select & left click again drop")
						if ms.lc && !drop3D {
							dropO3D = usr.o3d[i]
							drop3D = true
						}
					}
					drrlw(r, 1, SET.UIC)
					r.X += r.W + sbR.spc
					if r.X+r.W+sbR.spc > sbR.r.X+sbR.r.W {
						r.X = ox
						r.Y += r.H + sbR.spc
					}

				}
			}
		} else {
			dAddFile(r, "Add 3D models", 0)
		}
	case 0: //3D SHAPES
		r := R(x, y, sbR.thumbSiz, sbR.thumbSiz)
		//Draw thumbnails
		if len(prim3DwinOBJ) > 0 {
			for i := range len(prim3DwinOBJ) {
				drr(r, CA(SET.BGC, uiAlpha))
				dtexscale1(prim3DwinOBJ[i].thumb, r.X, r.Y)
				if cms(r) {
					mInfoTXT(prim3DwinOBJ[i].nm + ": Left click model select & left click again drop")
					if ms.lc && !drop3D {
						dropO3D = prim3DwinOBJ[i]
						drop3D = true
					}
				}
				drrlw(r, 1, SET.UIC)
				r.X += r.W + sbR.spc
				if r.X+r.W+sbR.spc > sbR.r.X+sbR.r.W {
					r.X = ox
					r.Y += r.H + sbR.spc
				}
			}
		}
	}

	sbR.ys = dscrollY(y, un8th, SET.SCRW-scrollBarSiz, 0, scrollBarSiz, SET.SCRH-uiIconsSiz)

}
func dSidebarL() { // MARK: LEFT SIDEBAR
	if cms(sbL.r) {
		inSBL = true
	} else {
		inSBL = false
	}
	drr(sbL.r, SET.UIC) //BG REC
	var xtx, ytx float32
	xtx = sbL.r.X + sbL.spc
	ytx = sbL.r.Y + sbL.spc
	xmid := sbL.r.X + sbL.r.W/2
	x3q := sbL.r.X + ((sbL.r.W / 4) * 3)
	x3f := sbL.r.X + ((sbL.r.W / 5) * 3)

	switch sbLmenuNum {

	case -1: //METRICS
		dtxtSpans(fs, xtx, ytx,
			TxtSpan{"GPU Backend: ", SET.UICtxt},
			TxtSpan{gpuDev.Driver(), SET.UIC4},
		)
		ytx += fonh(fs)
		dtxtSpans(fs, xtx, ytx,
			TxtSpan{"Mouse ", SET.UICtxt},
			TxtSpan{" X: " + fmt.Sprintf("%.0f", ms.X), SET.UIC4},
			TxtSpan{"  | ", SET.UICtxt},
			TxtSpan{" Y: " + fmt.Sprintf("%.0f", ms.Y), SET.UIC4},
		)
		ytx += fonh(fs)
		if len(lev3d.o3d) > 0 {
			dtxt(fs, "Drawn (Visible) 3D: "+fmt.Sprint(drawn3Dcount), xtx, ytx, SET.UICtxt)
			ytx += fonh(fs)
			dtxt(fs, "Culled (Invisible) 3D: "+fmt.Sprint(culled3Dcount), xtx, ytx, SET.UICtxt)
			ytx += fonh(fs)
		}
	case 0: //NEW GRID
		dtxt(fd, "New Grid", xtx, ytx, SET.UIC4)
		ytx += fonh(fd) + uispc
		dtxt(fs, "Rows: ", xtx, ytx, SET.UICtxt)
		keysInpInt(fs, &gridVal.rows, xmid, ytx, un8th, 2, 10000, "Number of horizontal rows")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Columns: ", xtx, ytx, SET.UICtxt)
		keysInpInt(fs, &gridVal.columns, xmid, ytx, un8th, 2, 10000, "Number of vertical colums")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Block Width: ", xtx, ytx, SET.UICtxt)
		keysInpF32(fs, &gridVal.blockW, xmid, ytx, un8th, 4, 1024, "Width of each grid block")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Block Height: ", xtx, ytx, SET.UICtxt)
		keysInpF32(fs, &gridVal.blockH, xmid, ytx, un8th, 4, 1024, "Height of each grid block")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Line Color: ", xtx, ytx, SET.UICtxt)
		paletteAlphaOpenBoxFOR(&gridVal.c, xmid, ytx, un8thx3, gridVal.c, "Grid line color")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Hover Color: ", xtx, ytx, SET.UICtxt)
		paletteAlphaOpenBoxFOR(&gridVal.cHover, xmid, ytx, un8thx3, gridVal.cHover, "Grid block mouse hover color")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Hover Visible: ", xtx, ytx, SET.UICtxt)
		onoff(&gridVal.hoverVis, xmid, ytx, un8thx3, "Highlight grid block on mouse hover")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Hidden: ", xtx, ytx, SET.UICtxt)
		onoff(&gridVal.hidden, xmid, ytx, un8thx3, "Grid is only visible when 'Draw Hidden' is active")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Name: ", xtx, ytx, SET.UICtxt)
		ytx += fonh(fs) + uispc
		keysInpTxt(fs, &gridVal.nm, xtx, ytx, sbL.r.W-(sbL.spc*2), un8th, false, 20, "Name the grid (can be left blank)")
		ytx += fonh(fs) + uispc
		dtxt(fs, "Draw from: ", xtx, ytx, SET.UICtxt)
		ytx += fonh(fs) + uispc
		if !dropGrid {
			dropGridFrom = chooseListVert(fs, xtx+unq, ytx, un8th, []string{"Center", "Top Left XY"}, []string{"Draw grid from center at mouse drop position", "Draw grid from top left at mouse drop position"})
			if dropGridFrom != -1 {
				dropTimer = int(timer.FPS / 4)
				dropGrid = true
			}
		}
		debugnum = dropGridFrom
		if dropGrid {
			if gridVal.blockW == 0 || gridVal.blockH == 0 || gridVal.rows == 0 || gridVal.columns == 0 {
				mWarnTXT("Zero values in grid, fill all values")
				dropGrid = false
			} else {
				isCursorActiv = true
				mInfoTXT("Left click to postion grid")
				if dropTimer == 0 && ms.lc {
					ms.lc = false
					switch dropGridFrom {
					case 0: //CENTER
						lev2d.gr = append(lev2d.gr, mGridfromGridVal(mslcP, true))
					case 1: //TOP LEFT
						lev2d.gr = append(lev2d.gr, mGridfromGridVal(mslcP, false))
					}
					dropGridFrom = -1
					isCursorActiv = false
					dropGrid = false
				}

			}
		}

	case 1: //3D OBJECTS
		if selOBJ3Dnum >= 0 && selOBJ3Dnum < len(lev3d.o3d) {
			obj := &lev3d.o3d[selOBJ3Dnum]

			lineH := fonh(fs) + un8th
			contentW := sbL.r.W - sbL.spc*2

			// 1. Name
			dtxt(fd, obj.nm, xtx, ytx, SET.UIC4)
			ytx += lineH

			// 2. Color Swatch
			dtxt(fs, "Color: ", xtx, ytx, SET.UICtxt)
			x2 := xtx + txW(fs, "Color: ") + un8th
			paletteOpenBoxFOR(&obj.Color, x2, ytx, fonh(fs), obj.Color, "Change object color")
			ytx += lineH

			// 3. Texture Picker & Remove Buttons
			texBtn := R(xtx, ytx, contentW, fonh(fs)+4)
			if cms(texBtn) {
				drr(texBtn, SET.UIC3)
				mInfoTXT("Select a 2D Texture (PNG/JPG) for this model")
				if ms.lc {
					openNativeFileDialog("object_texture", 0)
					ms.lc = false
				}
			} else {
				drr(texBtn, SET.UIC2)
				drrl(texBtn, SET.UIC3)
			}

			texLabel := "Texture: [ None ] -> Browse..."
			if obj.CustomTex != nil {
				texLabel = "Texture: Applied [ Change... ]"
			}
			dtxtreccnt(fs, texLabel, texBtn, SET.UICtxt)
			ytx += lineH + 2

			// If texture is applied, show UV Tiling stepper & Remove button
			if obj.CustomTex != nil {
				// UV Tiling Stepper (1.0 = Stretch/1x, 2.0 = Tile 2x2, etc.)
				numStepperFloat(&obj.UVScale, 0.25, 0.25, 16.0, xtx, ytx, contentW, fonh(fs)+4, false, "UV Tile:", "Size of texture tile")
				ytx += lineH + 2

				// Remove texture button
				rmTexBtn := R(xtx, ytx, contentW, fonh(fs)+4)
				if cms(rmTexBtn) {
					drr(rmTexBtn, Col.RedBright)
					if ms.lc {
						obj.CustomTex = nil
						obj.UVScale = 1.0
						ms.lc = false
					}
				} else {
					drr(rmTexBtn, SET.UIC2)
					drrl(rmTexBtn, Col.RedBright)
				}
				dtxtreccnt(fs, "Remove Texture", rmTexBtn, Col.White)
				ytx += lineH + 2
			}

			// 4. Object Scale Stepper
			numStepperFloat(&obj.Scale, 0.05, 0.05, 10.0, xtx, ytx, contentW, fonh(fs)+4, false, "Scale:", "Increase/decrease object size")
			ytx += lineH + 4

			// 5. Rotation Speed Steppers (X, Y, Z)
			numStepperFloat(&obj.RotSpdX, 0.1, -5.0, 5.0, xtx, ytx, contentW, fonh(fs)+4, false, "RotSpd X:", "Rotation speed in the X axis")
			ytx += lineH + 2
			numStepperFloat(&obj.RotSpdY, 0.1, -5.0, 5.0, xtx, ytx, contentW, fonh(fs)+4, false, "RotSpd Y:", "Rotation speed in the Y axis")
			ytx += lineH + 2
			numStepperFloat(&obj.RotSpdZ, 0.1, -5.0, 5.0, xtx, ytx, contentW, fonh(fs)+4, false, "RotSpd Z:", "Rotation speed in the Z axis")
			ytx += lineH + 4

			// 6. Stop Rotation Button
			stopBtn := R(xtx, ytx, contentW, fonh(fs)+4)
			if cms(stopBtn) {
				drr(stopBtn, SET.UIC3)
				if ms.lc {
					obj.RotSpdX = 0
					obj.RotSpdY = 0
					obj.RotSpdZ = 0
					ms.lc = false
				}
			} else {
				drr(stopBtn, SET.UIC2)
				drrl(stopBtn, SET.UIC3)
			}
			dtxtreccnt(fs, "Stop Auto Rotation", stopBtn, SET.UICtxt)
		}
	case 2: //TEXTURES
		if selTEXnum == -1 {
			sbLmenuNum = -1
		}
		if selTEXnum >= 0 && selTEXnum < len(lev2d.tex) {
			var t *TEX2D
			t = &lev2d.tex[selTEXnum]

			ytx = sbL.ys

			dtxt(fd, t.nm, xtx, ytx, SET.UIC4) //TEX NAME
			ytx += fonh(fd) + uispc
			dtxt(fs, "Width: "+fmt.Sprint(t.W)+" Height: "+fmt.Sprint(t.H)+" px", xtx, ytx, SET.UICtxt)
			ytx += fonh(fs) + uispc

			dSBqicons(xtx, ytx) //SB QUICK ICONS
			ytx += fonh(id2) + uispc*2

			dtxt(fs, "Rotation Step:", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &roStep, x3f, ytx, 4, 0.01, 10, "Change increment for setting texture rotation")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Scale Step: ", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &scaleStep, x3f, ytx, 4, 0.01, 10, "Change increment for increasing & decreasing scale value")
			ytx += fonh(fs) + uispc
			numStepperFloat(&t.ro, roStep, -360, 360, xtx, ytx, sbL.r.W-(uispc*2), fonh(fs)+4, false, "Rotation: ", "Rotation angle of texture around center point")
			ytx += fonh(fs) + uispc
			numStepperFloat(&t.scale, scaleStep, 0.01, 100, xtx, ytx, sbL.r.W-(uispc*2), fonh(fs)+4, false, "Scale: ", "Increase/decrease size")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Color: ", xtx, ytx, SET.UICtxt)
			paletteAlphaOpenBoxFOR(&t.c, xmid, ytx, un8thx3, t.c, "Texture draw color")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Shadow: ", xtx, ytx, SET.UICtxt)
			onoff(&t.shadow, xmid, ytx, un8thx3, "Draw shadow image of texture")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Shadow Color: ", xtx, ytx, SET.UICtxt)
			paletteAlphaOpenBoxFOR(&t.cShadow, xmid, ytx, un8thx3, t.cShadow, "Set texture shadow color")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Shadow X: ", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &t.shadowX, x3f, ytx, 4, -1000, 1000, "Change X position of shadow texture")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Shadow Y: ", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &t.shadowY, x3f, ytx, 4, -1000, 1000, "Change Y position of shadow texture")
			ytx += fonh(fs) + uispc
			dtxt(fs, "FlipX Left/Right: ", xtx, ytx, SET.UICtxt)
			onoff(&t.flipX, x3q, ytx, un8thx3, "Flips texture left/right")
			ytx += fonh(fs) + uispc
			dtxt(fs, "FlipY Top/Bottom: ", xtx, ytx, SET.UICtxt)
			onoff(&t.flipY, x3q, ytx, un8thx3, "Flips texture top/bottom")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Color Effect: ", xtx, ytx, SET.UICtxt)
			listStepperU8(&t.fx, xmid-un8th, ytx, un8thx3, []string{"None", "Grayscale", "Sepia", "Posterize", "Scanlines"})
			if t.fx == 3 {
				ytx += fonh(fs) + uispc
				dtxt(fs, "Colors Bands: ", xtx, ytx, SET.UICtxt)
				keysInpF32(fs, &t.fxParam, xmid, ytx, 4, 1, 10, "Adjust posterization level")
			}
			ytx += fonh(fs) + uispc
			dtxt(fs, "Auto Create Frames: ", xtx, ytx, SET.UICtxt)
			mFrames := clickOnceSwitch(x3q, ytx, un8thx3, "Generate frames for animations automatically")
			if mFrames {
				t.frames, t.colFrames = AutoSliceSpriteSheet(t.path, t.tex)
				t.dAllFrames = true
			}
			ytx += fonh(fs) + uispc
			if len(t.frames) > 1 { //MARK: SELECT FRAMES
				dtxt(fs, "Draw All Frames: ", xtx, ytx, SET.UICtxt)
				onoff(&t.dAllFrames, x3q, ytx, un8thx3, "Draw all frames of sprite sheet")
				ytx += fonh(fs) + uispc
				dtxt(fs, "Draw Collision Boxes: ", xtx, ytx, SET.UICtxt)
				onoff(&collisFramesVis, x3q, ytx, un8thx3, "Draw auto generated collision boxes")
				ytx += fonh(fs) + uispc
				if t.dAllFrames {
					dtxt(fs, "Select Anim Frames: ", xtx, ytx, SET.UICtxt)
					prevSelectFrame := selectFrame
					onoff(&selectFrame, x3q, ytx, un8thx3, "Select frames for animation")
					ytx += fonh(fs) + uispc
					if prevSelectFrame != selectFrame && selectFrame {
						frameListTemp = nil
						animTemp = ANIM2D{
							playing: true,
							loop:    true,
							fps:     10,
							cbg:     Col.BlackCharcoal,
						}
						isCursorActiv = true

					} else if prevSelectFrame != selectFrame && !selectFrame {
						isCursorActiv = false
					}
					if len(frameListTemp) > 0 {
						maxX := sbL.r.X + sbL.r.W
						lineH := fonh(fs) + uispc
						curX := xtx

						var line strings.Builder

						for _, v := range frameListTemp {
							txt := strconv.Itoa(v) + ", "
							w := txW(fs, txt)

							// If adding this word exceeds the border, draw the current line and wrap
							if curX+w > maxX && line.Len() > 0 {
								dtxt(fs, line.String(), xtx, ytx, SET.UICtxt) // 🚀 1 draw call for the whole line!
								ytx += lineH
								curX = xtx
								line.Reset()
							}

							line.WriteString(txt)
							curX += w
						}

						// Flush remaining text on the final line
						if line.Len() > 0 {
							dtxt(fs, line.String(), xtx, ytx, SET.UICtxt)
						}
						ytx += fonh(fs) + uispc
					}
					if len(frameListTemp) > 1 { //MARK: DRAW ANIM PREVIEW
						dtxt(fs, "New Animation Name: ", xtx, ytx, SET.UICtxt)
						ytx += fonh(fs) + uispc
						keysInpTxt(fs, &animTemp.nm, uispc, ytx, sbL.innerW, 4, false, 20, "Animation name, can be left blank")
						ytx += fonh(fs) + uispc
						dtxt(fs, "FPS: ", xtx, ytx, SET.UICtxt)
						keysInpF32(fs, &animTemp.fps, x3f, ytx, 4, 1, 120, "Animation speed, number of frames displayed per second")
						ytx += fonh(fs) + uispc
						dtxt(fs, "Loop: ", xtx, ytx, SET.UICtxt)
						onoff(&animTemp.loop, x3q, ytx, un8thx3, "Loop animation continuosly")
						ytx += fonh(fs) + uispc
						dtxt(fs, "Preview Animation: ", xtx, ytx, SET.UICtxt)
						prevAnimPreviewVis := animPreviewVis
						onoff(&animPreviewVis, x3q, ytx, un8thx3, "Show animation preview of selected frames")
						if prevAnimPreviewVis != animPreviewVis {
							animTemp.frames = nil
							for _, framenum := range frameListTemp {
								animTemp.frames = append(animTemp.frames, t.frames[framenum])
							}
							prevAnimPreviewVis = animPreviewVis
						}
						ytx += fonh(fs) + uispc
						if animPreviewVis {
							dtxt(fs, "Preview Background: ", xtx, ytx, SET.UICtxt)
							paletteAlphaOpenBoxFOR(&t.cShadow, x3q, ytx, un8thx3, t.cShadow, "Set animation preview background color")
							ytx += fonh(fs) + uispc
							if len(animTemp.frames) > 0 {
								rbg := recfromcntr(P(xmid, ytx+(sbL.innerW/2)+uispc), sbL.innerW, sbL.innerW)
								drr(rbg, animTemp.cbg)
								dAnimPreview(t, P(xmid, ytx+un2), &animTemp, sbL.innerW)
								ytx += sbL.innerW + uispc*2
							}
						}

					}

				}
			}

			if len(lev2d.tex) > 1 { //MOVE BACK FORWARD
				dtxt(fs, "Move to Front: ", xtx, ytx, SET.UICtxt)
				move2front := clickOnceSwitch(x3q, ytx, un8thx3, "Move draw texture to foreground")
				if move2front {
					selTEXnum = MoveTexToFront(lev2d.tex, selTEXnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move to Back: ", xtx, ytx, SET.UICtxt)
				move2back := clickOnceSwitch(x3q, ytx, un8thx3, "Move draw texture to background")
				if move2back {
					selTEXnum = MoveTexToEnd(lev2d.tex, selTEXnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move Forward 1 Place: ", xtx, ytx, SET.UICtxt)
				move1forward := clickOnceSwitch(x3q, ytx, un8thx3, "Move draw texture 1 position forward")
				if move1forward {
					selTEXnum = MoveTexStepForward(lev2d.tex, selTEXnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move Back 1 Place: ", xtx, ytx, SET.UICtxt)
				move1back := clickOnceSwitch(x3q, ytx, un8thx3, "Move draw texture 1 position backward")
				if move1back {
					selTEXnum = MoveTexStepBackward(lev2d.tex, selTEXnum)
				}
				ytx += fonh(fs) + uispc
			}

			if t.scale != 1 { //SCALED WIDTH HEIGHT
				dtxt(fs, "Scaled Width: "+fmt.Sprintf("%.2f", t.W*t.scale)+"px", xtx, ytx, SET.UICtxt)
				ytx += fonh(fs) + uispc
				dtxt(fs, "Scaled Height: "+fmt.Sprintf("%.2f", t.H*t.scale)+" px", xtx, ytx, SET.UICtxt)
				ytx += fonh(fs) + uispc
				if len(t.frames) > 0 {
					dtxt(fs, "Scaled Frame Width: "+fmt.Sprintf("%.2f", t.frames[0].W*t.scale)+"px", xtx, ytx, SET.UICtxt)
					ytx += fonh(fs) + uispc
					dtxt(fs, "Scaled Frame Height: "+fmt.Sprintf("%.2f", t.frames[0].H*t.scale)+" px", xtx, ytx, SET.UICtxt)
					ytx += fonh(fs) + uispc
				}
			}
			sbL.ys = dscrollY(sbL.ys, un8th, sbL.r.X+sbL.r.W-scrollBarSiz, 0, scrollBarSiz, SET.SCRH-uiIconsSiz) //MARK:LEFT SCROLL BAR
			if sbL.ys > 0 {
				sbL.ys = 0
			}
		}
	case 3: //2D LINE
		if selLINE2Dnum != -1 && len(lev2d.li) > 0 {
			if nameEdit {
				keysInpTxt(fd, &lev2d.li[selLINE2Dnum].nm, xtx, ytx, sbL.innerW, 12, false, 20, "Name the line")
			} else {
				if lev2d.li[selLINE2Dnum].nm != "" {
					dtxt(fd, "Name: "+lev2d.li[selLINE2Dnum].nm, xtx, ytx, SET.UIC4)
				} else {
					dtxt(fd, "Line: "+fmt.Sprint(selLINE2Dnum), xtx, ytx, SET.UIC4)
				}
				if cms(sbL.rEditName) {
					mInfoTXT("Edit line name")
					drr(sbL.rEditName, SET.UIC3)
					if ms.lc {
						nameEdit = true
					}
				}
				dtxtreccnt(id, iconPencil, sbL.rEditName, SET.UICtxt)
			}
		} else {
			dtxt(fd, "New Line", xtx, ytx, SET.UIC4)
		}
		ytx += fonh(fd) + uispc*2
		if selLINE2Dnum != -1 && len(lev2d.li) > 0 {
			l := &lev2d.li[selLINE2Dnum]
			dtxt(fs, "Color: ", xtx, ytx, SET.UICtxt)
			paletteAlphaOpenBoxFOR(&l.c, xmid, ytx, un8thx3, l.c, "Line draw color")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Line Width Step:", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &lineWstep, x3f, ytx, 4, 0.5, 2, "Change increment for setting adjusting line width")
			ytx += fonh(fs) + uispc
			numStepperFloat(&l.w, lineWstep, 0.5, 20, xtx, ytx, sbL.innerW, fonh(fs)+4, false, "Line Width:", "Line thickness")
			ytx += fonh(fs) + uispc
			dtxt2c(fs, "Total Length: ", fmt.Sprint(l.l), SET.UICtxt, SET.UIC4, xtx, ytx)
			ytx += fonh(fs) + uispc
			for i := range l.p {
				dtxtSpans(fs, xtx, ytx,
					TxtSpan{"Point " + fmt.Sprint(i) + ": ", SET.UICtxt},
					TxtSpan{" X: " + fmt.Sprintf("%.0f", l.p[i].X), SET.UIC4},
					TxtSpan{"  | ", SET.UICtxt},
					TxtSpan{" Y: " + fmt.Sprintf("%.0f", l.p[i].Y), SET.UIC4},
				)
				ytx += fonh(fs) + uispc
			}
			for i := range l.cnt {
				dtxtSpans(fs, xtx, ytx,
					TxtSpan{"Midpoint " + fmt.Sprint(i) + ": ", SET.UICtxt},
					TxtSpan{" X: " + fmt.Sprintf("%.0f", l.cnt[i].X), SET.UIC4},
					TxtSpan{"  | ", SET.UICtxt},
					TxtSpan{" Y: " + fmt.Sprintf("%.0f", l.cnt[i].Y), SET.UIC4},
				)
				ytx += fonh(fs) + uispc
			}
		}
	case 4: //2D RECTANGLES
		if selREC2Dnum != -1 && len(lev2d.rec) > 0 {
			if nameEdit {
				keysInpTxt(fd, &lev2d.rec[selREC2Dnum].nm, xtx, ytx, sbL.innerW, 12, false, 20, "Name the rectangle")
			} else {
				if lev2d.rec[selREC2Dnum].nm != "" {
					dtxt(fd, "Name: "+lev2d.rec[selREC2Dnum].nm, xtx, ytx, SET.UIC4)
				} else {
					dtxt(fd, "Rectangle: "+fmt.Sprint(selREC2Dnum), xtx, ytx, SET.UIC4)
				}
				if cms(sbL.rEditName) {
					mInfoTXT("Edit rectangle name")
					drr(sbL.rEditName, SET.UIC3)
					if ms.lc {
						nameEdit = true
					}
				}
				dtxtreccnt(id, iconPencil, sbL.rEditName, SET.UICtxt)
			}
		} else {
			dtxt(fd, "New Rectangle", xtx, ytx, SET.UIC4)
		}
		ytx += fonh(fd) + uispc*2
		if selREC2Dnum != -1 && len(lev2d.rec) > 0 {
			r := &lev2d.rec[selREC2Dnum]
			rEditDim := R(x3f, ytx, un8thx3, un8thx3)
			if widthEdit {
				wt, _ := dtxt(fs, "New Width: ", xtx, ytx, SET.UICtxt)
				keysInpF32(fs, &lev2d.rec[selREC2Dnum].r.W, xtx+wt+uispc, ytx, 8, 2, 50000, "New 2D rectangle width")
			} else {
				wt, _ := dtxt2c(fs, "Width: ", fmt.Sprintf("%.0f", r.r.W), SET.UICtxt, SET.UIC4, xtx, ytx)
				rEditDim.X = xtx + uispc + wt
				if cms(rEditDim) {
					drr(rEditDim, SET.UIC3)
					mInfoTXT("Left click edit 2D rectangle width")
					if ms.lc {
						widthEdit = true
						ms.lc = false
					}
				}

				dtxtreccnt(is, iconPencil, rEditDim, SET.UICtxt)
			}
			ytx += fonh(fs) + uispc
			if heightEdit {
				wt, _ := dtxt(fs, "New Height: ", xtx, ytx, SET.UICtxt)
				keysInpF32(fs, &lev2d.rec[selREC2Dnum].r.H, xtx+wt+uispc, ytx, 8, 2, 50000, "New 2D rectangle height")
			} else {
				wt, _ := dtxt2c(fs, "Height: ", fmt.Sprintf("%.0f", r.r.H), SET.UICtxt, SET.UIC4, xtx, ytx)
				rEditDim = R(xtx+uispc+wt, ytx, un8thx3, un8thx3)
				if cms(rEditDim) {
					drr(rEditDim, SET.UIC3)
					mInfoTXT("Left click edit 2D rectangle height")
					if ms.lc {
						heightEdit = true
						ms.lc = false
					}
				}
				dtxtreccnt(is, iconPencil, rEditDim, SET.UICtxt)
			}
			ytx += fonh(fs) + uispc

			dtxtSpans(fs, xtx, ytx,
				TxtSpan{"Width: ", SET.UICtxt},
				TxtSpan{fmt.Sprintf("%.0f", r.r.W), SET.UIC4},
				TxtSpan{"  |  Height: ", SET.UICtxt},
				TxtSpan{fmt.Sprintf("%.0f", r.r.H), SET.UIC4},
				TxtSpan{" px", SET.UICtxt},
			)
			ytx += fonh(fs) + uispc
			dtxt(fs, "Draw Outline: ", xtx, ytx, SET.UICtxt)
			onoff(&r.outline, x3f, ytx, un8thx3, "Draw perimeter outline")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Draw Fill: ", xtx, ytx, SET.UICtxt)
			onoff(&r.fill, x3f, ytx, un8thx3, "Draw solid rectangle")
			ytx += fonh(fs) + uispc
			if r.outline {
				dtxt(fs, "Outline Color: ", xtx, ytx, SET.UICtxt)
				paletteAlphaOpenBoxFOR(&r.cLine, xmid, ytx, un8thx3, r.cLine, "Rectangle fill color")
				ytx += fonh(fs) + uispc
				numStepperFloat(&r.lineW, 0.5, 1, 20, xtx, ytx, sbL.innerW, fonh(fs)+4, false, "Outline Width: ", "Thickness of rectangle outline")
				ytx += fonh(fs) + uispc
			}
			if r.fill {
				dtxt(fs, "Fill Type: ", xtx, ytx, SET.UICtxt)
				listStepperINT(&r.fillGrad, xmid-unh, ytx, un8thx3, []string{"Solid", "Gradient Horiz", "Gradient Vert", "Radial", "4 Corner"})
				ytx += fonh(fs) + uispc
				if r.fillGrad == 0 {
					dtxt(fs, "Fill Color: ", xtx, ytx, SET.UICtxt)
					paletteAlphaOpenBoxFOR(&r.cFill, xmid, ytx, un8thx3, r.cFill, "Rectangle draw color")
					ytx += fonh(fs) + uispc
				} else {
					dtxt(fs, "Gradient Color 1: ", xtx, ytx, SET.UICtxt)
					paletteAlphaOpenBoxFOR(&r.cGrad1, xmid, ytx, un8thx3, r.cGrad1, "Rectangle gradient fill color 1")
					ytx += fonh(fs) + uispc
					dtxt(fs, "Gradient Color 2: ", xtx, ytx, SET.UICtxt)
					paletteAlphaOpenBoxFOR(&r.cGrad2, xmid, ytx, un8thx3, r.cGrad2, "Rectangle gradient fill color 2")
					ytx += fonh(fs) + uispc
					if r.fillGrad == 4 {
						dtxt(fs, "Gradient Color 3: ", xtx, ytx, SET.UICtxt)
						paletteAlphaOpenBoxFOR(&r.cGrad3, xmid, ytx, un8thx3, r.cGrad3, "Rectangle gradient fill color 3")
						ytx += fonh(fs) + uispc
						dtxt(fs, "Gradient Color 4: ", xtx, ytx, SET.UICtxt)
						paletteAlphaOpenBoxFOR(&r.cGrad4, xmid, ytx, un8thx3, r.cGrad4, "Rectangle gradient fill color 4")
						ytx += fonh(fs) + uispc
					}
				}
			}
			dtxt(fs, "Rotation Step:", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &roStep, x3f, ytx, 4, 0.01, 10, "Change increment for rectangle rotation")
			ytx += fonh(fs) + uispc
			numStepperFloat(&r.ro, roStep, -360, 360, xtx, ytx, sbL.r.W-(uispc*2), fonh(fs)+4, false, "Rotation: ", "Rotation angle of rectangle around center point")
			ytx += fonh(fs) + uispc
			if len(lev2d.rec) > 1 { //MOVE BACK FORWARD
				dtxt(fs, "Move to Front: ", xtx, ytx, SET.UICtxt)
				move2front := clickOnceSwitch(x3q, ytx, un8thx3, "Move 2D rectangle to foreground")
				if move2front {
					selREC2Dnum = MoveRec2DtoFront(lev2d.rec, selREC2Dnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move to Back: ", xtx, ytx, SET.UICtxt)
				move2back := clickOnceSwitch(x3q, ytx, un8thx3, "Move 2D rectangle to background")
				if move2back {
					selREC2Dnum = MoveRec2DtoEnd(lev2d.rec, selREC2Dnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move Forward 1 Place: ", xtx, ytx, SET.UICtxt)
				move1forward := clickOnceSwitch(x3q, ytx, un8thx3, "Move 2D rectangle 1 position forward")
				if move1forward {
					selREC2Dnum = MoveRec2DstepForward(lev2d.rec, selREC2Dnum)
				}
				ytx += fonh(fs) + uispc
				dtxt(fs, "Move Back 1 Place: ", xtx, ytx, SET.UICtxt)
				move1back := clickOnceSwitch(x3q, ytx, un8thx3, "Move 2D rectangle 1 position backward")
				if move1back {
					selREC2Dnum = MoveRec2DstepBackward(lev2d.rec, selREC2Dnum)
				}
				ytx += fonh(fs) + uispc
			}
		}
	case 5: //2D TRIANGLES
		if selTRI2Dnum != -1 && len(lev2d.tri) > 0 {
			if nameEdit {
				keysInpTxt(fd, &lev2d.tri[selTRI2Dnum].nm, xtx, ytx, sbL.innerW, 12, false, 20, "Name the triangle")
			} else {
				if lev2d.tri[selTRI2Dnum].nm != "" {
					dtxt(fd, "Name: "+lev2d.tri[selTRI2Dnum].nm, xtx, ytx, SET.UIC4)
				} else {
					dtxt(fd, "Triangle: "+fmt.Sprint(selTRI2Dnum), xtx, ytx, SET.UIC4)
				}
				if cms(sbL.rEditName) {
					mInfoTXT("Edit triangle name")
					drr(sbL.rEditName, SET.UIC3)
					if ms.lc {
						nameEdit = true
					}
				}
				dtxtreccnt(id, iconPencil, sbL.rEditName, SET.UICtxt)
			}
		} else {
			dtxt(fd, "New Triangle", xtx, ytx, SET.UIC4)
		}
		ytx += fonh(fd) + uispc*2
		if selTRI2Dnum != -1 && len(lev2d.tri) > 0 {
			t := &lev2d.tri[selTRI2Dnum]
			rEditDim := R(x3f, ytx, un8thx3, un8thx3)
			if widthEdit {
				wt, _ := dtxt(fs, "New Side Width: ", xtx, ytx, SET.UICtxt)
				keysInpF32(fs, &t.sideW, xtx+wt+uispc, ytx, 8, 3, 50000, "New 2D triangle side")
			} else {
				wt, _ := dtxt2c(fs, "Side Width: ", fmt.Sprintf("%.0f", t.sideW), SET.UICtxt, SET.UIC4, xtx, ytx)
				rEditDim.X = xtx + uispc + wt
				if cms(rEditDim) {
					drr(rEditDim, SET.UIC3)
					mInfoTXT("Left click edit 2D triangle side width")
					if ms.lc {
						widthEdit = true
						ms.lc = false
					}
				}
				dtxtreccnt(is, iconPencil, rEditDim, SET.UICtxt)
			}
			ytx += fonh(fs) + uispc
			dtxt(fs, "Draw Outline: ", xtx, ytx, SET.UICtxt)
			onoff(&t.outline, x3f, ytx, un8thx3, "Draw triangle outline")
			ytx += fonh(fs) + uispc
			dtxt(fs, "Draw Fill: ", xtx, ytx, SET.UICtxt)
			onoff(&t.fill, x3f, ytx, un8thx3, "Draw triangle fill")
			ytx += fonh(fs) + uispc
			if t.outline {
				dtxt(fs, "Outline Color: ", xtx, ytx, SET.UICtxt)
				paletteAlphaOpenBoxFOR(&t.cLine, xmid, ytx, un8thx3, t.cLine, "Triangle fill color")
				ytx += fonh(fs) + uispc
				numStepperFloat(&t.lineW, 0.5, 1, 20, xtx, ytx, sbL.innerW, fonh(fs)+4, false, "Outline Width: ", "Thickness of triangle outline")
				ytx += fonh(fs) + uispc
			}
			if t.fill {
				dtxt(fs, "Fill Color: ", xtx, ytx, SET.UICtxt)
				paletteAlphaOpenBoxFOR(&t.cFill, xmid, ytx, un8thx3, t.cFill, "Rectangle draw color")
				ytx += fonh(fs) + uispc
			}
			ytx += fonh(fs) + uispc
			dtxt(fs, "Rotation Step:", xtx, ytx, SET.UICtxt)
			keysInpF32(fs, &roStep, x3f, ytx, 4, 0.01, 10, "Change increment for triangle rotation")
			ytx += fonh(fs) + uispc
			numStepperFloat(&t.ro, roStep, -360, 360, xtx, ytx, sbL.r.W-(uispc*2), fonh(fs)+4, false, "Rotation: ", "Rotation angle of triangle around center point")
			ytx += fonh(fs) + uispc
		}
	}

}

func dFooterIcons() { //MARK: DRAW FOOTER ICONS
	r := R(0, SET.SCRH-uiIconsSiz, uiIconsSiz, uiIconsSiz)
	qiconsBGR = R(0, SET.SCRH-uiIconsSiz, SET.SCRW, uiIconsSiz)
	drr(qiconsBGR, SET.UIC)
	for i := range qIcons {
		q := qIcons[i].nm
		for j, nm := range q {
			if cms(r) {
				drr(r, SET.UIC3)
				mInfoTXT(nm)
				if ms.lc {
					switch nm {
					case "Recyle Bin":
						recycleVis = true
					case "Draw Hidden":
						SET.DrawHidden = !SET.DrawHidden
					case "Settings":
						settingsVis = true
					case "Shaded":
						R3D.Mode = ModeShaded
					case "Wireframe":
						R3D.Mode = ModeWireframe
					case "Solid":
						R3D.Mode = ModeSolid
					case "Ambient Color":
						OpenPaletteFor(&R3D.AmbientColor)
					case "Lights":
						if sceneLightsActiv {
							LightsOff = make([]*Light3D, len(SceneLights))
							copy(LightsOff, SceneLights)
							SceneLights = nil
							sceneLightsActiv = false
						} else {
							SceneLights = make([]*Light3D, len(LightsOff))
							copy(SceneLights, LightsOff)
							LightsOff = nil
							sceneLightsActiv = true
						}
					}
					uQuickIcons()
				}
			}
			if qIcons[i].onoff[j] {
				dtxtreccnt(id, qIcons[i].ic[j], r, SET.UICgrn)
			} else {
				dtxtreccnt(id, qIcons[i].ic[j], r, SET.UICtxt)
			}
			r.X += r.W
		}
	}

}

// dAddFile: Draws the add file + rec, txInfo: hover text information, cat: filetype loaded
func dAddFile(r sdl.FRect, txInfo string, cat int) { //MARK: DRAW ADD FILE
	if cms(r) {
		drr(r, CA(Col.PinkHot, 100))
		mInfoTXT(txInfo)
		if ms.lc {
			switch cat {
			case 1: //TEXTURE
				openNativeFileDialog("texture", thumbWdef)
			case 0: //MODEL
				openNativeFileDialog("model", thumbWdef)
			}

		}
	}
	drrlw(r, 2, Col.PinkHot)
	dtxtreccnt(id, iconAdd, r, SET.UICtxt)
}
func dBGiconPattern(f *ttf.Font, ic string, pad float32, r sdl.FRect, c sdl.Color) { //MARK: DRAW BG ICON PATTERN
	var x, y = r.X, r.Y
	w, h := txWH(f, ic)
	w += pad * 2
	h += pad * 2
	r2 := R(x, y, w, h)
	for r2.Y < r.Y+r.H {
		dtxtreccnt(f, ic, r2, c)
		r2.X += r2.W
		if r2.X > r.X+r.W {
			r2.Y += r2.H
			r2.X = r.X
		}
	}
}
func dTexPreview(t *TEX2D) { //MARK: DRAW TEX PREVIEW

	if t.tex.W/t.tex.H > 10 || t.tex.H/t.tex.W > 10 {
		if t.tex.W > t.tex.H {
			w := viewportW - un2
			h := scaleH(t.tex.W, t.tex.H, w)
			if SET.TexPreviewBG {
				drcntrWH(CNT, w, h, SET.TexPreviewBGC)
			}
			dtexcenterpWH(t.tex, CNT, w, h, t.c)
		} else {
			h := viewportH - un2
			w := scaleW(t.tex.W, t.tex.H, h)
			if SET.TexPreviewBG {
				drcntrWH(CNT, w, h, SET.TexPreviewBGC)
			}
			dtexcenterpWH(t.tex, CNT, w, h, t.c)
		}
	} else if t.tex.H > viewportH/2 {
		h := viewportH - un2
		w := scaleW(t.tex.W, t.tex.H, h)
		if SET.TexPreviewBG {
			drcntrWH(CNT, w, h, SET.TexPreviewBGC)
		}
		dtexcenterpWH(t.tex, CNT, w, h, t.c)
	} else {
		scale := fitScale(t.tex.W, t.tex.H, un7)
		if SET.TexPreviewBG {
			drcntrWH(CNT, t.tex.W*scale, t.tex.H*scale, SET.TexPreviewBGC)
		}
		dtexcenterp(t.tex, CNT, scale)
	}

	/*
		if t.tex.W > viewportW || t.tex.H > viewportH {
			if t.tex.W > t.tex.H {
				w := viewportW - un2
				h := scaleH(t.tex.W, t.tex.H, w)
				dtexcenterpWH(t.tex, cntr, w, h)
			} else {
				h := viewportH - un2
				w := scaleW(t.tex.W, t.tex.H, h)
				dtexcenterpWH(t.tex, cntr, w, h)
			}
		} else {
			scale := fitScale(t.tex.W, t.tex.H, un7)
			dtexcenterp(t.tex, cntr, scale)
		}
	*/
}

// MARK: MAKE
func mUI() { //MARK:MAKE UI

	//SIDEBAR RIGHT
	siz := un4
	sbR.r = R(SET.SCRW-siz, 0, siz, SET.SCRH-uiIconsSiz)
	sbR.ic.siz = unh
	sbR.ic.nm = append(sbR.ic.nm, "3D Shapes", "3D Models", "Lights", "Textures", "2D")
	sbR.ic.ic = append(sbR.ic.ic, ic23d, ic23dmodels, ic2lights, ic2texture, ic22d)
	sbR.thumbSiz = un + unh
	sbR.spc = un16th
	sbR.ys = sbR.r.Y
	sbR.innerW = (sbR.r.W - sbR.ic.siz) - uispc*2
	sbR.tabHr = []sdl.FRect{R(sbR.r.X+sbR.ic.siz, sbR.r.Y, (sbR.innerW/2)-uispc, uiIconsSiz), R(sbR.r.X+sbR.ic.siz+((sbR.innerW/2)-uispc), sbR.r.Y, (sbR.innerW/2)-uispc, uiIconsSiz)}

	//sbR.tabVertNum = 4

	//2D GEOMETRY ICON LIST SIDEBAR RIGHT
	sbRGeomIcons.nm = append(sbRGeomIcons.nm, "Grid", "Line", "Polyline", "Rectangle", "Triangle", "Circle", "Line Shape", "Pentagon", "Hexagon", "Polygon")
	sbRGeomIcons.ic = append(sbRGeomIcons.ic, ic2grid, ic2line, ic2polyline, ic22d, ic2triangle, ic2circle, ic2lineshape, ic2pentagon, ic2hexagon, ic2polygon)
	sbRGeomIcons.siz = sbR.thumbSiz / 2

	//SIDEBAR LEFT
	siz = un4
	sbL.r = R(0, 0, siz, SET.SCRH-uiIconsSiz)
	sbL.rEditName = R(0+sbL.r.W-uiIconsSiz, 0, uiIconsSiz, uiIconsSiz)
	sbL.ic.siz = unh
	sbL.thumbSiz = un + unh
	sbL.spc = un16th
	sbL.ys = sbR.r.Y
	sbL.innerW = sbL.r.W - uispc*2

	//SIDEBAR ICONS
	sbQicons.ic = []string{ic2save, ic2delete}
	sbQicons.nm = []string{"Save", "Delete"}
	sbQicons.hovertxt = []string{"Saves current selection with edits", "Moves current selection to recycle bin"}
	sbQicons.siz = uiIconsSiz

	//QUICK ICONS
	//3D
	qIcons3D.ic = []string{iconDrawShaded, iconDrawSolid, iconDrawWire, iconSunlight, iconLightbulb, iconLightbulbSolid, iconAmbient, iconPalette}
	qIcons3D.nm = []string{"Shaded", "Solid", "Wireframe", "Sunlight", "Lights", "Light Visibility", "Ambient Color", "Interface Colors"}
	for range qIcons3D.ic {
		qIcons3D.onoff = append(qIcons3D.onoff, false)
	}
	//DEFAULT
	qIconsDEF.ic = []string{iconSettings, iconRecycle, iconVisible, iconlayers}
	qIconsDEF.nm = []string{"Settings", "Recycle Bin", "Draw Hidden", "Layers"}
	for range qIconsDEF.ic {
		qIconsDEF.onoff = append(qIconsDEF.onoff, false)
	}
	//ALL QICONS
	qIcons = append(qIcons, qIconsDEF, qIcons3D)

	uQuickIcons()

	//VIEWPORT
	viewportW, viewportH = SET.SCRW-(sbL.r.W+sbR.r.W), SET.SCRH-uiIconsSiz
}

// MARK: UPDATE
func uQuickIcons() {
	for i := range qIcons {
		q := qIcons[i]

		for j, nm := range q.nm {
			activ := false
			switch nm {
			case "Settings":
				if settingsVis {
					activ = true
				}
			case "Draw Hidden":
				if SET.DrawHidden {
					activ = true
				}
			case "Shaded":
				if R3D.Mode == ModeShaded {
					activ = true
				}
			case "Wireframe":
				if R3D.Mode == ModeWireframe {
					activ = true
				}
			case "Solid":
				if R3D.Mode == ModeSolid {
					activ = true
				}
			}
			if activ {
				q.onoff[j] = true
			} else {
				q.onoff[j] = false
			}
		}

	}
}

// MARK: UTILS
func listStepperINT(val *int, x, y, h float32, l []string) {
	w := txW(fs, l[0])
	for i := range l {
		if txW(fs, l[i]) > w {
			w = txW(fs, l[i])
		}
	}
	w += uispc
	rl := R(x, y, h, h)
	drr(rl, SET.UIC2)
	drrl(rl, SET.UIC3)
	dtxtreccnt(id2, ic2leftarrow, rl, SET.UICtxt)
	rm := R(rl.X+h, y, w, h)
	drr(rm, SET.BGC)
	rr := rl
	rr.X += h + w
	drr(rr, SET.UIC2)
	drrl(rr, SET.UIC3)
	dtxtreccnt(id2, ic2rightarrow, rr, SET.UICtxt)
	if cms(rl) {
		drrlw(rl, 2, SET.UIC4)
		if ms.lc {
			if *val > 0 {
				*val--
			}
			ms.lc = false
		}
	}
	if cms(rr) {
		drrlw(rr, 2, SET.UIC4)
		if ms.lc {
			if *val < len(l)-1 {
				*val++
			}
			ms.lc = false
		}
	}
	dtxtreccnt(fs, l[int(*val)], rm, SET.UICtxt)
}

func listStepperU8(val *uint8, x, y, h float32, l []string) {
	w := txW(fs, l[0])
	for i := range l {
		if txW(fs, l[i]) > w {
			w = txW(fs, l[i])
		}
	}
	w += uispc
	rl := R(x, y, h, h)
	drr(rl, SET.UIC2)
	drrl(rl, SET.UIC3)
	dtxtreccnt(id2, ic2leftarrow, rl, SET.UICtxt)
	rm := R(rl.X+h, y, w, h)
	drr(rm, SET.BGC)
	rr := rl
	rr.X += h + w
	drr(rr, SET.UIC2)
	drrl(rr, SET.UIC3)
	dtxtreccnt(id2, ic2rightarrow, rr, SET.UICtxt)
	if cms(rl) {
		drrlw(rl, 2, SET.UIC4)
		if ms.lc {
			if *val > 0 {
				*val--
			}
			ms.lc = false
		}
	}
	if cms(rr) {
		drrlw(rr, 2, SET.UIC4)
		if ms.lc {
			if *val < uint8(len(l)-1) {
				*val++
			}
			ms.lc = false
		}
	}
	dtxtreccnt(fs, l[int(*val)], rm, SET.UICtxt)
}

// chooseListVert: Input left click list select, returns int[i] of txt4nums[i]
func chooseListVert(f *ttf.Font, x, y, pad float32, listnames, hovertxt []string) int {
	choose := -1

	for i := range len(listnames) {
		w, h := txWH(f, listnames[i])
		w += pad * 2
		//h += pad * 2
		r := R(x, y, w, h)
		if cms(r) {
			drr(r, SET.UIC3)
			if ms.lc {
				choose = i
			}
			mInfoTXT(hovertxt[i])
		} else {
			drr(r, SET.BGC)
		}
		dtxtreccnt(f, listnames[i], r, SET.UICtxt)
		x += r.W + uispc

	}
	return choose
}
func dscrollY(yScroll, speedScroll, x, y, w, h float32) float32 {
	r2 := R(x, y, w, w)
	r3 := r2
	r3.Y = r2.Y + h - w
	if cms(r2) {
		if ms.lc || ms.ld {
			if yScroll < 0 {
				yScroll += speedScroll
				if yScroll > 0 {
					yScroll = 0
				}
			}
		}
		drr(r2, CA(SET.UIC3, 100))

	} else {
		drr(r2, CA(SET.UIC2, 100))
	}
	if cms(r3) {
		if ms.lc || ms.ld {
			yScroll -= speedScroll
		}
		drr(r3, CA(SET.UIC3, 100))
	} else {
		drr(r3, CA(SET.UIC2, 100))
	}

	dtxtreccnt(id2, ic2upArrow, r2, SET.UICtxt)
	dtxtreccnt(id2, ic2downArrow, r3, SET.UICtxt)
	return yScroll
}

func closeIcon(x, y, w float32) bool {
	close := false
	r := R(x, y, w, w)
	if cms(r) {
		drrl(r, SET.UICred)
		dtxtreccnt(id, iconClose, r, SET.UICred)
		if ms.lc {
			close = true
			ms.lc = false
		}
	} else {
		drrl(r, SET.UICtxt)
		dtxtreccnt(id, iconClose, r, SET.UICtxt)
	}
	return close
}

// numStepperFloat draws a label with [-] and [+] buttons to edit a float32 value
func numStepperFloat(val *float32, step, minVal, maxVal float32, x, y, w, h float32, nodecimal bool, label, hovertxt string) {
	if val == nil {
		return
	}

	btnSize := h
	spacing := float32(4)
	valBoxW := float32(60)

	// 1. Draw Label on the left
	dtxt(fs, label, x, y+(h-fonh(fs))/2, SET.UICtxt)

	// Layout positions on the right side of the window
	btnMinusX := x + w - (btnSize*2 + valBoxW + spacing*2)
	valBoxX := btnMinusX + btnSize + spacing
	btnPlusX := valBoxX + valBoxW + spacing

	btnMinus := R(btnMinusX, y, btnSize, btnSize)
	valBox := R(valBoxX, y, valBoxW, h)
	btnPlus := R(btnPlusX, y, btnSize, btnSize)

	// 2. Minus Button [-]
	drr(btnMinus, SET.BGC)
	if cms(btnMinus) {
		drrlw(btnMinus, 2, SET.UIC4)
		if ms.lc {
			*val -= step
			if *val < minVal {
				*val = minVal
			}
			ms.lc = false
		}
	} else {
		drrl(btnMinus, SET.UIC3)
	}
	dtxtreccnt(fs, "-", btnMinus, SET.UICtxt)

	// 3. Current Value Display
	drr(valBox, SET.BGC)
	// Round cleanly to 2 decimal places:
	valText := fmt.Sprintf("%.2f", *val)
	if nodecimal {
		valText = fmt.Sprintf("%.0f", *val)
	}
	dtxtreccnt(fs, valText, valBox, Col.White)

	// 4. Plus Button [+]
	drr(btnPlus, SET.BGC)
	if cms(btnPlus) {
		drrlw(btnPlus, 2, SET.UIC4)
		if ms.lc {
			*val += step
			if *val > maxVal {
				*val = maxVal
			}
			ms.lc = false
		}
	} else {
		drrl(btnPlus, SET.UIC3)
	}
	dtxtreccnt(fs, "+", btnPlus, SET.UICtxt)
	if cms(btnMinus) || cms(btnPlus) || cms(valBox) {
		mInfoTXT(hovertxt)
	}
}

// onoff draws a clickable toggle checkbox/switch and returns the toggled boolean value
func onoff(val *bool, x, y, w float32, hovertxt string) {
	r := R(x, y, w, w)
	drrl(r, SET.UIC3)
	r2 := recsmlr(r, 4)
	if *val {
		drr(r2, SET.UICgrn)
	} else {
		drr(r2, SET.BGC)
	}
	if cms(r) {
		drrlw(r, 2, SET.UIC4)
		if ms.lc {
			*val = !*val
		}
		mInfoTXT(hovertxt)
	}
}
func onoffcolor(val *bool, x, y, w, oultineW float32, cOn, cBG, cOutline sdl.Color, outline bool, hovertxt string) {
	r := R(x, y, w, w)
	r2 := recsmlr(r, 4)
	if *val {
		drr(r2, cOn)
	} else {
		drr(r2, cBG)
	}
	if cms(r) {
		drrlw(r, 2, SET.UIC4)
		if ms.lc {
			*val = !*val
		}
		mInfoTXT(hovertxt)
	} else {
		if outline {
			drrlw(r, oultineW, cOutline)
		}
	}
}

func clickOnceSwitch(x, y, w float32, hovertxt string) bool {
	val := false
	r := R(x, y, w, w)
	drrl(r, SET.UIC3)
	r2 := recsmlr(r, 4)
	if cms(r) {
		drr(r2, SET.UIC3)
		drrlw(r, 2, SET.UIC4)
		if ms.lc {
			drr(r2, SET.UICgrn)
			val = true
			ms.lc = false
		}
		mInfoTXT(hovertxt)
	} else {
		drr(r2, SET.UICgrn)
	}
	return val
}
