package main

import (
	"math"
	"slices"

	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	selectFade       uint8 = 200
	selectFadeSwitch bool
)

type TRI2D struct {
	nm                                  string
	p                                   [3]sdl.FPoint
	cnt                                 sdl.FPoint
	cFill, cLine                        sdl.Color
	outline, fill                       bool
	ro, roPrev, sideW, lineW, sideWprev float32
}
type GRID2D struct {
	nm                     string
	tags                   []string
	rs                     []sdl.FRect
	cnt                    sdl.FPoint
	cornersfromTL          []sdl.FPoint
	recBorder              sdl.FRect
	c, cHover              sdl.Color
	hidden, hoverVis       bool
	cat, id, rows, columns int
	W, H, blockW, blockH   float32
}
type LINE2D struct {
	pPrevious sdl.FPoint
	p, cnt    []sdl.FPoint
	w, l      float32
	c         sdl.Color
	nm        string
}

type REC2D struct {
	r                 sdl.FRect
	p                 []sdl.FPoint
	ro, roPrev, lineW float32
	nm                string
	fillGrad          int
	outline, fill     bool
	xyPrev            sdl.FPoint

	cLine, cFill, cGrad1, cGrad2, cGrad3, cGrad4 sdl.Color
}

// MARK: TRIANGLES ████
func triDef() TRI2D {
	return TRI2D{outline: true, fill: true, cFill: SET.Tri2DFillC, cLine: SET.Tri2DLineC, sideW: SET.Tri2DSideW, lineW: SET.Tri2DLineW, sideWprev: SET.Tri2DSideW, roPrev: 0}
}
func (t TRI2D) Draw() {
	if t.sideWprev != t.sideW {
		t.p = triPointsCntr(t.cnt, t.sideW)
		t.sideWprev = t.sideW
	}
	if t.ro != t.roPrev {
		t.p = triPointsRotated(t.cnt, t.sideW, t.ro)
		t.roPrev = t.ro
	}
	if t.fill {
		dtrislice(t.p, t.cFill)
	}
	if t.outline {
		if t.lineW == 1 {
			dtrilslice(t.p, t.cLine)
		} else {
			dtrilwslice(t.p, t.lineW, t.cLine)
		}
	}
}
func (t TRI2D) DrawSelected() {
	dtrilrgrOutline(t, 12, 2, CA(SET.Tri2DSelectC, selectFade))
}

// MARK: RECTANGLES ████
func recDef() REC2D {
	return REC2D{cLine: SET.Rec2DLineC, cFill: SET.Rec2DFillC, cGrad1: SET.Rec2DGradFill1C, cGrad2: SET.Rec2DGradFill2C, cGrad3: SET.Rec2DGradFill3C, cGrad4: SET.Rec2DGradFill4C, lineW: SET.Rec2DLineWdef, fill: true, outline: true}
}

// Contains checks if point 'pt' is inside the rectangle (handles any rotation angle).
func (rec *REC2D) Contains(pt sdl.FPoint) bool {
	// 🚀 Fast path: If rotation is 0, use standard instant AABB bounding box check
	if rec.ro == 0 {
		return pt.X >= rec.r.X && pt.X <= rec.r.X+rec.r.W &&
			pt.Y >= rec.r.Y && pt.Y <= rec.r.Y+rec.r.H
	}

	// Rotated check: must have 4 corner points
	if len(rec.p) < 4 {
		return false
	}

	var hasPos, hasNeg bool

	// Check point against all 4 perimeter edges (0->1, 1->2, 2->3, 3->0)
	for i := 0; i < 4; i++ {
		next := (i + 1) & 3 // 0->1, 1->2, 2->3, 3->0

		// 2D Cross product of edge vector and point vector
		cross := (rec.p[next].X-rec.p[i].X)*(pt.Y-rec.p[i].Y) - (rec.p[next].Y-rec.p[i].Y)*(pt.X-rec.p[i].X)

		if cross > 0 {
			hasPos = true
		} else if cross < 0 {
			hasNeg = true
		}

		// Early exit: if point is to the right of one line and left of another, it's outside!
		if hasPos && hasNeg {
			return false
		}
	}

	return true
}

// HasMouse returns true if the mouse cursor is currently hovering inside the rectangle.
func (rec *REC2D) HasMouse() bool {
	return rec.Contains(sdl.FPoint{X: ms.X, Y: ms.Y})
}
func (r REC2D) Draw() {
	if r.ro == 0 {
		if r.fill {
			switch r.fillGrad {
			case 4: //4 CORNERS
				drrgrad4(r.r, r.cGrad1, r.cGrad2, r.cGrad3, r.cGrad4)
			case 3: //RADIAL
				drrgradRadial(r.r, r.cGrad1, r.cGrad2)
			case 2: //VERT
				drrgradV(r.r, r.cGrad1, r.cGrad2)
			case 1: //HORIZ
				drrgradH(r.r, r.cGrad1, r.cGrad2)
			case 0: //SOLID
				drr(r.r, r.cFill)
			}
		}
		if r.outline {
			if r.lineW == 1 {
				drrl(r.r, r.cLine)
			} else if r.lineW > 0 && r.lineW < 1 || r.lineW > 1 {
				drrlw(r.r, r.lineW, r.cLine)
			}
		}
	} else {
		if r.roPrev != r.ro {
			r.p = recRotatedPoints(r.r, r.ro)
			r.roPrev = r.ro
		}
		if r.fill {
			switch r.fillGrad {
			case 4: //4 CORNERS
				drgrad4p(r.p, r.cGrad1, r.cGrad2, r.cGrad3, r.cGrad4)
			case 3: //RADIAL
				drgradRadialp(r.p, r.cGrad1, r.cGrad2)
			case 2: //VERT
				drgradVp(r.p, r.cGrad1, r.cGrad2)
			case 1: //HORIZ
				drgradHp(r.p, r.cGrad1, r.cGrad2)
			case 0: //SOLID
				drro(r.p, r.cFill)
			}
		}
		if r.outline {
			drlro(r.p, r.lineW, r.cLine)
		}

	}
}
func (r REC2D) DrawSelected() {
	if r.ro == 0 {
		r2 := reclrgr(r.r, 8)
		drrlw(r2, 2, CA(SET.Rec2DSelectC, selectFade))
	} else {
		drlrolrgr(r.r, 2, r.ro, 8, CA(SET.Rec2DSelectC, selectFade))
	}
}

// MARK: LINES ████
func (l LINE2D) Draw() {
	for i := range len(l.p) - 1 {
		dliw(l.p[i].X, l.p[i].Y, l.p[i+1].X, l.p[i+1].Y, l.w, l.c)
	}
}
func lineDef() LINE2D {
	return LINE2D{c: SET.LineCdef, w: SET.LineWdef}
}
func (l LINE2D) Mid() []sdl.FPoint {
	var p []sdl.FPoint
	for i := range len(l.p) - 1 {
		p = append(p, sdl.FPoint{X: (l.p[i].X + l.p[i+1].X) * 0.5, Y: (l.p[i].Y + l.p[i+1].Y) * 0.5})
	}
	return p
}

// Length calculates the total length across all segments of the polyline.
func (l LINE2D) Length() float32 {
	if len(l.p) < 2 {
		return 0
	}

	var total float32
	for i := 0; i < len(l.p)-1; i++ {
		dx := l.p[i+1].X - l.p[i].X
		dy := l.p[i+1].Y - l.p[i].Y
		total += float32(math.Sqrt(float64(dx*dx + dy*dy)))
	}

	return total
}
func (l LINE2D) Clone() LINE2D {
	cp := l // 1. Shallow copy primitive values (w, l, c, nm, pPrevious)

	// 2. Deep copy the slices so they don't share memory backing arrays
	cp.p = slices.Clone(l.p)
	cp.cnt = slices.Clone(l.cnt)

	return cp
}

// 3. Move / Translate (Pointer Receiver '*'): MODIFIES the actual line!
func (l *LINE2D) Move(dx, dy float32) {
	// 1. Shift all line vertices
	for i := range l.p {
		l.p[i].X += dx
		l.p[i].Y += dy
	}

	// 2. Shift midpoints in place (0 garbage collection pressure)
	for i := range l.cnt {
		l.cnt[i].X += dx
		l.cnt[i].Y += dy
	}
}

// MARK: GRIDS ████
func gridDef() GRID2D {
	return GRID2D{c: SET.GridLineCdef, cHover: SET.GridHoverCdef}
}
func mGridfromGridVal(p sdl.FPoint, cntrTRUE bool) GRID2D {
	g2 := GRID2D{}
	if cntrTRUE {
		g2 = mGridfromCenter(p, gridVal.blockW, gridVal.blockH, gridVal.columns, gridVal.rows)
	} else {
		g2 = mGridfromTopLeft(p, gridVal.blockW, gridVal.blockH, gridVal.columns, gridVal.rows)
	}
	g2.c = gridVal.c
	g2.cHover = gridVal.cHover
	g2.hidden = gridVal.hidden
	g2.hoverVis = gridVal.hoverVis
	g2.nm = gridVal.nm
	g2.hidden = gridVal.hidden
	return g2
}
func mGridfromCenter(cnt sdl.FPoint, blockW, blockH float32, rows, columns int) GRID2D {
	gr := GRID2D{
		cnt:           cnt,
		W:             f32(rows) * blockW,
		H:             f32(columns) * blockH,
		recBorder:     R(cnt.X-(f32(rows)*blockW)/2, cnt.Y-(f32(columns)*blockH)/2, f32(rows)*blockW, f32(columns)*blockH),
		cornersfromTL: recCorners(R(cnt.X-(f32(rows)*blockW)/2, cnt.Y-(f32(columns)*blockH)/2, f32(rows)*blockW, f32(columns)*blockH)),
	}
	x, y := gr.cornersfromTL[0].X, gr.cornersfromTL[0].Y
	ox := x
	a := rows * columns
	c := 0
	for a > 0 {
		gr.rs = append(gr.rs, R(x, y, blockW, blockH))
		x += blockW
		c++
		a--
		if c == rows {
			c = 0
			x = ox
			y += blockH
		}
	}
	return gr
}
func mGridfromTopLeft(tl sdl.FPoint, blockW, blockH float32, rows, columns int) GRID2D {
	cnt := P(tl.X+(f32(rows)*blockW)/2, tl.Y+(f32(columns)*blockH)/2)
	gr := mGridfromCenter(cnt, blockW, blockH, rows, columns)
	return gr
}
