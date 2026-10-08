package main

import (
	"slices"

	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	lev3d       *LEV3DOBJ
	lev2d       *LEV2DOBJ
	lev         LEVEL
	settingsVis bool

	prevselTEXnum, prevselOBJ3Dnum int
	prevREC2Dpos                   sdl.FPoint

	//NEW OBJ
	lineChangeActiv, recMoveActiv = -1, -1
	lineVal                       LINE2D
	gridVal                       GRID2D
	recVal                        REC2D
	triVal                        TRI2D

	triNewActiv, lineNewActiv, polylineNewActiv, recNewActiv bool

	//UNITS
	un float32 = 64

	unh, unq, un8th, un16th, un3q, un8thx3 = un / 2, un / 4, un / 8, un / 16, (un / 4) * 3, (un / 8) * 3

	un2, un3, un4, un5, un6, un7, un8, un9, un10 = un * 2, un * 3, un * 4, un * 5, un * 6, un * 7, un * 8, un * 9, un * 10
)

const (
	FxNone      uint8 = iota
	FxGrayscale       // 1. Black & White
	FxSepia           // 2. Vintage Warm Tone
	FxPosterize       // 3. Retro Comic Color Bands (fxParam = bands, e.g. 4.0)
	FxScanlines       // 4. Animated CRT / Hologram
)

type LEVEL struct {
	l3d LEV3DOBJ
	l2d LEV2DOBJ
}
type V3 struct {
	X, Y, Z float32
}

type ANIM2D struct {
	nm        string
	curFrame  int
	fps       float32
	animTimer float32
	playing   bool
	loop      bool
	frames    []sdl.FRect
	cbg       sdl.Color
}
type TEX2D struct {
	nm, path                          string
	cat, id                           int
	tags                              []string
	tex                               *Texture
	frames                            []sdl.FRect // 1. Source UV rectangles on the spritesheet
	colFrames                         []sdl.FRect // 2. Tight collision boxes (relative to frame 0,0)
	shadowX, shadowY, W, H, scale, ro float32
	c, cShadow                        sdl.Color
	cnt                               sdl.FPoint
	wide, high                        bool
	flipX, flipY                      bool
	shadow, dAllFrames                bool
	fx                                uint8
	fxParam                           float32

	curAnim int
	anims   []ANIM2D

	// 🚀 Animation Control
	curFrame      int     // Current active frame index
	fps           float32 // Playback speed in frames per second (e.g. 10.0)
	animTimer     float32 // Delta-time accumulator
	playing       bool    // True if animation is running
	loop          bool    // True to loop continuously, false to play once
	partialFrames []int   // Optional frame sequence (e.g. [0, 1, 2, 3] for walk)

	recBorder    sdl.FRect // Visual screen bounds
	recCollision sdl.FRect // 🚀 Tight world/screen collision box
	recPoints    []sdl.FPoint
}
type OBJ3D struct {
	nm, path string
	cat, id  int
	tags     []string
	mod      *GPUModel
	thumb    *Texture
	dropV2   sdl.FPoint
	cnt      V3

	Color     sdl.Color
	CustomTex *Texture
	UVScale   float32 //Texture tile scale

	RotSpdX, RotSpdY, RotSpdZ float32
	WorldX, WorldY, WorldZ    float32
	RotX, RotY, RotZ          float32
	Scale                     float32
}

type LEV2DOBJ struct {
	gr  []GRID2D
	tex []TEX2D
	li  []LINE2D
	rec []REC2D
	tri []TRI2D
}

type LEV3DOBJ struct {
	o3d []OBJ3D
}

func UP() { //MARK: UPDATE ████
	uINP()
	B4()
	uTimers()
	// 🚀 UPDATE 2D SPRITE ANIMATIONS & COLLISION HITBOXES
	if lev2d != nil {
		for i := range lev2d.tex {
			t := &lev2d.tex[i]

			// Advance animation if playing and has multiple frames
			if len(t.anims) > 0 {
				if t.anims[t.curAnim].playing && len(t.anims[t.curAnim].frames) > 1 && t.anims[t.curAnim].fps > 0 {
					t.anims[t.curAnim].animTimer += timer.DT
					frameDuration := 1.0 / t.fps

					totalFrames := len(t.anims[t.curAnim].frames)

					for t.anims[t.curAnim].animTimer >= frameDuration {
						t.anims[t.curAnim].animTimer -= frameDuration
						t.anims[t.curAnim].curFrame++

						if t.anims[t.curAnim].curFrame >= totalFrames {
							if t.anims[t.curAnim].loop {
								t.anims[t.curAnim].curFrame = 0
							} else {
								t.anims[t.curAnim].curFrame = totalFrames - 1
								t.anims[t.curAnim].playing = false
								break
							}
						}
					}
				}
			}

			/*
				if t.playing && len(t.frames) > 1 && t.fps > 0 {
					t.animTimer += timer.DT
					frameDuration := 1.0 / t.fps

					totalFrames := len(t.frames)
					if len(t.partialFrames) > 0 {
						totalFrames = len(t.partialFrames)
					}

					for t.animTimer >= frameDuration {
						t.animTimer -= frameDuration
						t.curFrame++

						if t.curFrame >= totalFrames {
							if t.loop {
								t.curFrame = 0
							} else {
								t.curFrame = totalFrames - 1
								t.playing = false
								break
							}
						}
					}
				}
			*/

			// Update tight collision box in screen space
			t.recCollision = t.GetCollisionRect()
		}
	}
	// UPDATE OBJ3D
	for i := range lev3d.o3d {
		uLevOBJ3D(&lev3d.o3d[i])
	}
	UpdateLights(timer.DT)
}

// Update object motion and auto-rotation
func uLevOBJ3D(o *OBJ3D) {
	if o.RotSpdX != 0 {
		o.RotX += o.RotSpdX * timer.DT
	}
	if o.RotSpdY != 0 {
		o.RotY += o.RotSpdY * timer.DT
	}
	if o.RotSpdZ != 0 {
		o.RotZ += o.RotSpdZ * timer.DT
	}
}
func B4() { //MARK: BEFORE ████
	prevselTEXnum = selTEXnum
	prevselOBJ3Dnum = selOBJ3Dnum
}
func AFTER() { //MARK: AFTER ████
	if prevselTEXnum != selTEXnum {
		sbLmenuNum = 2
	}
	if prevselOBJ3Dnum != selOBJ3Dnum {
		sbLmenuNum = 1
	}
}

func addTEX(t TEX2D) {
	t.cnt = mslcP
	t.recBorder.X = t.cnt.X - t.recBorder.W/2
	t.recBorder.Y = t.cnt.Y - t.recBorder.H/2
	lev2d.tex = append(lev2d.tex, t)
}
func addOBJ3D(o []OBJ3D) {
	for i := range o {
		if o[i].Color.A == 0 {
			o[i].Color = Col.White
		}
		if o[i].UVScale <= 0 {
			o[i].UVScale = 1.0 // Default 1x scale
		}
		if o[i].Scale <= 0 {
			o[i].Scale = 0.5
		}
		lev3d.o3d = append(lev3d.o3d, o[i])
	}
}

func mLevDefault() {
	lev3d = &lev.l3d
	lev2d = &lev.l2d
}

// MARK: UTILS ████
func MoveRec2DstepForward(rs []REC2D, idx int) int {
	// Can't move forward if already at the top (index 0) or out of bounds
	if idx <= 0 || idx >= len(rs) {
		return idx
	}

	// In-place swap with neighbor towards 0
	rs[idx], rs[idx-1] = rs[idx-1], rs[idx]

	return idx - 1
}

func MoveRec2DstepBackward(rs []REC2D, idx int) int {
	// Can't move backward if already at the bottom (last index) or out of bounds
	if idx < 0 || idx >= len(rs)-1 {
		return idx
	}

	// In-place swap with neighbor towards the end
	rs[idx], rs[idx+1] = rs[idx+1], rs[idx]

	return idx + 1
}
func MoveRec2DtoFront(rs []REC2D, idx int) int {
	if idx <= 0 || idx >= len(rs) {
		return idx // Already at index 0 or invalid
	}

	// 1. Save the selected item
	item := rs[idx]

	// 2. Shift all elements before idx one position to the right
	copy(rs[1:idx+1], rs[0:idx])

	// 3. Place item at the beginning
	rs[0] = item

	return 0 // New index is 0
}
func MoveRec2DtoEnd(rs []REC2D, idx int) int {
	lastIdx := len(rs) - 1
	if idx < 0 || idx >= lastIdx {
		return idx // Already at the end or invalid
	}

	// 1. Save the selected item
	item := rs[idx]

	// 2. Shift all elements after idx one position to the left
	copy(rs[idx:], rs[idx+1:])

	// 3. Place item at the very end
	rs[lastIdx] = item

	return lastIdx // New index is the end
}

// MoveTexToFront: moves TEX2D struct to front (beginning) of slice
func MoveTexToFront(texSlice []TEX2D, idx int) int {
	if idx <= 0 || idx >= len(texSlice) {
		return idx // Already at index 0 or invalid
	}

	// 1. Save the selected item
	item := texSlice[idx]

	// 2. Shift all elements before idx one position to the right
	copy(texSlice[1:idx+1], texSlice[0:idx])

	// 3. Place item at the beginning
	texSlice[0] = item

	return 0 // New index is 0
}

// MoveTexToEnd: moves TEX2D struct to back (end) of slice
func MoveTexToEnd(texSlice []TEX2D, idx int) int {
	lastIdx := len(texSlice) - 1
	if idx < 0 || idx >= lastIdx {
		return idx // Already at the end or invalid
	}

	// 1. Save the selected item
	item := texSlice[idx]

	// 2. Shift all elements after idx one position to the left
	copy(texSlice[idx:], texSlice[idx+1:])

	// 3. Place item at the very end
	texSlice[lastIdx] = item

	return lastIdx // New index is the end
}

// MoveTexStepBackward moves the texture 1 step closer to the end (pushes it down 1 visual layer).
// Returns the new index of the texture (idx + 1).
func MoveTexStepBackward(texSlice []TEX2D, idx int) int {
	// Can't move backward if already at the bottom (last index) or out of bounds
	if idx < 0 || idx >= len(texSlice)-1 {
		return idx
	}

	// In-place swap with neighbor towards the end
	texSlice[idx], texSlice[idx+1] = texSlice[idx+1], texSlice[idx]

	return idx + 1
}

// MoveTexStepForward moves the texture 1 step closer to index 0 (brings it up 1 visual layer).
// Returns the new index of the texture (idx - 1).
func MoveTexStepForward(texSlice []TEX2D, idx int) int {
	// Can't move forward if already at the top (index 0) or out of bounds
	if idx <= 0 || idx >= len(texSlice) {
		return idx
	}

	// In-place swap with neighbor towards 0
	texSlice[idx], texSlice[idx-1] = texSlice[idx-1], texSlice[idx]

	return idx - 1
}

// isSameTEX2D checks whether two TEX2D structs are duplicates.
func isSameTEX2D(a, b *TEX2D) bool {
	// 1. Asset & Texture Identity (Fails immediately if different images)
	if a.path != b.path || a.nm != b.nm || a.tex != b.tex || a.cat != b.cat || a.id != b.id {
		return false
	}

	// 2. Visual Transforms
	if a.scale != b.scale || a.ro != b.ro || a.flipX != b.flipX || a.flipY != b.flipY {
		return false
	}

	// 3. Color & Effects
	if a.c != b.c || a.fx != b.fx || a.fxParam != b.fxParam {
		return false
	}

	// 4. Shadow Properties
	if a.shadow != b.shadow || a.cShadow != b.cShadow || a.shadowX != b.shadowX || a.shadowY != b.shadowY {
		return false
	}

	// 5. Animation Settings
	if a.curAnim != b.curAnim || a.fps != b.fps || a.loop != b.loop {
		return false
	}

	// 6. Tags Slice Check
	if len(a.tags) != len(b.tags) {
		return false
	}
	for i := range a.tags {
		if a.tags[i] != b.tags[i] {
			return false
		}
	}

	// ⚠️ NOTE ON POSITION:
	// If you want duplicate checking to IGNORE screen position (i.e. saving a preset template),
	// leave 'cnt' and 'recBorder' out.
	// If you want exact screen position to matter, uncomment this line:
	// if a.cnt != b.cnt { return false }

	return true
}
func addUniqueTEX2D(list []TEX2D, target TEX2D) ([]TEX2D, bool) {
	for i := range list {
		if isSameTEX2D(&list[i], &target) {
			return list, false // Duplicate found; do not append
		}
	}

	return append(list, target), true // Successfully added
}

// DeleteTexAt removes the texture at index 'idx' while preserving layer order.
// Returns the updated slice and the new valid selected index.
func DeleteTEX2Dat(texSlice []TEX2D, idx int) ([]TEX2D, int) {
	if idx < 0 || idx >= len(texSlice) {
		return texSlice, idx // Out of bounds, do nothing
	}

	// Remove item (slices.Delete also zeroes the trailing slot to prevent memory leaks)
	texSlice = slices.Delete(texSlice, idx, idx+1)

	// Adjust selection index
	newIdx := idx
	if len(texSlice) == 0 {
		newIdx = -1
	} else if newIdx >= len(texSlice) {
		newIdx = len(texSlice) - 1
	}

	return texSlice, newIdx
}
