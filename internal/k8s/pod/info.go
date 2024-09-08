package pod

import (
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
)

type Info struct {
	Name              string
	Namespace         string
	NodeName          string
	Containers        container.InfoList
	CreationTimestamp time.Time
	Labels            map[string]string
	CPURequests       float64
	CPULimits         float64
	MemoryRequests    float64
	MemoryLimits      float64
	IP                string
}

type InfoList []*Info
