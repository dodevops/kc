package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	tea "github.com/charmbracelet/bubbletea"
	"kc/internal"
	"kc/pkg"
	"os"
)

func contextDescription(context pkg.Context) string {
	if context.Namespace == "" {
		return context.Cluster
	}
	return fmt.Sprintf("%s/%s", context.Cluster, context.Namespace)
}

func main() {
	parser := argparse.NewParser("kc", "Kubernetes Context Switcher")
	namespace := parser.Flag("n", "namespace", &argparse.Options{
		Help:    "Switch current context's namespace instead",
		Default: false,
	})
	contextOrNamespace := parser.StringPositional(&argparse.Options{
		Help:     "The context or namespace to switch to. If empty, a selection list will be shown.",
		Required: false,
	})
	if err := parser.Parse(os.Args); err != nil {
		println(parser.Usage(err.Error()))
		os.Exit(1)
	}

	c := *contextOrNamespace

	if c == "" {
		if *namespace {
			if ns, err := pkg.GetNamespaces(); err != nil {
				println(parser.Usage(err.Error()))
				os.Exit(1)
			} else {
				var choices []internal.SelectionItem
				for _, namespace := range ns {
					choices = append(choices, internal.NewSelectionItem(namespace, namespace, namespace))
				}
				selector := tea.NewProgram(internal.NewSelectionList("Please select a namespace", choices, false))
				if m, err := selector.Run(); err != nil {
					println(parser.Usage(err.Error()))
					os.Exit(1)
				} else {
					c = m.(internal.SelectionList).SelectedChoice.ID
				}
			}
		} else {
			var choices []internal.SelectionItem
			for _, context := range pkg.GetContexts() {
				choices = append(choices, internal.NewSelectionItem(context.Name, context.Name, contextDescription(context)))
			}
			selector := tea.NewProgram(internal.NewSelectionList("Please select a context", choices, true))
			if m, err := selector.Run(); err != nil {
				println(parser.Usage(err.Error()))
				os.Exit(1)
			} else {
				c = m.(internal.SelectionList).SelectedChoice.ID
			}
		}
	}

	if c == "" {
		os.Exit(0)
	}

	if *namespace {
		if err := pkg.ChangeNamespaceOfCurrentContext(c); err != nil {
			println(parser.Usage(err.Error()))
			os.Exit(1)
		}
	} else {
		if err := pkg.SetCurrentContext(c); err != nil {
			println(parser.Usage(err.Error()))
			os.Exit(1)
		}
	}
}
