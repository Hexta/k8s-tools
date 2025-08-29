package service

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
		l, err := clientset.CoreV1().Services(v1.NamespaceAll).List(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list services: %v", err)
		}
		for idx := range l.Items {
			item := &l.Items[idx]
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
		return l.Continue, nil
	})

	return infoList, err
}
