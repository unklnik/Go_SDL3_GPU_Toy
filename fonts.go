package main

import (
	"fmt"

	"github.com/Zyko0/go-sdl3/ttf"
)

var (
	fx, fd, fs, fl, is, id, il, id2, il2, ixx2 *ttf.Font
)

// Global namespace for all engine fonts (Type Fon. in your editor)
var Fon struct {
	// Standard Text (Rubik-Regular)
	Sml *ttf.Font // 14px - Tooltips, badges, timestamps
	Def *ttf.Font // 18px - Standard UI buttons, menu items, labels
	Lrg *ttf.Font // 24px - Panel headers, section titles
	Xl  *ttf.Font // 36px - Banners, main window titles

	// UI Icons (RemixIcon)
	IconSml *ttf.Font // 18px - Small button icons
	IconDef *ttf.Font // 24px - Standard toolbar icons
	IconLrg *ttf.Font // 32px - Large header icons

	Icon2Def *ttf.Font
	Icon2Lrg *ttf.Font
	Icon2XXL *ttf.Font
}

// InitFonts loads all standard engine fonts at startup
func InitFonts() error {
	var err error

	// 1. Load Text Font Tiers (Rubik-Regular)
	Fon.Sml, err = ttf.OpenFont("fon/BarlowSemiCondensed-Medium.ttf", 18)
	if err != nil {
		return fmt.Errorf("failed loading Fon.Sml: %w", err)
	}

	Fon.Def, err = ttf.OpenFont("fon/BarlowSemiCondensed-Medium.ttf", 20)
	if err != nil {
		return fmt.Errorf("failed loading Fon.Def: %w", err)
	}

	Fon.Lrg, err = ttf.OpenFont("fon/BarlowSemiCondensed-Medium.ttf", 24)
	if err != nil {
		return fmt.Errorf("failed loading Fon.Lrg: %w", err)
	}

	Fon.Xl, err = ttf.OpenFont("fon/BarlowSemiCondensed-Medium.ttf", 36)
	if err != nil {
		return fmt.Errorf("failed loading Fon.Xl: %w", err)
	}

	// 2. Load Icon Font Tiers (RemixIcon) - Optional/When available
	Fon.IconSml, _ = ttf.OpenFont("fon/remixicon.ttf", 18)
	Fon.IconDef, _ = ttf.OpenFont("fon/remixicon.ttf", 24)
	Fon.IconLrg, _ = ttf.OpenFont("fon/remixicon.ttf", 32)

	Fon.Icon2Def, _ = ttf.OpenFont("fon/tabler-icons.ttf", 24)
	Fon.Icon2Lrg, _ = ttf.OpenFont("fon/tabler-icons.ttf", 32)
	Fon.Icon2XXL, _ = ttf.OpenFont("fon/tabler-icons.ttf", 64)

	fd = Fon.Def
	fs = Fon.Sml
	fl = Fon.Lrg
	fx = Fon.Xl
	is = Fon.IconSml
	id = Fon.IconDef
	il = Fon.IconLrg
	id2 = Fon.Icon2Def
	il2 = Fon.Icon2Lrg
	ixx2 = Fon.Icon2XXL

	return nil
}

// CloseFonts frees all font handles at engine exit
func CloseFonts() {
	if Fon.Sml != nil {
		Fon.Sml.Close()
	}
	if Fon.Def != nil {
		Fon.Def.Close()
	}
	if Fon.Lrg != nil {
		Fon.Lrg.Close()
	}
	if Fon.Xl != nil {
		Fon.Xl.Close()
	}

	if Fon.IconSml != nil {
		Fon.IconSml.Close()
	}
	if Fon.IconDef != nil {
		Fon.IconDef.Close()
	}
	if Fon.IconLrg != nil {
		Fon.IconLrg.Close()
	}
}
