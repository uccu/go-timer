package timer

import (
	"context"
	"sync"
	"time"
)

type State int

const (
	OPEN State = iota
	CLOSED
)

type TimerFunc struct {
	run      func()
	groupIds []string
	time     time.Time
	timer    *Timer
	state    State
	next     *TimerFunc
}

func (f *TimerFunc) Delete() {
	if f.state == CLOSED {
		return
	}

	f.state = CLOSED
	f.timer.delTimerFunc(f)
}

func (f *TimerFunc) Run() {
	if f.state == CLOSED {
		return
	}
	f.state = CLOSED
	f.run()
}

func NewTimerFunc(t time.Time, r func(), g ...string) *TimerFunc {
	return &TimerFunc{
		run:      r,
		groupIds: g,
		time:     t,
	}
}

type Timer struct {
	group      map[string][]*TimerFunc
	unix       *TimerFunc
	mu         sync.RWMutex
	ctxCancel  context.CancelFunc
	errHandler func(err interface{})
}

func (c *Timer) run(t *TimerFunc) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if c.errHandler != nil {
					c.errHandler(err)
				}
			}
		}()
		t.Run()
	}()

}

func (c *Timer) AddTimerFunc(t *TimerFunc) {

	c.mu.Lock()
	defer c.mu.Unlock()

	t.timer = c

	if t.time.Unix() <= time.Now().Unix() {
		c.run(t)
		return
	}

	var prev *TimerFunc
	next := c.unix

	for {
		if next == nil || next.time.After(t.time) {
			t.next = next
			if prev != nil {
				prev.next = t
			} else {
				c.unix = t
			}
			break
		}
		prev = next
		next = next.next
	}

	for _, g := range t.groupIds {
		if timers, ok := c.group[g]; ok {
			c.group[g] = append(timers, t)
		} else {
			c.group[g] = []*TimerFunc{t}
		}
	}

}
func (c *Timer) delUnixTimerFunc(f *TimerFunc) {

	c.mu.Lock()
	defer c.mu.Unlock()

	var prev *TimerFunc
	next := c.unix
	for {
		if next == f {
			if prev == nil {
				c.unix = next.next
			} else {
				prev.next = next.next
			}
			break
		}
		prev = next
		next = next.next
	}
}

func (c *Timer) delGroupTimerFunc(f *TimerFunc) {

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, g := range f.groupIds {
		timers, ok := c.group[g]
		if !ok {
			return
		}

		for k, t := range timers {
			if t == f {
				timers = append(timers[:k], timers[k+1:]...)
				break
			}
		}

		c.group[g] = timers
		if len(timers) == 0 {
			delete(c.group, g)
		}

	}

	f.groupIds = []string{}
}

func (c *Timer) delTimerFunc(f *TimerFunc) {

	c.delUnixTimerFunc(f)

	c.delGroupTimerFunc(f)

}

func (c *Timer) DelGroup(groupId string) []*TimerFunc {
	c.mu.RLock()
	timers, ok := c.group[groupId]
	timersCopy := make([]*TimerFunc, len(timers))
	copy(timersCopy, timers)
	c.mu.RUnlock()
	if !ok {
		return nil
	}

	for _, t := range timersCopy {
		t.Delete()
	}
	return timersCopy
}

func (c *Timer) Start() {

	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for {
			<-ticker.C
			if !c.exec(ctx) {
				ticker.Stop()
				break
			}
		}
	}()

	c.ctxCancel = cancel
}

func (c *Timer) SetErrhandler(i func(interface{})) {
	c.errHandler = i
}

func (c *Timer) Stop() {
	c.ctxCancel()
}

func (c *Timer) exec(ctx context.Context) bool {

	select {
	case <-ctx.Done():
		return false
	default:
	}

	unix := time.Now()
	timer := c.unix
	for {
		if timer == nil {
			break
		}

		if !timer.time.After(unix) {
			c.run(timer)
			c.delTimerFunc(timer)
		} else {
			break
		}

		timer = timer.next
	}

	return true
}
