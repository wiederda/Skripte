package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"gopkg.in/yaml.v2"
)

// ContainerInfo repräsentiert die Struktur der Docker-Containerinformationen
type ContainerInfo struct {
	Name       string          `json:"Name"`
	Config     ContainerConfig `json:"Config"`
	HostConfig HostConfig      `json:"HostConfig"`
	Mounts     []Mount         `json:"Mounts"`
}

// ContainerConfig enthält Konfigurationsinformationen für den Container
type ContainerConfig struct {
	Image  string            `json:"Image"`
	Labels map[string]string `json:"Labels"`
	Env    []string          `json:"Env"`
}

// HostConfig enthält Host-spezifische Konfigurationen
type HostConfig struct {
	PortBindings map[string][]PortBinding `json:"PortBindings"`
}

// PortBinding beschreibt die Portbindung zwischen Host und Container
type PortBinding struct {
	HostPort string `json:"HostPort"`
}

// Mount beschreibt die Volumes, die im Container verwendet werden
type Mount struct {
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
}

// ComposeFile repräsentiert die Struktur für die Docker Compose-Datei
type ComposeFile struct {
	Version  string             `yaml:"version"`
	Services map[string]Service `yaml:"services"`
}

// Service beschreibt einen einzelnen Service in der Docker Compose-Datei
type Service struct {
	Image         string            `yaml:"image"`
	ContainerName string            `yaml:"container_name"`
	Ports         []string          `yaml:"ports"`
	Environment   []string          `yaml:"environment"`
	Volumes       []string          `yaml:"volumes"`
	Labels        map[string]string `yaml:"labels,omitempty"`
}

// convertToCompose konvertiert ContainerInfo in ComposeFile-Format
func convertToCompose(container ContainerInfo) ComposeFile {
	// Container-Name ohne den führenden Slash
	containerName := strings.TrimPrefix(container.Name, "/")

	// Port-Mappings
	var portMappings []string
	for containerPort, bindings := range container.HostConfig.PortBindings {
		for _, binding := range bindings {
			portMappings = append(portMappings, fmt.Sprintf("%s:%s", binding.HostPort, containerPort))
		}
	}

	// Volume-Mappings
	var volumeMappings []string
	for _, mount := range container.Mounts {
		volumeMappings = append(volumeMappings, fmt.Sprintf("%s:%s", mount.Source, mount.Destination))
	}

	// Docker Compose-Struktur aufbauen
	return ComposeFile{
		Version: "3",
		Services: map[string]Service{
			containerName: {
				Image:         container.Config.Image,
				ContainerName: containerName,
				Ports:         portMappings,
				Environment:   container.Config.Env,
				Volumes:       volumeMappings,
				Labels:        container.Config.Labels,
			},
		},
	}
}

func main() {
	// Kommandozeilenargumente definieren
	containerName := flag.String("container", "", "Name des Docker-Containers (optional mit führendem /)")
	outputFile := flag.String("output", "docker-compose.yml", "Pfad zur Ausgabedatei")

	// Hilfetext zur Verwendung von Flaggen
	executableName := "docker-compose-converter" // Default name for Linux and macOS
	if runtime.GOOS == "windows" {
		executableName += ".exe" // Append .exe for Windows
	}

	flag.Usage = func() {
		fmt.Printf("Verwendung: %s -container container_name -output output_file\n", executableName)
		fmt.Println("Flags:")
		flag.PrintDefaults()
	}

	// Argumente parsen
	flag.Parse()

	// Wenn kein Container-Name angegeben ist, zeige Fehlermeldung an
	if *containerName == "" {
		fmt.Println("Fehler: Container-Name muss angegeben werden.")
		flag.Usage()
		os.Exit(1)
	}

	// Sicherstellen, dass der Container-Name mit "/" beginnt
	if (*containerName)[0] != '/' {
		*containerName = "/" + *containerName
	}

	// Prüfen, ob Docker installiert ist
	if _, err := exec.LookPath("docker"); err != nil {
		fmt.Println("Fehler: Docker ist nicht installiert oder nicht im PATH verfügbar.")
		os.Exit(1)
	}

	// Docker-Inspect-Befehl ausführen
	jsonData, err := exec.Command("docker", "inspect", *containerName).Output()
	if err != nil {
		fmt.Printf("Fehler beim Ausführen des Docker-Befehls: %v\n", err)
		os.Exit(1)
	}

	// JSON-Daten in ContainerInfo-Struktur umwandeln
	var containerInfos []ContainerInfo
	if err := json.Unmarshal(jsonData, &containerInfos); err != nil {
		fmt.Printf("Fehler beim Parsen der JSON-Daten: %v\n", err)
		os.Exit(1)
	}

	// Nur den ersten Container nehmen
	container := containerInfos[0]

	// In Docker Compose umwandeln
	compose := convertToCompose(container)

	// YAML-Output generieren
	output, err := yaml.Marshal(compose)
	if err != nil {
		fmt.Printf("Fehler beim Marshaling der YAML-Daten: %v\n", err)
		os.Exit(1)
	}

	// In die Ausgabedatei schreiben
	if err := ioutil.WriteFile(*outputFile, output, 0644); err != nil {
		fmt.Printf("Fehler beim Schreiben in die Ausgabedatei: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Docker Compose-Datei erfolgreich in %s geschrieben.\n", *outputFile)
}
