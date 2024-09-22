package ds

import "time"

type Info struct {
	Name                   string            `db:"name"`
	Namespace              string            `db:"namespace"`
	CreationTimestamp      time.Time         `db:"creation_ts"`
	Labels                 map[string]string `db:"labels"`
	CurrentNumberScheduled int32             `db:"current_number_scheduled"`
	DesiredNumberScheduled int32             `db:"desired_number_scheduled"`
	NumberAvailable        int32             `db:"number_available"`
	NumberMisscheduled     int32             `db:"number_misscheduled"`
	NumberReady            int32             `db:"number_ready"`
	NumberUnavailable      int32             `db:"number_unavailable"`
}

type InfoList []*Info
