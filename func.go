package timer

func New() *Timer {
	return &Timer{
		group: map[interface{}][]TimerFunc{},
		uinx:  map[int64][]TimerFunc{},
	}
}
