package pv

import (
	"context"
	"fmt"

	"github.com/Hexta/k8s-tools/internal/k8sutil"
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

				AccessModes:                   k8sutil.ConvertPVAccessModesToStrings(item.Spec.AccessModes),
				Capacity:                      k8sutil.ConvertResourceListToStringMap(item.Spec.Capacity),
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
