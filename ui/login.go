package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginSubmit struct {
	un string
	pw string
	lt loginType
}

func SelectSubmit(un string, pw string) tea.Cmd {
	return func() tea.Msg {
		return LoginSubmit{un, pw, login}
	}
}

type Login struct {
	focusIndex int
	inputs     []textinput.Model
	text       string
}

func NewLogin() Login {
	l := Login{
		inputs: make([]textinput.Model, 2),
		text:   "Please enter your Bitwarden Login",
	}

	var t textinput.Model
	for i := range l.inputs {
		t = textinput.New()
		t.CursorStyle = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
			t.CharLimit = 64
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		l.inputs[i] = t
	}

	return l
}

func (l Login) Init() tea.Cmd {
	return textinput.Blink
}

func (l Login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadingLoginFailed:
		l.text = "Login Failed. Please try again or press 'esc' to exit"
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return l, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && l.focusIndex == len(l.inputs) {
				return l, SelectSubmit(l.inputs[0].Value(), l.inputs[1].Value())
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				l.focusIndex--
			} else {
				l.focusIndex++
			}

			if l.focusIndex > len(l.inputs) {
				l.focusIndex = 0
			} else if l.focusIndex < 0 {
				l.focusIndex = len(l.inputs)
			}

			cmds := make([]tea.Cmd, len(l.inputs))
			for i := 0; i <= len(l.inputs)-1; i++ {
				if i == l.focusIndex {
					// Set focused state
					cmds[i] = l.inputs[i].Focus()
					l.inputs[i].PromptStyle = focusedStyle
					l.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				l.inputs[i].Blur()
				l.inputs[i].PromptStyle = noStyle
				l.inputs[i].TextStyle = noStyle
			}

			return l, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := l.updateInputs(msg)

	return l, cmd
}

func (l *Login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (l Login) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s ", logo)))
	b.WriteString("\n\n")
	b.WriteString(l.text)
	b.WriteString("\n\n")
	for i := range l.inputs {
		b.WriteString(l.inputs[i].View())
		if i < len(l.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if l.focusIndex == len(l.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	return docStyle.Render(b.String())
}
