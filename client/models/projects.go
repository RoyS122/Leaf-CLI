package models

import (
	"encoding/json"
	"fmt"
	"leafcli/utils"
	"os"
	"path/filepath"
)

type Project struct {
	Name        string
	Description string
	Directory   string
	Version     Version
	Rooms       []Room
}

func (p Project) CreateProjectDirectory(basePath string) error {
	projectPath := filepath.Join(basePath, p.Name)

	if err := utils.EnsureDir(projectPath); err != nil {
		return err
	}
	p.Directory = projectPath
	file, err := os.Create(filepath.Join(projectPath, "projectdata.ldat"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(p)
	if err != nil {
		panic(err)
	}

	fmt.Println("Project directory ready at:", projectPath)
	return nil
}

func LoadProjectFromFile(filePath string) (*Project, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var p Project
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&p); err != nil {
		return nil, err
	}

	// S'assurer que le champ Directory est correct
	p.Directory = filepath.Dir(filePath)

	return &p, nil
}
