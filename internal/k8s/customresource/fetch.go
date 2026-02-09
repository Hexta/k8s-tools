package customresource

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/k8sutil"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

func Fetch(ctx context.Context, dynamicClient *dynamic.DynamicClient, apiExtClient apiextensionsclient.Interface) (InfoList, error) {
	crdList, err := apiExtClient.ApiextensionsV1().CustomResourceDefinitions().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list custom resource definitions: %w", err)
	}

	infoList := make(InfoList, 0, 10000)
	for idx := range crdList.Items {
		crd := &crdList.Items[idx]
		if crd.Spec.Names.Plural == "" {
			continue
		}

		resources, err := fetchCRDResources(ctx, dynamicClient, crd)
		if err != nil {
			return nil, err
		}

		infoList = append(infoList, resources...)
	}

	return infoList, nil
}

func fetchCRDResources(ctx context.Context, dynamicClient dynamic.Interface, crd *apiextensionsv1.CustomResourceDefinition) (InfoList, error) {
	var infoList InfoList
	isNamespaced := crd.Spec.Scope == apiextensionsv1.NamespaceScoped

	for _, version := range crd.Spec.Versions {
		if version.Name == "" {
			continue
		}

		gvr := schema.GroupVersionResource{
			Group:    crd.Spec.Group,
			Version:  version.Name,
			Resource: crd.Spec.Names.Plural,
		}

		resourceClient := dynamicClient.Resource(gvr)

		err := k8sutil.Paginate(ctx, func(opts metav1.ListOptions) (string, error) {
			list, err := listResources(ctx, resourceClient, isNamespaced, opts)
			if err != nil {
				return "", fmt.Errorf("failed to list custom resources for %s/%s/%s: %w",
					crd.Spec.Group, version.Name, crd.Spec.Names.Plural, err)
			}

			for itemIdx := range list.Items {
				info, err := convertToInfo(&list.Items[itemIdx], crd.Spec.Group, crd.Spec.Names.Kind, version.Name)
				if err != nil {
					return "", err
				}

				infoList = append(infoList, info)
			}

			return list.GetContinue(), nil
		})
		if err != nil {
			return nil, err
		}
	}

	return infoList, nil
}

func listResources(ctx context.Context, client dynamic.NamespaceableResourceInterface, namespaced bool, opts metav1.ListOptions) (*unstructured.UnstructuredList, error) {
	if namespaced {
		return client.Namespace(metav1.NamespaceAll).List(ctx, opts)
	}

	return client.List(ctx, opts)
}

func convertToInfo(item *unstructured.Unstructured, group, kind, version string) (*Info, error) {
	payload, err := json.Marshal(item.Object)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal custom resource %s/%s: %w",
			item.GetNamespace(), item.GetName(), err)
	}

	return &Info{
		Name:              item.GetName(),
		Namespace:         item.GetNamespace(),
		CreationTimestamp: item.GetCreationTimestamp().Time,
		Labels:            item.GetLabels(),
		ResourceGroup:     group,
		ResourceKind:      kind,
		ResourceVersion:   version,
		Data:              string(payload),
	}, nil
}
