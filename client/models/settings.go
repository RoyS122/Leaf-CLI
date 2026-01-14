package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type EngineSettings struct {
	CodeEditor string `json:"codeeditor"`
	Language   string `json:"language"`

	// ajoute d'autres settings ici si besoin
}

func GetSettingsPath() string {
	configDir, _ := os.UserConfigDir()
	appDir := filepath.Join(configDir, "LeafCLI")
	os.MkdirAll(appDir, 0755)
	return filepath.Join(appDir, "settings.json")
}

func LoadSettings() EngineSettings {
	path := GetSettingsPath()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return EngineSettings{
			CodeEditor: "notepad.exe",
			Language:   "en",
		}
	}

	file, err := os.Open(path)
	if err != nil {
		return EngineSettings{}
	}
	defer file.Close()

	var s EngineSettings
	if err := json.NewDecoder(file).Decode(&s); err != nil {
		return EngineSettings{}
	}
	return s
}

func SaveSettings(s *EngineSettings) error {
	path := GetSettingsPath()
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(s)
}
