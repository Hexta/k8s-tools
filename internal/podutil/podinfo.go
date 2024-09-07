package podutil

import "time"

type PodInfo struct {
	Name              string
	Namespace         string
	NodeName          string
	CreationTimestamp time.Time
	Labels            map[string]string
	CPURequests       float64
	CPULimits         float64
	MemoryRequests    float64
	MemoryLimits      float64
}

type PodInfoList []*PodInfo
