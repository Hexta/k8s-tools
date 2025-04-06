package endpoints

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Fetch retrieves a list of endpoint information from all namespaces using the provided Kubernetes clientset.
func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	infoList := make(InfoList, 0, 10000)

	for {
		list, err := clientset.CoreV1().Endpoints(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list endpoints: %w", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]

			infoList = append(infoList, &Info{
				Annotations:       item.Annotations,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				Name:              item.Name,
				Namespace:         item.Namespace,
				Subsets:           getSubsets(item),
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}

// getSubsets extracts subsets from the provided Endpoints object and returns a slice of Subset objects.
func getSubsets(endpoints *corev1.Endpoints) []*Subset {
	result := make([]*Subset, 0, len(endpoints.Subsets))

	for _, subset := range endpoints.Subsets {
		result = appendAddresses(result, endpoints, subset.Addresses, nil, true)
		result = appendAddresses(result, endpoints, subset.NotReadyAddresses, nil, false)
	}

	return result
}

// appendAddresses adds address information to the result slice
func appendAddresses(
	result []*Subset,
	endpoints *corev1.Endpoints,
	addresses []corev1.EndpointAddress,
	ports []int32,
	isReady bool,
) []*Subset {
	for _, address := range addresses {
		subset := &Subset{
			EndpointsName: endpoints.Name,
			Hostname:      address.Hostname,
			IP:            address.IP,
			IsReady:       isReady,
			Namespace:     endpoints.Namespace,
			NodeName:      address.NodeName,
			Ports:         ports,
		}

		targetRef := address.TargetRef
		if targetRef != nil {
			subset.HasTargetRef = true
			subset.TargetRefAPIVersion = targetRef.APIVersion
			subset.TargetRefFieldPath = targetRef.FieldPath
			subset.TargetRefKind = targetRef.Kind
			subset.TargetRefName = targetRef.Name
			subset.TargetRefNamespace = targetRef.Namespace
			subset.TargetRefResourceVersion = targetRef.ResourceVersion
			subset.TargetRefUID = (string)(targetRef.UID)
		}

		result = append(result, subset)
	}
	return result
}
