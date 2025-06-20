package container

import (
	"os/exec"
	"strings"
)

type Container struct {
	ID    string
	Name  string
	State string
	Image string
}

func ListAll() ([]Container, error) {
	output, err := exec.Command("container", "list", "--all").Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var containers []Container

	// skip header
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		container := Container{
			ID:    fields[0],
			Name:  fields[1],
			State: fields[2],
			Image: fields[3],
		}
		containers = append(containers, container)
	}

	return containers, nil
}
