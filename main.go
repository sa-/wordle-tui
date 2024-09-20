package main

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	words = []string{
		"WORLD", "GLOBE", "EARTH", "OCEAN", "RIVER",
	}

	correctStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#3a3")).Foreground(lipgloss.Color("#fff"))
	wrongPosStyle = lipgloss.NewStyle().Background(lipgloss.Color("#aa3")).Foreground(lipgloss.Color("#fff"))
	wrongStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#777")).Foreground(lipgloss.Color("#fff"))
	emptyStyle    = lipgloss.NewStyle().Width(15)
)

type model struct {
	word         string
	guesses      []string
	input        textinput.Model
	gameOver     bool
	win          bool
	windowWidth  int
	windowHeight int
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "guess"
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 20

	return model{
		word:        words[rand.Intn(len(words))],
		guesses:     make([]string, 0),
		input:       ti,
		windowWidth: 20,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if len(m.input.Value()) == 5 {
				m.guesses = append(m.guesses, strings.ToUpper(m.input.Value()))
				m.input.SetValue("")
				if m.guesses[len(m.guesses)-1] == m.word {
					m.gameOver = true
					m.win = true
				} else if len(m.guesses) == 6 {
					m.gameOver = true
				}
			}
		}
	}

	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) View() string {

	title := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Width(15).PaddingLeft(1).PaddingRight(1).Border(lipgloss.NormalBorder()).Render("Wordle")

	s := strings.Builder{}
	for i := 0; i < 6; i++ {
		if i < len(m.guesses) {
			s.WriteString(renderGuess(m.guesses[i], m.word))
		} else {
			s.WriteString(renderEmptyRow())
		}
		s.WriteString("\n")
	}

	s.WriteString("\n")
	s.WriteString(m.input.View())
	s.WriteString("\n\n")

	letters := "\nQ W E R T Y U I O P\n A S D F G H J K L\nZ X C V B N M"
	letters = lipgloss.NewStyle().Bold(true).Align(lipgloss.Left).Render(letters)

	if m.gameOver {
		if m.win {
			s.WriteString("\n\nCongratulations! You guessed the word!")
		} else {
			s.WriteString(fmt.Sprintf("Game over! The word was: %s", m.word))
		}
	}

	style := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight).AlignVertical(lipgloss.Center).Align(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, title, s.String(), letters)
	return style.Render(content)
}

func renderGuess(guess, word string) string {
	s := make([]string, 5)
	for i, ch := range guess {
		style := wrongStyle
		if ch == rune(word[i]) {
			style = correctStyle
		} else if strings.ContainsRune(word, ch) {
			style = wrongPosStyle
		}
		s[i] = style.Render(fmt.Sprintf(" %c ", ch))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, s...)
}

func renderEmptyRow() string {
	return emptyStyle.Render(" _  _  _  _  _ ")
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
	}
}

/*
┌───┐
│ A │
└───┘
*/
