package models

import (
	"encoding/json"
	"fmt"
	"io"
	"leafcli/utils"
	"os"
	"path/filepath"
)

type Sprite struct {
	Name, ImagePath, Directory                        string
	Columns, Rows, ImageWidth, ImageHeight, AnimSpeed uint
}

func (s Sprite) Save() error {
	fmt.Println("save called for sprite: ", s)
	if err := utils.EnsureDir(s.Directory); err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(s.Directory, "info.lsp"))
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

func (s Sprite) GetCompiledPath(buildPath string) string {
	return filepath.Join(buildPath, "Sprites", s.Name, filepath.Base(s.ImagePath))
}

func CopySpritesAssets(paths []string, destPath string) error {
	sprites := LoadSprites(paths)
	for _, spr := range sprites {
		var relativeSpritePath string = filepath.Join(destPath, spr.Name)
		var relativeSpriteFilePath string = filepath.Join(relativeSpritePath, filepath.Base(spr.ImagePath))
		if err := utils.EnsureDir(relativeSpritePath); err != nil {
			return err
		}
		// var fileIn, fileOut *os.File
		fileIn, err := os.Open(spr.ImagePath)
		if err != nil {
			return err
		}
		defer fileIn.Close()

		fileOut, err := os.Create(relativeSpriteFilePath)
		if err != nil {
			return err
		}
		defer fileOut.Close()

		if _, err := io.Copy(fileOut, fileIn); err != nil {
			return err
		}
	}
	return nil
}

func LoadSpriteFromFile(filePath string) (Sprite, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening sprite file: ", filePath)
		return Sprite{}, err
	}
	defer file.Close()

	var s Sprite
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&s); err != nil {
		return Sprite{}, err
	}

	// S'assurer que le champ Directory est correct
	s.Directory = filepath.Dir(filePath)
	return s, nil
}
