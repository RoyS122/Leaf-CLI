package models

import (
	"encoding/json"
	"leafcli/utils"
	"os"
	"path/filepath"
)

type Room struct {
	Name, Directory string
	GOInstances     []Instance
}

func (r Room) Save() error {
	if err := utils.EnsureDir(r.Directory); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(r.Directory, r.Name+".lsp"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(r)
	if err != nil {
		panic(err)
	}
	return nil
}

func LoadRooms(paths []string) (arr []Room) {
	for _, k := range paths {
		file, err := os.Open(k)
		if err != nil {
			break
		}

		var r Room
		decoder := json.NewDecoder(file)
		if err := decoder.Decode(&r); err != nil {
			break
		}

		// S'assurer que le champ Directory est correct
		r.Directory = filepath.Dir(k)
		arr = append(arr, r)
		file.Close()
	}

	return arr
}
