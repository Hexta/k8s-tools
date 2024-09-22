package container

type Info struct {
	Name           string  `db:"name"`
	Namespace      string  `db:"namespace"`
	PodName        string  `db:"pod_name"`
	CPURequests    float64 `db:"cpu_requests"`
	CPULimits      float64 `db:"cpu_limits"`
	MemoryRequests float64 `db:"memory_requests"`
	MemoryLimits   float64 `db:"memory_limits"`
}

type InfoList []*Info
