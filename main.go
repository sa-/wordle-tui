package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	words "github.com/sa-/wordle-tui/words"
)

var (
	correctStyle  = lipgloss.NewStyle().Background(lipgloss.Color("#40a02b")).Foreground(lipgloss.Color("#fff"))
	wrongPosStyle = lipgloss.NewStyle().Background(lipgloss.Color("#df8e1d")).Foreground(lipgloss.Color("#fff"))
	wrongStyle    = lipgloss.NewStyle().Background(lipgloss.Color("#777")).Foreground(lipgloss.Color("#fff"))
	emptyStyle    = lipgloss.NewStyle().Width(15)
)

const (
	LetterStatusUnguessed = iota
	LetterStatusGreen
	LetterStatusYellow
	LetterStatusWrong
)

type model struct {
	word           string
	guesses        []string
	input          textinput.Model
	gameOver       bool
	win            bool
	windowWidth    int
	windowHeight   int
	guessedLetters map[rune]int
	cheats         bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "guess"
	ti.Focus()
	ti.CharLimit = 5
	ti.Width = 5

	wordToGuess := os.Getenv("WORDLE_WORD")
	if len(wordToGuess) == 0 {
		wordToGuess = words.AnswerWords[rand.Intn(len(words.AnswerWords))]
	}
	wordToGuess = strings.ToUpper(wordToGuess)

	return model{
		word:           wordToGuess,
		guesses:        make([]string, 0),
		input:          ti,
		windowWidth:    20,
		guessedLetters: make(map[rune]int),
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
		case tea.KeyCtrlY:
			m.cheats = !m.cheats
			return m, cmd
		case tea.KeyEnter:
			if m.gameOver {
				newModel := initialModel()
				newModel.windowWidth = m.windowWidth
				newModel.windowHeight = m.windowHeight
				return newModel, cmd
			}
			if len(m.input.Value()) == 5 {
				currentGuess := strings.ToUpper(m.input.Value())

				if !words.IsValidGuess(currentGuess) {
					return m, cmd
				}

				// process the guess
				for i, ch := range currentGuess {
					if ch == rune(m.word[i]) {
						m.guessedLetters[ch] = LetterStatusGreen
					} else if strings.ContainsRune(m.word, ch) {
						m.guessedLetters[ch] = LetterStatusYellow
					} else {
						m.guessedLetters[ch] = LetterStatusWrong
					}
				}
				m.guesses = append(m.guesses, currentGuess)
				//

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
	if !m.gameOver {
		m.input, cmd = m.input.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	titleText := "WORDLE!"
	if m.cheats {
		titleText = m.word
	}
	renderedTitle := lipgloss.NewStyle().AlignHorizontal(lipgloss.Center).Height(2).Width(21).Render(titleText)

	s := strings.Builder{}
	for i := 0; i < 6; i++ {
		if i < len(m.guesses) {
			s.WriteString(renderGuess(m.guesses[i], m))
		} else {
			s.WriteString(renderEmptyRow())
		}
		s.WriteString("\n")
	}

	if !m.gameOver {
		s.WriteString("\n")
		s.WriteString(m.input.View())
		s.WriteString("\n")
	} else {
		if m.win {
			s.WriteString("\n")
			s.WriteString("🎊 Congratulations 😎")
			s.WriteString("\n")
		} else {
			s.WriteString(fmt.Sprintf("Game over!\nThe word was:\n%s", m.word))
		}
	}

	letters := "Q W E R T Y U I O P\n A S D F G H J K L \n   Z X C V B N M   "
	lettersColored := strings.Builder{}
	for _, ch := range letters {
		if unicode.IsLetter(ch) {
			letterStatus := m.guessedLetters[ch]
			switch letterStatus {
			case LetterStatusUnguessed:
				lettersColored.WriteString(lipgloss.NewStyle().Render(string(ch)))
			case LetterStatusGreen:
				lettersColored.WriteString(correctStyle.Render(string(ch)))
			case LetterStatusYellow:
				lettersColored.WriteString(wrongPosStyle.Render(string(ch)))
			case LetterStatusWrong:
				lettersColored.WriteString(wrongStyle.Foreground(lipgloss.Color("#ccc")).Render(string(ch)))
			}
		} else {
			lettersColored.WriteRune(ch)
		}
	}
	letters = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).PaddingLeft(1).PaddingRight(1).Align(lipgloss.Center).Render(lettersColored.String())

	style := lipgloss.NewStyle().Width(m.windowWidth).Height(m.windowHeight).AlignVertical(lipgloss.Center).Align(lipgloss.Center)
	content := lipgloss.JoinVertical(lipgloss.Center, renderedTitle, s.String(), letters)
	return style.Render(content)
}

func renderGuess(guess string, m model) string {
	letterCount := make(map[rune]int)
	for _, ch := range m.word {
		letterCount[ch]++
	}

	letterStatus := [5]int{}

	// green pass
	for i, ch := range guess {
		if ch == rune(m.word[i]) {
			letterCount[ch]--
			letterStatus[i] = LetterStatusGreen
		}
	}

	// yellow and grey pass
	for i, ch := range guess {
		if letterStatus[i] == LetterStatusGreen {
			continue
		} else if letterCount[ch] > 0 {
			letterStatus[i] = LetterStatusYellow
			letterCount[ch]--
		} else {
			letterStatus[i] = LetterStatusWrong
		}
	}

	s := make([]string, 5)
	for i, ch := range guess {
		style := wrongStyle
		switch letterStatus[i] {
		case LetterStatusGreen:
			style = correctStyle
		case LetterStatusYellow:
			style = wrongPosStyle
		}
		s[i] = style.Render(fmt.Sprintf(" %c ", ch))
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, s...)
}

func processGuess(currentGuess string, m model) {

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
