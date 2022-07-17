package main

import (
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TickMsg time.Time

type KeyMap struct {
	Up   key.Binding
	Quit key.Binding
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Quit}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

var DefaultKeyMap = KeyMap{
	Up:   key.NewBinding(key.WithKeys("k", "up", "space"), key.WithHelp("↑/k/space", "jump")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
}

type location struct {
	x, y int
}

type model struct {
	obstacles []location
	cursor    location
	score     int
	help      help.Model
}

var viewport = location{
	x: 80,
	y: 9,
}

func initialModel() model {
	return model{
		obstacles: []location{},
		cursor:    location{0, 0},
		score:     0,
		help:      help.New(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.Up):
		}
	case tea.WindowSizeMsg:
		// Terminal resized.
		m.help.Width = msg.Width
	case TickMsg:
		// Tick. Update the game.
	}
	return m, nil
}

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
var viewStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(2).MarginBottom(2)

func (m model) View() string {
	s := titleStyle.Render("Flapioca")
	s += "\n"

	for y := 0; y < viewport.y; y++ {
		for x := 0; x < viewport.x; x++ {
			if m.cursor.x == x && m.cursor.y == y {
				s += "•"
			} else {
				s += " "
			}
		}
		s += "\n"
	}

	helpView := m.help.View(DefaultKeyMap)
	s += helpView

	// Send the UI for rendering
	return viewStyle.Render(s)
}

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
