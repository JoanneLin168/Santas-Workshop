package util

import (
	"time"
)

type TimeHandler struct {
	Start time.Time
}

var startTime time.Time

func (t *TimeHandler) SetStartTime() {
	startTime = time.Now()
}

func (t *TimeHandler) GetStartTime() time.Time {
	return startTime
}

func (t *TimeHandler) GetTime() time.Duration {
	return time.Now().Sub(startTime)
}
