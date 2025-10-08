package node

import (
	"time"
)

type Info struct {
	Name   string            `db:"name"`
	Labels map[string]string `db:"labels"`

	Address                     map[string]string `db:"address"`
	Age                         time.Duration
	AllocatableCPU              float64           `db:"allocatable_cpu"`
	AllocatableEphemeralStorage int64             `db:"allocatable_ephemeral_storage"`
	AllocatableMemory           float64           `db:"allocatable_memory"`
	AllocatablePods             int64             `db:"allocatable_pods"`
	Annotations                 map[string]string `db:"annotations"`
	Architecture                string            `db:"architecture"`
	CapacityCPU                 float64           `db:"capacity_cpu"`
	CapacityMemory              float64           `db:"capacity_memory"`
	ContainerRuntimeVersion     string            `db:"container_runtime_version"`
	CPUUtilisation              float64           `db:"cpu_utilisation"`
	CreationTimestamp           time.Time         `db:"creation_ts"`
	Images                      ImageList
	InstanceType                string  `db:"instance_type"`
	KernelVersion               string  `db:"kernel_version"`
	KubeletVersion              string  `db:"kubelet_version"`
	MemoryUtilisation           float64 `db:"memory_utilisation"`
	OperatingSystem             string  `db:"operating_system"`
	OSImage                     string  `db:"os_image"`
	Taints                      TaintList
}

type InfoList []*Info

type Utilisation struct {
	CPU    float64
	Memory float64
}

type Taint struct {
	Effect string
	Key    string
	Value  string
}

type TaintList []*Taint

type Image struct {
	Names []string `db:"names"`
	Node  string   `db:"node_name"`
	Size  int      `db:"size"`
}

type ImageList []*Image
