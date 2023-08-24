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

type itemShowRow struct {
	label       string
	value       string
	hidden      bool
	blockRender bool
	marginTop   int
}

func (row itemShowRow) render(selected bool) string {
	value := row.value
	if row.hidden && !selected && len(value) > 0 {
		value = "•••"
	}
	spacer := "\t"
	if len(row.label) < 4 {
		spacer += "\t"
	}
	if selected {
		value = focusedStyle.Render(value)
	} else {
		value = mutedStyle.Render(value)
	}
	marginTop := ""
	for i := 0; i < row.marginTop; i++ {
		marginTop += "\n"
	}
	if row.blockRender {
		var label string
		if selected {
			label = selectedRowStyle.Render(row.label + ":")
		} else {
			label = rowStyle.Render(row.label + ":")
		}
		label += "\n"
		if row.value == "" {
			value = mutedStyle.Render(docStyle.Render("-"))
		} else {
			value = docStyle.Render(value)
			if selected {
				value = focusedStyle.Render(value)
			} else {
				value = mutedStyle.Render(value)
			}
		}
		return marginTop + label + value + "\n"
	}
	intermediate := fmt.Sprintf("%s:%s%s", row.label, spacer, value)
	if selected {
		return marginTop + selectedRowStyle.Render(intermediate) + "\n"
	}
	return marginTop + rowStyle.Render(intermediate) + "\n"
}

type ItemShow struct {
	bwm        *bw.Manager
	item       bw.Item
	selected   int
	keys       itemShowKeyBindings
	help       help.Model
	flashMsg   string
	flashTimer timer.Model
	rows       []itemShowRow
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

func (c ItemShow) setItem(listItem BWListItem) tea.Model {
	c.item = listItem.Item
	c.rows = make([]itemShowRow, 0)
	c.rows = append(c.rows, itemShowRow{
		label:     "Object",
		value:     c.item.Object,
		marginTop: 1,
	})
	c.rows = append(c.rows, itemShowRow{
		label: "ID",
		value: c.item.ID,
	})
	c.rows = append(c.rows, itemShowRow{
		label: "Type",
		value: c.item.Type.String(),
	})
	if c.item.OrganizationID != "" {
		c.rows = append(c.rows, itemShowRow{
			label: "Org ID",
			value: c.item.OrganizationID,
		})
	}
	if c.item.FolderID != "" {
		c.rows = append(c.rows, itemShowRow{
			label: "Folder ID",
			value: c.item.FolderID,
		})
	}
	c.rows = append(c.rows, itemShowRow{
		label:     "Username",
		value:     c.item.Login.Username,
		marginTop: 1,
	})
	c.rows = append(c.rows, itemShowRow{
		label:  "Password",
		value:  c.item.Login.Password,
		hidden: true,
	})
	c.rows = append(c.rows, itemShowRow{
		label:       "Notes",
		value:       c.item.Notes,
		marginTop:   1,
		blockRender: true,
	})
	return c
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
		return c.setItem(listItem), tick
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
			if c.selected > len(c.rows)-1 {
				c.selected = len(c.rows) - 1
			}
		case key.Matches(msg, c.keys.Copy):
			data := ""
			row := c.rows[c.selected]
			data = row.value
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

func (c ItemShow) View() string {
	sections := make([]string, 2)
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s Item | %s", logo, c.item.Name)))
	b.WriteString("\n")
	for i, row := range c.rows {
		b.WriteString(row.render(i == c.selected))
	}
	sections[0] = b.String()
	sections[1] = lipgloss.JoinVertical(lipgloss.Bottom, c.flashMsg, c.help.View(c.keys))
	return docStyle.Render(lipgloss.JoinVertical(lipgloss.Left, sections...))
}
