package utils

import (
	"fmt"
	"os"
)

func EnsureDir(path string) error {
	info, err := os.Stat(path)

	if err == nil {
		// Le chemin existe
		if info.IsDir() {
			return nil // dossier OK
		}
		return fmt.Errorf("path exists but is not a directory")
	}

	if os.IsNotExist(err) {
		// Le dossier n'existe pas → on le crée
		return os.MkdirAll(path, 0755)
	}

	// Autre erreur (permissions, etc.)
	return err
}
