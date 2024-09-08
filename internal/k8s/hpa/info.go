package hpa

import "time"

type Info struct {
	Name              string
	Namespace         string
	CreationTimestamp time.Time
	Labels            map[string]string
	CurrentReplicas   int32
	DesiredReplicas   int32
}

type InfoList []*Info
