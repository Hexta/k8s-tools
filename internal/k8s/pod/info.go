package pod

import (
	"time"

	"github.com/Hexta/k8s-tools/internal/k8s/container"
	"github.com/Hexta/k8s-tools/internal/k8s/volumesource"
	apicorev1 "k8s.io/api/core/v1"
)

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`

	Affinity                      *apicorev1.Affinity
	AutomountServiceAccountToken  *bool   `db:"automount_service_account_token"`
	CPULimits                     float64 `db:"cpu_limits"`
	CPURequests                   float64 `db:"cpu_requests"`
	Containers                    container.InfoList
	DNSConfig                     *apicorev1.PodDNSConfig
	DNSPolicy                     apicorev1.DNSPolicy `db:"dns_policy"`
	EnableServiceLinks            *bool               `db:"enable_service_links"`
	HostIPC                       bool                `db:"host_ipc"`
	HostNetwork                   bool                `db:"host_network"`
	HostPID                       bool                `db:"host_pid"`
	HostUsers                     *bool               `db:"host_users"`
	Hostname                      string              `db:"hostname"`
	IP                            string              `db:"ip"`
	InitContainers                container.InfoList
	MemoryLimits                  float64           `db:"memory_limits"`
	MemoryRequests                float64           `db:"memory_requests"`
	NodeName                      string            `db:"node_name"`
	NodeSelector                  map[string]string `db:"node_selector"`
	PreemptionPolicy              *apicorev1.PreemptionPolicy
	Priority                      *int32 `db:"priority"`
	PriorityClassName             string `db:"priority_class_name"`
	Tolerations                   TolerationList
	RestartPolicy                 string  `db:"restart_policy"`
	RuntimeClassName              *string `db:"runtime_class_name"`
	SchedulerName                 string  `db:"scheduler_name"`
	ServiceAccountName            string  `db:"service_account_name"`
	SetHostnameAsFQDN             *bool   `db:"set_hostname_as_fqdn"`
	ShareProcessNamespace         *bool   `db:"share_process_namespace"`
	Subdomain                     string  `db:"subdomain"`
	TerminationGracePeriodSeconds *int64  `db:"termination_grace_period_seconds"`
	Volumes                       VolumeList
}

type InfoList []*Info

type Toleration struct {
	Effect            string
	Key               string
	Operator          string
	TolerationSeconds *int64
	Value             string
}

type TolerationList []*Toleration

type Volume struct {
	Name         string
	VolumeSource volumesource.VolumeSource
}

type VolumeList []*Volume
