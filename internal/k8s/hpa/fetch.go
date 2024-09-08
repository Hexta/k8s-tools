package hpa

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	deployments := make(InfoList, 0, 10000)

	for {
		log.Debugf("Listing deployments with continue token %q", continueToken)
		list, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			deployments = append(deployments, &Info{
				Name:              item.Name,
				Namespace:         item.Namespace,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				CurrentReplicas:   item.Status.CurrentReplicas,
				DesiredReplicas:   item.Status.DesiredReplicas,
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return deployments, nil
}
