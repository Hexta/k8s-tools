package service

import (
	"context"
	"fmt"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	infoList := make(InfoList, 0, 10000)

	for {
		list, err := clientset.CoreV1().Services(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list services: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			infoList = append(infoList, &Info{
				Annotations:       item.Annotations,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				Name:              item.Name,
				Namespace:         item.Namespace,

				ClusterIP:                item.Spec.ClusterIP,
				ClusterIPs:               item.Spec.ClusterIPs,
				ExternalIPs:              item.Spec.ExternalIPs,
				ExternalTrafficPolicy:    item.Spec.ExternalTrafficPolicy,
				HealthCheckNodePort:      item.Spec.HealthCheckNodePort,
				IPFamilies:               item.Spec.IPFamilies,
				IPFamilyPolicy:           item.Spec.IPFamilyPolicy,
				InternalTrafficPolicy:    item.Spec.InternalTrafficPolicy,
				LoadBalancerClass:        item.Spec.LoadBalancerClass,
				LoadBalancerSourceRanges: item.Spec.LoadBalancerSourceRanges,
				Ports:                    item.Spec.Ports,
				PublishNotReadyAddresses: item.Spec.PublishNotReadyAddresses,
				Selector:                 item.Spec.Selector,
				SessionAffinity:          item.Spec.SessionAffinity,
				Type:                     item.Spec.Type,
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}
