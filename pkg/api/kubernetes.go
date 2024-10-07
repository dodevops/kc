package api

import (
	"context"
	"github.com/thoas/go-funk"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetNamespaces returns all available namespaces
func GetNamespaces(c *kubernetes.Clientset) ([]string, error) {
	var timeout int64 = 5
	if n, err := c.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{TimeoutSeconds: &timeout}); err != nil {
		return []string{}, err
	} else {
		return funk.Map(n.Items, func(n v12.Namespace) string {
			return n.Name
		}).([]string), nil
	}
}
