package fetch

import "time"

type Options struct {
	LabelSelector        string
	RetryInitialInterval time.Duration
	RetryJitterPercent   uint64
	RetryMaxAttempts     uint64
	RetryMaxInterval     time.Duration
}
