package container

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Container struct {
	ID    string
	Name  string
	State string
	Image string
}

type ContainerDetails struct {
	ID          string
	Image       string
	CPU         int64
	Memory      int64
	Networks    []string
	Environment []string
}

type NetworkInfo struct {
	Address  string `json:"address"`
	Network  string `json:"network"`
	Hostname string `json:"hostname"`
	Gateway  string `json:"gateway"`
}

type NetworksField []string
type containerInspectRaw struct {
	Configuration struct {
		ID    string `json:"id"`
		Image struct {
			Descriptor struct {
				Digest string `json:"digest"`
			} `json:"descriptor"`
			Reference string `json:"reference"`
		} `json:"image"`
		Resources struct {
			CPUs          int64 `json:"cpus"`
			MemoryInBytes int64 `json:"memoryInBytes"`
		} `json:"resources"`
		InitProcess struct {
			Environment []string `json:"environment"`
		} `json:"initProcess"`
		Networks NetworksField `json:"networks"`
	} `json:"configuration"`
	Networks NetworksField `json:"networks"`
}

func ListAll() ([]Container, error) {
	var containers []Container

	output, err := exec.Command("container", "list", "--all").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		// skip malformed lines
		if len(fields) < 4 {
			continue
		}

		// ID, IMAGE, OS, ARCH, STATE, ADDR
		container := Container{
			ID:    fields[0],
			Name:  "",
			State: fields[4],
			Image: fields[1],
		}
		containers = append(containers, container)
	}

	return containers, nil
}

func GetDetails(id string) (ContainerDetails, error) {
	rawJSON, err := inspect(id)
	if err != nil {
		return ContainerDetails{}, err
	}

	var entries []containerInspectRaw
	if err := json.Unmarshal([]byte(rawJSON), &entries); err != nil {
		return ContainerDetails{}, fmt.Errorf("failed to parse container JSON: %w", err)
	}
	if len(entries) == 0 {
		return ContainerDetails{}, fmt.Errorf("no container entries found in inspect output")
	}

	e := entries[0]

	// show base configuration.networks if not set in base struct
	networks := e.Networks
	if len(e.Networks) == 0 {
		networks = e.Configuration.Networks
	}

	containerDetails := ContainerDetails{
		ID:          e.Configuration.ID,
		Image:       e.Configuration.Image.Reference,
		CPU:         e.Configuration.Resources.CPUs,
		Memory:      e.Configuration.Resources.MemoryInBytes,
		Networks:    networks,
		Environment: e.Configuration.InitProcess.Environment,
	}
	return containerDetails, nil
}

func GetLogs(id string) (string, error) {
	output, err := exec.Command("container", "logs", id).Output()
	if err != nil {
		return "Error reading container logs", err
	}

	return string(output), nil
}

func inspect(id string) (string, error) {
	output, err := exec.Command("container", "inspect", id).Output()
	if err != nil {
		return "Error inspecting container", err
	}

	return string(output), nil
}

func (nf *NetworksField) UnmarshalJSON(data []byte) error {
	var strArr []string
	if err := json.Unmarshal(data, &strArr); err == nil {
		*nf = strArr
		return nil
	}

	var objArr []NetworkInfo
	if err := json.Unmarshal(data, &objArr); err == nil {
		var result []string
		for _, n := range objArr {
			result = append(result, n.Address)
		}
		*nf = result
		return nil
	}

	return fmt.Errorf("networks field is neither []string nor []NetworkInfo")
}
