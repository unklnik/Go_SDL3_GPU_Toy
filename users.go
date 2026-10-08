package main

var (
	users []USER
	usr   *USER
)

type USER struct {
	o3d, o3dSave, o3dTrash       []OBJ3D
	tex2d, tex2dSave, tex2dTrash []TEX2D
	lights3d                     []Light3D
	grids2d                      []GRID2D
	rec2d                        []REC2D
	line2d                       []LINE2D
}

func mUsers() {
	u := USER{}
	users = append(users, u)
	usr = &users[0]
}
