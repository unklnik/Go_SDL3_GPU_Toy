package main

import (
	"fmt"
	"math"

	"github.com/Zyko0/go-sdl3/sdl"
)

// MARK: GRIDS ███████ CHANGE TO METHOD ███████
func dGrid(gr *GRID2D, lineW float32) int {
	candraw := true
	if gr.hidden && !SET.DrawHidden {
		candraw = false
	}
	num := -1
	for i, r := range gr.rs {
		if cms(r) {
			if gr.hoverVis {
				if candraw {
					drr(r, gr.cHover)
				}
			}
			num = i
			mInfoTXT("Grid block number " + fmt.Sprint(i))
		}
		if candraw {
			drrlw(r, lineW, gr.c)
		}
	}
	if cms(gr.recBorder) {
		if gr.nm != "" {
			mInfoTXT2(gr.nm)
		}
	}

	return num
}

// MARK:TEXTURES ████
func dtex2Dcntr(t *TEX2D, p sdl.FPoint) {
	if t == nil || t.tex == nil || t.tex.GPU == nil {
		return
	}
	t.cnt = p

	// 1. Determine active source frame
	src := sdl.FRect{X: 0, Y: 0, W: t.tex.W, H: t.tex.H}
	if !t.dAllFrames {
		if len(t.frames) > 0 {
			frameIdx := t.curFrame
			if len(t.partialFrames) > 0 {
				if frameIdx >= len(t.partialFrames) {
					frameIdx = 0
				}
				frameIdx = t.partialFrames[frameIdx]
			}
			if frameIdx >= 0 && frameIdx < len(t.frames) {
				src = t.frames[frameIdx]
			}
		}
	}

	// 2. Scale factor
	scale := t.scale
	if scale <= 0 {
		scale = 1.0
	}

	// 🚀 Compute destination rectangle using the FRAME dimensions, not the entire sheet!
	w := src.W * scale
	h := src.H * scale

	dst := sdl.FRect{
		X: p.X - w*0.5,
		Y: p.Y - h*0.5,
		W: w,
		H: h,
	}
	t.recBorder = dst

	// 3. Fallback color to White if empty
	c := t.c
	if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
		c = Col.White
	}
	if t.shadow {
		dst2 := dst
		dst2.X += t.shadowX
		dst2.Y += t.shadowY
		dtexrectEx(t.tex, src, dst2, t.ro, t.cnt, t.flipX, t.flipY, 0, 0, t.cShadow)
	}

	if t.dAllFrames { //DRAW FRAME BORDERS COLLISION RECS SELECT FRAMES
		rd, rc := t.GetCollisionRectAll()
		inFrame := false
		for i := range rd {
			if selectFrame {
				if cms(rd[i]) {
					inFrame = true
					drr(rd[i], CA(Col.PinkHot, 100))
					drrlw(rd[i], 2, RandomColor())
					if ms.lc {
						frameListTemp = addremSlice(frameListTemp, i)
					}
					mInfoTXT2("Frame Number: " + fmt.Sprint(i))
				} else {
					drrl(rd[i], Col.BlueCyan)
				}
			} else {
				drrl(rd[i], Col.BlueCyan)
			}
			if len(frameListTemp) > 0 {
				for j := range frameListTemp {
					if frameListTemp[j] == i {
						drr(rd[i], CA(Col.OrangeBlaze, 100))
					}
				}
			}
		}
		if selectFrame && !inFrame {
			if ms.lc || ms.rc {
				selectFrame = false
				isCursorActiv = false
			}
		}
		if collisFramesVis {
			for i := range rc {
				if rc[i] != rd[i] {
					drrl(rc[i], Col.PinkHot)
				}
			}
		} else {
			drrl(dst, Col.CyanBright)
			if collisFramesVis {
				r := t.GetCollisionRect()
				drrl(r, Col.PinkHot)
			}
		}
	}

	// 4. Draw rotated, flipped, with active shader effect
	t.recPoints = dtexrectEx(t.tex, src, dst, t.ro, t.cnt, t.flipX, t.flipY, t.fx, t.fxParam, c)

}

func dtex2D(t *TEX2D) {
	if t == nil || t.tex == nil || t.tex.GPU == nil {
		return
	}

	// 1. Determine active source frame
	src := sdl.FRect{X: 0, Y: 0, W: t.tex.W, H: t.tex.H}
	if !t.dAllFrames {
		if len(t.frames) > 0 {
			frameIdx := t.curFrame
			if len(t.partialFrames) > 0 {
				if frameIdx >= len(t.partialFrames) {
					frameIdx = 0
				}
				frameIdx = t.partialFrames[frameIdx]
			}
			if frameIdx >= 0 && frameIdx < len(t.frames) {
				src = t.frames[frameIdx]
			}
		}
	}

	scale := t.scale
	if scale <= 0 {
		scale = 1.0
	}

	// Destination rectangle based on frame size
	dst := t.recBorder
	if dst.W == 0 && dst.H == 0 {
		dst.W = src.W * scale
		dst.H = src.H * scale
	}

	c := t.c
	if c.A == 0 && c.R == 0 && c.G == 0 && c.B == 0 {
		c = Col.White
	}

	pivot := t.cnt
	if t.shadow {
		dst2 := dst
		dst2.X += t.shadowX
		dst2.Y += t.shadowY
		dtexrectEx(t.tex, src, dst2, t.ro, pivot, t.flipX, t.flipY, 0, 0, t.cShadow)
	}
	t.recPoints = dtexrectEx(t.tex, src, dst, t.ro, pivot, t.flipX, t.flipY, t.fx, t.fxParam, c)
}

// dtexthumb renders a texture fitted proportionally and centered inside 'r'
func dtexthumb(t *TEX2D, r sdl.FRect) {
	if t == nil || t.tex == nil || t.tex.W <= 0 || t.tex.H <= 0 || r.W <= 0 || r.H <= 0 {
		return
	}

	origW := t.tex.W
	origH := t.tex.H

	// 1. Calculate scale factor to fit inside the box (preserving aspect ratio)
	scaleW := r.W / origW
	scaleH := r.H / origH

	// The smaller scale factor ensures the whole image fits without clipping:
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	// 2. Proportional dimensions
	fitW := origW * scale
	fitH := origH * scale

	// 3. Center the thumbnail inside the destination box 'r'
	posX := r.X + (r.W-fitW)/2
	posY := r.Y + (r.H-fitH)/2

	dst := sdl.FRect{
		X: posX,
		Y: posY,
		W: fitW,
		H: fitH,
	}
	src := sdl.FRect{
		X: 0,
		Y: 0,
		W: origW,
		H: origH,
	}

	// 4. Default to White tint so colors are unaffected
	tint := Col.White
	if t.c.A > 0 {
		tint = t.c
	}

	// Draw using your 2D batcher
	dtexrect(t.tex, src, dst, tint)
}

// dtex: Draw entire texture at native pixel size (x, y)
func dtexscale1(t *Texture, x, y float32) {
	dtexscale(t, x, y, 1.0)
}
func dtexXYWH(t *Texture, x, y, w, h float32, c sdl.Color) {
	r := sdl.FRect{x, y, w, h}
	dtexr(t, r, c)
}

// dtexscale: Draw texture with scale multiplier
func dtexscale(t *Texture, x, y, scale float32) {
	if t == nil || t.GPU == nil {
		return
	}
	w := t.W * scale
	h := t.H * scale
	dtexrect(t, sdl.FRect{X: 0, Y: 0, W: t.W, H: t.H}, sdl.FRect{X: x, Y: y, W: w, H: h}, Col.White)
}
func dtexcenterpWH(t *Texture, p sdl.FPoint, w, h float32, c sdl.Color) {
	if t == nil || t.GPU == nil {
		return
	}
	x := p.X - w/2
	y := p.Y - h/2
	dtexXYWH(t, x, y, w, h, c)
}
func dtexcenterp(t *Texture, p sdl.FPoint, scale float32) {
	if t == nil || t.GPU == nil {
		return
	}
	w := t.W * scale
	h := t.H * scale
	dtexscale(t, p.X-w/2, p.Y-h/2, scale)
}

// dtexcenter: Draw texture centered at (cx, cy)
func dtexcenter(t *Texture, cx, cy, scale float32) {
	if t == nil || t.GPU == nil {
		return
	}
	w := t.W * scale
	h := t.H * scale
	dtexscale(t, cx-w/2, cy-h/2, scale)
}
func dtexmscntr(t *Texture, r sdl.FRect, c sdl.Color) {
	if t == nil || t.GPU == nil {
		return
	}
	r.X = ms.X - r.W/2
	r.Y = ms.Y - r.H/2
	dtexr(t, r, c)
}
func dtexmscntrscale(t *Texture, scale float32, c sdl.Color) {
	if t == nil || t.GPU == nil {
		return
	}
	w := t.W * scale
	h := t.H * scale
	x := ms.X - w/2
	y := ms.Y - h/2
	dtexXYWH(t, x, y, w, h, c)
}

// dtexr: Draw texture in a rectangle using tex.W,tex.H for source rectangle
func dtexr(t *Texture, r sdl.FRect, c sdl.Color) {
	srcrec := sdl.FRect{0, 0, t.W, t.H}
	dtexrect(t, srcrec, r, c)
}

func dtexrect(t *Texture, src, dst sdl.FRect, c sdl.Color) {
	if R2D == nil || t == nil || t.GPU == nil {
		return
	}

	u0 := src.X / t.W
	v0 := src.Y / t.H
	u1 := (src.X + src.W) / t.W
	v1 := (src.Y + src.H) / t.H

	rf, gf, bf, af := color2float(c)

	quad := [6]SpriteVertex{
		{X: dst.X, Y: dst.Y, U: u0, V: v0, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},
		{X: dst.X + dst.W, Y: dst.Y, U: u1, V: v0, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},
		{X: dst.X, Y: dst.Y + dst.H, U: u0, V: v1, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},

		{X: dst.X + dst.W, Y: dst.Y, U: u1, V: v0, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},
		{X: dst.X + dst.W, Y: dst.Y + dst.H, U: u1, V: v1, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},
		{X: dst.X, Y: dst.Y + dst.H, U: u0, V: v1, R: rf, G: gf, B: bf, A: af, Fx: 0, FxParam: 0},
	}

	R2D.addSpriteVertices(t.GPU, quad[:])
}

// dtexrectEx draws a textured quad with rotation, center pivot, flipping, and shader effects.
func dtexrectEx(t *Texture, src, dst sdl.FRect, angleDeg float32, pivot sdl.FPoint, flipX, flipY bool, fx uint8, fxParam float32, c sdl.Color) []sdl.FPoint {
	if R2D == nil || t == nil || t.GPU == nil {
		return nil
	}

	// 1. Calculate UV texture coordinates
	u0 := src.X / t.W
	v0 := src.Y / t.H
	u1 := (src.X + src.W) / t.W
	v1 := (src.Y + src.H) / t.H

	// 2. Flip UVs
	if flipX {
		u0, u1 = u1, u0
	}
	if flipY {
		v0, v1 = v1, v0
	}

	rf, gf, bf, af := color2float(c)
	fxf := float32(fx)

	// 3. Determine Pivot point in screen space
	px := dst.X + dst.W*0.5
	py := dst.Y + dst.H*0.5
	if pivot.X != 0 || pivot.Y != 0 {
		px = pivot.X
		py = pivot.Y
	}

	// 4. Calculate corner positions (accounting for rotation)
	var x0, y0, x1, y1, x2, y2, x3, y3 float32

	if angleDeg != 0 {
		rad := float64(angleDeg) * (math.Pi / 180.0)
		cos := float32(math.Cos(rad))
		sin := float32(math.Sin(rad))

		rotatePoint := func(x, y float32) (float32, float32) {
			dx := x - px
			dy := y - py
			return px + (dx*cos - dy*sin), py + (dx*sin + dy*cos)
		}

		x0, y0 = rotatePoint(dst.X, dst.Y)
		x1, y1 = rotatePoint(dst.X+dst.W, dst.Y)
		x2, y2 = rotatePoint(dst.X+dst.W, dst.Y+dst.H)
		x3, y3 = rotatePoint(dst.X, dst.Y+dst.H)
	} else {
		x0, y0 = dst.X, dst.Y
		x1, y1 = dst.X+dst.W, dst.Y
		x2, y2 = dst.X+dst.W, dst.Y+dst.H
		x3, y3 = dst.X, dst.Y+dst.H
	}

	recpoints := [4]sdl.FPoint{
		{X: x0, Y: y0},
		{X: x1, Y: y1},
		{X: x2, Y: y2},
		{X: x3, Y: y3},
	}

	// 5. Build stack array with Fx and FxParam (0 heap allocations)
	quad := [6]SpriteVertex{
		// Triangle 1: TL -> TR -> BL
		{X: x0, Y: y0, U: u0, V: v0, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},
		{X: x1, Y: y1, U: u1, V: v0, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},
		{X: x3, Y: y3, U: u0, V: v1, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},

		// Triangle 2: TR -> BR -> BL
		{X: x1, Y: y1, U: u1, V: v0, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},
		{X: x2, Y: y2, U: u1, V: v1, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},
		{X: x3, Y: y3, U: u0, V: v1, R: rf, G: gf, B: bf, A: af, Fx: fxf, FxParam: fxParam},
	}

	R2D.addSpriteVertices(t.GPU, quad[:])
	return recpoints[:]
}

// dtexrectRotated draws a textured quad rotated by 'angleDeg' (in degrees) around a pivot point.
// If pivot is (0, 0), it rotates around the center of dst.
func dtexrectRotated(t *Texture, src, dst sdl.FRect, angleDeg float32, pivot sdl.FPoint, c sdl.Color) []sdl.FPoint {
	if R2D == nil || t == nil || t.GPU == nil {
		return nil
	}

	// 1. Calculate UV texture coordinates
	u0 := src.X / t.W
	v0 := src.Y / t.H
	u1 := (src.X + src.W) / t.W
	v1 := (src.Y + src.H) / t.H

	rf, gf, bf, af := color2float(c)

	// 2. Determine Pivot point in screen space
	px := dst.X + dst.W*0.5
	py := dst.Y + dst.H*0.5
	if pivot.X != 0 || pivot.Y != 0 {
		px = pivot.X
		py = pivot.Y
	}

	// 3. Fast-path: If rotation is 0, draw normal quad (no sin/cos overhead)
	if angleDeg == 0 {
		dtexrect(t, src, dst, c)
		return nil
	}

	// 4. Calculate Sine and Cosine (angles in degrees)
	rad := float64(angleDeg) * (math.Pi / 180.0)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	// Helper to rotate a point (x, y) around pivot (px, py)
	rotatePoint := func(x, y float32) (float32, float32) {
		dx := x - px
		dy := y - py
		rx := px + (dx*cos - dy*sin)
		ry := py + (dx*sin + dy*cos)
		return rx, ry
	}

	// 5. Rotate the 4 corners of the quad
	// Top-Left (TL)
	x0, y0 := rotatePoint(dst.X, dst.Y)
	// Top-Right (TR)
	x1, y1 := rotatePoint(dst.X+dst.W, dst.Y)
	// Bottom-Right (BR)
	x2, y2 := rotatePoint(dst.X+dst.W, dst.Y+dst.H)
	// Bottom-Left (BL)
	x3, y3 := rotatePoint(dst.X, dst.Y+dst.H)

	recpoints := [4]sdl.FPoint{
		{X: x0, Y: y0},
		{X: x1, Y: y1},
		{X: x2, Y: y2},
		{X: x3, Y: y3},
	}
	// 6. Build the 2 triangles (6 vertices)
	// Using a fixed-size array [6] on the stack prevents heap allocation/GC pressure!
	quad := [6]SpriteVertex{
		// Triangle 1: TL -> TR -> BL
		{X: x0, Y: y0, U: u0, V: v0, R: rf, G: gf, B: bf, A: af},
		{X: x1, Y: y1, U: u1, V: v0, R: rf, G: gf, B: bf, A: af},
		{X: x3, Y: y3, U: u0, V: v1, R: rf, G: gf, B: bf, A: af},

		// Triangle 2: TR -> BR -> BL
		{X: x1, Y: y1, U: u1, V: v0, R: rf, G: gf, B: bf, A: af},
		{X: x2, Y: y2, U: u1, V: v1, R: rf, G: gf, B: bf, A: af},
		{X: x3, Y: y3, U: u0, V: v1, R: rf, G: gf, B: bf, A: af},
	}

	R2D.addSpriteVertices(t.GPU, quad[:])
	return recpoints[:]
}

// MARK: RECTANGLES ████
// drppl draws a rectangle line outline between two arbitrary corner points.
// Handles p2 being above, below, left, or right of p1.
func drppl(p1, p2 sdl.FPoint, c sdl.Color) {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)

	w := maxX - minX
	h := maxY - minY

	drl(minX, minY, w, h, c)
}
func drpplw(p1, p2 sdl.FPoint, lineW float32, c sdl.Color) {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)

	drlw(minX, minY, maxX-minX, maxY-minY, lineW, c)
}
func drlwcntrWH(p sdl.FPoint, w, h, lineW float32, c sdl.Color) {
	x, y := p.X-w/2, p.Y-h/2
	drlw(x, y, w, h, lineW, c)
}
func drcntrWH(p sdl.FPoint, w, h float32, c sdl.Color) {
	x, y := p.X-w/2, p.Y-h/2
	dr(x, y, w, h, c)
}
func drlro(p []sdl.FPoint, linew float32, c sdl.Color) {
	dliwp(p[0], p[1], linew, c)
	dliwp(p[1], p[2], linew, c)
	dliwp(p[2], p[3], linew, c)
	dliwp(p[3], p[0], linew, c)
}
func drro(p []sdl.FPoint, c sdl.Color) {
	dtrip(p[0], p[1], p[2], c)
	dtrip(p[0], p[2], p[3], c)
}
func drlrolrgr(r sdl.FRect, lineW, ro, offset float32, c sdl.Color) {
	r2 := reclrgr(r, offset)
	p := recRotatedPoints(r2, ro)
	drlro(p, lineW, c)
}

// dr: Draw 2D solid fill rectangle
func dr(x, y, w, h float32, c sdl.Color) {
	if R2D == nil {
		return
	}
	R2D.addShapeVertices(
		V2D(x, y, c), V2D(x+w, y, c), V2D(x, y+h, c),
		V2D(x+w, y, c), V2D(x+w, y+h, c), V2D(x, y+h, c),
	)
}

// drl: Draw 2D rectangle outline (1px default thickness)
func drl(x, y, w, h float32, c sdl.Color) {
	drlw(x, y, w, h, 1.0, c)
}

// drlw: Draw 2D rectangle outline with custom thickness
func drlw(x, y, w, h, thickness float32, c sdl.Color) {
	if R2D == nil || thickness <= 0 {
		return
	}
	dr(x, y, w, thickness, c)
	dr(x, y+h-thickness, w, thickness, c)
	dr(x, y+thickness, thickness, h-thickness*2, c)
	dr(x+w-thickness, y+thickness, thickness, h-thickness*2, c)
}

// drgradH draws a horizontal 2-color gradient rectangle (left to right).
func drgradH(x, y, w, h float32, cLeft, cRight sdl.Color) {
	if R2D == nil || w <= 0 || h <= 0 {
		return
	}

	tl := V2D(x, y, cLeft)
	tr := V2D(x+w, y, cRight)
	bl := V2D(x, y+h, cLeft)
	br := V2D(x+w, y+h, cRight)

	// Two triangles forming the quad
	R2D.addShapeVertices(
		tl, tr, bl,
		tr, br, bl,
	)
}

// drgradHp draws a horizontal gradient quad across 4 corner points (TL, TR, BR, BL).
// Rotates automatically with the points.
func drgradHp(p []sdl.FPoint, cLeft, cRight sdl.Color) {
	if R2D == nil || len(p) < 4 {
		return
	}

	tl := V2D(p[0].X, p[0].Y, cLeft)
	tr := V2D(p[1].X, p[1].Y, cRight)
	br := V2D(p[2].X, p[2].Y, cRight)
	bl := V2D(p[3].X, p[3].Y, cLeft)

	// Triangle 1: TL -> TR -> BL | Triangle 2: TR -> BR -> BL
	R2D.addShapeVertices(
		tl, tr, bl,
		tr, br, bl,
	)
}

// drgradV draws a vertical 2-color gradient rectangle (top to bottom).
func drgradV(x, y, w, h float32, cTop, cBottom sdl.Color) {
	if R2D == nil || w <= 0 || h <= 0 {
		return
	}

	tl := V2D(x, y, cTop)
	tr := V2D(x+w, y, cTop)
	bl := V2D(x, y+h, cBottom)
	br := V2D(x+w, y+h, cBottom)

	R2D.addShapeVertices(
		tl, tr, bl,
		tr, br, bl,
	)
}

// drgradVp draws a vertical gradient quad across 4 corner points (TL, TR, BR, BL).
// Rotates automatically with the points.
func drgradVp(p []sdl.FPoint, cTop, cBottom sdl.Color) {
	if R2D == nil || len(p) < 4 {
		return
	}

	tl := V2D(p[0].X, p[0].Y, cTop)
	tr := V2D(p[1].X, p[1].Y, cTop)
	br := V2D(p[2].X, p[2].Y, cBottom)
	bl := V2D(p[3].X, p[3].Y, cBottom)

	R2D.addShapeVertices(
		tl, tr, bl,
		tr, br, bl,
	)
}

// drgradRadial draws an elliptical/circular radial gradient from center to edge.
// drgradRadial fills the ENTIRE rectangle with a smooth radial gradient that
// extends all the way into the edges and 4 corners with zero empty space.
func drgradRadial(x, y, w, h float32, cCenter, cOuter sdl.Color) {
	if R2D == nil || w <= 0 || h <= 0 {
		return
	}

	cx := x + w*0.5
	cy := y + h*0.5

	// Distance from center to the farthest corner
	dxCorner := w * 0.5
	dyCorner := h * 0.5
	maxRadius := float32(math.Sqrt(float64(dxCorner*dxCorner + dyCorner*dyCorner)))
	if maxRadius <= 0 {
		return
	}

	const res = 8 // 8x8 grid (128 triangles) is buttery smooth on GPU
	stepX := w / float32(res)
	stepY := h / float32(res)

	// Helper to calculate vertex position and radial smoothstep color
	calcVertex := func(px, py float32) Vertex2D {
		dx := px - cx
		dy := py - cy
		dist := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		// Normalized distance (0.0 at center, reaches 1.0 at corners)
		t := dist / (maxRadius * 0.95)
		if t > 1.0 {
			t = 1.0
		}

		// Smoothstep curve for natural soft light falloff: t * t * (3 - 2*t)
		t = t * t * (3.0 - 2.0*t)

		// Lerp between center color and outer color
		rf := (float32(cCenter.R) + t*(float32(cOuter.R)-float32(cCenter.R))) / 255.0
		gf := (float32(cCenter.G) + t*(float32(cOuter.G)-float32(cCenter.G))) / 255.0
		bf := (float32(cCenter.B) + t*(float32(cOuter.B)-float32(cCenter.B))) / 255.0
		af := (float32(cCenter.A) + t*(float32(cOuter.A)-float32(cCenter.A))) / 255.0

		return Vertex2D{X: px, Y: py, Z: 0, R: rf, G: gf, B: bf, A: af}
	}

	// 1. Precalculate the (res+1) x (res+1) grid points (81 vertices total)
	var grid [res + 1][res + 1]Vertex2D
	for j := 0; j <= res; j++ {
		py := y + float32(j)*stepY
		for i := 0; i <= res; i++ {
			px := x + float32(i)*stepX
			grid[j][i] = calcVertex(px, py)
		}
	}

	// 2. Build quads directly on the stack (0 heap allocations)
	const totalVerts = res * res * 6
	var verts [totalVerts]Vertex2D
	vIdx := 0

	for j := 0; j < res; j++ {
		for i := 0; i < res; i++ {
			tl := grid[j][i]
			tr := grid[j][i+1]
			bl := grid[j+1][i]
			br := grid[j+1][i+1]

			// Triangle 1: TL -> TR -> BL
			verts[vIdx] = tl
			verts[vIdx+1] = tr
			verts[vIdx+2] = bl

			// Triangle 2: TR -> BR -> BL
			verts[vIdx+3] = tr
			verts[vIdx+4] = br
			verts[vIdx+5] = bl

			vIdx += 6
		}
	}

	R2D.addShapeVertices(verts[:]...)
}

// drgradRadialp draws a smooth radial gradient across 4 arbitrary points (TL, TR, BR, BL).
// Handles rotation and fills 100% of the quad without empty corners.
func drgradRadialp(p []sdl.FPoint, cCenter, cOuter sdl.Color) {
	if R2D == nil || len(p) < 4 {
		return
	}

	// 1. Center of the 4 points
	cx := (p[0].X + p[1].X + p[2].X + p[3].X) * 0.25
	cy := (p[0].Y + p[1].Y + p[2].Y + p[3].Y) * 0.25

	// 2. Maximum radius (distance to farthest corner)
	distCorner := func(pt sdl.FPoint) float32 {
		dx := pt.X - cx
		dy := pt.Y - cy
		return float32(math.Sqrt(float64(dx*dx + dy*dy)))
	}

	maxRadius := max(distCorner(p[0]), distCorner(p[1]), distCorner(p[2]), distCorner(p[3]))
	if maxRadius <= 0 {
		return
	}

	const res = 8

	// Helper to calculate vertex position via bilinear interpolation and radial color
	calcVertex := func(u, v float32) Vertex2D {
		// Bilinear interpolation between the 4 corners:
		// p[0] = TL, p[1] = TR, p[2] = BR, p[3] = BL
		px := (1-u)*(1-v)*p[0].X + u*(1-v)*p[1].X + u*v*p[2].X + (1-u)*v*p[3].X
		py := (1-u)*(1-v)*p[0].Y + u*(1-v)*p[1].Y + u*v*p[2].Y + (1-u)*v*p[3].Y

		dx := px - cx
		dy := py - cy
		d := float32(math.Sqrt(float64(dx*dx + dy*dy)))

		t := d / (maxRadius * 0.95)
		if t > 1.0 {
			t = 1.0
		}
		// Smoothstep curve for natural soft falloff
		t = t * t * (3.0 - 2.0*t)

		rf := (float32(cCenter.R) + t*(float32(cOuter.R)-float32(cCenter.R))) / 255.0
		gf := (float32(cCenter.G) + t*(float32(cOuter.G)-float32(cCenter.G))) / 255.0
		bf := (float32(cCenter.B) + t*(float32(cOuter.B)-float32(cCenter.B))) / 255.0
		af := (float32(cCenter.A) + t*(float32(cOuter.A)-float32(cCenter.A))) / 255.0

		return Vertex2D{X: px, Y: py, Z: 0, R: rf, G: gf, B: bf, A: af}
	}

	// 3. Precalculate 81 grid points on stack
	var grid [res + 1][res + 1]Vertex2D
	for j := 0; j <= res; j++ {
		v := float32(j) / float32(res)
		for i := 0; i <= res; i++ {
			u := float32(i) / float32(res)
			grid[j][i] = calcVertex(u, v)
		}
	}

	// 4. Build triangles (0 heap allocations)
	const totalVerts = res * res * 6
	var verts [totalVerts]Vertex2D
	vIdx := 0

	for j := 0; j < res; j++ {
		for i := 0; i < res; i++ {
			tl := grid[j][i]
			tr := grid[j][i+1]
			bl := grid[j+1][i]
			br := grid[j+1][i+1]

			verts[vIdx] = tl
			verts[vIdx+1] = tr
			verts[vIdx+2] = bl

			verts[vIdx+3] = tr
			verts[vIdx+4] = br
			verts[vIdx+5] = bl

			vIdx += 6
		}
	}

	R2D.addShapeVertices(verts[:]...)
}

// drgrad4 gives each corner its own independent color.
func drgrad4(x, y, w, h float32, cTL, cTR, cBR, cBL sdl.Color) {
	if R2D == nil || w <= 0 || h <= 0 {
		return
	}

	tl := V2D(x, y, cTL)
	tr := V2D(x+w, y, cTR)
	br := V2D(x+w, y+h, cBR)
	bl := V2D(x, y+h, cBL)

	R2D.addShapeVertices(
		tl, tr, bl,
		tr, br, bl,
	)
}

// drgrad4p gives each of the 4 corner points an individual color.
func drgrad4p(p []sdl.FPoint, c0, c1, c2, c3 sdl.Color) {
	if R2D == nil || len(p) < 4 {
		return
	}

	v0 := V2D(p[0].X, p[0].Y, c0)
	v1 := V2D(p[1].X, p[1].Y, c1)
	v2 := V2D(p[2].X, p[2].Y, c2)
	v3 := V2D(p[3].X, p[3].Y, c3)

	R2D.addShapeVertices(
		v0, v1, v3,
		v1, v2, v3,
	)
}
func drrgradRadial(r sdl.FRect, cCenter, cOuter sdl.Color) {
	drgradRadial(r.X, r.Y, r.W, r.H, cCenter, cOuter)
}
func drrgrad4(r sdl.FRect, cTL, cTR, cBR, cBL sdl.Color) {
	drgrad4(r.X, r.Y, r.W, r.H, cTL, cTR, cBR, cBL)
}
func drrgradV(r sdl.FRect, cTop, cBottom sdl.Color) {
	drgradV(r.X, r.Y, r.W, r.H, cTop, cBottom)
}
func drrgradH(r sdl.FRect, cLeft, cRight sdl.Color) {
	drgradH(r.X, r.Y, r.W, r.H, cLeft, cRight)
}

// MARK: ROUNDED RECTANGLES
// drround: Draw a filled rectangle with rounded corners
func drround(x, y, w, h, radius float32, c sdl.Color) {
	if R2D == nil {
		return
	}
	// Clamp radius so it doesn't exceed half the width or height
	maxRadius := min(w, h) / 2
	if radius > maxRadius {
		radius = maxRadius
	}
	if radius <= 0 {
		dr(x, y, w, h, c)
		return
	}

	// 1. Inner Cross Body (3 Fill Rectangles)
	dr(x+radius, y, w-2*radius, h, c)               // Center vertical & horizontal spine
	dr(x, y+radius, radius, h-2*radius, c)          // Left edge block
	dr(x+w-radius, y+radius, radius, h-2*radius, c) // Right edge block

	// 2. Corner Centers
	corners := [4][2]float32{
		{x + w - radius, y + radius},     // Top-Right
		{x + radius, y + radius},         // Top-Left
		{x + radius, y + h - radius},     // Bottom-Left
		{x + w - radius, y + h - radius}, // Bottom-Right
	}

	baseAngles := [4]float32{
		0,
		float32(math.Pi) / 2,
		float32(math.Pi),
		float32(math.Pi) * 3 / 2,
	}

	segments := 8 // 8 segments per corner for smooth curves
	step := (float32(math.Pi) / 2) / float32(segments)

	var verts []Vertex2D

	// 3. Draw 4 Corner Triangle Fans
	for i := 0; i < 4; i++ {
		cx, cy := corners[i][0], corners[i][1]
		startA := baseAngles[i]

		for s := 0; s < segments; s++ {
			a1 := startA + float32(s)*step
			a2 := a1 + step

			px1 := cx + float32(math.Cos(float64(a1)))*radius
			py1 := cy - float32(math.Sin(float64(a1)))*radius
			px2 := cx + float32(math.Cos(float64(a2)))*radius
			py2 := cy - float32(math.Sin(float64(a2)))*radius

			verts = append(verts,
				V2D(cx, cy, c),
				V2D(px1, py1, c),
				V2D(px2, py2, c),
			)
		}
	}

	R2D.addShapeVertices(verts...)
}

// drlround: Draw a rectangle outline with rounded corners and custom thickness
func drlround(x, y, w, h, radius, thickness float32, c sdl.Color) {
	if R2D == nil || thickness <= 0 {
		return
	}
	maxRadius := min(w, h) / 2
	if radius > maxRadius {
		radius = maxRadius
	}
	if radius <= thickness {
		drlw(x, y, w, h, thickness, c)
		return
	}

	innerRadius := radius - thickness

	// 1. Four Straight Edge Bars
	dr(x+radius, y, w-2*radius, thickness, c)             // Top
	dr(x+radius, y+h-thickness, w-2*radius, thickness, c) // Bottom
	dr(x, y+radius, thickness, h-2*radius, c)             // Left
	dr(x+w-thickness, y+radius, thickness, h-2*radius, c) // Right

	// 2. Corner Centers
	corners := [4][2]float32{
		{x + w - radius, y + radius},     // Top-Right
		{x + radius, y + radius},         // Top-Left
		{x + radius, y + h - radius},     // Bottom-Left
		{x + w - radius, y + h - radius}, // Bottom-Right
	}

	baseAngles := [4]float32{
		0,
		float32(math.Pi) / 2,
		float32(math.Pi),
		float32(math.Pi) * 3 / 2,
	}

	segments := 8
	step := (float32(math.Pi) / 2) / float32(segments)

	var verts []Vertex2D

	// 3. Draw 4 Curved Ribbon Arcs
	for i := 0; i < 4; i++ {
		cx, cy := corners[i][0], corners[i][1]
		startA := baseAngles[i]

		for s := 0; s < segments; s++ {
			a1 := startA + float32(s)*step
			a2 := a1 + step

			cos1, sin1 := float32(math.Cos(float64(a1))), float32(math.Sin(float64(a1)))
			cos2, sin2 := float32(math.Cos(float64(a2))), float32(math.Sin(float64(a2)))

			// Outer arc points
			ox1, oy1 := cx+cos1*radius, cy-sin1*radius
			ox2, oy2 := cx+cos2*radius, cy-sin2*radius

			// Inner arc points
			ix1, iy1 := cx+cos1*innerRadius, cy-sin1*innerRadius
			ix2, iy2 := cx+cos2*innerRadius, cy-sin2*innerRadius

			// 2 Triangles for this arc segment
			verts = append(verts,
				V2D(ox1, oy1, c), V2D(ox2, oy2, c), V2D(ix1, iy1, c),
				V2D(ox2, oy2, c), V2D(ix2, iy2, c), V2D(ix1, iy1, c),
			)
		}
	}

	R2D.addShapeVertices(verts...)
}

// Helpers for REC and sdl.FRect
func drrround(r sdl.FRect, radius float32, c sdl.Color) { drround(r.X, r.Y, r.W, r.H, radius, c) }
func drrlround(r sdl.FRect, radius, thickness float32, c sdl.Color) {
	drlround(r.X, r.Y, r.W, r.H, radius, thickness, c)
}
func drectround(r sdl.FRect, radius float32, c sdl.Color) { drround(r.X, r.Y, r.W, r.H, radius, c) }
func drectlround(r sdl.FRect, radius, thickness float32, c sdl.Color) {
	drlround(r.X, r.Y, r.W, r.H, radius, thickness, c)
}
func drrlw(r sdl.FRect, thickness float32, c sdl.Color) { drlw(r.X, r.Y, r.W, r.H, thickness, c) }
func drr(r sdl.FRect, c sdl.Color)                      { dr(r.X, r.Y, r.W, r.H, c) }
func drrl(r sdl.FRect, c sdl.Color)                     { drl(r.X, r.Y, r.W, r.H, c) }

// MARK: TRIANGLES ████
func dtrilrgrOutline(t TRI2D, offset, lineW float32, c sdl.Color) {
	p := trilrgrPoints(t.cnt, t.sideW, offset)
	if t.ro != 0 {
		p = triAnyPointRotate(t.cnt, p, t.ro)
	}
	dtrilwslice(p, lineW, c)
}

// dtri: Draw 2D solid fill triangle
func dtrilslice(p [3]sdl.FPoint, c sdl.Color) {
	dtrilp(p[0], p[1], p[2], c)
}
func dtrilwslice(p [3]sdl.FPoint, lineW float32, c sdl.Color) {
	dtrilwp(p[0], p[1], p[2], lineW, c)
}
func dtrislice(p [3]sdl.FPoint, c sdl.Color) {
	dtrip(p[0], p[1], p[2], c)
}
func dtri(x1, y1, x2, y2, x3, y3 float32, c sdl.Color) {
	if R2D == nil {
		return
	}
	R2D.addShapeVertices(V2D(x1, y1, c), V2D(x2, y2, c), V2D(x3, y3, c))
}
func dtrip(p1, p2, p3 sdl.FPoint, c sdl.Color) {
	if R2D == nil {
		return
	}
	x1, y1, x2, y2, x3, y3 := p1.X, p1.Y, p2.X, p2.Y, p3.X, p3.Y
	R2D.addShapeVertices(V2D(x1, y1, c), V2D(x2, y2, c), V2D(x3, y3, c))
}
func dtrilp(p1, p2, p3 sdl.FPoint, c sdl.Color) {
	dlip(p1, p2, c)
	dlip(p2, p3, c)
	dlip(p1, p3, c)
}
func dtrilwp(p1, p2, p3 sdl.FPoint, lineW float32, c sdl.Color) {
	dliwp(p1, p2, lineW, c)
	dliwp(p2, p3, lineW, c)
	dliwp(p1, p3, lineW, c)
}

// dtricntr draws an equilateral triangle centered at 'cntr' where each side equals 'side'.
func dtricntr(cntr sdl.FPoint, side float32, c sdl.Color) {
	if R2D == nil || side <= 0 {
		return
	}

	// 1 / sqrt(3) ≈ 0.57735027
	const invSqrt3 = 0.57735027
	topDist := side * invSqrt3  // Distance from center to top vertex
	bottomDist := topDist * 0.5 // Distance from center to bottom base
	halfSide := side * 0.5

	p1 := sdl.FPoint{X: cntr.X, Y: cntr.Y - topDist}               // Top
	p2 := sdl.FPoint{X: cntr.X + halfSide, Y: cntr.Y + bottomDist} // Bottom-Right
	p3 := sdl.FPoint{X: cntr.X - halfSide, Y: cntr.Y + bottomDist} // Bottom-Left

	dtrip(p1, p2, p3, c)
}
func dtrilcntr(cntr sdl.FPoint, side float32, c sdl.Color) {
	if R2D == nil || side <= 0 {
		return
	}

	// 1 / sqrt(3) ≈ 0.57735027
	const invSqrt3 = 0.57735027
	topDist := side * invSqrt3  // Distance from center to top vertex
	bottomDist := topDist * 0.5 // Distance from center to bottom base
	halfSide := side * 0.5

	p1 := sdl.FPoint{X: cntr.X, Y: cntr.Y - topDist}               // Top
	p2 := sdl.FPoint{X: cntr.X + halfSide, Y: cntr.Y + bottomDist} // Bottom-Right
	p3 := sdl.FPoint{X: cntr.X - halfSide, Y: cntr.Y + bottomDist} // Bottom-Left

	dtrilp(p1, p2, p3, c)
}

// dtricntrRotated draws an equilateral triangle centered at 'cntr' rotated by angleDeg (0 = pointing UP).
func dtricntrRotated(cntr sdl.FPoint, side float32, angleDeg float32, c sdl.Color) {
	if R2D == nil || side <= 0 {
		return
	}

	const invSqrt3 = 0.57735027
	topDist := side * invSqrt3
	bottomDist := topDist * 0.5
	halfSide := side * 0.5

	// Fast-path: no rotation math needed if angle is 0
	if angleDeg == 0 {
		dtrip(
			sdl.FPoint{X: cntr.X, Y: cntr.Y - topDist},
			sdl.FPoint{X: cntr.X + halfSide, Y: cntr.Y + bottomDist},
			sdl.FPoint{X: cntr.X - halfSide, Y: cntr.Y + bottomDist},
			c,
		)
		return
	}

	rad := float64(angleDeg) * (math.Pi / 180.0)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rotatePoint := func(ox, oy float32) sdl.FPoint {
		return sdl.FPoint{
			X: cntr.X + (ox*cos - oy*sin),
			Y: cntr.Y + (ox*sin + oy*cos),
		}
	}

	p1 := rotatePoint(0, -topDist)
	p2 := rotatePoint(halfSide, bottomDist)
	p3 := rotatePoint(-halfSide, bottomDist)

	dtrip(p1, p2, p3, c)
}

// dtricntrWH draws a triangle centered at 'cntr' fitting within box width and height.
func dtricntrWH(cntr sdl.FPoint, w, h float32, c sdl.Color) {
	if R2D == nil || w <= 0 || h <= 0 {
		return
	}

	halfW := w * 0.5
	halfH := h * 0.5

	p1 := sdl.FPoint{X: cntr.X, Y: cntr.Y - halfH}
	p2 := sdl.FPoint{X: cntr.X + halfW, Y: cntr.Y + halfH}
	p3 := sdl.FPoint{X: cntr.X - halfW, Y: cntr.Y + halfH}

	dtrip(p1, p2, p3, c)
}

// MARK: LINES ████
// dli: Draw 2D line (1px default thickness)
func dli(x1, y1, x2, y2 float32, c sdl.Color) {
	dliw(x1, y1, x2, y2, 1.0, c)
}
func dlip(p1, p2 sdl.FPoint, c sdl.Color) {
	dliwp(p1, p2, 1.0, c)
}

// dliwRecSides: Draw selected rectangle sides as lines clockwise from top: sides[1(T),2(R),3(B),4(L)]
func drliwSides(r sdl.FRect, sides []int, w float32, c sdl.Color) {
	for i := range sides {
		switch sides[i] {
		case 1: //TOP
			dliw(r.X, r.Y, r.X+r.W, r.Y, w, c)
		case 2: //RIGHT
			dliw(r.X+r.W, r.Y, r.X+r.W, r.Y+r.H, w, c)
		case 3: //BOTTOM
			dliw(r.X, r.Y+r.H, r.X+r.W, r.Y+r.H, w, c)
		case 4: //LEFT
			dliw(r.X, r.Y, r.X, r.Y+r.H, w, c)
		}
	}
}
func dliwp(p1, p2 sdl.FPoint, linew float32, c sdl.Color) {
	dliw(p1.X, p1.Y, p2.X, p2.Y, linew, c)
}

// dliw: Draw 2D line with custom thickness (Quad)
func dliw(x1, y1, x2, y2, thickness float32, c sdl.Color) {
	if R2D == nil || thickness <= 0 {
		return
	}
	dx := x2 - x1
	dy := y2 - y1
	length := float32(math.Sqrt(float64(dx*dx + dy*dy)))
	if length == 0 {
		return
	}

	nx := (-dy / length) * (thickness / 2)
	ny := (dx / length) * (thickness / 2)

	R2D.addShapeVertices(
		V2D(x1-nx, y1-ny, c), V2D(x1+nx, y1+ny, c), V2D(x2+nx, y2+ny, c),
		V2D(x1-nx, y1-ny, c), V2D(x2+nx, y2+ny, c), V2D(x2-nx, y2-ny, c),
	)
}
