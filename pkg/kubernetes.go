package pkg

import (
	"context"
	"github.com/thoas/go-funk"
	v12 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetNamespaces() ([]string, error) {
	if c, err := GetKubernetesClient(); err != nil {
		return []string{}, err
	} else {
		if n, err := c.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{}); err != nil {
			return []string{}, err
		} else {
			return funk.Map(n.Items, func(n v12.Namespace) string {
				return n.Name
			}).([]string), nil
		}
	}
}
