package image

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

type Image struct {
	Name   string
	Tag    string
	Digest string
}

type ImageDetails struct {
	Name    string
	Id      string
	Size    int64
	Created string
}

type imageInspectRaw struct {
	Name  string `json:"name"`
	Index struct {
		Digest string `json:"digest"`
		Size   int64  `json:"size"`
	} `json:"index"`
	Variants []struct {
		Config struct {
			Created string `json:"created"`
		} `json:"config"`
	} `json:"variants"`
}

func ListAll() ([]Image, error) {
	var images []Image

	output, err := exec.Command("container", "images", "list").Output()
	if err != nil {
		return nil, err
	}

	// NAME, TAG, DIGEST
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines[1:] {
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

func GetDetails(name string) (ImageDetails, error) {
	rawJSON, err := inspect(name)
	if err != nil {
		return ImageDetails{}, err
	}

	var entries []imageInspectRaw
	if err := json.Unmarshal([]byte(rawJSON), &entries); err != nil {
		return ImageDetails{}, fmt.Errorf("failed to parse image JSON: %w", err)
	}
	if len(entries) == 0 {
		return ImageDetails{}, fmt.Errorf("no image entries found in inspect output")
	}

	e := entries[0]
	created := ""
	if len(e.Variants) > 0 {
		created = e.Variants[0].Config.Created
	}

	imageDetails := ImageDetails{
		Name:    e.Name,
		Id:      e.Index.Digest,
		Size:    e.Index.Size,
		Created: created,
	}
	return imageDetails, nil
}

func inspect(name string) (string, error) {
	output, err := exec.Command("container", "images", "inspect", name).Output()
	if err != nil {
		return "Error inspecting image", err
	}

	return string(output), nil
}
