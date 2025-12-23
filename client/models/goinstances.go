package models

import "path/filepath"

type Instance struct {
	Parent string // Parent path
	X, Y   int
}

func NewInstanceFromGO(gO GameObject, x, y int) Instance {
	return Instance{
		Parent: filepath.Join(gO.Directory, gO.Name+".lgo"),
		X:      x,
		Y:      y,
	}
}
