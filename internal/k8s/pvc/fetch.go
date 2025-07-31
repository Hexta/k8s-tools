package pvc

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
		list, err := clientset.CoreV1().PersistentVolumeClaims(v1.NamespaceAll).List(ctx, v1.ListOptions{
			Continue: continueToken,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to list persistent volumes: %v", err)
		}

		for idx := range list.Items {
			item := &list.Items[idx]
			info := &Info{
				Name:              item.Name,
				Namespace:         item.Namespace,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,

				AccessModes:                      k8sutil.ConvertPVAccessModesToStrings(item.Status.AccessModes),
				AllocatedResourceStatuses:        k8sutil.ConvertResourceStatusesToStringMap(item.Status.AllocatedResourceStatuses),
				AllocatedResources:               k8sutil.ConvertResourceListToStringMap(item.Status.Capacity),
				Capacity:                         k8sutil.ConvertResourceListToStringMap(item.Status.Capacity),
				CurrentVolumeAttributesClassName: item.Status.CurrentVolumeAttributesClassName,
				DesiredAccessModes:               k8sutil.ConvertPVAccessModesToStrings(item.Spec.AccessModes),
				Phase:                            string(item.Status.Phase),
				ResourceLimits:                   k8sutil.ConvertResourceListToStringMap(item.Spec.Resources.Limits),
				ResourceRequests:                 k8sutil.ConvertResourceListToStringMap(item.Spec.Resources.Requests),
				StorageClassName:                 item.Spec.StorageClassName,
				VolumeAttributesClassName:        item.Spec.VolumeAttributesClassName,
				VolumeMode:                       (*string)(item.Spec.VolumeMode),
				VolumeName:                       item.Spec.VolumeName,
			}

			if item.Status.ModifyVolumeStatus != nil {
				info.ModifyVolumeStatus = (string)(item.Status.ModifyVolumeStatus.Status)
				info.TargetVolumeAttributesClassName = item.Status.ModifyVolumeStatus.TargetVolumeAttributesClassName
			}

			infoList = append(infoList, info)
		}

		if continueToken = list.Continue; continueToken == "" {
			break
		}
	}

	return infoList, nil
}
