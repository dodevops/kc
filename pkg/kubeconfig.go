package pkg

import (
	"fmt"
	"github.com/thoas/go-funk"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func getClientConfig() clientcmd.ClientConfig {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{}
	// if you want to change override values or bind them to flags, there are methods to help you

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
}

func getKubeConfig() api.Config {
	clientConfig := getClientConfig()
	if kubeConfig, err := clientConfig.RawConfig(); err != nil {
		panic(fmt.Sprintf("Can not read Kubeconfig: %s", err.Error()))
	} else {
		return kubeConfig
	}
}

type Context struct {
	Name      string
	Cluster   string
	Namespace string
}

func GetContexts() []Context {
	return funk.Map(getKubeConfig().Contexts, func(name string, context *api.Context) Context {
		return Context{
			Name:      name,
			Cluster:   context.Cluster,
			Namespace: context.Namespace,
		}
	}).([]Context)
}

func SetCurrentContext(context string) error {
	config := getKubeConfig()
	if !funk.ContainsString(funk.Keys(config.Contexts).([]string), context) {
		return fmt.Errorf("context %s not found", context)
	}
	config.CurrentContext = context
	return clientcmd.ModifyConfig(clientcmd.NewDefaultClientConfigLoadingRules(), config, true)
}

func ChangeNamespaceOfCurrentContext(namespace string) error {
	config := getKubeConfig()
	config.Contexts[config.CurrentContext].Namespace = namespace
	return clientcmd.ModifyConfig(clientcmd.NewDefaultClientConfigLoadingRules(), config, true)
}

func GetKubernetesClient() (*kubernetes.Clientset, error) {
	if c, err := getClientConfig().ClientConfig(); err != nil {
		return nil, err
	} else {
		return kubernetes.NewForConfig(c)
	}
}
