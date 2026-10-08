package main

import (
	"github.com/Zyko0/go-sdl3/sdl"
)

var (
	paletteTarget                 *sdl.Color // Pointer to the variable being edited
	paletteAlpha                  uint8      = 255
	PaletteModeAlpha, PaletteMode bool
)

// MARK: DRAW
// dPalette: Draws the colors in PaletteList slice
func dPalette(siz float32) (sdl.Color, bool) {
	txt := "Hover over a color for name"
	var pickedColor sdl.Color
	hasPicked := false

	// 1. Dark Fullscreen Overlay (Block background view)
	dr(0, 0, SET.SCRW, SET.SCRH, Col.BlackPure)
	dBGiconPattern(ixx2, ic2patternboxes, 0, scrREC, CA(Col.BlackCharcoal, 100))

	x, y, spc := unq, un+unq, un8th
	ox := x
	r := R(x, y, siz, siz)

	// 2. Draw Color Swatches
	for _, pc := range PaletteList {
		drr(r, pc.Color)

		if cms(r) {
			// Highlight hovered swatch with animated/contrasting outline
			drrlw(r, 2.0, Col.White)
			txt = pc.Name

			// Left-click to pick color
			if ms.lc {
				pickedColor = pc.Color
				if paletteTarget != nil {
					*paletteTarget = pc.Color // 🚀 Instantly modifies the target variable!
				}
				hasPicked = true
				break // Stop checking other swatches this frame
			}
		}

		r.X += r.W + spc
		if r.X+r.W+spc > SET.SCRW-unq {
			r.X = ox
			r.Y += r.H + spc
		}
	}

	// 3. Hovered Color Name Tooltip
	if txt != "" {
		dtxt(Fon.Def, txt, un8th, un8th, Col.White)
	}

	// 4. Close button or picked selection closes the palette
	if closeIcon(SET.SCRW-(unh+un16th), un16th, unh) || hasPicked {
		PaletteMode = false
		ms.lc = false // 🚀 Consume click so it doesn't click items behind the palette!
	}

	return pickedColor, hasPicked
}
func dPaletteAlpha(siz float32) (sdl.Color, bool) {
	txt := "Hover over a color for name"
	var pickedColor sdl.Color
	hasPicked := false

	// 1. Dark Fullscreen Overlay (Block background view)
	dr(0, 0, SET.SCRW, SET.SCRH, Col.BlackPure)
	dBGiconPattern(ixx2, ic2patternboxes, 0, scrREC, CA(Col.BlackCharcoal, 100))

	x, y, spc := unq, un+unq, un8th
	ox := x
	r := R(x, y, siz, siz)

	// 2. Draw Color Swatches
	for _, pc := range PaletteList {
		drr(r, CA(pc.Color, paletteAlpha))
		drrl(r, pc.Color)

		if cms(r) {
			// Highlight hovered swatch with animated/contrasting outline
			drrlw(r, 2.0, Col.White)
			txt = pc.Name

			// Left-click to pick color
			if ms.lc {
				pickedColor = CA(pc.Color, paletteAlpha)
				if paletteTarget != nil {
					*paletteTarget = CA(pc.Color, paletteAlpha) // 🚀 Instantly modifies the target variable!
				}
				hasPicked = true
				break // Stop checking other swatches this frame
			}
		}

		r.X += r.W + spc
		if r.X+r.W+spc > SET.SCRW-unq {
			r.X = ox
			r.Y += r.H + spc
		}
	}

	// 3. Hovered Color Name Tooltip
	if txt != "" {
		dtxt(fd, txt, un8th, un8th, Col.White)
	}

	//Alpha txt
	xtx := CNT.X - txW(fd, "Change alpha (fade):")
	dtxt(fd, "Change alpha (fade):", xtx, un8th, SET.UICtxt)
	keysInpU8(fd, &paletteAlpha, CNT.X+unh, un8th, unq, 1, 255, "Change transparency of color")
	// 4. Close button or picked selection closes the palette
	if closeIcon(SET.SCRW-(unh+un16th), un16th, unh) || hasPicked {
		PaletteModeAlpha = false
		paletteAlpha = 255
		ms.lc = false // 🚀 Consume click so it doesn't click items behind the palette!
	}

	return pickedColor, hasPicked
}

// MARK:UTILS
func paletteAlphaOpenBoxFOR(targetColor *sdl.Color, x, y, w float32, c sdl.Color, hovertxt string) {
	r := R(x, y, w, w)
	drr(r, c)
	drrl(r, SET.UIC3)
	if cms(r) {
		drrlw(r, 2, SET.UIC4)
		if ms.lc {
			OpenPaletteAlphaFor(targetColor)
			ms.lc = false
		}
		mInfoTXT(hovertxt)
	}
}
func paletteOpenBoxFOR(targetColor *sdl.Color, x, y, w float32, c sdl.Color, hovertxt string) {
	r := R(x, y, w, w)
	drr(r, c)
	drrl(r, SET.UIC3)
	if cms(r) {
		drrlw(r, 2, SET.UIC4)
		if ms.lc {
			OpenPaletteFor(targetColor)
			ms.lc = false
		}
		mInfoTXT(hovertxt)
	} else {
		drrlw(r, 2, SET.UIC2)
	}
}
func paletteOpenBox(x, y, w float32, c sdl.Color) {
	r := R(x, y, w, w)
	drr(r, c)
	if cms(r) {
		drrlw(r, 2, SET.UIC3)
		if ms.lc {
			PaletteMode = true
			ms.lc = false
		}
	} else {
		drrlw(r, 2, SET.UIC2)
	}
}

// Open the palette for ANY color variable in your engine:
func OpenPaletteFor(target *sdl.Color) {
	paletteTarget = target
	PaletteMode = true
}
func OpenPaletteAlphaFor(target *sdl.Color) {
	paletteTarget = target
	PaletteModeAlpha = true
}
func RandomColor() sdl.Color {
	return PaletteList[RINT(0, len(PaletteList)-1)].Color
}

// color2float: Returns sdl.Color as RGBA float32 values
func color2float(c sdl.Color) (float32, float32, float32, float32) {
	return float32(c.R) / 255.0, float32(c.G) / 255.0, float32(c.B) / 255.0, float32(c.A) / 255.0
}

func CA(c sdl.Color, a uint8) sdl.Color {
	c.A = a
	return c
}

// MARK: PALETTE
func Hex(rgb uint32) sdl.Color {
	return sdl.Color{
		R: uint8((rgb >> 16) & 0xFF),
		G: uint8((rgb >> 8) & 0xFF),
		B: uint8(rgb & 0xFF),
		A: 255,
	}
}
func HexA(rgba uint32) sdl.Color {
	return sdl.Color{
		R: uint8((rgba >> 24) & 0xFF),
		G: uint8((rgba >> 16) & 0xFF),
		B: uint8((rgba >> 8) & 0xFF),
		A: uint8(rgba & 0xFF),
	}
}

type ColorItem struct {
	Name  string
	Color sdl.Color
	Hex   uint32
}

// ==========================================
// ENDESGA 64 PALETTE https://lospec.com/palette-list/endesga-64
// ==========================================

var Col = struct {
	// Accent
	PinkHot sdl.Color

	// Monochromes & Neutral Grays
	BlackPure     sdl.Color
	BlackCharcoal sdl.Color
	GrayDarkest   sdl.Color
	GrayDark      sdl.Color
	GrayMid       sdl.Color
	GraySteel     sdl.Color
	GrayLight     sdl.Color
	White         sdl.Color

	// Slate & Cold Blues
	SlateWhite sdl.Color
	SlateLight sdl.Color
	SlateBlue  sdl.Color
	SlateDark  sdl.Color
	NavyDark   sdl.Color
	Midnight   sdl.Color
	VoidBlack  sdl.Color

	// Earth, Browns & Creams
	BrownDarkest   sdl.Color
	BrownChocolate sdl.Color
	BrownRedwood   sdl.Color
	BrownClay      sdl.Color
	BrownTan       sdl.Color
	BrownPeach     sdl.Color
	Sand           sdl.Color
	Cream          sdl.Color

	// Amber, Terracotta & Rust
	AmberLight   sdl.Color
	Terracotta   sdl.Color
	RedRust      sdl.Color
	RedDarkBlood sdl.Color

	// Fire Oranges & Yellows
	OrangeBlaze  sdl.Color
	OrangeBright sdl.Color
	GoldWarm     sdl.Color
	Gold         sdl.Color
	YellowCyber  sdl.Color

	// Greens & Teals
	GreenPale      sdl.Color
	GreenLime      sdl.Color
	GreenMeadow    sdl.Color
	GreenLeaf      sdl.Color
	GreenForest    sdl.Color
	GreenTeal      sdl.Color
	BlueDeepMarine sdl.Color

	// Ocean Blues & Cyans
	BlueNavy   sdl.Color
	BlueCobalt sdl.Color
	BlueSky    sdl.Color
	BlueCyan   sdl.Color
	CyanBright sdl.Color
	CyanIce    sdl.Color

	// Purples, Violets & Indigos
	PinkLavender   sdl.Color
	PinkBubblegum  sdl.Color
	MagentaBright  sdl.Color
	PurpleViolet   sdl.Color
	BlueElectric   sdl.Color
	BlueDeepIndigo sdl.Color
	BlueAbyss      sdl.Color

	// Plums & Orchids
	PlumDarkest   sdl.Color
	PlumDeep      sdl.Color
	PurplePlum    sdl.Color
	MagentaOrchid sdl.Color

	// Roses, Corals & Crimson Reds
	RoseMuted  sdl.Color
	PinkSalmon sdl.Color
	RedCoral   sdl.Color
	RedBright  sdl.Color
	RedCrimson sdl.Color
	RedRuby    sdl.Color
	RedMaroon  sdl.Color
}{
	// Accent
	PinkHot: Hex(0xff0040),

	// Monochromes
	BlackPure:     Hex(0x131313),
	BlackCharcoal: Hex(0x1b1b1b),
	GrayDarkest:   Hex(0x272727),
	GrayDark:      Hex(0x3d3d3d),
	GrayMid:       Hex(0x5d5d5d),
	GraySteel:     Hex(0x858585),
	GrayLight:     Hex(0xb4b4b4),
	White:         Hex(0xffffff),

	// Slate & Cold Blues
	SlateWhite: Hex(0xc7cfdd),
	SlateLight: Hex(0x92a1b9),
	SlateBlue:  Hex(0x657392),
	SlateDark:  Hex(0x424c6e),
	NavyDark:   Hex(0x2a2f4e),
	Midnight:   Hex(0x1a1932),
	VoidBlack:  Hex(0x0e071b),

	// Earth, Browns & Creams
	BrownDarkest:   Hex(0x1c121c),
	BrownChocolate: Hex(0x391f21),
	BrownRedwood:   Hex(0x5d2c28),
	BrownClay:      Hex(0x8a4836),
	BrownTan:       Hex(0xbf6f4a),
	BrownPeach:     Hex(0xe69c69),
	Sand:           Hex(0xf6ca9f),
	Cream:          Hex(0xf9e6cf),

	// Amber, Terracotta & Rust
	AmberLight:   Hex(0xedab50),
	Terracotta:   Hex(0xe07438),
	RedRust:      Hex(0xc64524),
	RedDarkBlood: Hex(0x8e251d),

	// Fire Oranges & Yellows
	OrangeBlaze:  Hex(0xff5000),
	OrangeBright: Hex(0xed7614),
	GoldWarm:     Hex(0xffa214),
	Gold:         Hex(0xffc825),
	YellowCyber:  Hex(0xffeb57),

	// Greens & Teals
	GreenPale:      Hex(0xd3fc7e),
	GreenLime:      Hex(0x99e65f),
	GreenMeadow:    Hex(0x5ac54f),
	GreenLeaf:      Hex(0x33984b),
	GreenForest:    Hex(0x1e6f50),
	GreenTeal:      Hex(0x134c4c),
	BlueDeepMarine: Hex(0x0c2e44),

	// Ocean Blues & Cyans
	BlueNavy:   Hex(0x00396d),
	BlueCobalt: Hex(0x0069aa),
	BlueSky:    Hex(0x0098dc),
	BlueCyan:   Hex(0x00cdf9),
	CyanBright: Hex(0x0cf1ff),
	CyanIce:    Hex(0x94fdff),

	// Purples, Violets & Indigos
	PinkLavender:   Hex(0xfdd2ed),
	PinkBubblegum:  Hex(0xf389f5),
	MagentaBright:  Hex(0xdb3ffd),
	PurpleViolet:   Hex(0x7a09fa),
	BlueElectric:   Hex(0x3003d9),
	BlueDeepIndigo: Hex(0x0c0293),
	BlueAbyss:      Hex(0x03193f),

	// Plums & Orchids
	PlumDarkest:   Hex(0x3b1443),
	PlumDeep:      Hex(0x622461),
	PurplePlum:    Hex(0x93388f),
	MagentaOrchid: Hex(0xca52c9),

	// Roses, Corals & Crimson Reds
	RoseMuted:  Hex(0xc85086),
	PinkSalmon: Hex(0xf68187),
	RedCoral:   Hex(0xf5555d),
	RedBright:  Hex(0xea323c),
	RedCrimson: Hex(0xc42430),
	RedRuby:    Hex(0x891e2b),
	RedMaroon:  Hex(0x571c27),
}

// PaletteList is an ordered slice of all 64 colors for UI grids and browsers
// PaletteList: Ordered into 4 harmonious rows of 16 colors (64 total)
var PaletteList = []ColorItem{
	// =========================================================================
	// ROW 1: NEUTRAL GRAYS (8) -> SLATE COLD BLUES (8)
	// =========================================================================
	{"BlackPure", Col.BlackPure, 0x131313},
	{"BlackCharcoal", Col.BlackCharcoal, 0x1b1b1b},
	{"GrayDarkest", Col.GrayDarkest, 0x272727},
	{"GrayDark", Col.GrayDark, 0x3d3d3d},
	{"GrayMid", Col.GrayMid, 0x5d5d5d},
	{"GraySteel", Col.GraySteel, 0x858585},
	{"GrayLight", Col.GrayLight, 0xb4b4b4},
	{"White", Col.White, 0xffffff},
	{"SlateWhite", Col.SlateWhite, 0xc7cfdd},
	{"SlateLight", Col.SlateLight, 0x92a1b9},
	{"SlateBlue", Col.SlateBlue, 0x657392},
	{"SlateDark", Col.SlateDark, 0x424c6e},
	{"NavyDark", Col.NavyDark, 0x2a2f4e},
	{"Midnight", Col.Midnight, 0x1a1932},
	{"VoidBlack", Col.VoidBlack, 0x0e071b},
	{"BlueAbyss", Col.BlueAbyss, 0x03193f},

	// =========================================================================
	// ROW 2: CREAMS & BROWNS (8) -> AMBERS, RUST & ORANGES (8)
	// =========================================================================
	{"Cream", Col.Cream, 0xf9e6cf},
	{"Sand", Col.Sand, 0xf6ca9f},
	{"BrownPeach", Col.BrownPeach, 0xe69c69},
	{"BrownTan", Col.BrownTan, 0xbf6f4a},
	{"BrownClay", Col.BrownClay, 0x8a4836},
	{"BrownRedwood", Col.BrownRedwood, 0x5d2c28},
	{"BrownChocolate", Col.BrownChocolate, 0x391f21},
	{"BrownDarkest", Col.BrownDarkest, 0x1c121c},
	{"AmberLight", Col.AmberLight, 0xedab50},
	{"Terracotta", Col.Terracotta, 0xe07438},
	{"RedRust", Col.RedRust, 0xc64524},
	{"RedDarkBlood", Col.RedDarkBlood, 0x8e251d},
	{"GoldWarm", Col.GoldWarm, 0xffa214},
	{"OrangeBright", Col.OrangeBright, 0xed7614},
	{"OrangeBlaze", Col.OrangeBlaze, 0xff5000},
	{"PinkHot", Col.PinkHot, 0xff0040},

	// =========================================================================
	// ROW 3: YELLOWS (2) -> GREENS & TEALS (8) -> OCEAN CYANS (6)
	// =========================================================================
	{"YellowCyber", Col.YellowCyber, 0xffeb57},
	{"Gold", Col.Gold, 0xffc825},
	{"GreenPale", Col.GreenPale, 0xd3fc7e},
	{"GreenLime", Col.GreenLime, 0x99e65f},
	{"GreenMeadow", Col.GreenMeadow, 0x5ac54f},
	{"GreenLeaf", Col.GreenLeaf, 0x33984b},
	{"GreenForest", Col.GreenForest, 0x1e6f50},
	{"GreenTeal", Col.GreenTeal, 0x134c4c},
	{"BlueDeepMarine", Col.BlueDeepMarine, 0x0c2e44},
	{"BlueNavy", Col.BlueNavy, 0x00396d},
	{"BlueCobalt", Col.BlueCobalt, 0x0069aa},
	{"BlueSky", Col.BlueSky, 0x0098dc},
	{"BlueCyan", Col.BlueCyan, 0x00cdf9},
	{"CyanBright", Col.CyanBright, 0x0cf1ff},
	{"CyanIce", Col.CyanIce, 0x94fdff},
	{"PinkLavender", Col.PinkLavender, 0xfdd2ed},

	// =========================================================================
	// ROW 4: INDIGOS & VIOLETS (4) -> MAGENTAS & PLUMS (4) -> CRIMSON REDS (8)
	// =========================================================================
	{"BlueDeepIndigo", Col.BlueDeepIndigo, 0x0c0293},
	{"BlueElectric", Col.BlueElectric, 0x3003d9},
	{"PurpleViolet", Col.PurpleViolet, 0x7a09fa},
	{"MagentaBright", Col.MagentaBright, 0xdb3ffd},
	{"PinkBubblegum", Col.PinkBubblegum, 0xf389f5},
	{"MagentaOrchid", Col.MagentaOrchid, 0xca52c9},
	{"PurplePlum", Col.PurplePlum, 0x93388f},
	{"PlumDeep", Col.PlumDeep, 0x622461},
	{"PlumDarkest", Col.PlumDarkest, 0x3b1443},
	{"RoseMuted", Col.RoseMuted, 0xc85086},
	{"PinkSalmon", Col.PinkSalmon, 0xf68187},
	{"RedCoral", Col.RedCoral, 0xf5555d},
	{"RedBright", Col.RedBright, 0xea323c},
	{"RedCrimson", Col.RedCrimson, 0xc42430},
	{"RedRuby", Col.RedRuby, 0x891e2b},
	{"RedMaroon", Col.RedMaroon, 0x571c27},
}
