package service

import (
	"time"
)

var startedAt = time.Now()

type Health struct {
	Uptime    time.Duration `json:"uptime"`
	StartedAt time.Time     `json:"started_at"`
	Status    string        `json:"status"`
}

func GetHealth() *Health {
	return &Health{
		Uptime:    time.Since(startedAt),
		StartedAt: startedAt,
		Status:    "OK",
	}
}
