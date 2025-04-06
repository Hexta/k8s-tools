package endpoints

import (
	"time"
)

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`
	Subsets           SubsetList
}

type InfoList []*Info

type Subset struct {
	EndpointsName            string  `db:"endpoints_name"`
	HasTargetRef             bool    `db:"has_target_ref"`
	Hostname                 string  `db:"hostname"`
	IP                       string  `db:"ip"`
	IsReady                  bool    `db:"is_ready"`
	Namespace                string  `db:"namespace"`
	NodeName                 *string `db:"node_name"`
	Ports                    []int32 `db:"ports"`
	TargetRefAPIVersion      string  `db:"target_ref_api_version"`
	TargetRefFieldPath       string  `db:"target_ref_field_path"`
	TargetRefKind            string  `db:"target_ref_kind"`
	TargetRefName            string  `db:"target_ref_name"`
	TargetRefNamespace       string  `db:"target_ref_namespace"`
	TargetRefResourceVersion string  `db:"target_ref_resource_version"`
	TargetRefUID             string  `db:"target_ref_uid"`
}

type SubsetList []*Subset
