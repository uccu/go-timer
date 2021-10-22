package timer

func New() *Timer {
	return &Timer{
		group: map[string][]*TimerFunc{},
	}
}
