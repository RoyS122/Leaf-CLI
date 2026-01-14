package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Instance struct {
	Parent string // Parent path
	X, Y   int
}

func NewInstanceFromGO(gO GameObject, x, y int) Instance {
	return Instance{
		Parent: filepath.Join(gO.Directory, "info.lgo"),
		X:      x,
		Y:      y,
	}
}

func (i Instance) GetGameObject() GameObject {
	file, err := os.Open(i.Parent)
	if err != nil {
		fmt.Println("Error opening GameObject file for instance:", err)
		return GameObject{}
	}

	var gO GameObject
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&gO); err != nil {
		return GameObject{}
	}

	// S'assurer que le champ Directory est correct
	gO.Directory = filepath.Dir(i.Parent)
	fmt.Println("Loaded GameObject for instance:", gO.Name)
	file.Close()
	return gO
}
