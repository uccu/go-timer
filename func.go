package timer

import "time"

func New() Timer {
	return &timer{
		group: map[string][]TimerFunc{},
		state: CLOSED,
	}
}

func NewTimerFunc(t time.Time, r func(), g ...string) TimerFunc {
	return &timerFunc{
		run:      r,
		groupIds: g,
		time:     t,
	}
}
