package ds

import "time"

type Info struct {
	Name                   string
	Namespace              string
	CreationTimestamp      time.Time
	Labels                 map[string]string
	CurrentNumberScheduled int32
	DesiredNumberScheduled int32
	NumberAvailable        int32
	NumberMisscheduled     int32
	NumberReady            int32
	NumberUnavailable      int32
}

type InfoList []*Info
