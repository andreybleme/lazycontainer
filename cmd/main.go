package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	container "lazycontainer/pkg/container"
	image "lazycontainer/pkg/image"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type model struct {
	containersTable table.Model
	containers      []container.Container
	imageTable      table.Model
	images          []image.Image
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			if m.containersTable.Focused() {
				m.containersTable.Blur()
				m.imageTable.Focus()
			} else if m.imageTable.Focused() {
				m.imageTable.Blur()
				m.containersTable.Focus()
			}
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			// containers table actions
			if m.containersTable.Focused() {
				index := m.containersTable.Cursor()
				containerSelected := m.containers[index]
				containerLogs, err := container.GetLogs(containerSelected.ID)

				if err != nil {
					return m, tea.Batch(
						tea.Printf("Error reading logs for container %s: %v", containerSelected.ID, err),
					)
				}

				if containerLogs == "" {
					return m, tea.Batch(
						tea.Printf("No logs available"),
					)
				}

				return m, tea.Batch(
					tea.Printf("Logs for container %s", containerLogs),
				)
			}
			// images table actions
			if m.imageTable.Focused() {
				imageDetails, err := image.GetInspect(m.imageTable.SelectedRow()[0])
				if err != nil {
					return m, tea.Batch(
						tea.Printf("Error inspecting image %s: %v", m.imageTable.SelectedRow()[0], err),
					)
				}

				return m, tea.Batch(
					tea.Printf("Image details: %s", imageDetails),
				)
			}
		}
	}

	if m.containersTable.Focused() {
		m.containersTable, cmd = m.containersTable.Update(msg)
	} else if m.imageTable.Focused() {
		m.imageTable, cmd = m.imageTable.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	return baseStyle.Render(m.containersTable.View()) + "\n" + baseStyle.Render(m.imageTable.View()) + "\n"
}

func main() {
	// containers table
	containers, err := container.ListAll()
	if err != nil {
		fmt.Println("Error listing containers:", err)
	}

	containerRows := []table.Row{}
	for _, c := range containers {
		containerRows = append(containerRows, table.Row{c.State, c.Image})
	}

	containerColumns := []table.Column{
		{Title: "Containers", Width: 10},
		{Title: "", Width: 20},
	}

	containersTable := table.New(
		table.WithColumns(containerColumns),
		table.WithRows(containerRows),
		table.WithFocused(true),
		table.WithHeight(5),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	containersTable.SetStyles(s)

	// Images table
	images, err := image.ListAll()
	if err != nil {
		fmt.Println("Error listing images:", err)
	}

	imageRows := []table.Row{}
	for _, image := range images {
		imageRows = append(imageRows, table.Row{image.Name, image.Tag})
	}

	imageColumns := []table.Column{
		{Title: "Images", Width: 10},
		{Title: "", Width: 20},
	}

	imageTable := table.New(
		table.WithColumns(imageColumns),
		table.WithRows(imageRows),
		table.WithFocused(false),
		table.WithHeight(5),
	)
	imageTable.SetStyles(s)

	m := model{containersTable, containers, imageTable, nil}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
