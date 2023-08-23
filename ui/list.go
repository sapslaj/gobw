package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sapslaj/gobw/bw"
)

type ListSelectedEntry struct {
	item list.Item
}

func SelectListSelectedEntry(id list.Item) tea.Cmd {
	return func() tea.Msg {
		return ListSelectedEntry{id}
	}
}

type BWListItem struct {
	Item       bw.Item
	ID         string
	ObjectName string
	UserName   string
}

func NewBWListItem(bwi bw.Item) BWListItem {
	return BWListItem{
		Item:       bwi,
		ID:         bwi.ID,
		ObjectName: bwi.Name,
		UserName:   bwi.Login.Username,
	}
}

func (bwl BWListItem) Title() string       { return bwl.ObjectName }
func (bwl BWListItem) Description() string { return bwl.UserName }
func (bwl BWListItem) FilterValue() string { return bwl.ObjectName }

type UIList struct {
	list list.Model
	bwm  *bw.BWManager
}

func NewUIList(h int, v int, bwm *bw.BWManager) UIList {
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(selectedColor).BorderLeftForeground(selectedColor)
	d.Styles.SelectedDesc = d.Styles.SelectedTitle.Copy()
	width, height := docStyle.GetFrameSize()
	l := list.New(nil, d, h-width, v-height)
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("Enter", "view item"),
			),
		}
	}
	l.Styles.Title = titleStyle

	return UIList{
		list: l,
		bwm:  bwm,
	}
}

func (m UIList) Init() tea.Cmd {
	return nil
}

func (m *UIList) GetEntries() {
	listItems := []list.Item{}
	items, err := m.bwm.GetList()
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	for _, v := range items {
		listItems = append(listItems, NewBWListItem(v))
	}
	if err != nil {
		fmt.Println("Error: " + err.Error())
		os.Exit(1)
	}
	m.list.Title = fmt.Sprintf(" %s Vault | %s ", logo, m.bwm.VaultStatus.UserEmail)
	m.list.SetItems(listItems)
}

func (m UIList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoadingDone:
		m.GetEntries()
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			return m, SelectListSelectedEntry(m.list.SelectedItem())
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, tea.ClearScreen
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m UIList) View() string {
	return docStyle.Render(m.list.View())
}
