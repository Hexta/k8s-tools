package deployment

import "time"

type Info struct {
	Name              string
	Namespace         string
	CreationTimestamp time.Time
	Labels            map[string]string
}

type InfoList []*Info
