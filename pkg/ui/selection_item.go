package ui

import "github.com/charmbracelet/bubbles/list"

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

func NewSelectionItem(id string, title string) SelectionItem {
	return SelectionItem{
		ID:    id,
		title: title,
	}
}
