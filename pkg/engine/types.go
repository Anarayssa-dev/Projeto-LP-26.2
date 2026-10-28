package engine

import "time"

type PingResult struct {
	Target    string
	RTT       time.Duration
	Success   bool
	Timestamp time.Time
}

type MetricsStorage interface {
	Save(result PingResult) error
	GetAll(target string) ([]PingResult, error)
	Close() error
}