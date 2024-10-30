// build/build.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// BuildConfig definiert die Build-Konfiguration für jede Plattform
type BuildConfig struct {
	OS      string
	Arch    string
	Output  string
	DirName string
}

func main() {
	// Name der Ausgabedatei ohne Erweiterung
	outputName := "docker-compose-converter"

	// Erstelle Build-Konfigurationen für jede Plattform
	buildConfigs := []BuildConfig{
		{"windows", "amd64", outputName + ".exe", "windows"},
		{"linux", "amd64", outputName, "linux"},
		{"darwin", "amd64", outputName, "macos"},
	}

	// Baue das Programm für jede Konfiguration
	for _, config := range buildConfigs {
		if err := build(config); err != nil {
			fmt.Printf("%s build failed: %v\n", config.OS, err)
			os.Exit(1)
		} else {
			fmt.Printf("%s build completed successfully!\n", config.OS)
		}
	}

	fmt.Println("All builds completed successfully!")
}

// Funktion zum Bauen des Programms für eine bestimmte Konfiguration
func build(config BuildConfig) error {
	// Erstelle das Ausgabeverzeichnis
	if err := os.MkdirAll(config.DirName, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", config.DirName, err)
	}

	// Definiere den Build-Befehl mit OS, Arch und Output-Optionen
	// Hier geben wir den Pfad zur main.go-Datei an
	cmd := exec.Command("go", "build", "-o", filepath.Join(config.DirName, config.Output), "../main.go")
	cmd.Env = append(os.Environ(), "GOOS="+config.OS, "GOARCH="+config.Arch)

	// Führe den Befehl aus und gib ggf. Fehler zurück
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed for %s/%s: %v", config.OS, config.Arch, err)
	}

	return nil
}
