package entity

import (
	"sync"
	"time"
)

type RequestStatus map[string]int

type Report struct {
	Mutex          sync.Mutex
	TimeStarted    time.Time
	TotalTime      time.Duration
	RequestsCount  uint64
	SuccesssCount  uint64
	RequestsStatus RequestStatus
}

func NewReport() *Report {
	return &Report{
		TimeStarted:    time.Now(),
		RequestsCount:  0,
		SuccesssCount:  0,
		RequestsStatus: make(RequestStatus),
	}
}

func (r *Report) ValidateTotalTime() {
	r.TotalTime = time.Since(r.TimeStarted)
}
