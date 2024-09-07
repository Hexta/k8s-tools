package containerutil

type ContainerInfo struct {
	Name           string
	Namespace      string
	PodName        string
	CPURequests    float64
	CPULimits      float64
	MemoryRequests float64
	MemoryLimits   float64
}
