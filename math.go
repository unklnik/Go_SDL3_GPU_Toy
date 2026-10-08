package main

import (
	"math"
	"math/rand/v2"
	"slices"
	"strconv"
	"strings"

	"github.com/Zyko0/go-sdl3/sdl"
)

// MATH

// MARK: TRIANGLES
func triAnyPointRotate(cnt sdl.FPoint, p [3]sdl.FPoint, ro float32) [3]sdl.FPoint {
	// Fast-path: no rotation needed
	if ro == 0 {
		return p
	}

	rad := float64(ro) * (math.Pi / 180.0)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	var rotated [3]sdl.FPoint

	for i := 0; i < 3; i++ {
		dx := p[i].X - cnt.X
		dy := p[i].Y - cnt.Y

		rotated[i] = sdl.FPoint{
			X: cnt.X + (dx*cos - dy*sin),
			Y: cnt.Y + (dx*sin + dy*cos),
		}
	}

	return rotated
}
func triPointsRotated(cnt sdl.FPoint, sideW, ro float32) [3]sdl.FPoint {
	const invSqrt3 = 0.57735027
	topDist := sideW * invSqrt3
	bottomDist := topDist * 0.5
	halfSide := sideW * 0.5

	rad := float64(ro) * (math.Pi / 180.0)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rotatePoint := func(ox, oy float32) sdl.FPoint {
		return sdl.FPoint{
			X: cnt.X + (ox*cos - oy*sin),
			Y: cnt.Y + (ox*sin + oy*cos),
		}
	}

	p1 := rotatePoint(0, -topDist)
	p2 := rotatePoint(halfSide, bottomDist)
	p3 := rotatePoint(-halfSide, bottomDist)
	return [3]sdl.FPoint{p1, p2, p3}
}
func triPointsCntr(cnt sdl.FPoint, sideW float32) [3]sdl.FPoint {
	// 1 / sqrt(3) ≈ 0.57735027
	const invSqrt3 = 0.57735027
	topDist := sideW * invSqrt3 // Distance from center to top vertex
	bottomDist := topDist * 0.5 // Distance from center to bottom base
	halfSide := sideW * 0.5

	p1 := sdl.FPoint{X: cnt.X, Y: cnt.Y - topDist}               // Top
	p2 := sdl.FPoint{X: cnt.X + halfSide, Y: cnt.Y + bottomDist} // Bottom-Right
	p3 := sdl.FPoint{X: cnt.X - halfSide, Y: cnt.Y + bottomDist} // Bottom-Left
	return [3]sdl.FPoint{p1, p2, p3}
}
func trilrgrPoints(cnt sdl.FPoint, sideW, offset float32) [3]sdl.FPoint {
	sideW += offset
	return triPointsCntr(cnt, sideW)
}

// MARK: TRIANGLES
func tricntr(p [3]sdl.FPoint) sdl.FPoint {
	const inv3 = float32(1.0 / 3.0)
	return sdl.FPoint{
		X: (p[0].X + p[1].X + p[2].X) * inv3,
		Y: (p[0].Y + p[1].Y + p[2].Y) * inv3,
	}
}
func tricntrp(p1, p2, p3 sdl.FPoint) sdl.FPoint {
	const inv3 = float32(1.0 / 3.0)

	return sdl.FPoint{
		X: (p1.X + p2.X + p3.X) * inv3,
		Y: (p1.Y + p2.Y + p3.Y) * inv3,
	}
}
func ctripp(pt, p1, p2, p3 sdl.FPoint) bool {
	d1 := (p2.X-p1.X)*(pt.Y-p1.Y) - (p2.Y-p1.Y)*(pt.X-p1.X)
	d2 := (p3.X-p2.X)*(pt.Y-p2.Y) - (p3.Y-p2.Y)*(pt.X-p2.X)
	d3 := (p1.X-p3.X)*(pt.Y-p3.Y) - (p1.Y-p3.Y)*(pt.X-p3.X)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	return !(hasNeg && hasPos)
}
func ctrip(pt sdl.FPoint, p [3]sdl.FPoint) bool {
	// 2D cross product for each of the 3 directed edges
	d1 := (p[1].X-p[0].X)*(pt.Y-p[0].Y) - (p[1].Y-p[0].Y)*(pt.X-p[0].X)
	d2 := (p[2].X-p[1].X)*(pt.Y-p[1].Y) - (p[2].Y-p[1].Y)*(pt.X-p[1].X)
	d3 := (p[0].X-p[2].X)*(pt.Y-p[2].Y) - (p[0].Y-p[2].Y)*(pt.X-p[2].X)

	hasNeg := (d1 < 0) || (d2 < 0) || (d3 < 0)
	hasPos := (d1 > 0) || (d2 > 0) || (d3 > 0)

	// If the signs are mixed, the point is outside.
	// If all are >= 0 or all are <= 0, the point is inside.
	return !(hasNeg && hasPos)
}

// MARK: RECS
func recRotatedPoints(r sdl.FRect, ro float32) []sdl.FPoint {
	cnt := reccntr(r)
	px, py := cnt.X, cnt.Y
	rad := float64(ro) * (math.Pi / 180.0)
	cos := float32(math.Cos(rad))
	sin := float32(math.Sin(rad))

	rotatePoint := func(x, y float32) (float32, float32) {
		dx := x - px
		dy := y - py
		rx := px + (dx*cos - dy*sin)
		ry := py + (dx*sin + dy*cos)
		return rx, ry
	}

	// Top-Left (TL)
	x0, y0 := rotatePoint(r.X, r.Y)
	// Top-Right (TR)
	x1, y1 := rotatePoint(r.X+r.W, r.Y)
	// Bottom-Right (BR)
	x2, y2 := rotatePoint(r.X+r.W, r.Y+r.H)
	// Bottom-Left (BL)
	x3, y3 := rotatePoint(r.X, r.Y+r.H)

	return []sdl.FPoint{
		{X: x0, Y: y0},
		{X: x1, Y: y1},
		{X: x2, Y: y2},
		{X: x3, Y: y3},
	}
}
func recFromPoints(p1, p2 sdl.FPoint) sdl.FRect {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)

	return sdl.FRect{
		X: minX,
		Y: minY,
		W: maxX - minX,
		H: maxY - minY,
	}
}
func recPointsFrom2Points(p1, p2 sdl.FPoint) []sdl.FPoint {
	minX := min(p1.X, p2.X)
	maxX := max(p1.X, p2.X)
	minY := min(p1.Y, p2.Y)
	maxY := max(p1.Y, p2.Y)
	p3 := P(minX, minY)
	p4 := P(minX+(maxX-minX), minY)
	p5 := P(minX+(maxX-minX), minY+(maxY-minY))
	p6 := P(minX, minY+(maxY-minY))
	return []sdl.FPoint{p3, p4, p5, p6}
}
func recfromcntr(p sdl.FPoint, w, h float32) sdl.FRect {
	return R(p.X-w/2, p.Y-h/2, w, h)
}
func recmovedcntr(r sdl.FRect, newcntr sdl.FPoint) sdl.FRect {
	return R(newcntr.X-r.W/2, newcntr.Y-r.H/2, r.W, r.H)
}
func R(x, y, w, h float32) sdl.FRect {
	return sdl.FRect{x, y, w, h}
}
func reccntr(r sdl.FRect) sdl.FPoint {
	return sdl.FPoint{r.X + r.W/2, r.Y + r.H/2}
}
func recPoints(r sdl.FRect) (sdl.FPoint, sdl.FPoint, sdl.FPoint, sdl.FPoint) {
	v1 := P(r.X, r.Y)
	v2 := v1
	v2.X += r.W
	v3 := v2
	v3.Y += r.H
	v4 := v1
	v4.Y += r.H
	return v1, v2, v3, v4
}
func recsmlr(r sdl.FRect, inset float32) sdl.FRect {
	return R(r.X+inset, r.Y+inset, r.W-(inset*2), r.H-(inset*2))
}
func reclrgr(r sdl.FRect, offset float32) sdl.FRect {
	return R(r.X-offset, r.Y-offset, r.W+offset*2, r.H+offset*2)
}

// MARK: NUMBERS
func divmulti[T1, T2, T3 numConvert](num T1, divby T2, multiby T3) float32 {
	return (float32(num) / float32(divby)) * float32(multiby)
}

// MARK:COLLISIONS
// cmsrro checks if mouse is inside a 4-sided polygon 'r'. Points in 'r' must be in consecutive order around the perimeter (e.g., TL -> TR -> BR -> BL)
func cmsrro4slice(r [4]sdl.FPoint) bool {
	var hasPos, hasNeg bool

	for i := 0; i < 4; i++ {
		next := (i + 1) & 3
		cross := (r[next].X-r[i].X)*(ms.Y-r[i].Y) - (r[next].Y-r[i].Y)*(ms.X-r[i].X)

		if cross > 0 {
			hasPos = true
		} else if cross < 0 {
			hasNeg = true
		}

		if hasPos && hasNeg {
			return false
		}
	}

	return true
}
func cmsrro(r []sdl.FPoint) bool {
	var hasPos, hasNeg bool

	for i := 0; i < 4; i++ {
		next := (i + 1) & 3
		cross := (r[next].X-r[i].X)*(ms.Y-r[i].Y) - (r[next].Y-r[i].Y)*(ms.X-r[i].X)

		if cross > 0 {
			hasPos = true
		} else if cross < 0 {
			hasNeg = true
		}

		if hasPos && hasNeg {
			return false
		}
	}

	return true
}

// cprro checks if point 'p' is inside a 4-sided polygon 'r'. Points in 'r' must be in consecutive order around the perimeter (e.g., TL -> TR -> BR -> BL)
func cprro(p sdl.FPoint, r [4]sdl.FPoint) bool {
	var hasPos, hasNeg bool

	for i := 0; i < 4; i++ {
		next := (i + 1) & 3
		cross := (r[next].X-r[i].X)*(p.Y-r[i].Y) - (r[next].Y-r[i].Y)*(p.X-r[i].X)

		if cross > 0 {
			hasPos = true
		} else if cross < 0 {
			hasNeg = true
		}

		if hasPos && hasNeg {
			return false
		}
	}

	return true
}

// crrslice: Checks if sdl.FRect collides with any other rectangles in []sdl.FRect
func crrslice(r sdl.FRect, rs []sdl.FRect) bool {
	for i := range rs {
		b := &rs[i]
		if r.X < b.X+b.W &&
			r.X+r.W > b.X &&
			r.Y < b.Y+b.H &&
			r.Y+r.H > b.Y {
			return true
		}
	}
	return false
}

// crrsliceReturnNum: Checks if sdl.FRect collides with any other rectangles in []sdl.FRect, return num ([i]) in slice of sdl.FRect that it collides with
func crrsliceReturnNum(r sdl.FRect, list []sdl.FRect) (int, bool) {
	for i := range list {
		b := &list[i]
		if r.X < b.X+b.W &&
			r.X+r.W > b.X &&
			r.Y < b.Y+b.H &&
			r.Y+r.H > b.Y {
			return i, true
		}
	}
	return -1, false
}

// crrslice: Checks if sdl.FRect (also in slice) collides with any other rectangles in the slice []sdl.FRect, skipIdx ([i]) is the index of r sdl.FRect in slice
func crrsliceInSlice(r sdl.FRect, list []sdl.FRect, skipIdx int) bool {
	for i := range list {
		if i == skipIdx {
			continue
		}
		b := &list[i]
		if r.X < b.X+b.W &&
			r.X+r.W > b.X &&
			r.Y < b.Y+b.H &&
			r.Y+r.H > b.Y {
			return true
		}
	}
	return false
}
func crpxy(x, y float32, r sdl.FRect) bool {
	return x >= r.X && x <= r.X+r.W &&
		y >= r.Y && y <= r.Y+r.H
}
func crp(p sdl.FPoint, r sdl.FRect) bool {
	return p.X >= r.X && p.X <= r.X+r.W &&
		p.Y >= r.Y && p.Y <= r.Y+r.H
}
func cms(r sdl.FRect) bool {
	return crp(msp, r)
}
func crr(a, b sdl.FRect) bool {
	return a.X < b.X+b.W &&
		a.X+a.W > b.X &&
		a.Y < b.Y+b.H &&
		a.Y+a.H > b.Y
}
func GetIntersection(a, b sdl.FRect) (sdl.FRect, bool) {
	if !crr(a, b) {
		return sdl.FRect{}, false
	}

	minX := max(a.X, b.X)
	maxX := min(a.X+a.W, b.X+b.W)
	minY := max(a.Y, b.Y)
	maxY := min(a.Y+a.H, b.Y+b.H)

	return sdl.FRect{
		X: minX,
		Y: minY,
		W: maxX - minX,
		H: maxY - minY,
	}, true
}

// MARK:SCALE
func fitScale(width, height, maxSize float32) float32 {
	// Guard against zero, negative values, or invalid inputs
	if width <= 0 || height <= 0 || maxSize <= 0 {
		return 0
	}

	// Determine the dominant dimension
	maxDim := width
	if height > maxDim {
		maxDim = height
	}

	// Calculate uniform scale
	return maxSize / maxDim
}
func scale2wh(origW, origH, maxW, maxH float32) (float32, float32) {
	if origW == 0 || origH == 0 {
		return 0, 0
	}
	scaleW := maxW / origW
	scaleH := maxH / origH

	// Use smaller scale factor so it fits entirely inside
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	return origW * scale, origH * scale
}
func scaleH(origW, origH, newW float32) float32 {
	if origW == 0 {
		return 0
	}
	return (newW / origW) * origH
}
func scaleW(origW, origH, newH float32) float32 {
	if origH == 0 {
		return 0
	}
	return (newH / origH) * origW
}

// MARK: RANDOM
func RINT8(min, max uint8) uint8 {
	return min + uint8(rand.N(int(max-min)+1))
}
func RINT(min, max int) int {
	return min + rand.IntN(max-min+1)
}
func RI32(min, max int) int32 {
	return int32(min + rand.IntN(max-min+1))
}
func RF32(min, max float32) float32 {
	return min + rand.Float32()*(max-min)
}
func ROLL(sides int) int {
	if sides < 1 {
		return 1
	}
	return 1 + rand.N(sides)
}
func FLIP() bool {
	return rand.IntN(2) == 0
}

// MARK:POINTS
func recCorners(r sdl.FRect) []sdl.FPoint {
	var p []sdl.FPoint
	p = append(p, P(r.X, r.Y), P(r.X+r.W, r.Y), P(r.X+r.W, r.Y+r.H), P(r.X, r.Y+r.H))
	return p
}
func P[T1, T2 numConvert](x T1, y T2) sdl.FPoint {
	return sdl.FPoint{
		X: float32(x),
		Y: float32(y),
	}
}

// MARK: SLICES
func addPointUnique(points []sdl.FPoint, p sdl.FPoint) []sdl.FPoint {
	for _, pt := range points {
		if pt == p { // Exact match (pt.X == p.X && pt.Y == p.Y)
			return points // Duplicate found, do not append
		}
	}
	return append(points, p) // Added successfully
}
func addremSlice[T comparable](slice []T, val T) []T {
	if idx := slices.Index(slice, val); idx != -1 {
		return slices.Delete(slice, idx, idx+1) // Remove
	}
	return append(slice, val) // Add
}
func adduniqueSliceSlice[T comparable](list [][]T, target []T) [][]T {
	if len(target) == 0 {
		return list
	}
	for _, item := range list {
		if slices.Equal(item, target) {
			return list // Already exists, do nothing
		}
	}
	return append(list, slices.Clone(target))
}
func printIntSliceSeperator(nums []int, seperator string) string {
	if len(nums) == 0 {
		return ""
	}
	var sb strings.Builder
	for i, n := range nums {
		if i > 0 {
			sb.WriteString(seperator)
		}
		sb.WriteString(strconv.Itoa(n))
	}
	return sb.String()
}

// MARK:CONVERSIONS
func f32[T numConvert](val T) float32 {
	return float32(val)
}

type numConvert interface { //CONVERT ANY NUM TO F32
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}
