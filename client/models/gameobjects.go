package models

import (
	"encoding/json"
	"leafcli/utils"
	"os"
	"path/filepath"
)

type GameObject struct {
	Name      string
	Sprite    string
	Scripts   map[string]string
	Directory string
}

func (gO GameObject) Save() error {
	if err := utils.EnsureDir(gO.Directory); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(gO.Directory, "info.lgo"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(gO)
	if err != nil {
		panic(err)
	}
	return nil
}

type GOList []GameObject

func (goList GOList) GetGOByName(name string) GameObject {
	for i, gO := range goList {
		if gO.Name == name {
			return goList[i]
		}
	}
	return GameObject{}
}

func LoadGameObjects(paths []string) (arr []GameObject) {
	for _, k := range paths {
		file, err := os.Open(k)
		if err != nil {
			break
		}

		var gO GameObject
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&gO); err != nil {
			break
		}

		// S'assurer que le champ Directory est correct
		gO.Directory = filepath.Dir(k)
		arr = append(arr, gO)
		file.Close()
	}

	return arr
}
