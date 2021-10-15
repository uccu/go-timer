package timer

import (
	"context"
	"sync"
	"time"
)

type TimerFunc struct {
	Run     func()
	GroupId interface{}
	Unix    int64
}

type Timer struct {
	group      map[interface{}][]TimerFunc
	uinx       map[int64][]TimerFunc
	mu         sync.RWMutex
	ctxCancel  context.CancelFunc
	errHandler func(err interface{})
}

func (c *Timer) run(t TimerFunc) {
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

func (c *Timer) AddTimerFunc(t TimerFunc) {

	if t.Run == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	unix := t.Unix
	if unix <= time.Now().Unix() {
		c.run(t)
		return
	}

	if timers, ok := c.uinx[unix]; ok {
		c.uinx[unix] = append(timers, t)
	} else {
		c.uinx[unix] = []TimerFunc{t}
	}

	groupId := t.GroupId
	if groupId != nil {
		if timers, ok := c.group[groupId]; ok {
			c.group[groupId] = append(timers, t)
		} else {
			c.group[groupId] = []TimerFunc{t}
		}
	}
}

func (c *Timer) DelGroup(groupId interface{}) {

	c.mu.Lock()
	defer c.mu.Unlock()

	timers, ok := c.group[groupId]
	if !ok {
		return
	}

	for _, timer := range timers {

		unix := timer.Unix
		timers, ok := c.uinx[unix]
		if !ok {
			continue
		}

		ntimers := []TimerFunc{}
		for _, timer := range timers {
			if timer.GroupId != groupId {
				ntimers = append(ntimers, timer)
			}
		}

		if len(ntimers) == 0 {
			delete(c.uinx, unix)
			continue
		}

		c.uinx[unix] = ntimers
	}

	delete(c.group, groupId)

}

func (c *Timer) Start() {

	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		for {
			time.Sleep(time.Millisecond * 100)
			m := time.Now().UnixMilli()
			if m/100%10 == 5 {
				break
			}
		}
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

	unix := time.Now().Unix()

	c.mu.Lock()
	defer c.mu.Unlock()

	timers, ok := c.uinx[unix]
	if !ok {
		return true
	}

	for _, timer := range timers {
		c.run(timer)
		groupId := timer.GroupId
		timers, ok := c.group[groupId]
		if !ok {
			continue
		}

		ntimers := []TimerFunc{}
		for _, timer := range timers {
			if timer.Unix != unix {
				ntimers = append(ntimers, timer)
			}
		}

		if len(ntimers) == 0 {
			delete(c.group, groupId)
			continue
		}

		c.group[groupId] = ntimers
	}

	delete(c.uinx, unix)
	return true

}
