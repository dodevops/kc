package api

import (
	"github.com/thoas/go-funk"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"log"
	"slices"
)

type KubeConfig struct {
	configFlags      *genericclioptions.ConfigFlags
	config           api.Config
	kubernetesClient *kubernetes.Clientset
	counter          int
}

func NewKubeConfigAPI(configFlags *genericclioptions.ConfigFlags, config api.Config) KubeConfig {
	return KubeConfig{
		configFlags: configFlags,
		config:      config,
	}
}

func (r *KubeConfig) updateKubernetesClient() error {
	if configBytes, err := clientcmd.Write(r.config); err != nil {
		return err
	} else {
		if restConfig, err := clientcmd.RESTConfigFromKubeConfig(configBytes); err != nil {
			return err
		} else {
			if client, err := kubernetes.NewForConfig(restConfig); err != nil {
				return err
			} else {
				r.kubernetesClient = client
			}
		}
	}
	return nil
}

func (r *KubeConfig) GetContexts() []string {
	contexts := funk.Keys(r.config.Contexts).([]string)
	slices.Sort(contexts)
	return contexts
}

func (r *KubeConfig) GetCurrentContext() string {
	return r.config.CurrentContext
}

func (r *KubeConfig) GetCurrentNamespace() string {
	return r.config.Contexts[r.config.CurrentContext].Namespace
}

func (r *KubeConfig) GetNamespacesInContext(context string) ([]string, error) {
	log.Printf("Getting namespaces in context %s", context)
	var namespaces []string
	oldContext := r.config.CurrentContext

	if oldContext != context {
		if err := r.SwitchContext(context); err != nil {
			return namespaces, err
		}
	}

	if err := r.updateKubernetesClient(); err != nil {
		return namespaces, err
	}

	if ns, err := GetNamespaces(r.kubernetesClient); err != nil {
		return namespaces, err
	} else {
		namespaces = ns
	}

	log.Printf("Found %d namespaces", len(namespaces))

	if oldContext != context {
		log.Printf("Switching back to old context")
		if err := r.SwitchContext(context); err != nil {
			return namespaces, err
		}
	}

	return namespaces, nil
}

func (r *KubeConfig) SwitchContext(context string) error {
	log.Printf("Switching to context %s", context)
	r.config.CurrentContext = context
	return r.updateConfig(r.config)
}

func (r *KubeConfig) SwitchNamespace(context string, namespace string) error {
	r.config.Contexts[context].Namespace = namespace
	r.config.CurrentContext = context
	return r.updateConfig(r.config)
}

func (r *KubeConfig) updateConfig(config api.Config) error {
	configAccess := clientcmd.NewDefaultPathOptions()
	return clientcmd.ModifyConfig(configAccess, config, true)
}
