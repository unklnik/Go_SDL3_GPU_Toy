package main

import (
	"encoding/json"
	"os"

	"github.com/Zyko0/go-sdl3/sdl"
)

const settingsFilePath = "settings.json"

// Settings holds all persistent engine and editor configurations
type Settings struct {
	// Window & Display
	SCRW       float32 `json:"scr_w"`
	SCRH       float32 `json:"scr_h"`
	Fullscreen bool    `json:"fullscreen"`
	VSync      bool    `json:"vsync"`
	TargetFPS  float64 `json:"target_fps"`

	// Editor Theme & UI Colors
	UIC         sdl.Color `json:"ui_color"`
	UIC2        sdl.Color `json:"ui_color2"`
	UIC3        sdl.Color `json:"ui_color3"`
	UIC4        sdl.Color `json:"ui_color4"`
	UICtxt      sdl.Color `json:"ui_text_color"`
	UICgrn      sdl.Color `json:"ui_green_color"`
	UICred      sdl.Color `json:"ui_red_color"`
	BGC         sdl.Color `json:"bg_color"`
	CursorColor sdl.Color `json:"cursor_color"`

	// Editor Preferences
	DrawHidden      bool      `json:"draw_hidden_objects"`
	TexPreviewHover bool      `json:"texture_preview_hover"`
	TexPreviewBG    bool      `json:"texture_preview_background"`
	TexPreviewBGC   sdl.Color `json:"texture_preview_background_color"`
	Rec2DLineWdef   float32   `json:"rec2d_line_width_def"`
	Rec2DLineC      sdl.Color `json:"rec2d_line_color_def"`
	Rec2DFillC      sdl.Color `json:"rec2d_fill_color_def"`
	Rec2DGradFill1C sdl.Color `json:"rec2d_gradient_fill_color1_def"`
	Rec2DGradFill2C sdl.Color `json:"rec2d_gradient_fill_color2_def"`
	Rec2DGradFill3C sdl.Color `json:"rec2d_gradient_fill_color3_def"`
	Rec2DGradFill4C sdl.Color `json:"rec2d_gradient_fill_color4_def"`
	Rec2DSelectC    sdl.Color `json:"rec2d_select_color_def"`
	Tri2DLineC      sdl.Color `json:"tri2d_line_color_def"`
	Tri2DFillC      sdl.Color `json:"tri2d_fill_color_def"`
	Tri2DSideW      float32   `json:"tri2d_side_width_def"`
	Tri2DLineW      float32   `json:"tri2d_line_width_def"`
	Tri2DSelectC    sdl.Color `json:"tri2d_select_color_def"`
	LineWdef        float32   `json:"line_width_def"`
	LineCdef        sdl.Color `json:"line_color_def"`
	GridLineCdef    sdl.Color `json:"gridline_color_def"`
	GridHoverCdef   sdl.Color `json:"grid_hover_color_def"`
	GridSnap        bool      `json:"grid_snap"`
	GridSize        float32   `json:"grid_size"`
}

// Global instance accessible from ANY file/function in your engine
var SET Settings

// DefaultSettings returns sensible defaults on first launch
func DefaultSettings() Settings {
	return Settings{
		SCRW:       1920,
		SCRH:       1080,
		Fullscreen: false,
		VSync:      true,
		TargetFPS:  60.0,

		UIC:         Hex(0x0e071b),
		UIC2:        Hex(0x03193f),
		UIC3:        Hex(0xff5000),
		UIC4:        Hex(0xffc825),
		UICtxt:      Hex(0xffffff),
		UICgrn:      Hex(0x5ac54f),
		UICred:      Hex(0xea323c),
		BGC:         Hex(0x1a1932),
		CursorColor: Hex(0xffffff),

		DrawHidden:      true,
		TexPreviewHover: true,
		TexPreviewBG:    true,
		TexPreviewBGC:   Hex(0x1b1b1b),
		Rec2DLineWdef:   1,
		Rec2DLineC:      Hex(0x0c0293),
		Rec2DFillC:      HexA(0x0c029340),
		Rec2DGradFill1C: Hex(0xff5000),
		Rec2DGradFill2C: Hex(0x0c0293),
		Rec2DGradFill3C: Hex(0xff0040),
		Rec2DGradFill4C: Hex(0x5ac54f),
		Rec2DSelectC:    Hex(0xff0040),
		Tri2DLineC:      Hex(0x99e65f),
		Tri2DFillC:      HexA(0x99e65f40),
		Tri2DSelectC:    Hex(0xff0040),
		Tri2DSideW:      un4,
		Tri2DLineW:      1,
		LineWdef:        1,
		LineCdef:        Hex(0x99e65f),
		GridLineCdef:    Hex(0x99e65f),
		GridHoverCdef:   HexA(0x99e65f40),
		GridSnap:        true,
		GridSize:        32.0,
	}
}

// LoadSettings loads settings from disk or creates default if missing
func LoadSettings() {
	// 1. Start with full default values in memory:
	SET = DefaultSettings()

	file, err := os.Open(settingsFilePath)
	if err != nil {
		// If file doesn't exist, create it with all defaults
		SaveSettings()
		return
	}
	defer file.Close()

	// 2. Decode JSON on top of defaults (new fields won't become zero/transparent!)
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&SET); err != nil {
		SET = DefaultSettings()
	}
}

// SaveSettings writes the current settings to disk as formatted JSON
func SaveSettings() {
	data, err := json.MarshalIndent(SET, "", "    ")
	if err != nil {
		return
	}
	_ = os.WriteFile(settingsFilePath, data, 0644)
}

// MARK: DRAW USER SETTINGS
func dSettings() {
	dr(0, 0, SET.SCRW, SET.SCRH, SET.BGC)
	var x, y = un8th, un8th
	var x2 float32 = 270
	_, spc := dtxt(fs, "Textures 2D", x, y, SET.UIC4) //MARK: TEXTURES 2D ████
	y += spc
	_, spc = dtxt(fs, "Preview on hover", x, y, SET.UICtxt)
	onoffcolor(&SET.TexPreviewHover, x2, y+2, un8thx3, 2, SET.UICgrn, SET.UICred, SET.UIC, true, "Draws texture preview in screen center on mouse hover")
	y += spc + uispc
	_, spc = dtxt(fs, "Preview background rectangle", x, y, SET.UICtxt)
	onoffcolor(&SET.TexPreviewBG, x2, y+2, un8thx3, 2, SET.UICgrn, SET.UICred, SET.UIC, true, "Draws texture preview with background rectangle")
	y += spc + uispc
	_, spc = dtxt(fs, "Preview background color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.TexPreviewBGC, x2, y+2, un8thx3, SET.TexPreviewBGC, "Color of texture preview background")
	y += spc + uispc
	_, spc = dtxt(fs, "Grids 2D", x, y, SET.UIC4) //MARK: GRIDS 2D ████
	y += spc
	_, spc = dtxt(fs, "Default grid line color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.GridLineCdef, x2, y+2, un8thx3, SET.GridLineCdef, "Default color of 2D grid lines")
	y += spc + uispc
	_, spc = dtxt(fs, "Default grid hover color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.GridHoverCdef, x2, y+2, un8thx3, SET.GridHoverCdef, "Default color of 2D grid lines")
	y += spc + uispc
	_, spc = dtxt(fs, "Lines 2D", x, y, SET.UIC4) //MARK: LINES 2D ████
	y += spc
	numStepperFloat(&SET.LineWdef, 0.5, 0.5, 20, x, y, 285, fonh(fs)+4, false, "Default line width", "Thickness of new 2D lines")
	y += spc + uispc
	_, spc = dtxt(fs, "Default line color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.LineCdef, x2, y+2, un8thx3, SET.LineCdef, "Default color of 2D lines")
	y += spc + uispc
	_, spc = dtxt(fs, "Rectangles 2D", x, y, SET.UIC4) //MARK: RECTANGLES 2D ████
	y += spc
	numStepperFloat(&SET.Rec2DLineWdef, 0.5, 0.5, 20, x, y, 285, fonh(fs)+4, false, "Default outline width", "Thickness of new 2D rectangle outlines")
	y += spc + uispc
	_, spc = dtxt(fs, "Default outline color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.Rec2DLineC, x2, y+2, un8thx3, SET.Rec2DLineC, "Default color of 2D rectangle outlines")
	y += spc + uispc
	_, spc = dtxt(fs, "Default fill color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.Rec2DFillC, x2, y+2, un8thx3, SET.Rec2DFillC, "Default color of 2D rectangle fill")
	y += spc + uispc
	_, spc = dtxt(fs, "Select outline color", x, y, SET.UICtxt)
	paletteAlphaOpenBoxFOR(&SET.Rec2DSelectC, x2, y+2, un8thx3, SET.Rec2DSelectC, "Default color of selected 2D rectangle outlines")
	y += spc + uispc

	//CLOSE
	settingsVis = !closeIcon(SET.SCRW-(unh+un16th), un16th, unh)
	if !settingsVis {
		uQuickIcons()
	}
}
