package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/sapslaj/gobw/bw"
)

type loginType int

const (
	login loginType = iota
	unlock
)

type LoadingLoginFailed struct {
	lt loginType
}

func SelectLoadingFailed(lt loginType) tea.Cmd {
	return func() tea.Msg {
		return LoadingLoginFailed{lt}
	}
}

type LoadingDone struct{}

func SelectLoadingDone() tea.Cmd {
	return func() tea.Msg {
		return LoadingDone{}
	}
}

type Loading struct {
	un  string
	pw  string
	lt  loginType
	bwm *bw.Manager
}

func NewLoading(bwm *bw.Manager) Loading {
	return Loading{
		bwm: bwm,
	}
}

func (m Loading) Init() tea.Cmd {
	return nil
}

func (m Loading) Login() error {
	var err error
	switch m.lt {
	case login:
		err := m.bwm.Login(m.un, m.pw)
		if err != nil {
			return err
		}
	case unlock:
		err := m.bwm.Unlock(m.pw)
		if err != nil {
			return err
		}
	}
	err = m.bwm.UpdateList()
	return err
}

func (m Loading) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginSubmit:
		m.un = msg.un
		m.pw = msg.pw
		m.lt = msg.lt
		return m, tick
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	default:
		err := m.Login()
		if err != nil {
			return m, SelectLoadingFailed(m.lt)
		}
		return m, SelectLoadingDone()
	}
}

func (m Loading) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s ", logo)))
	b.WriteString("\n\n Logging in. Please wait\n\n")
	return docStyle.Render(b.String())
}
