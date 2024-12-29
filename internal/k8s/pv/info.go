package pv

import "time"

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`

	AccessModes                   []string         `db:"access_modes"`
	Capacity                      map[string]int64 `db:"capacity"`
	ClaimRefKind                  string           `db:"claim_ref_kind"`
	ClaimRefName                  string           `db:"claim_ref_name"`
	ClaimRefNamespace             string           `db:"claim_ref_namespace"`
	PersistentVolumeReclaimPolicy string           `db:"persistent_volume_reclaim_policy"`
	Phase                         string           `db:"phase"`
	StorageClassName              string           `db:"storage_class_name"`
	VolumeMode                    *string          `db:"volume_mode"`
}

type InfoList []*Info
