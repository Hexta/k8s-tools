package sts

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	infoList := make(InfoList, 0, 10000)

	for {
		log.Debugf("Listing STS with continue token %q", continueToken)
		list, err := clientset.AppsV1().StatefulSets(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list deployments: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			infoList = append(infoList, &Info{
				Name:              item.Name,
				Namespace:         item.Namespace,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				Replicas:          item.Status.Replicas,
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}
