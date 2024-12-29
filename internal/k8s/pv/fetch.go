package pv

import (
	"context"
	"fmt"

	v2 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, clientset *kubernetes.Clientset) (InfoList, error) {
	var continueToken string
	infoList := make(InfoList, 0, 10000)

	for {
		list, err := clientset.CoreV1().PersistentVolumes().List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list persistent volumes: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			infoList = append(infoList, &Info{
				Name:              item.Name,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,

				AccessModes:                   convertAccessModesToStringSlice(item.Spec.AccessModes),
				Capacity:                      convertResourceListToStringMap(item.Spec.Capacity),
				ClaimRefKind:                  item.Spec.ClaimRef.Kind,
				ClaimRefName:                  item.Spec.ClaimRef.Name,
				ClaimRefNamespace:             item.Spec.ClaimRef.Namespace,
				PersistentVolumeReclaimPolicy: string(item.Spec.PersistentVolumeReclaimPolicy),
				Phase:                         string(item.Status.Phase),
				StorageClassName:              item.Spec.StorageClassName,
				VolumeMode:                    (*string)(item.Spec.VolumeMode),
			})
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}

func convertResourceListToStringMap(list v2.ResourceList) map[string]int64 {
	capacity := make(map[string]int64, len(list))
	for k, v := range list {
		capacity[string(k)] = v.Value()
	}
	return capacity
}

func convertAccessModesToStringSlice(modes []v2.PersistentVolumeAccessMode) []string {
	modesSlice := make([]string, len(modes))
	for i, mode := range modes {
		modesSlice[i] = string(mode)
	}
	return modesSlice
}
