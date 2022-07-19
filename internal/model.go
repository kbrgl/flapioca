package internal

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type TickMsg time.Time

type Model struct {
	Keys      KeyMap
	Obstacles Obstacles
	Cursor    Location
	Score     int
	Help      help.Model
	Viewport  Location
	Over      bool
	Pressed   bool
}

func NewModel() Model {
	viewport := Location{
		X: 60,
		Y: 9,
	}
	return Model{
		Keys: KeyMap{
			Up:   key.NewBinding(key.WithKeys("k", "up", " ", "w"), key.WithHelp("↑/k/w/space", "jump")),
			Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		},
		Obstacles: NewObstacles(NewObstacle(2, &Location{viewport.X / 2, viewport.Y / 2})),
		Cursor:    Location{4, viewport.Y / 2},
		Score:     0,
		Help:      help.New(),
		Over:      false,
		Viewport:  viewport,
	}
}

func (m Model) tick() tea.Cmd {
	return tea.Tick(time.Second/4, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

func (m Model) Init() tea.Cmd {
	return m.tick()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.Keys.Quit):
			fmt.Println()
			return m, tea.Quit

		case key.Matches(msg, m.Keys.Up):
			if !m.Pressed && m.Cursor.Y > 0 {
				m.Cursor.Y--
			}
			// Disable the key until the next tick.
			// Since the view does not update in real time, this prevents
			// hidden states in the game that are invisible to the user.
			m.Pressed = true
		}
	case tea.WindowSizeMsg:
		// Terminal resized.
		m.Help.Width = msg.Width
	case TickMsg:
		if !m.Pressed {
			m.Cursor.Y++
		}
		m.Pressed = false
		for _, obst := range m.Obstacles {
			if (obst.Collides(m.Cursor)) || m.Cursor.Y >= m.Viewport.Y {
				m.Over = true
				return m, tea.Quit
			}
			if obst.Location.X == m.Cursor.X {
				m.Score++
			} else if obst.Location.X < 0 {
				m.Obstacles.Remove()
			}
			obst.Location.X--
		}
		rightmost := m.Obstacles.Rightmost()
		if rightmost == nil {
			return m, tea.Quit
		}
		// Create a new obstacle some percent of the time.
		gap := m.Viewport.X - rightmost.Location.X
		if gap > 5 || (rand.Intn(100) > 90 && gap > 2) {
			// Select a y that makes the obstacle possible to avoid.
			x := m.Viewport.X
			var y int
			yDelta := rand.Intn(x - rightmost.Location.X)
			if rand.Intn(100) > 50 {
				y = rightmost.Location.Y + yDelta
			} else {
				y = rightmost.Location.Y - yDelta
			}
			if y < 0 {
				y = 0
			} else if y >= m.Viewport.Y {
				y = m.Viewport.Y - 1
			}
			m.Obstacles.Add(NewObstacle(2, &Location{x, y}))
		}
		return m, m.tick()
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Flapioca"))
	sb.WriteByte('\n')

	viewport := make([]string, 0, m.Viewport.Y)
	for y := 0; y < m.Viewport.Y; y++ {
		var line strings.Builder
		// Store the index of the leftmost obstacle encountered.
		// This is used to slice the obstacle list to avoid checking obstacles
		// we've already seen.
		leftmost := 0
		for x := 0; x < m.Viewport.X; x++ {
			// Check if any obstacles collide with this cell.
			cellValue := ' '
			for _, o := range m.Obstacles[leftmost:] {
				if o.Collides(Location{x, y}) {
					cellValue = '#'
					leftmost++
					break
				}
			}
			if m.Cursor.X == x && m.Cursor.Y == y {
				if cellValue == '#' {
					cellValue = '@'
				} else {
					cellValue = '*'
				}
			}
			line.WriteRune(cellValue)
		}
		viewport = append(viewport, line.String())
	}

	sb.WriteString(ViewportStyle.Render(strings.Join(viewport, "\n")))
	sb.WriteString(fmt.Sprintf("\n%d point(s) ", m.Score))
	sb.WriteString(m.Help.View(m.Keys))

	if m.Over {
		sb.WriteString(GameOverStyle.Render("\n\n> Game over! <"))
	}

	// Send the UI for rendering
	return ViewStyle.Render(sb.String())
}
