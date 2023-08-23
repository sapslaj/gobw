package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sapslaj/gobw/bw"
)

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}

type property int

const (
	copyUsername property = iota
	copyPassword
)

type UIItemShow struct {
	prop property
	bwm  *bw.BWManager
	item bw.Item
}

func NewUIItemShow(bwm *bw.BWManager) tea.Model {
	return UIItemShow{
		bwm: bwm,
	}
}

func (c UIItemShow) Init() tea.Cmd {
	return nil
}

func (c UIItemShow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ListSelectedEntry:
		listItem, ok := msg.item.(BWListItem)
		if !ok {
			panic("Could not get BWListItem")
		}
		switch msg.prop {
		case copyPassword:
			data, err := c.bwm.GetPassword(listItem.ID)
			if err != nil {
				panic("Error getting Password")
			}
			err = clipboard.WriteAll(data)
			if err != nil {
				panic(fmt.Errorf("Error copying password to clipboard: %w", err))
			}
			c.prop = copyPassword
		case copyUsername:
			data := listItem.UserName
			err := clipboard.WriteAll(data)
			if err != nil {
				panic(fmt.Errorf("Error copying username to clipboard: %w", err))
			}
			c.prop = copyUsername
		}
		c.item = listItem.Item
		return c, tick
	case tea.KeyMsg:
		if msg.String() == "q" {
			return c, SelectLoadingDone()
		}
	}
	var cmd tea.Cmd
	return c, cmd
}

func (c UIItemShow) View() string {
	var b strings.Builder
	b.WriteString("  ")
	b.WriteString(titleStyle.Render(fmt.Sprintf(" %s Item | %s", logo, c.item.Name)))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("Object:\t%s\n", focusedStyle.Render(c.item.Object)))
	b.WriteString(fmt.Sprintf("ID:\t\t%s\n", focusedStyle.Render(c.item.ID)))
	b.WriteString(fmt.Sprintf("Type:\t\t%s\n", focusedStyle.Render(c.item.Type.String())))
	if c.item.OrganizationID != "" {
		b.WriteString(fmt.Sprintf("Organization ID:\t%s\n", focusedStyle.Render(c.item.OrganizationID)))
	}
	if c.item.FolderID != "" {
		b.WriteString(fmt.Sprintf("Folder ID:\t%s\n", focusedStyle.Render(c.item.FolderID)))
	}
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Username:\t%s\n", focusedStyle.Render(c.item.Login.Username)))
	b.WriteString(fmt.Sprintf("Password:\t%s\n", focusedStyle.Render(c.item.Login.Password)))
	b.WriteString("\nNotes:\n")
	if c.item.Notes != "" {
		b.WriteString(c.item.Notes)
	} else {
		b.WriteString("-")
	}
	return docStyle.Render(b.String())
}
