package event

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
		l, err := clientset.CoreV1().Events(v1.NamespaceAll).List(ctx, opts)
		if err != nil {
			return "", fmt.Errorf("failed to list events: %v", err)
		}
		for idx := range l.Items {
			item := &l.Items[idx]
			info := &Info{
				Annotations:       item.Annotations,
				CreationTimestamp: item.CreationTimestamp.Time,
				Labels:            item.Labels,
				Name:              item.Name,
				Namespace:         item.Namespace,

				Action:         item.Action,
				ApiVersion:     item.APIVersion,
				Count:          item.Count,
				EventTime:      item.EventTime.Time,
				FirstTimestamp: item.FirstTimestamp.Time,

				InvolvedObjectAPIVersion:      item.InvolvedObject.APIVersion,
				InvolvedObjectFieldPath:       item.InvolvedObject.FieldPath,
				InvolvedObjectKind:            item.InvolvedObject.Kind,
				InvolvedObjectName:            item.InvolvedObject.Name,
				InvolvedObjectNamespace:       item.InvolvedObject.Namespace,
				InvolvedObjectResourceVersion: item.InvolvedObject.ResourceVersion,
				InvolvedObjectUID:             string(item.InvolvedObject.UID),

				Kind:          item.Kind,
				LastTimestamp: item.LastTimestamp.Time,
				Message:       item.Message,
				Reason:        item.Reason,

				ReportingComponent: item.ReportingController,
				ReportingInstance:  item.ReportingInstance,

				SourceComponent: item.Source.Component,
				SourceHost:      item.Source.Host,

				Type: item.Type,
			}

			if item.Related != nil {
				info.RelatedObjectAPIVersion = item.Related.APIVersion
				info.RelatedObjectFieldPath = item.Related.FieldPath
				info.RelatedObjectKind = item.Related.Kind
				info.RelatedObjectName = item.Related.Name
				info.RelatedObjectNamespace = item.Related.Namespace
				info.RelatedObjectResourceVersion = item.Related.ResourceVersion
				info.RelatedObjectUID = string(item.Related.UID)
			}

			if item.Series != nil {
				info.SeriesCount = item.Series.Count
				info.SeriesLastObservedTime = item.Series.LastObservedTime.Time
			}

			infoList = append(infoList, info)
		}
		return l.Continue, nil
	})

	return infoList, err
}
