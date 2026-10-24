package main

import (
	"testing"

	"github.com/charmbracelet/bubbles/table"
)

func TestLoadImageInfoWithoutSelectedImage(t *testing.T) {
	m := model{
		imageTable: table.New(
			table.WithColumns([]table.Column{{Title: "Images", Width: 10}}),
			table.WithRows([]table.Row{}),
		),
	}

	m.loadImageInfo()

	if m.infoBox != "No image selected" {
		t.Fatalf("infoBox = %q, want %q", m.infoBox, "No image selected")
	}
}
