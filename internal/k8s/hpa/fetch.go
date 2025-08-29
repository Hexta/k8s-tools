package hpa

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/k8sutil"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	infoList := make(InfoList, 0, 10000)

	err := k8sutil.Paginate(ctx, func(opts v1.ListOptions) (string, error) {
		list, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(v1.NamespaceAll).List(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list deployments: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			infoList = append(infoList, &Info{
				Name:              item.Name,
				Namespace:         item.Namespace,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				CurrentReplicas:   item.Status.CurrentReplicas,
				DesiredReplicas:   item.Status.DesiredReplicas,
			})
		}

		return list.Continue, nil
	})

	return infoList, err
}
