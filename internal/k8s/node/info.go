package node

import (
	"time"
)

type Info struct {
	Name              string
	Age               time.Duration
	CreationTimestamp time.Time
	InstanceType      string
	Utilisation       NodeUtilisation
	Labels            map[string]string
}

type NodeUtilisation struct {
	CPU    float64
	Memory float64
}

type InfoList []*Info
