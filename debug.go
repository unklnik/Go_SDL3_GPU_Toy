package main

import "fmt"

var (
	debugstring string
	debugnum    int
	DebugMode   bool
)

func dDebug() {
	wid := f32(320)
	r := R(SET.SCRW-wid, 0, wid, SET.SCRH)
	drr(r, CA(Col.RedRuby, 200))
	x, y, spc := r.X+un8th, r.Y+un8th, fonh(Fon.Sml)
	dtxtsmlwhite("debugnum "+fmt.Sprint(debugnum)+"  |  debugstring "+fmt.Sprint(debugstring), x, y)
	y += spc
	dtxtsmlwhite("SCRW "+fmt.Sprint(SET.SCRW)+"  |  SCRH "+fmt.Sprint(SET.SCRH)+"  |  msinUI "+fmt.Sprint(msinUI), x, y)
	y += spc
	dtxtsmlwhite("recNewActiv "+fmt.Sprint(recNewActiv)+"len(recVal.p)"+fmt.Sprint(len(recVal.p)), x, y)
	y += spc
	dtxtsmlwhite("recNewActiv "+fmt.Sprint(recNewActiv)+"len(recVal.p)"+fmt.Sprint(len(recVal.p)), x, y)
	y += spc
}
