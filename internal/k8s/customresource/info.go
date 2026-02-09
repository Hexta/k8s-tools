package customresource

import "time"

type Info struct {
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Data              string            `db:"data"`
	ResourceGroup     string            `db:"resource_group"`
	ResourceKind      string            `db:"resource_kind"`
	ResourceVersion   string            `db:"resource_version"`
}

type InfoList []*Info
