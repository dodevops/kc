package internal

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectionItem struct {
	ID          string
	title       string
	description string
}

var _ list.DefaultItem = SelectionItem{}

func (s SelectionItem) Title() string {
	return s.title
}

func (s SelectionItem) Description() string {
	return s.description
}

func (s SelectionItem) FilterValue() string {
	return s.title
}

func NewSelectionItem(id string, title string, description string) SelectionItem {
	return SelectionItem{
		ID:          id,
		title:       title,
		description: description,
	}
}

type SelectionList struct {
	list           list.Model
	SelectedChoice SelectionItem
}

var _ tea.Model = SelectionList{}

func NewSelectionList(title string, choices []SelectionItem, showDescription bool) SelectionList {
	d := list.NewDefaultDelegate()
	d.SetSpacing(0)
	d.ShowDescription = showDescription

	var c []list.Item

	for _, choice := range choices {
		c = append(c, choice)
	}

	l := list.New(c, d, 0, 0)
	l.Title = title

	return SelectionList{
		list: l,
	}
}

func (s SelectionList) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (s SelectionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			s.SelectedChoice = s.list.SelectedItem().(SelectionItem)
			return s, tea.Quit
		}
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width, msg.Height)
	}
	lm, cmd := s.list.Update(msg)
	s.list = lm
	return s, cmd
}

func (s SelectionList) View() string {
	return s.list.View()
}
