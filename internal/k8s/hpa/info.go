package hpa

import "time"

type Info struct {
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	CurrentReplicas   int32             `db:"current_replicas"`
	DesiredReplicas   int32             `db:"desired_replicas"`
}

type InfoList []*Info
