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

func GetLogs(id string) (string, error) {
	output, err := exec.Command("container", "logs", id).Output()
	if err != nil {
		return "Error reading container logs", err
	}

	return string(output), nil
}
