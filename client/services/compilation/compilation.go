package compilation

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func BuildEngine(engineDir, outputDir string) {

	fmt.Println("🔨 Building engine...")

	// 1. Exécuter make dans engine
	fmt.Println("test make in :", engineDir)

	cmd := exec.Command("make", "-f", "Makefile.windows")
	cmd.Dir = engineDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("⚠️  'make build' failed, trying 'make'...")
		cmd = exec.Command("make", "-f", "Makefile.linux")
		cmd.Dir = engineDir
		if err := cmd.Run(); err != nil {
			log.Fatal("❌ Failed to build engine: ", err)
		}
	}

	// 2. Chercher l'output
	possibleDirs := []string{"output", "bin", "build", "dist"}
	var sourceDir string

	for _, dir := range possibleDirs {
		checkPath := filepath.Join(engineDir, dir)
		if info, err := os.Stat(checkPath); err == nil && info.IsDir() {
			sourceDir = checkPath
			break
		}
	}

	if sourceDir == "" {
		log.Fatal("❌ No output directory found")
	}

	// 3. Nettoyer et créer le dossier de destination
	os.RemoveAll(outputDir)
	os.MkdirAll(outputDir, 0755)

	// 4. Copier les fichiers
	fmt.Println("📋 Copying files...")

	entries, err := os.ReadDir(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		src := filepath.Join(sourceDir, entry.Name())
		dst := filepath.Join(outputDir, entry.Name())

		if entry.IsDir() {
			// Copier récursivement
			cmd = exec.Command("cp", "-r", src, dst)
		} else {
			// Copier fichier
			data, err := os.ReadFile(src)
			if err != nil {
				fmt.Printf("⚠️  Failed to copy %s: %v\n", entry.Name(), err)
				continue
			}
			os.WriteFile(dst, data, 0644)
		}
	}

	fmt.Printf("✅ Engine deployed to: %s\n", outputDir)

	// Lister les fichiers copiés
	entries, _ = os.ReadDir(outputDir)
	fmt.Println("📁 Contents:")
	for _, entry := range entries {
		fmt.Printf("  - %s\n", entry.Name())
	}
}
