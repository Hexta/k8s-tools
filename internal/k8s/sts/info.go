package sts

import "time"

type Info struct {
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Replicas          int32             `db:"replicas"`
}

type InfoList []*Info
