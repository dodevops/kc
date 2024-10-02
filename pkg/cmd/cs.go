package cmd

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/thoas/go-funk"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd/api"
	api2 "kc/pkg/api"
	"kc/pkg/ui"
	"strings"
)

var (
	csExample = `
	# select a context and namespace
	%[1]s cs

	# switch to context foo, but keep its configured namespace
	%[1]s cs foo:

	# switch to namespace bar in the current context
	%[1]s cs :bar

	# switch to context foo and namespace bar
	%[1]s cs foo:bar
`
)

const (
	TaskNone = iota
	TaskSelect
	TaskSwitchContext
	TaskSwitchNamespace
	TaskSwitchBoth
)

type ContextSwitcherOptions struct {
	configFlags        *genericclioptions.ConfigFlags
	context            string
	namespace          string
	task               int
	rawConfig          api.Config
	kubernetesClient   *kubernetes.Clientset
	kubeConfigApi      api2.KubeConfig
	onlyCurrentContext bool
}

// NewContextSwitcherOptions provides an instance of NamespaceOptions with default values
func NewContextSwitcherOptions() *ContextSwitcherOptions {
	return &ContextSwitcherOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
	}
}

func NewCmdContextSwitcher() *cobra.Command {
	o := NewContextSwitcherOptions()

	cmd := &cobra.Command{
		Use:          "cs [context:namespace]",
		Short:        "Switch or select a context and/or namespace",
		Example:      fmt.Sprintf(csExample, "kubectl"),
		SilenceUsage: true,
		Annotations: map[string]string{
			cobra.CommandDisplayNameAnnotation: "kubectl cs",
		},
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(
		&o.onlyCurrentContext,
		"only-current-context",
		"c",
		false,
		"if true will only gather namespaces of the current context instead of all contexts",
	)

	return cmd
}

func (o *ContextSwitcherOptions) Complete(args []string) error {
	if rc, err := o.configFlags.ToRawKubeConfigLoader().RawConfig(); err != nil {
		return err
	} else {
		o.rawConfig = rc
	}

	if len(funk.Keys(o.rawConfig.Contexts).([]string)) == 0 {
		return fmt.Errorf("kubeconfig doesn't contain any contexts")
	}

	o.kubeConfigApi = api2.NewKubeConfigAPI(o.configFlags, o.rawConfig)

	o.task = TaskNone
	if len(args) == 0 {
		o.task = TaskSelect
	} else {
		o.interpretTask(args[0])
	}
	return nil
}

func (o *ContextSwitcherOptions) interpretTask(value string) {
	components := strings.Split(value, ":")
	if components[0] != "" {
		o.context = components[0]
		o.task = TaskSwitchContext
	}
	if components[1] != "" {
		o.namespace = components[1]
		o.task = TaskSwitchNamespace
	}
	if components[0] != "" && components[1] != "" {
		o.task = TaskSwitchBoth
	}
}

func (o *ContextSwitcherOptions) Validate() error {
	if o.task == TaskNone {
		return fmt.Errorf("wrong argument specified")
	}

	if o.task == TaskSwitchBoth || o.task == TaskSwitchContext &&
		!funk.ContainsString(funk.Keys(o.rawConfig.Contexts).([]string), o.context) {
		return fmt.Errorf("selected context does not exist")
	}

	return nil
}

func (o *ContextSwitcherOptions) Run() error {
	if o.task == TaskSelect {
		title := fmt.Sprintf(
			"Please select a namespace | Current context: %s | Current namespace: %s",
			o.kubeConfigApi.GetCurrentContext(),
			o.kubeConfigApi.GetCurrentNamespace(),
		)
		selector := tea.NewProgram(ui.NewSelectionList(title, o.kubeConfigApi, o.onlyCurrentContext))
		if m, err := selector.Run(); err != nil {
			return err
		} else {
			if err := m.(ui.SelectionList).Error; err != nil {
				return err
			}
			if m.(ui.SelectionList).WasCanceled {
				return nil
			}
			o.interpretTask(m.(ui.SelectionList).SelectedChoice.ID)
		}
	}
	switch o.task {
	default:
		return fmt.Errorf("invalid task")
	case TaskSwitchContext:
		return o.kubeConfigApi.SwitchContext(o.context)
	case TaskSwitchNamespace:
		if err := o.checkNamespace(o.rawConfig.CurrentContext, o.namespace); err != nil {
			return err
		}
		return o.kubeConfigApi.SwitchNamespace(o.rawConfig.CurrentContext, o.namespace)
	case TaskSwitchBoth:
		if err := o.checkNamespace(o.context, o.namespace); err != nil {
			return err
		}
		if err := o.kubeConfigApi.SwitchNamespace(o.context, o.namespace); err != nil {
			return err
		}
		return o.kubeConfigApi.SwitchContext(o.context)
	}
}

func (o *ContextSwitcherOptions) checkNamespace(context string, namespace string) error {
	if namespaces, err := o.kubeConfigApi.GetNamespacesInContext(context); err != nil {
		return err
	} else {
		if !funk.Contains(namespaces, namespace) {
			return fmt.Errorf("context %s doesn't contain namespace %s", o.rawConfig.CurrentContext, o.namespace)
		}
	}
	return nil
}
