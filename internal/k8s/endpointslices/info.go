package endpointslices

import (
	"time"
)

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`

	AddressType string `db:"address_type"`
	ApiVersion  string `db:"api_version"`
	Endpoints   EndpointList
	Kind        string `db:"kind"`
	Ports       PortList
}

type InfoList []*Info

type Endpoint struct {
	EndpointSliceName string `db:"endpoint_slice_name"`
	Namespace         string `db:"namespace"`

	Addresses                []string `db:"addresses"`
	Hostname                 *string  `db:"hostname"`
	NodeName                 *string  `db:"node_name"`
	HasTargetRef             bool     `db:"has_target_ref"`
	TargetRefAPIVersion      string   `db:"target_ref_api_version"`
	TargetRefFieldPath       string   `db:"target_ref_field_path"`
	TargetRefKind            string   `db:"target_ref_kind"`
	TargetRefName            string   `db:"target_ref_name"`
	TargetRefNamespace       string   `db:"target_ref_namespace"`
	TargetRefResourceVersion string   `db:"target_ref_resource_version"`
	TargetRefUID             string   `db:"target_ref_uid"`
}

type EndpointList []*Endpoint

type Port struct {
	EndpointSliceName string  `db:"endpoint_slice_name"`
	Name              *string `db:"name"`
	Namespace         string  `db:"namespace"`

	AppProtocol *string `db:"app_protocol"`
	Port        *int32  `db:"port"`
	Protocol    *string `db:"protocol"`
}

type PortList []*Port
