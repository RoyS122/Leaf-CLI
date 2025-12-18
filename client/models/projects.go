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

	// Indexes
	Rooms       []string
	GameObjects []string
	Sprites     []string
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

func (p *Project) UpdateProject() error {
	if p.Directory == "" {
		return fmt.Errorf("project directory is not set")
	}

	projectFilePath := filepath.Join(p.Directory, "projectdata.ldat")

	file, err := os.Create(projectFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // JSON lisible (optionnel)

	if err := encoder.Encode(p); err != nil {
		return err
	}

	fmt.Println("Project updated:", projectFilePath)
	return nil
}
