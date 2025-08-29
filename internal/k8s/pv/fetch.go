package pv

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
		list, err := clientset.CoreV1().PersistentVolumes().List(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list persistent volumes: %v", err)
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
		return list.Continue, nil
	})

	return infoList, err
}
