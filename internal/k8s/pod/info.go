package pod

import (
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
)

type Info struct {
	Name              string `db:"name"`
	Namespace         string `db:"namespace"`
	NodeName          string `db:"node_name"`
	Containers        container.InfoList
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	CPURequests       float64           `db:"cpu_requests"`
	CPULimits         float64           `db:"cpu_limits"`
	MemoryRequests    float64           `db:"memory_requests"`
	MemoryLimits      float64           `db:"memory_limits"`
	IP                string            `db:"ip"`
}

type InfoList []*Info
