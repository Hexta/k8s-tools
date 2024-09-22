package node

import (
	"time"
)

type Info struct {
	Name              string `db:"name"`
	Age               time.Duration
	CreationTimestamp time.Time         `db:"creation_ts"`
	InstanceType      string            `db:"instance_type"`
	CPUUtilisation    float64           `db:"cpu_utilisation"`
	MemoryUtilisation float64           `db:"memory_utilisation"`
	Labels            map[string]string `db:"labels"`
	Address           map[string]string `db:"address"`
}

type InfoList []*Info

type Utilisation struct {
	CPU    float64
	Memory float64
}
