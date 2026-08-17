package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
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

	imageTable table.Model
	images     []image.Image

	infoBox string

	containerTab string
	logsViewport viewport.Model
}

// container tab menu constants
const (
	containerInfoTab = "Info"
	containerLogsTab = "Logs"
	infoBoxWidth     = 60
	infoBoxRows      = 10
)

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
		case "left", "right":
			// container tab menu
			if m.containersTable.Focused() {
				if m.containerTab == containerInfoTab {
					m.containerTab = containerLogsTab
				} else {
					m.containerTab = containerInfoTab
				}
				m.loadContainerTab()
			}
		case "enter":
			// containers table actions
			if m.containersTable.Focused() {
				m.loadContainerTab()
			}

			// images table actions
			if m.imageTable.Focused() {
				m.loadImageInfo()
			}
		}
	}

	if m.containersTable.Focused() && m.containerTab == containerLogsTab {
		m.logsViewport, cmd = m.logsViewport.Update(msg)
	} else if m.containersTable.Focused() {
		m.containersTable, cmd = m.containersTable.Update(msg)
	} else if m.imageTable.Focused() {
		m.imageTable, cmd = m.imageTable.Update(msg)
	}
	return m, cmd
}

func (m *model) loadContainerTab() {
	index := m.containersTable.Cursor()
	if index < 0 || index >= len(m.containers) {
		m.infoBox = "No container selected"
		return
	}

	containerSelected := m.containers[index]

	if m.containerTab == containerLogsTab {
		logs, err := container.GetLogs(containerSelected.ID)
		if err != nil {
			m.infoBox = fmt.Sprintf("Error reading logs for container %s: %v", containerSelected.ID, err)
		} else {
			m.infoBox = logs
		}
		m.logsViewport.SetContent(m.infoBox)
		m.logsViewport.GotoTop()
		return
	}

	containerDetails, err := container.GetDetails(containerSelected.ID)
	if err != nil {
		m.infoBox = fmt.Sprintf("Error inspecting container %s: %v", containerSelected.ID, err)
		return
	}

	m.infoBox = fmt.Sprintf("ID: %s \nImage: %s \nCPU: %d \nMemory: %d \nNetworks: %s \nEnvironment: %s", containerDetails.ID, containerDetails.Image, containerDetails.CPU, containerDetails.Memory,
		lipgloss.JoinVertical(lipgloss.Left, containerDetails.Networks...),
		lipgloss.JoinVertical(lipgloss.Left, containerDetails.Environment...),
	)
}

func (m *model) loadImageInfo() {
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

func (m model) View() string {
	tables := lipgloss.JoinVertical(lipgloss.Left,
		baseStyle.Render(m.containersTable.View()),
		baseStyle.Render(m.imageTable.View()),
	)

	infoBoxStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Width(infoBoxWidth).
		Height(14).
		Padding(0, 2)

	infoBox := m.infoBox
	if m.containersTable.Focused() {
		content := lipgloss.NewStyle().MaxHeight(20).Render(m.infoBox)
		if m.containerTab == containerLogsTab {
			content = m.logsViewport.View()
		}
		infoBox = lipgloss.JoinVertical(lipgloss.Left,
			m.containerTabs(),
			content,
		)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		tables,
		infoBoxStyle.Render(infoBox),
	)
}

func (m model) containerTabs() string {
	activeTab := lipgloss.NewStyle().
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(true).
		Padding(0, 1)
	inactiveTab := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Padding(0, 1)

	infoHeader := inactiveTab.Render(containerInfoTab)
	logsHeader := inactiveTab.Render(containerLogsTab)
	if m.containerTab == containerLogsTab {
		logsHeader = activeTab.Render(containerLogsTab)
	} else {
		infoHeader = activeTab.Render(containerInfoTab)
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, infoHeader, logsHeader)
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
		{Title: "", Width: 25},
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
		{Title: "", Width: 25},
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

	m := model{
		containersTable: containersTable,
		containers:      containers,
		imageTable:      imageTable,
		images:          images,
		containerTab:    containerInfoTab,
		logsViewport:    viewport.New(infoBoxWidth, infoBoxRows),
	}
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
