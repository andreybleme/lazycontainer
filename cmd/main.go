package main

import (
	"fmt"
	"os"
	"time"

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
	infoBox         string
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
				containerDetails, err := container.GetDetails(containerSelected.ID)
				if err != nil {
					m.infoBox = fmt.Sprintf("Error inspecting container %s: %v", containerSelected.ID, err)
				} else {
					m.infoBox = fmt.Sprintf("ID: %s \nImage: %s \nCPU: %d \nMemory: %d \nNetworks: %s \nEnvironment: %s", containerDetails.ID, containerDetails.Image, containerDetails.CPU, containerDetails.Memory,
						lipgloss.JoinVertical(lipgloss.Left, containerDetails.Networks...),
						lipgloss.JoinVertical(lipgloss.Left, containerDetails.Environment...),
					)
				}
			}

			// images table actions
			if m.imageTable.Focused() {
				imageDetails, err := image.GetDetails(m.imageTable.SelectedRow()[0])
				if err != nil {
					m.infoBox = fmt.Sprintf("Error inspecting image %s: %v", m.imageTable.SelectedRow()[0], err)
				} else {
					createdDataTime, _ := time.Parse(time.RFC3339, imageDetails.Created)
					// adjust to readable local date time (2025-05-29T16:02:07Z)
					localTime := createdDataTime.Local()
					formattedDateTime := localTime.Format("Mon, 02 Jan 2006 15:04:05 -07")
					// convert bytes to megabytes
					sizeMB := float64(imageDetails.Size) / (1024 * 1024)
					m.infoBox = fmt.Sprintf("Name: %s \nID: %s \nSize: %.2fMB \nCreated: %s", imageDetails.Name, imageDetails.Id, sizeMB, formattedDateTime)
				}
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
	tables := lipgloss.JoinVertical(lipgloss.Left,
		baseStyle.Render(m.containersTable.View()),
		baseStyle.Render(m.imageTable.View()),
	)

	infoBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(60).
		Height(14).
		Padding(1, 2)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		tables,
		infoBoxStyle.Render(m.infoBox),
	)
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

	styleContainers := table.DefaultStyles()
	styleContainers.Header = styleContainers.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	styleContainers.Selected = styleContainers.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	containersTable.SetStyles(styleContainers)

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

	styleImages := table.DefaultStyles()
	styleImages.Header = styleImages.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(true)
	styleImages.Selected = styleImages.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("201")).
		Bold(false)
	imageTable.SetStyles(styleImages)

	m := model{containersTable, containers, imageTable, images, ""}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
