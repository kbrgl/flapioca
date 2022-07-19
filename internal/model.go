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

type ObstacleList []*Obstacle

func NewObstacleList(obstacles ...*Obstacle) ObstacleList {
	return ObstacleList(obstacles)
}

func (ol *ObstacleList) Remove() {
	(*ol)[0] = nil
	*ol = (*ol)[1:]
}

func (ol *ObstacleList) Add(obst *Obstacle) {
	*ol = append(*ol, obst)
}

func (ol ObstacleList) Index(i int) *Obstacle {
	if i < 0 || i >= len(ol) {
		return nil
	}
	return (ol)[i]
}

func (ol ObstacleList) Rightmost() *Obstacle {
	if len(ol) == 0 {
		return nil
	}
	return (ol)[len(ol)-1]
}

type Model struct {
	keys      KeyMap
	obstacles ObstacleList
	cursor    Location
	score     int
	help      help.Model
	viewport  Location
	over      bool
	pressed   bool
}

func NewModel() Model {
	viewport := Location{
		x: 60,
		y: 9,
	}
	return Model{
		keys: KeyMap{
			Up:   key.NewBinding(key.WithKeys("k", "up", " ", "w"), key.WithHelp("↑/k/w/space", "jump")),
			Quit: key.NewBinding(key.WithKeys("q", "ctrl+c"), key.WithHelp("q/ctrl+c", "quit")),
		},
		obstacles: NewObstacleList(NewObstacle(2, &Location{viewport.x / 2, viewport.y / 2})),
		cursor:    Location{4, viewport.y / 2},
		score:     0,
		help:      help.New(),
		over:      false,
		viewport:  viewport,
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
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			if !m.pressed && m.cursor.y > 0 {
				m.cursor.y--
			}
			// Disable the key until the next tick.
			// Since the view does not update in real time, this prevents
			// hidden states in the game that are invisible to the user.
			m.pressed = true
		}
	case tea.WindowSizeMsg:
		// Terminal resized.
		m.help.Width = msg.Width
	case TickMsg:
		if !m.pressed {
			m.cursor.y++
		}
		m.pressed = false
		for _, obst := range m.obstacles {
			if (obst.Collides(m.cursor)) || m.cursor.y >= m.viewport.y {
				m.over = true
				return m, tea.Quit
			}
			if obst.Location.x == m.cursor.x {
				m.score++
			} else if obst.Location.x < 0 {
				m.obstacles.Remove()
			}
			obst.Location.x--
		}
		rightmost := m.obstacles.Rightmost()
		if rightmost == nil {
			return m, tea.Quit
		}
		// Create a new obstacle some percent of the time.
		gap := m.viewport.x - rightmost.Location.x
		if gap > 5 || (rand.Intn(100) > 90 && gap > 2) {
			// Select a y that makes the obstacle possible to avoid.
			x := m.viewport.x
			var y int
			yDelta := rand.Intn(x - rightmost.Location.x)
			if rand.Intn(100) > 50 {
				y = rightmost.Location.y + yDelta
			} else {
				y = rightmost.Location.y - yDelta
			}
			if y < 0 {
				y = 0
			} else if y >= m.viewport.y {
				y = m.viewport.y - 1
			}
			m.obstacles.Add(NewObstacle(2, &Location{x, y}))
		}
		return m, m.tick()
	}
	return m, nil
}

func (m Model) View() string {
	var sb strings.Builder
	sb.WriteString(TitleStyle.Render("Flapioca"))
	sb.WriteByte('\n')

	viewport := make([]string, 0, m.viewport.y)
	for y := 0; y < m.viewport.y; y++ {
		var line strings.Builder
		// Store the index of the leftmost obstacle encountered.
		// This is used to slice the obstacle list to avoid checking obstacles
		// we've already seen.
		leftmost := 0
		for x := 0; x < m.viewport.x; x++ {
			// Check if any obstacles collide with this cell.
			cellValue := ' '
			for _, o := range m.obstacles[leftmost:] {
				if o.Collides(Location{x, y}) {
					cellValue = '#'
					leftmost++
					break
				}
			}
			if m.cursor.x == x && m.cursor.y == y {
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
	sb.WriteString(fmt.Sprintf("\n%d point(s) ", m.score))
	sb.WriteString(m.help.View(m.keys))

	if m.over {
		sb.WriteString(GameOverStyle.Render("\n\n> Game over! <"))
	}

	// Send the UI for rendering
	return ViewStyle.Render(sb.String())
}
