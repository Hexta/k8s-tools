package endpointslices

import (
	"context"
	"fmt"

	v2 "k8s.io/api/discovery/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Fetch retrieves a list of endpoint information from all namespaces using the provided Kubernetes clientset.
func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	infoList := make(InfoList, 0, 10000)

	for {
		list, err := clientset.DiscoveryV1().EndpointSlices(v1.NamespaceAll).List(ctx, v1.ListOptions{
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

				AddressType: string(item.AddressType),
				ApiVersion:  item.APIVersion,
				Endpoints:   getEndpoints(item),
				Kind:        item.Kind,
				Ports:       getPorts(item),
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}

func getPorts(es *v2.EndpointSlice) PortList {
	result := make(PortList, 0, len(es.Ports))

	for _, item := range es.Ports {
		port := &Port{
			EndpointSliceName: es.Name,
			Name:              item.Name,
			Namespace:         es.Namespace,

			AppProtocol: item.AppProtocol,
			Port:        item.Port,
			Protocol:    (*string)(item.Protocol),
		}
		result = append(result, port)
	}

	return result
}

func getEndpoints(endpointSlice *v2.EndpointSlice) EndpointList {
	result := make(EndpointList, 0, len(endpointSlice.Endpoints))

	for _, item := range endpointSlice.Endpoints {
		endpoint := &Endpoint{
			EndpointSliceName: endpointSlice.Name,
			Namespace:         endpointSlice.Namespace,

			Addresses: item.Addresses,
			Hostname:  item.Hostname,
			NodeName:  item.NodeName,
		}

		targetRef := item.TargetRef
		if targetRef != nil {
			endpoint.HasTargetRef = true
			endpoint.TargetRefAPIVersion = targetRef.APIVersion
			endpoint.TargetRefFieldPath = targetRef.FieldPath
			endpoint.TargetRefKind = targetRef.Kind
			endpoint.TargetRefName = targetRef.Name
			endpoint.TargetRefNamespace = targetRef.Namespace
			endpoint.TargetRefResourceVersion = targetRef.ResourceVersion
			endpoint.TargetRefUID = (string)(targetRef.UID)
		}

		result = append(result, endpoint)
	}

	return result
}
