package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"kc/pkg/api"
	"slices"
)

type SelectionList struct {
	list               list.Model
	spinner            spinner.Model
	SelectedChoice     SelectionItem
	WasCanceled        bool
	Loading            bool
	Error              error
	LoadingContext     string
	ContextsToLoad     []string
	kubeConfigAPI      api.KubeConfig
	onlyCurrentContext bool
}

var _ tea.Model = SelectionList{}

func NewSelectionList(title string, kubeConfigAPI api.KubeConfig, onlyCurrentContext bool) SelectionList {
	_, _ = tea.LogToFile("/tmp/cslog", "debug")
	d := list.NewDefaultDelegate()
	d.SetSpacing(0)
	d.ShowDescription = false

	l := list.New(make([]list.Item, 0), d, 0, 0)
	l.Title = title

	return SelectionList{
		list:               l,
		spinner:            spinner.New(spinner.WithSpinner(spinner.Dot)),
		WasCanceled:        true,
		Loading:            true,
		kubeConfigAPI:      kubeConfigAPI,
		onlyCurrentContext: onlyCurrentContext,
	}
}

func (s SelectionList) Init() tea.Cmd {
	return tea.Sequence(tea.EnterAltScreen, s.spinner.Tick, func() tea.Msg {
		return loadContextsMessage{}
	})
}

func (s SelectionList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case loadContextsMessage:
		s.ContextsToLoad = []string{s.kubeConfigAPI.GetCurrentContext()}
		for _, context := range s.kubeConfigAPI.GetContexts() {
			if context != s.kubeConfigAPI.GetCurrentContext() {
				s.ContextsToLoad = append(s.ContextsToLoad, context)
			}
		}
		return s, func() tea.Msg {
			return loadNextContext{}
		}
	case loadNextContext:
		if len(s.ContextsToLoad) == 0 {
			return s, func() tea.Msg {
				return loadingFinished{}
			}
		}
		s.LoadingContext = s.ContextsToLoad[0]
		s.ContextsToLoad = slices.Delete(s.ContextsToLoad, 0, 1)
		return s, func() tea.Msg {
			if s.onlyCurrentContext && s.LoadingContext != s.kubeConfigAPI.GetCurrentContext() {
				return addItemsMessage{items: []list.Item{NewSelectionItem(fmt.Sprintf("%s:", s.LoadingContext), s.LoadingContext)}}
			} else {
				return loadNamespace(s.kubeConfigAPI, s.LoadingContext)
			}
		}
	case addItemsMessage:
		s.list.SetItems(slices.Concat(s.list.Items(), msg.items))
		cmds = append(cmds, func() tea.Msg {
			return loadNextContext{}
		})
	case loadingFinished:
		s.Loading = false
		if len(s.list.Items()) == 0 {
			return s, tea.Quit
		} else {
			s.Error = nil
		}
	case errorMessage:
		s.Error = msg.err
		return s, tea.Quit
	case tea.KeyMsg:
		if s.list.FilterState() != list.Filtering {
			switch msg.Type {
			case tea.KeyEnter:
				s.SelectedChoice = s.list.SelectedItem().(SelectionItem)
				s.WasCanceled = false
				return s, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		s.list.SetSize(msg.Width, msg.Height)
	}
	if s.Loading {
		sm, scmd := s.spinner.Update(msg)
		cmds = append(cmds, scmd)
		s.spinner = sm
	} else {
		lm, lcmd := s.list.Update(msg)
		cmds = append(cmds, lcmd)
		s.list = lm
	}
	return s, tea.Sequence(cmds...)
}

func (s SelectionList) View() string {
	if s.Loading {
		errorText := ""
		if s.Error != nil {
			errorText = s.Error.Error()
		}
		return fmt.Sprintf("%s Please wait. Loading namespaces from context %s...\n\n%s", s.spinner.View(), s.LoadingContext, errorText)
	} else {
		return s.list.View()
	}
}
