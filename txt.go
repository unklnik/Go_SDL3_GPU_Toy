package main

import (
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"github.com/Zyko0/go-sdl3/sdl"
	"github.com/Zyko0/go-sdl3/ttf"
)

var (
	infotxt, infotxt2, warntxt string

	infotxtTimer, infotxt2Timer, warntxtTimer int

	//KEYS INPUT
	activeInputPtr uintptr // Memory address of the variable currently being typed into
	inputBuffer    string  // Current text being typed
)

const (
	iconClose          = "\uF4C6"
	iconMin            = "\uF4D9"
	iconTick           = "\uEB7A"
	iconCursor         = "\uEF89"
	iconMove           = "\uEC5F"
	iconMax            = "\uF515"
	iconAdd            = "\uF4B1"
	iconDrawWire       = "\uF2F5"
	iconDrawShaded     = "\uF2F4"
	iconDrawSolid      = "\uF2F3"
	iconHand           = "\uF444"
	iconSunlight       = "\uF1BF"
	iconLightbulb      = "\uEEA9"
	iconLightbulbSolid = "\uEEA6"
	iconAmbient        = "\uF42E"
	iconPalette        = "\uEFC5"
	icon3D             = "\uF2F5"
	iconTexture        = "\uEE4B"
	icon3DPrimitives   = "\uF3DA"
	iconSettings       = "\uF0E7"
	iconVisible        = "\uECB5"
	iconRecycle        = "\uF05D"
	iconPencil         = "\uEFDF"
	iconlayers         = "\uF181"

	ic23d           = "\ufa97"
	ic22d           = "\ueb2c"
	ic23dprimitives = "\ufaae"
	ic23dmodels     = "\ufdb2"
	ic2lights       = "\uf09e"
	ic2texture      = "\ueb0a"
	ic2upArrow      = "\uea62"
	ic2downArrow    = "\uea5f"
	ic2grid         = "\ufca5"
	ic2patterndiag  = "\uf51b"
	ic2patternboxes = "\uedba"
	ic2patterncubes = "\uf2c9"
	ic2hidden       = "\uf7ec"
	ic2leftarrow    = "\ueb5e"
	ic2rightarrow   = "\ueb5f"
	ic2delete       = "\ueb41"
	ic2save         = "\ueb62"
	ic2line         = "\uec40"
	ic2triangle     = "\ueb44"
	ic2circle       = "\uea6b"
	ic2lineshape    = "\uefd0"
	ic2pentagon     = "\uefe3"
	ic2hexagon      = "\uec02"
	ic2polygon      = "\uecbd"
	ic2pencil       = "\ueb04"
	ic2polyline     = "\ueed9"
)

type TxtSpan struct {
	Text string
	C    sdl.Color
}

// MARK: DRAW
// dtxt2c draws two colored segments side by side (e.g. "Score: " in Gray, "100" in Yellow)
func dtxt2c(f *ttf.Font, t1, t2 string, c1, c2 sdl.Color, x, y float32) (float32, float32) {
	return dtxtSpans(f, x, y, TxtSpan{t1, c1}, TxtSpan{t2, c2})
}

// dtxtTopBottomWin: Draws text centered at top or bottom of screen with offset
func dtxtTopBottomWin(f *ttf.Font, t string, offset float32, topBottom bool, c sdl.Color) {
	x := CNT.X - txW(f, t)/2
	h := fonh(f)
	y := h + offset
	if topBottom {
		y = SET.SCRH - (h + offset)
	}
	dtxt(f, t, x, y, c)
}

// dtxtinrecfromright: Draws text inside rectangle from right side with offset center Y
func dtxtinrecfromright(f *ttf.Font, t string, r sdl.FRect, offset float32, c sdl.Color) {
	x := (r.X + r.W) - (txW(f, t) + offset)
	y := (r.Y + r.H/2) - fonh(f)/2
	dtxt(f, t, x, y, c)
}

// dtxtreccnt: Draw text in the center of rectangle
func dtxtreccnt(f *ttf.Font, t string, r sdl.FRect, c sdl.Color) {
	w, h := txWH(f, t)
	v2 := reccntr(r)
	x, y := v2.X-w/2, v2.Y-h/2
	dtxt(f, t, x, y, c)
}

// dtxtrecside: Draws text left/right (lr) of rectangle with center Y
func dtxtrecside(f *ttf.Font, t string, r sdl.FRect, lr bool, offset float32, c sdl.Color) {
	x := r.X + r.W + offset
	if lr {
		x = r.X - (txW(f, t) + offset)
	}
	y := (r.Y + r.W/2) - fonh(f)/2
	dtxt(f, t, x, y, c)
}

// dtxtsmlwhite: Draws text at X,y in white with small font
func dtxtsmlwhite(t string, x, y float32) {
	dtxt(Fon.Sml, t, x, y, Col.White)
}

// dtxtcntrp: Draws text centered on a point
func dtxtcntrp(f *ttf.Font, t string, p sdl.FPoint, c sdl.Color) {
	w := txW(f, t)
	x := p.X - w/2
	h := fonh(f)
	y := p.Y - h/2
	dtxt(f, t, x, y, c)
}

// dtxt: Draw text anywhere using TTF font and sdl.Color (returns width, height)
func dtxt(f *ttf.Font, t string, x, y float32, c sdl.Color) (float32, float32) {
	if R2D == nil || f == nil || t == "" {
		return 0, 0
	}

	textObj, err := R2D.textEngine.CreateText(f, t)
	if err != nil || textObj == nil {
		return 0, 0
	}
	defer textObj.Destroy()

	w, h, _ := textObj.Size()

	seq, err := textObj.GPUDrawData()
	if err != nil || seq == nil {
		return 0, 0
	}

	rf, gf, bf, af := color2float(c)

	var batch []TextVertex

	for curr := seq; curr != nil; curr = curr.Next {
		if curr.NumVertices == 0 || curr.NumIndices == 0 {
			continue
		}

		R2D.textAtlas = curr.AtlasTexture

		xyList := unsafe.Slice(curr.Xy, curr.NumVertices)
		uvList := unsafe.Slice(curr.Uv, curr.NumVertices)
		idxList := unsafe.Slice(curr.Indices, curr.NumIndices)

		for i := 0; i < int(curr.NumIndices); i++ {
			idx := idxList[i]
			batch = append(batch, TextVertex{
				X: x + xyList[idx].X,
				Y: y - xyList[idx].Y, // Inverts right-side up and starts at y
				U: uvList[idx].X,
				V: uvList[idx].Y,
				R: rf, G: gf, B: bf, A: af,
			})
		}
	}

	// Append text batch into the chronological draw command queue
	R2D.addTextVertices(batch)

	return float32(w), float32(h)
}

// dtxtSpans draws multiple colored text segments seamlessly on a single line in a SINGLE GPU draw call.
func dtxtSpans(f *ttf.Font, x, y float32, spans ...TxtSpan) (float32, float32) {
	if R2D == nil || f == nil || len(spans) == 0 {
		return 0, 0
	}

	curX := x
	var maxH float32
	var totalBatch []TextVertex

	for _, span := range spans {
		if span.Text == "" {
			continue
		}

		// 1. Create text object for this segment
		textObj, err := R2D.textEngine.CreateText(f, span.Text)
		if err != nil || textObj == nil {
			continue
		}

		w, h, _ := textObj.Size()
		if float32(h) > maxH {
			maxH = float32(h)
		}

		seq, err := textObj.GPUDrawData()
		if err != nil || seq == nil {
			textObj.Destroy()
			continue
		}

		// 2. Color for this specific segment
		rf, gf, bf, af := color2float(span.C)

		for curr := seq; curr != nil; curr = curr.Next {
			if curr.NumVertices == 0 || curr.NumIndices == 0 {
				continue
			}

			R2D.textAtlas = curr.AtlasTexture

			xyList := unsafe.Slice(curr.Xy, curr.NumVertices)
			uvList := unsafe.Slice(curr.Uv, curr.NumVertices)
			idxList := unsafe.Slice(curr.Indices, curr.NumIndices)

			for i := 0; i < int(curr.NumIndices); i++ {
				idx := idxList[i]
				totalBatch = append(totalBatch, TextVertex{
					X: curX + xyList[idx].X,
					Y: y - xyList[idx].Y,
					U: uvList[idx].X,
					V: uvList[idx].Y,
					R: rf, G: gf, B: bf, A: af,
				})
			}
		}

		textObj.Destroy()

		// Advance x by the width of this segment
		curX += float32(w)
	}

	// 3. Append the combined multi-colored vertices in ONE draw call
	if len(totalBatch) > 0 {
		R2D.addTextVertices(totalBatch)
	}

	return curX - x, maxH
}

// MARK: UTILS
// txlen: Measure pixel width of string without rendering
func txW(f *ttf.Font, t string) float32 {
	if f == nil || t == "" {
		return 0
	}
	w, _, _ := f.StringSize(t)
	return float32(w)
}
func txWH(f *ttf.Font, t string) (float32, float32) {
	w, h, _ := f.StringSize(t)
	return f32(w), f32(h)
}

// fonh: Measure font line height
func txH(f *ttf.Font, t string) float32 {
	if f == nil || t == "" {
		return 0
	}
	_, h, _ := f.StringSize(t)
	return float32(h)
}
func fonh(f *ttf.Font) float32 {
	h, _ := f.Size()
	return h
}

// MARK: UTILS
func mInfoTXT(t string) {
	infotxtTimer = int(SET.TargetFPS) * 2
	infotxt = t
}
func mInfoTXT2(t string) {
	infotxt2Timer = int(SET.TargetFPS) * 2
	infotxt2 = t
}
func mWarnTXT(t string) {
	warntxtTimer = int(SET.TargetFPS) * 2
	warntxt = t
}
func listRec(f *ttf.Font, spc, pad float32, u []UIITEM) sdl.FRect {
	h := listLen(f, spc, u) + pad
	var w float32
	for i := range u {
		if w < txW(f, u[i].nm) {
			w = txW(f, u[i].nm)
		}
	}
	w += pad
	return R(0, 0, w, h)
}
func listLen(f *ttf.Font, spc float32, u []UIITEM) float32 {
	return float32(len(u)-1) * (fonh(f) + spc/2)
}

// Helper: Converts an SDL keycode + Shift state into a typed character string
func getTypedChar(k sdl.Keycode, shift, lettersonly bool) string {
	// 1. Letters (a - z / A - Z)
	if k >= sdl.K_A && k <= sdl.K_Z {
		char := string(rune(k))
		if shift {
			return strings.ToUpper(char)
		}
		return strings.ToLower(char)
	}

	// If lettersonly is enabled, ignore all numbers and symbols!
	if lettersonly {
		return ""
	}

	// 2. Numbers (0 - 9)
	if k >= sdl.K_0 && k <= sdl.K_9 {
		if !shift {
			return string(rune(k))
		}
		// Shifted symbols on number row
		switch k {
		case sdl.K_1:
			return "!"
		case sdl.K_2:
			return "@"
		case sdl.K_3:
			return "#"
		case sdl.K_4:
			return "$"
		case sdl.K_5:
			return "%"
		case sdl.K_6:
			return "^"
		case sdl.K_7:
			return "&"
		case sdl.K_8:
			return "*"
		case sdl.K_9:
			return "("
		case sdl.K_0:
			return ")"
		}
	}

	// 3. Keypad & Special Characters
	switch k {
	case sdl.K_KP_0:
		return "0"
	case sdl.K_KP_1:
		return "1"
	case sdl.K_KP_2:
		return "2"
	case sdl.K_KP_3:
		return "3"
	case sdl.K_KP_4:
		return "4"
	case sdl.K_KP_5:
		return "5"
	case sdl.K_KP_6:
		return "6"
	case sdl.K_KP_7:
		return "7"
	case sdl.K_KP_8:
		return "8"
	case sdl.K_KP_9:
		return "9"

	case sdl.K_SPACE:
		return " "
	case sdl.K_MINUS, sdl.K_KP_MINUS:
		if shift {
			return "_"
		}
		return "-"
	case sdl.K_EQUALS:
		if shift {
			return "+"
		}
		return "="
	case sdl.K_PERIOD, sdl.K_KP_PERIOD:
		if shift {
			return ">"
		}
		return "."
	case sdl.K_COMMA:
		if shift {
			return "<"
		}
		return ","
	case sdl.K_SLASH, sdl.K_KP_DIVIDE:
		if shift {
			return "?"
		}
		return "/"
	}

	return ""
}

// keysInpTxt draws a text input box for strings with length limits and optional letters-only filtering
func keysInpTxt(f *ttf.Font, val *string, x, y, w, pad float32, lettersonly bool, maxLen int, hovertxt string) string {
	if val == nil {
		return ""
	}

	if maxLen <= 0 {
		maxLen = 32
	}

	h := fonh(f) + pad/2

	// Calculate a comfortable width based on max length
	sampleLen := maxLen
	if sampleLen > 18 {
		sampleLen = 18 // Keep width reasonable for long strings
	}

	r := R(x, y, w, h)
	id := uintptr(unsafe.Pointer(val))
	isFocused := activeInputPtr == id

	// 1. Mouse Focus / Unfocus
	if ms.lc {
		if cms(r) {
			activeInputPtr = id
			inputBuffer = *val
			isFocused = true
			ms.lc = false // Consume click
		} else if isFocused {
			*val = inputBuffer
			activeInputPtr = 0
			isFocused = false
		}
	}

	// 2. Keyboard Typing
	if isFocused && LastKey != 0 {
		isShift := Keys[sdl.K_LSHIFT] || Keys[sdl.K_RSHIFT] || Keys[sdl.K_CAPSLOCK]

		switch LastKey {
		case sdl.K_BACKSPACE:
			// UTF-8 safe backspace (removes by rune, not byte)
			runes := []rune(inputBuffer)
			if len(runes) > 0 {
				inputBuffer = string(runes[:len(runes)-1])
			}

		case sdl.K_RETURN, sdl.K_KP_ENTER, sdl.K_ESCAPE:
			*val = inputBuffer
			activeInputPtr = 0
			isFocused = false
			nameEdit = false

		default:
			// Check max length by rune count
			runes := []rune(inputBuffer)
			if len(runes) < maxLen {
				char := getTypedChar(LastKey, isShift, lettersonly)
				if char != "" {
					inputBuffer += char
				}
			}
		}
	}

	// 3. Draw Background & Border
	drr(r, SET.BGC)
	if isFocused {
		drrlw(r, 2.0, SET.UIC3)
	} else if cms(r) {
		drrlw(r, 2.0, SET.UIC4)
		if hovertxt != "" {
			mInfoTXT(hovertxt)
		}
	}

	// 4. Draw Text / Blinking Cursor
	var displayText string
	if isFocused {
		displayText = inputBuffer
		if (sdl.TicksNS()/500_000_000)%2 == 0 {
			displayText += "|"
		}
	} else {
		displayText = *val
	}

	dtxtreccnt(f, displayText, r, SET.UICtxt)
	return *val
}

// Helper: Formats a float32 cleanly without trailing zeros (e.g. 1.25 -> "1.25", 2.0 -> "2")
func formatF32(v float32) string {
	return strconv.FormatFloat(float64(v), 'f', -1, 32)
}

func keysInpF32(f *ttf.Font, val *float32, x, y, pad float32, minNum, maxNum float32, hovertxt string) float32 {
	if val == nil {
		return 0
	}

	h := fonh(f) + pad/2

	// Measure box width based on min/max possibilities with a little extra space for decimals
	maxLenText := formatF32(maxNum) + ".00"
	minLenText := formatF32(minNum) + ".00"
	if len(minLenText) > len(maxLenText) {
		maxLenText = minLenText
	}

	w := txW(f, maxLenText) + pad*2
	if w < 70 {
		w = 70
	}

	r := R(x, y, w, h)
	id := uintptr(unsafe.Pointer(val))
	isFocused := activeInputPtr == id

	// 1. Mouse Focus / Unfocus
	if ms.lc {
		if cms(r) {
			activeInputPtr = id
			if *val == 0 {
				inputBuffer = ""
			} else {
				// Cleanly format existing float into buffer with decimals
				inputBuffer = formatF32(*val)
			}
			isFocused = true
			ms.lc = false
		} else if isFocused {
			commitInputF32(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
		}
	}

	// 2. Typing Controls
	if isFocused && LastKey != 0 {
		appendNum := func(digit string) {
			if inputBuffer == "0" {
				inputBuffer = digit // Replace lonely leading zero
			} else if inputBuffer == "-0" {
				inputBuffer = "-" + digit
			} else {
				inputBuffer += digit
			}
		}

		switch LastKey {
		case sdl.K_0, sdl.K_KP_0:
			if inputBuffer == "0" || inputBuffer == "-0" {
				// Avoid "00"
			} else if inputBuffer == "" {
				inputBuffer = "0"
			} else if inputBuffer == "-" {
				inputBuffer = "-0"
			} else {
				// Allows typing zeroes after a decimal point (e.g. "1.0")
				inputBuffer += "0"
			}

		case sdl.K_1, sdl.K_KP_1:
			appendNum("1")
		case sdl.K_2, sdl.K_KP_2:
			appendNum("2")
		case sdl.K_3, sdl.K_KP_3:
			appendNum("3")
		case sdl.K_4, sdl.K_KP_4:
			appendNum("4")
		case sdl.K_5, sdl.K_KP_5:
			appendNum("5")
		case sdl.K_6, sdl.K_KP_6:
			appendNum("6")
		case sdl.K_7, sdl.K_KP_7:
			appendNum("7")
		case sdl.K_8, sdl.K_KP_8:
			appendNum("8")
		case sdl.K_9, sdl.K_KP_9:
			appendNum("9")

		// 🚀 Added: Decimal Point support
		case sdl.K_PERIOD, sdl.K_KP_PERIOD:
			if !strings.Contains(inputBuffer, ".") {
				if inputBuffer == "" || inputBuffer == "-" {
					inputBuffer += "0."
				} else {
					inputBuffer += "."
				}
			}

		case sdl.K_MINUS, sdl.K_KP_MINUS:
			if minNum < 0 && len(inputBuffer) == 0 {
				inputBuffer = "-"
			}

		case sdl.K_BACKSPACE:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
			}

		case sdl.K_RETURN, sdl.K_KP_ENTER, sdl.K_ESCAPE:
			commitInputF32(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
			widthEdit = false
			heightEdit = false
		}
	}

	drr(r, SET.BGC)
	if isFocused {
		drrlw(r, 2.0, SET.UIC3)
	} else if cms(r) {
		drrlw(r, 2, SET.UIC4)
		mInfoTXT(hovertxt)
	}

	// 4. Draw Text / Blinking Cursor
	var displayText string
	if isFocused {
		displayText = inputBuffer
		if (sdl.TicksNS()/500_000_000)%2 == 0 {
			displayText += "|"
		}
	} else {
		// Display float formatted with its decimal places
		displayText = formatF32(*val)
	}

	dtxtreccnt(f, displayText, r, SET.UICtxt)
	return *val
}
func commitInputF32(val *float32, minNum, maxNum float32) {
	if val == nil {
		return
	}

	// If empty or left with just symbols, default to 0 (clamped)
	if len(inputBuffer) == 0 || inputBuffer == "-" || inputBuffer == "." || inputBuffer == "-." {
		var f float32 = 0
		if f < minNum {
			f = minNum
		}
		if f > maxNum {
			f = maxNum
		}
		*val = f
		return
	}

	// 🚀 Use ParseFloat (with 32-bit precision) instead of Atoi
	n, err := strconv.ParseFloat(inputBuffer, 32)
	if err == nil {
		f := float32(n)
		if f < minNum {
			f = minNum
		}
		if f > maxNum {
			f = maxNum
		}
		*val = f
	}
}

func keysInpInt(f *ttf.Font, val *int, x, y, pad float32, minNum, maxNum int, hovertxt string) int {
	if val == nil {
		return 0
	}

	h := fonh(f) + +pad/2
	maxLenText := fmt.Sprint(maxNum)
	if len(fmt.Sprint(minNum)) > len(maxLenText) {
		maxLenText = fmt.Sprint(minNum)
	}
	w := txW(f, maxLenText) + pad*2
	if w < 60 {
		w = 60
	}

	r := R(x, y, w, h)
	id := uintptr(unsafe.Pointer(val))
	isFocused := activeInputPtr == id

	// 1. Mouse Focus / Unfocus
	if ms.lc {
		if cms(r) {
			activeInputPtr = id
			// 🚀 FIX 1: If 0, start with an EMPTY buffer so only the cursor blinks
			if *val == 0 {
				inputBuffer = ""
			} else {
				inputBuffer = strconv.Itoa(*val)
			}
			isFocused = true
			ms.lc = false
		} else if isFocused {
			commitInputInt(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
		}
	}

	// 2. Typing Controls
	if isFocused && LastKey != 0 {
		// Helper to prevent leading zeros
		appendNum := func(digit string) {
			if inputBuffer == "0" {
				inputBuffer = digit // Replace leading zero
			} else {
				inputBuffer += digit
			}
		}

		switch LastKey {
		case sdl.K_0, sdl.K_KP_0:
			if inputBuffer != "" && inputBuffer != "0" {
				inputBuffer += "0"
			}
		case sdl.K_1, sdl.K_KP_1:
			appendNum("1")
		case sdl.K_2, sdl.K_KP_2:
			appendNum("2")
		case sdl.K_3, sdl.K_KP_3:
			appendNum("3")
		case sdl.K_4, sdl.K_KP_4:
			appendNum("4")
		case sdl.K_5, sdl.K_KP_5:
			appendNum("5")
		case sdl.K_6, sdl.K_KP_6:
			appendNum("6")
		case sdl.K_7, sdl.K_KP_7:
			appendNum("7")
		case sdl.K_8, sdl.K_KP_8:
			appendNum("8")
		case sdl.K_9, sdl.K_KP_9:
			appendNum("9")

		case sdl.K_MINUS, sdl.K_KP_MINUS:
			if minNum < 0 && len(inputBuffer) == 0 {
				inputBuffer = "-"
			}

		case sdl.K_BACKSPACE:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
			}

		case sdl.K_RETURN, sdl.K_KP_ENTER, sdl.K_ESCAPE:
			commitInputInt(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
		}
	}

	drr(r, SET.BGC)
	if isFocused {
		drrlw(r, 2.0, SET.UIC3)
	} else if cms(r) {
		drrlw(r, 2, SET.UIC4)
		mInfoTXT(hovertxt)
	}

	// 4. Draw Text / Blinking Cursor
	var displayText string
	if isFocused {
		displayText = inputBuffer
		// Blinks cursor every 500ms
		if (sdl.TicksNS()/500_000_000)%2 == 0 {
			displayText += "|"
		}
	} else {
		displayText = fmt.Sprint(*val)
	}

	dtxtreccnt(f, displayText, r, SET.UICtxt)
	return *val
}

func commitInputInt(val *int, minNum, maxNum int) {
	// If left empty on exit, default to 0
	if len(inputBuffer) == 0 || inputBuffer == "-" {
		n := 0
		if n < minNum {
			n = minNum
		}
		if n > maxNum {
			n = maxNum
		}
		*val = n
		return
	}
	n, err := strconv.Atoi(inputBuffer)
	if err == nil {
		if n < minNum {
			n = minNum
		}
		if n > maxNum {
			n = maxNum
		}
		*val = n
	}
}

// keysInpU8 draws an input box for uint8 values (0 - 255, no negative numbers or decimals)
func keysInpU8(f *ttf.Font, val *uint8, x, y, pad float32, minNum, maxNum uint8, hovertxt string) uint8 {
	if val == nil {
		return 0
	}

	h := fonh(f) + pad/2
	maxLenText := fmt.Sprint(maxNum)
	if len(fmt.Sprint(minNum)) > len(maxLenText) {
		maxLenText = fmt.Sprint(minNum)
	}
	w := txW(f, maxLenText) + pad*2
	if w < 60 {
		w = 60
	}

	r := R(x, y, w, h)
	id := uintptr(unsafe.Pointer(val))
	isFocused := activeInputPtr == id

	// 1. Mouse Focus / Unfocus
	if ms.lc {
		if cms(r) {
			activeInputPtr = id
			// Start with an empty buffer if 0 so only the cursor blinks
			if *val == 0 {
				inputBuffer = ""
			} else {
				inputBuffer = strconv.Itoa(int(*val))
			}
			isFocused = true
			ms.lc = false
		} else if isFocused {
			commitInputU8(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
		}
	}

	// 2. Typing Controls
	if isFocused && LastKey != 0 {
		appendNum := func(digit string) {
			if inputBuffer == "0" {
				inputBuffer = digit // Replace leading zero
			} else {
				inputBuffer += digit
			}
		}

		switch LastKey {
		case sdl.K_0, sdl.K_KP_0:
			if inputBuffer != "" && inputBuffer != "0" {
				inputBuffer += "0"
			}
		case sdl.K_1, sdl.K_KP_1:
			appendNum("1")
		case sdl.K_2, sdl.K_KP_2:
			appendNum("2")
		case sdl.K_3, sdl.K_KP_3:
			appendNum("3")
		case sdl.K_4, sdl.K_KP_4:
			appendNum("4")
		case sdl.K_5, sdl.K_KP_5:
			appendNum("5")
		case sdl.K_6, sdl.K_KP_6:
			appendNum("6")
		case sdl.K_7, sdl.K_KP_7:
			appendNum("7")
		case sdl.K_8, sdl.K_KP_8:
			appendNum("8")
		case sdl.K_9, sdl.K_KP_9:
			appendNum("9")

		case sdl.K_BACKSPACE:
			if len(inputBuffer) > 0 {
				inputBuffer = inputBuffer[:len(inputBuffer)-1]
			}

		case sdl.K_RETURN, sdl.K_KP_ENTER, sdl.K_ESCAPE:
			commitInputU8(val, minNum, maxNum)
			activeInputPtr = 0
			isFocused = false
		}
	}

	// 3. Draw Box
	drr(r, SET.BGC)
	if isFocused {
		drrlw(r, 2.0, SET.UIC3)
	} else if cms(r) {
		drrlw(r, 2, SET.UIC4)
		mInfoTXT(hovertxt)
	}

	// 4. Draw Text / Blinking Cursor
	var displayText string
	if isFocused {
		displayText = inputBuffer
		if (sdl.TicksNS()/500_000_000)%2 == 0 {
			displayText += "|"
		}
	} else {
		displayText = fmt.Sprint(*val)
	}

	dtxtreccnt(f, displayText, r, SET.UICtxt)
	return *val
}

// Internal commit helper for uint8
func commitInputU8(val *uint8, minNum, maxNum uint8) {
	if len(inputBuffer) == 0 {
		*val = minNum
		return
	}
	n, err := strconv.Atoi(inputBuffer)
	if err == nil {
		if n < int(minNum) {
			n = int(minNum)
		}
		if n > int(maxNum) {
			n = int(maxNum)
		}
		*val = uint8(n)
	}
}
