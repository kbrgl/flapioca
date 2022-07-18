package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func absdiff(a, b int) int {
	if a > b {
		return a - b
	}
	return b - a
}

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
	Up:   key.NewBinding(key.WithKeys("k", "up", " ", "w"), key.WithHelp("↑/k/w/space", "jump")),
	Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
}

type location struct {
	x, y int
}

type obstacle struct {
}

type model struct {
	obstacles map[*location]*obstacle
	cursor    location
	score     int
	help      help.Model
	over      bool
}

var viewport = location{
	x: 60,
	y: 9,
}

func initialModel() model {
	return model{
		obstacles: map[*location]*obstacle{},
		cursor:    location{4, viewport.y / 2},
		score:     0,
		help:      help.New(),
		over:      false,
	}
}

func (m model) tick() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return m.tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, DefaultKeyMap.Up):
			m.cursor.y--
		}
	case tea.WindowSizeMsg:
		// Terminal resized.
		m.help.Width = msg.Width
	case TickMsg:
		// Tick. Update the game.
		var rightmost location
		for loc := range m.obstacles {
			if loc.x > rightmost.x {
				rightmost = *loc
			}
			loc.x--
			if (loc.x == m.cursor.x && absdiff(loc.y, m.cursor.y) >= 2) || m.cursor.y >= viewport.y {
				m.over = true
				return m, tea.Quit
			}
			if loc.x == m.cursor.x {
				m.score++
			} else if loc.x < 0 {
				delete(m.obstacles, loc)
			}
		}
		// If the rightmost obstacle is 2 or more units away from the right edge,
		// create a new obstacle in 20% of the cases.
		if absdiff(viewport.x, rightmost.x) >= 2 && rand.Intn(100) > 80 {
			// Select a y that makes the obstacle possible to avoid.
			var y int
			for {
				y = rand.Intn(viewport.y)
				if absdiff(viewport.x, rightmost.x) > absdiff(y, rightmost.y) {
					break
				}
			}
			m.obstacles[&location{x: viewport.x, y: y}] = &obstacle{}
		}
		m.cursor.y++
		return m, m.tick()
	}
	return m, nil
}

var titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
var viewStyle = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).MarginBottom(2)
var canvasStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())

func (m model) View() string {
	s := titleStyle.Render("Flapioca")
	s += "\n"

	canvas := make([]string, 0, viewport.y)
	for y := 0; y < viewport.y; y++ {
		line := ""
	Character:
		for x := 0; x < viewport.x; x++ {
			if m.cursor.x == x && m.cursor.y == y {
				line += "•"
				continue
			}
			for loc := range m.obstacles {
				if loc.x == x && absdiff(loc.y, y) >= 2 {
					line += "#"
					continue Character
				}
			}
			line += " "
		}
		canvas = append(canvas, line)
	}

	s += canvasStyle.Render(strings.Join(canvas, "\n"))
	s += fmt.Sprintf("\n%d point(s) ", m.score)
	helpView := m.help.View(DefaultKeyMap)
	s += helpView

	if m.over {
		s += "\n\nGame over!\n"
	}

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
