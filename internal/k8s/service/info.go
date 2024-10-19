package service

import (
	"time"

	apicorev1 "k8s.io/api/core/v1"
)

type Info struct {
	Annotations       map[string]string `db:"annotations"`
	CreationTimestamp time.Time         `db:"creation_ts"`
	Labels            map[string]string `db:"labels"`
	Name              string            `db:"name"`
	Namespace         string            `db:"namespace"`

	ClusterIPs               []string                               `db:"cluster_ips"`
	ClusterIP                string                                 `db:"cluster_ip"`
	ExternalIPs              []string                               `db:"external_ips"`
	ExternalTrafficPolicy    apicorev1.ServiceExternalTrafficPolicy `db:"external_traffic_policy"`
	HealthCheckNodePort      int32                                  `db:"health_check_node_port"`
	IPFamilies               []apicorev1.IPFamily
	IPFamilyPolicy           *apicorev1.IPFamilyPolicy
	InternalTrafficPolicy    *apicorev1.ServiceInternalTrafficPolicy
	LoadBalancerClass        *string  `db:"load_balancer_class"`
	LoadBalancerSourceRanges []string `db:"load_balancer_source_ranges"`
	Ports                    []apicorev1.ServicePort
	PublishNotReadyAddresses bool                      `db:"publish_not_ready_addresses"`
	Reason                   string                    `db:"reason"`
	Selector                 map[string]string         `db:"selector"`
	SessionAffinity          apicorev1.ServiceAffinity `db:"session_affinity"`
	Type                     apicorev1.ServiceType     `db:"type"`
}

type InfoList []*Info
