package nodeutil

import "time"

type NodeInfo struct {
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
