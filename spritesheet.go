package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/Zyko0/go-sdl3/sdl"
)

// FastAlpha reads pixel alpha directly from raw memory slices (15x faster than img.At)
func getPixelAlpha(img image.Image, x, y int) uint8 {
	bounds := img.Bounds()
	if x < bounds.Min.X || x >= bounds.Max.X || y < bounds.Min.Y || y >= bounds.Max.Y {
		return 0
	}

	switch m := img.(type) {
	case *image.NRGBA:
		return m.Pix[(y-m.Rect.Min.Y)*m.Stride+(x-m.Rect.Min.X)*4+3]
	case *image.RGBA:
		return m.Pix[(y-m.Rect.Min.Y)*m.Stride+(x-m.Rect.Min.X)*4+3]
	default:
		_, _, _, a := m.At(x, y).RGBA()
		return uint8(a >> 8)
	}
}

// AutoSliceSpriteSheet inspects a PNG file, determines frame grid cells,
// and extracts pixel-tight collision bounding boxes.
func AutoSliceSpriteSheet(path string, tex *Texture) (frames []sdl.FRect, colFrames []sdl.FRect) {
	if tex == nil || tex.W <= 0 || tex.H <= 0 {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		full := sdl.FRect{X: 0, Y: 0, W: tex.W, H: tex.H}
		return []sdl.FRect{full}, []sdl.FRect{full}
	}
	defer file.Close()

	imgData, _, err := image.Decode(file)
	if err != nil {
		full := sdl.FRect{X: 0, Y: 0, W: tex.W, H: tex.H}
		return []sdl.FRect{full}, []sdl.FRect{full}
	}

	bounds := imgData.Bounds()
	wInt := int(tex.W)
	hInt := int(tex.H)

	numCols := 1
	numRows := 1
	frameW := tex.W
	frameH := tex.H

	// 1. Single Horizontal Strip Check (e.g. 600x100 -> 6 columns of 100x100)
	if hInt > 0 && wInt >= hInt && (wInt%hInt == 0) {
		numCols = wInt / hInt
		numRows = 1
		frameW = tex.H
		frameH = tex.H
	} else {
		// 2. 2D Grid Check: test standard square cell sizes in ascending order
		candidates := []int{16, 24, 32, 48, 64, 80, 96, 128, 256}
		detected := false

		for _, sz := range candidates {
			if wInt%sz == 0 && hInt%sz == 0 && (wInt/sz)*(hInt/sz) > 1 {
				cuts := countGridCutPixels(imgData, bounds, wInt, hInt, sz)
				// If grid lines don't slice solid pixels (tolerance <= 15 for compression noise)
				if cuts <= 15 {
					numCols = wInt / sz
					numRows = hInt / sz
					frameW = float32(sz)
					frameH = float32(sz)
					detected = true
					break
				}
			}
		}

		if !detected {
			numCols = 1
			numRows = 1
			frameW = tex.W
			frameH = tex.H
		}
	}

	// 3. Slice Frames and calculate tight bounds per frame
	for r := 0; r < numRows; r++ {
		for c := 0; c < numCols; c++ {
			fx := float32(c) * frameW
			fy := float32(r) * frameH

			frames = append(frames, sdl.FRect{X: fx, Y: fy, W: frameW, H: frameH})
			tightBox := getTightBounds(imgData, int(fx), int(fy), int(frameW), int(frameH))
			colFrames = append(colFrames, tightBox)
		}
	}

	return frames, colFrames
}

func countGridCutPixels(imgData image.Image, bounds image.Rectangle, wInt, hInt, sz int) int {
	cutPixels := 0
	cols := wInt / sz
	rows := hInt / sz

	// Horizontal grid lines
	for r := 1; r < rows; r++ {
		y := bounds.Min.Y + r*sz
		for x := 0; x < wInt; x++ {
			if getPixelAlpha(imgData, bounds.Min.X+x, y) > 0 {
				cutPixels++
			}
		}
	}

	// Vertical grid lines
	for c := 1; c < cols; c++ {
		x := bounds.Min.X + c*sz
		for y := 0; y < hInt; y++ {
			if getPixelAlpha(imgData, x, bounds.Min.Y+y) > 0 {
				cutPixels++
			}
		}
	}

	return cutPixels
}

func getTightBounds(imgData image.Image, fx, fy, fw, fh int) sdl.FRect {
	minX, maxX := fw, 0
	minY, maxY := fh, 0
	hasPixels := false

	for y := 0; y < fh; y++ {
		for x := 0; x < fw; x++ {
			if getPixelAlpha(imgData, fx+x, fy+y) > 0 {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
				hasPixels = true
			}
		}
	}

	if !hasPixels {
		return sdl.FRect{X: 0, Y: 0, W: float32(fw), H: float32(fh)}
	}

	return sdl.FRect{
		X: float32(minX),
		Y: float32(minY),
		W: float32(maxX - minX + 1),
		H: float32(maxY - minY + 1),
	}
}

// GetCollisionRect transforms the current frame's tight bounds into screen/world coordinates.
// Accurately accounts for scale and horizontal/vertical flipping.
func (t *TEX2D) GetCollisionRect() sdl.FRect {
	if t == nil || len(t.frames) == 0 || len(t.colFrames) == 0 {
		return t.recBorder
	}

	frameIdx := t.curFrame
	if len(t.partialFrames) > 0 {
		if frameIdx >= len(t.partialFrames) {
			frameIdx = 0
		}
		frameIdx = t.partialFrames[frameIdx]
	}

	if frameIdx < 0 || frameIdx >= len(t.colFrames) {
		frameIdx = 0
	}

	scale := t.scale
	if scale <= 0 {
		scale = 1.0
	}

	c := t.colFrames[frameIdx]
	srcFrame := t.frames[frameIdx]

	colX := c.X
	colY := c.Y

	// Mirror tight bounds if flipped horizontally
	if t.flipX {
		colX = srcFrame.W - (c.X + c.W)
	}

	// Mirror tight bounds if flipped vertically
	if t.flipY {
		colY = srcFrame.H - (c.Y + c.H)
	}

	return sdl.FRect{
		X: t.recBorder.X + (colX * scale),
		Y: t.recBorder.Y + (colY * scale),
		W: c.W * scale,
		H: c.H * scale,
	}
}

func (t *TEX2D) GetCollisionRectAll() ([]sdl.FRect, []sdl.FRect) {
	if t == nil || len(t.frames) == 0 || len(t.colFrames) == 0 {
		return []sdl.FRect{t.recBorder}, []sdl.FRect{t.recBorder}
	}

	frameIdx := t.curFrame
	if len(t.partialFrames) > 0 {
		if frameIdx >= len(t.partialFrames) {
			frameIdx = 0
		}
		frameIdx = t.partialFrames[frameIdx]
	}

	if frameIdx < 0 || frameIdx >= len(t.colFrames) {
		frameIdx = 0
	}

	scale := t.scale
	if scale <= 0 {
		scale = 1.0
	}

	var df []sdl.FRect
	var cf []sdl.FRect

	for i := range t.colFrames {
		c := t.colFrames[i]
		srcFrame := t.frames[i]
		c.X += srcFrame.X
		c.Y += srcFrame.Y
		colX := c.X
		colY := c.Y

		// Mirror tight bounds if flipped horizontally
		if t.flipX {
			colX += srcFrame.W - (c.X + c.W) //MARK: FIX FLIPXY COLLIS RECS
		}

		// Mirror tight bounds if flipped vertically
		if t.flipY {
			colY += srcFrame.H - (c.Y + c.H)
		}
		df = append(df, sdl.FRect{
			X: t.recBorder.X + (srcFrame.X * scale),
			Y: t.recBorder.Y + (srcFrame.Y * scale),
			W: srcFrame.W * scale,
			H: srcFrame.H * scale,
		})
		cf = append(cf, sdl.FRect{
			X: t.recBorder.X + (colX * scale),
			Y: t.recBorder.Y + (colY * scale),
			W: c.W * scale,
			H: c.H * scale,
		})
	}

	return df, cf
}

// Helper methods on TEX2D
func (t *TEX2D) Play(animNum int, fps float32, loop bool) {
	t.anims[animNum].fps = fps
	t.anims[animNum].loop = loop
	t.anims[animNum].playing = true
}

func (t *TEX2D) Stop() {
	t.playing = false
	t.curFrame = 0
	t.animTimer = 0
}

// TEXTURE ANIMATION PREVIEW
func dAnimPreview(t *TEX2D, p sdl.FPoint, a *ANIM2D, maxW float32) {
	if t == nil || t.tex == nil || t.tex.GPU == nil {
		return
	}
	//t.cnt = p

	// 1. Determine active source frame
	src := sdl.FRect{X: 0, Y: 0, W: t.tex.W, H: t.tex.H}

	if len(a.frames) > 0 {
		frameIdx := a.curFrame
		if frameIdx >= 0 && frameIdx < len(a.frames) {
			src = a.frames[frameIdx]
		}
	}

	// 🚀 Compute destination rectangle using the FRAME dimensions, not the entire sheet!
	var w, h float32
	if src.W > src.H {
		h = scaleH(src.W, src.H, maxW)
		w = maxW
	} else {
		w = scaleW(src.W, src.H, maxW)
		h = maxW
	}

	dst := sdl.FRect{
		X: p.X - w*0.5,
		Y: p.Y - h*0.5,
		W: w,
		H: h,
	}
	//t.recBorder = dst

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
	// 4. Draw rotated, flipped, with active shader effect
	dtexrectEx(t.tex, src, dst, t.ro, t.cnt, t.flipX, t.flipY, t.fx, t.fxParam, c)
	if a.playing && len(a.frames) > 1 && a.fps > 0 {
		a.animTimer += timer.DT
		frameDuration := 1.0 / a.fps

		totalFrames := len(a.frames)

		for a.animTimer >= frameDuration {
			a.animTimer -= frameDuration
			a.curFrame++

			if a.curFrame >= totalFrames {
				if a.loop {
					a.curFrame = 0
				} else {
					a.curFrame = totalFrames - 1
					a.playing = false
					break
				}
			}
		}
	}
}
