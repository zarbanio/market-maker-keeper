package domain

import (
	"time"
)

type CycleStatus string

const (
	CycleStatusRunning CycleStatus = "running"
	CycleStatusSuccess CycleStatus = "success"
	CycleStatusFailed  CycleStatus = "failed"
)

type Cycle struct {
	ID     int64       `json:"id"`
	Start  time.Time   `json:"start"`
	End    time.Time   `json:"end"`
	Status CycleStatus `json:"status"`
}
