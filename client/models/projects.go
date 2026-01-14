package models

import (
	"encoding/json"
	"fmt"
	"leafcli/services/compilation"
	"leafcli/utils"
	"os"
	"path/filepath"
	"strings"
)

const ENGINE_SOURCE_DIR = "C:/Users/morga/git/Leaf-CLI/vendor/engine"

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

// func (p Project) CompileGameObject(gI Instance) string {
// 	var builder strings.Builder

// 	builder.WriteString(fmt.Sprintf("-- Auto-generated GameObject: %s\n\n", gO.Name))

// 	// Création directe avec l'API Engine
// 	builder.WriteString(fmt.Sprintf("local %s = Engine.create_gameobject(%s, 0, 0)\n",
// 		strings.ToLower(gO.Name), strings.ToLower()))

// 	// Configuration du sprite
// 	if gO.Sprite != "" {
// 		builder.WriteString(fmt.Sprintf("%s:setSprite(\"%s\", 64, 64, 4, 4, 10)\n",
// 			strings.ToLower(gO.Name), gO.Sprite))
// 	}

// 	// Ajout des scripts
// 	// for scriptType, scriptContent := range gO.Scripts {
// 	// 	if scriptContent != "" {
// 	// 		builder.WriteString(fmt.Sprintf("\n-- %s script\n", scriptType))

// 	// 	}
// 	// }

// 	builder.WriteString(fmt.Sprintf("\nreturn %s\n", strings.ToLower(gO.Name)))

// 	return builder.String()
// }

func (p Project) Compile() error {
	fmt.Println(p)
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\n-- Auto-generated Project: %s\n\n", p.Name))
	// builder.WriteString(fmt.Sprintf("local Engine = require(\"engine\")\n"))
	builder.WriteString(fmt.Sprintf("-- Auto-generated Rooms: %d\n\n", len(p.Rooms)))

	for _, r := range p.Rooms {
		builder.WriteString(p.CompileRoom(LoadRoom(r)))
	}

	buildPath := filepath.Join(p.Directory, "build")
	compilation.BuildEngine(ENGINE_SOURCE_DIR, buildPath)
	if err := os.MkdirAll(buildPath, 0755); err != nil {
		return fmt.Errorf("impossible de créer le dossier build: %w", err)
	}

	filePath := filepath.Join(buildPath, "game.lua")
	err := os.WriteFile(filePath, []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("erreur lors de l'écriture de game.lua: %w", err)
	}

	fmt.Printf("Compilation réussie : %s\n", filePath)
	return nil
}

func (p Project) CompileRoom(r Room) string {
	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("\n-- Auto-generated Room: %s\n\n", r.Name))
	builder.WriteString(fmt.Sprintf("local %s = Engine.create_room()\n", r.Name))
	builder.WriteString(fmt.Sprintf("\n-- Auto-generated Instances: %d\n\n", len(r.GOInstances)))
	for ins_id, GOi := range r.GOInstances {
		builder.WriteString(p.CompileInstance(r, GOi, uint(ins_id)))
	}
	return builder.String()
}

func (p Project) CompileInstance(r Room, i Instance, insID uint) string {
	var builder strings.Builder
	Gobj := LoadGameObject(i.Parent)
	var gobj_id string = fmt.Sprintf("ins_%s_%d", Gobj.Name, insID)
	builder.WriteString(fmt.Sprintf("local %s = Engine.create_gameobject(%s, %d, %d)\n", gobj_id, r.Name, i.X, i.Y))
	if Gobj.Sprite != "" {
		spr, err := LoadSpriteFromFile(Gobj.Sprite)
		if err != nil {
			fmt.Println(err)
			return builder.String()
		}
		safeSpritePath := strings.ReplaceAll(spr.ImagePath, "\\", "/") // On remplace \ par /

		if spr.Rows == 0 {
			spr.Rows = 1
		}

		if spr.Columns == 0 {
			spr.Columns = 1
		}

		builder.WriteString(fmt.Sprintf(
			"%s:setSprite(\"%s\", %d, %d, %d, %d, %d)\n",
			gobj_id,
			safeSpritePath,
			spr.ImageWidth,
			spr.ImageHeight,
			spr.Columns,
			spr.Rows,
			spr.AnimSpeed))
	}
	return builder.String()
}
