package event

import (
	"time"
)

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`

	Action         string    `db:"action"`
	ApiVersion     string    `db:"api_version"`
	Count          int32     `db:"count"`
	EventTime      time.Time `db:"event_time"`
	FirstTimestamp time.Time `db:"first_timestamp"`

	InvolvedObjectAPIVersion      string `db:"involved_object_api_version"`
	InvolvedObjectFieldPath       string `db:"involved_object_field_path"`
	InvolvedObjectKind            string `db:"involved_object_kind"`
	InvolvedObjectName            string `db:"involved_object_name"`
	InvolvedObjectNamespace       string `db:"involved_object_namespace"`
	InvolvedObjectResourceVersion string `db:"involved_object_resource_version"`
	InvolvedObjectUID             string `db:"involved_object_uid"`

	Kind          string    `db:"kind"`
	LastTimestamp time.Time `db:"last_timestamp"`
	Message       string    `db:"message"`
	Reason        string    `db:"reason"`

	RelatedObjectAPIVersion      string `db:"related_object_api_version"`
	RelatedObjectFieldPath       string `db:"related_object_field_path"`
	RelatedObjectKind            string `db:"related_object_kind"`
	RelatedObjectName            string `db:"related_object_name"`
	RelatedObjectNamespace       string `db:"related_object_namespace"`
	RelatedObjectResourceVersion string `db:"related_object_resource_version"`
	RelatedObjectUID             string `db:"related_object_uid"`

	ReportingComponent string `db:"reporting_component"`
	ReportingInstance  string `db:"reporting_instance"`

	SeriesCount            int32     `db:"series_count"`
	SeriesLastObservedTime time.Time `db:"series_last_observed_time"`

	SourceComponent string `db:"source_component"`
	SourceHost      string `db:"source_host"`

	Type string `db:"type"`
}

type InfoList []*Info
