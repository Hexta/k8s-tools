package pvc

import "time"

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`

	AccessModes                      []string          `db:"access_modes"`
	AllocatedResourceStatuses        map[string]string `db:"allocated_resource_statuses"`
	AllocatedResources               map[string]int64  `db:"allocated_resources"`
	Capacity                         map[string]int64  `db:"capacity"`
	CurrentVolumeAttributesClassName *string           `db:"current_volume_attributes_class_name"`
	DesiredAccessModes               []string          `db:"desired_access_modes"`
	ModifyVolumeStatus               string            `db:"modify_volume_status"`
	Phase                            string            `db:"phase"`
	ResourceLimits                   map[string]int64  `db:"resource_limits"`
	ResourceRequests                 map[string]int64  `db:"resource_requests"`
	StorageClassName                 *string           `db:"storage_class_name"`
	TargetVolumeAttributesClassName  string            `db:"target_volume_attributes_class_name"`
	VolumeAttributesClassName        *string           `db:"volume_attributes_class_name"`
	VolumeMode                       *string           `db:"volume_mode"`
	VolumeName                       string            `db:"volume_name"`
}

type InfoList []*Info
