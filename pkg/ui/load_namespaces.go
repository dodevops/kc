package ui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"kc/pkg/api"
)

type addItemsMessage struct{ items []list.Item }

type loadNextContext struct{}

type loadContextsMessage struct{}

type loadingFinished struct{}

func loadNamespace(kubeConfigAPI api.KubeConfig, context string) tea.Msg {
	if ns, err := kubeConfigAPI.GetNamespacesInContext(context); err != nil {
		return errorMessage{err: err}
	} else {
		var items []list.Item
		for _, namespace := range ns {
			title := fmt.Sprintf("%s:%s", context, namespace)
			items = append(items, NewSelectionItem(title, title))
		}
		return addItemsMessage{items: items}
	}
}
