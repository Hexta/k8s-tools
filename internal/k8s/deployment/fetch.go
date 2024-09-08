package deployment

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset, labelSelector string) (InfoList, error) {
	var continueToken string
	deployments := make(InfoList, 0, 10000)

	for {
		log.Debugf("Listing deployments with label selector %q, continue token %q", labelSelector, continueToken)
		list, err := clientset.AppsV1().Deployments(v1.NamespaceAll).List(ctx, v1.ListOptions{
			LabelSelector: labelSelector,
			Continue:      continueToken,
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
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return deployments, nil
}
