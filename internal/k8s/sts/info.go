package sts

import "time"

type Info struct {
	Name              string
	Namespace         string
	CreationTimestamp time.Time
	Labels            map[string]string
	Replicas          int32
}

type InfoList []*Info
