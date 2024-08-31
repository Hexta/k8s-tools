package nodeutil

import "time"

type NodeInfo struct {
	Name         string
	Age          time.Duration
	InstanceType string
	Utilisation  NodeUtilisation
}

type NodeUtilisation struct {
	CPU    float64
	Memory float64
}
