package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/sapslaj/gobw/bw"
)

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}

type itemShowKeyBindings struct {
	CursorUp     key.Binding
	CursorDown   key.Binding
	Copy         key.Binding
	CopyPassword key.Binding
	CopyUsername key.Binding
	Quit         key.Binding
}

func newItemShowKeyBindings() itemShowKeyBindings {
	return itemShowKeyBindings{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Copy: key.NewBinding(
			key.WithKeys("enter", "c", "y"),
			key.WithHelp("enter/c/y", "copy"),
		),
		CopyPassword: key.NewBinding(
			key.WithKeys("p"),
			key.WithHelp("p", "copy password"),
		),
		CopyUsername: key.NewBinding(
			key.WithKeys("u"),
			key.WithHelp("u", "copy username"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

func (k itemShowKeyBindings) ShortHelp() []key.Binding {
	return []key.Binding{k.CursorUp, k.CursorDown, k.Copy, k.CopyUsername, k.CopyPassword, k.Quit}
}

func (k itemShowKeyBindings) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}

type ItemShow struct {
	bwm        *bw.Manager
	item       bw.Item
	selected   int
	keys       itemShowKeyBindings
	help       help.Model
	flashMsg   string
	flashTimer timer.Model
}

func NewItemShow(bwm *bw.Manager) tea.Model {
	return ItemShow{
		bwm:  bwm,
		keys: newItemShowKeyBindings(),
		help: help.New(),
	}
}

func (c ItemShow) Init() tea.Cmd {
	return c.flashTimer.Init()
}

func (c ItemShow) flash(msg string) (tea.Model, tea.Cmd) {
	c.flashMsg = msg
	c.flashTimer = timer.NewWithInterval(5*time.Second, time.Second)
	return c, c.flashTimer.Start()
}

func (c ItemShow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		c.flashTimer, cmd = c.flashTimer.Update(msg)
		return c, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		c.flashTimer, cmd = c.flashTimer.Update(msg)
		return c, cmd
	case timer.TimeoutMsg:
		c.flashMsg = ""
	case ListSelectedEntry:
		listItem, ok := msg.item.(BWListItem)
		if !ok {
			panic("Could not get BWListItem")
		}
		c.item = listItem.Item
		return c, tick
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, c.keys.Quit):
			return c, SelectLoadingDone()
		case key.Matches(msg, c.keys.CursorUp):
			c.selected--
			if c.selected < 0 {
				c.selected = 0
			}
		case key.Matches(msg, c.keys.CursorDown):
			c.selected++
			if c.selected > 7 {
				c.selected = 7
			}
		case key.Matches(msg, c.keys.Copy):
			data := ""
			switch c.selected {
			case 0:
				data = c.item.Object
			case 1:
				data = c.item.ID
			case 2:
				data = c.item.Type.String()
			case 3:
				data = c.item.OrganizationID
			case 4:
				data = c.item.FolderID
			case 5:
				data = c.item.Login.Username
			case 6:
				data = c.item.Login.Password
			case 7:
				data = c.item.Notes
			}
			err := clipboard.WriteAll(data)
			if err != nil {
				panic(fmt.Errorf("error copying data to clipboard: %w", err))
			}
			return c.flash("copied to clipboard")
		case key.Matches(msg, c.keys.CopyPassword):
			err := clipboard.WriteAll(c.item.Login.Password)
			if err != nil {
				panic(fmt.Errorf("error copying password to clipboard: %w", err))
			}
			return c.flash("copied password to clipboard")
		case key.Matches(msg, c.keys.CopyUsername):
			err := clipboard.WriteAll(c.item.Login.Username)
			if err != nil {
				panic(fmt.Errorf("error copying username to clipboard: %w", err))
			}
			return c.flash("copied username to clipboard")
		}
	}
	var cmd tea.Cmd
	return c, cmd
}

func (c ItemShow) renderRow(i int, label string, value string, hidden bool) string {
	// this is a bad way to do this but will fix later
	if hidden && i != c.selected && len(value) > 0 {
		value = "•••"
	}
	spacer := "\t"
	if len(label) < 4 {
		spacer += "\t"
	}
	if i == c.selected {
		value = focusedStyle.Render(value)
	} else {
		value = mutedStyle.Render(value)
	}
	row := fmt.Sprintf("%s:%s%s", label, spacer, value)

	if i == c.selected {
		return selectedRowStyle.Render(row) + "\n"
	}
	return rowStyle.Render(row) + "\n"
}

func (c ItemShow) View() string {
	sections := make([]string, 2)
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s Item | %s", logo, c.item.Name)))
	b.WriteString("\n\n")
	b.WriteString(c.renderRow(0, "Object", c.item.Object, false))
	b.WriteString(c.renderRow(1, "ID", c.item.ID, false))
	b.WriteString(c.renderRow(2, "Type", c.item.Type.String(), false))
	b.WriteString(c.renderRow(3, "Org ID", c.item.OrganizationID, false))
	b.WriteString(c.renderRow(4, "Folder ID", c.item.FolderID, false))
	b.WriteString("\n")
	b.WriteString(c.renderRow(5, "Username", c.item.Login.Username, false))
	b.WriteString(c.renderRow(6, "Password", c.item.Login.Password, true))
	b.WriteString("\n")
	if c.selected == 7 {
		b.WriteString(selectedRowStyle.Render("Notes:"))
	} else {
		b.WriteString(rowStyle.Render("Notes:"))
	}
	b.WriteString("\n")
	notes := docStyle.Render(c.item.Notes)
	if c.item.Notes != "" {
		if c.selected == 7 {
			b.WriteString(focusedStyle.Render(notes))
		} else {
			b.WriteString(mutedStyle.Render(notes))
		}
	} else {
		b.WriteString(mutedStyle.Render(docStyle.Render("-")))
	}
	sections[0] = b.String()
	sections[1] = lipgloss.JoinVertical(lipgloss.Bottom, c.flashMsg, c.help.View(c.keys))
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
