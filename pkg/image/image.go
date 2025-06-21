package image

import (
	"os/exec"
	"strings"
)

type Image struct {
	Name   string
	Tag    string
	Digest string
}

func ListAll() ([]Image, error) {
	var images []Image

	output, err := exec.Command("container", "image", "list").Output()
	if err != nil {
		return nil, err
	}

	// NAME, TAG, DIGEST
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines[:1] {
		fields := strings.Fields(line)
		// skip malformed lines
		if len(fields) < 3 {
			continue
		}

		image := Image{
			Name:   fields[0],
			Tag:    fields[1],
			Digest: fields[2],
		}

		images = append(images, image)
	}

	return images, nil
}
