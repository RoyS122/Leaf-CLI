package models

import (
	"encoding/json"
	"leafcli/utils"
	"os"
	"path/filepath"
)

type Sprite struct {
	Name, ImagePath, Directory             string
	Columns, Rows, ImageWidth, ImageHeight uint
}

func (s Sprite) Save() error {
	if err := utils.EnsureDir(s.Directory); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(s.Directory, s.Name+".lsp"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(s)
	if err != nil {
		panic(err)
	}
	return nil
}

func LoadSprites(paths []string) (arr []Sprite) {
	for _, k := range paths {
		file, err := os.Open(k)
		if err != nil {
			break
		}

		var s Sprite
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&s); err != nil {
			break
		}

		// S'assurer que le champ Directory est correct
		s.Directory = filepath.Dir(k)
		arr = append(arr, s)
		file.Close()
	}

	return arr
}
