package deployment

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/k8sutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	deployments := make(InfoList, 0, 10000)

	err := k8sutil.Paginate(ctx, func(opts v1.ListOptions) (string, error) {
		list, err := clientset.AppsV1().Deployments(v1.NamespaceAll).List(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list deployments: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			deployments = append(deployments, &Info{
				Name:              item.Name,
				Namespace:         item.Namespace,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				Replicas:          item.Spec.Replicas,
			})
		}

		return list.Continue, nil
	})

	return deployments, err
}
