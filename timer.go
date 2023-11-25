package timer

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Timer interface {
	Start()                          // 运行定时器
	Stop()                           // 停止定时器
	DelGroup(string) []TimerFunc     // 删除分组
	SetErrhandler(func(interface{})) // 设置错误处理函数
	AddTimerFunc(TimerFunc)          // 添加运行函数
	GetCountData() *CountData        // 获取计数数据

	run(TimerFunc)
	delUnixTimerFunc(TimerFunc)
	delGroupTimerFunc(TimerFunc)
	delTimerFunc(TimerFunc)
	exec()
}

type timer struct {
	state      State
	group      map[string][]TimerFunc
	next       TimerFunc
	mu         sync.RWMutex
	ctxCancel  context.CancelFunc
	errHandler func(interface{})
}

func (t *timer) run(f TimerFunc) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				if t.errHandler != nil {
					t.errHandler(err)
				}
			}
		}()
		f.Run()
	}()
}

type CountData struct {
	GroupCount     int
	GroupMapCount  map[string]int
	TimerFuncCount int
}

func (t *timer) GetCountData() *CountData {

	t.mu.RLock()
	defer t.mu.RUnlock()

	fmt.Println("GetCountData", time.Now().Format("2006-01-02 15:04:05.999999"))

	timerFuncCount := 0
	next := t.next
	for next != nil {
		timerFuncCount++
		next = next.getNext()
	}

	groupMapCount := make(map[string]int)
	for k, v := range t.group {
		groupMapCount[k] = len(v)
	}
	return &CountData{
		GroupCount:     len(t.group),
		GroupMapCount:  groupMapCount,
		TimerFuncCount: timerFuncCount,
	}
}

func (t *timer) AddTimerFunc(f TimerFunc) {

	t.mu.Lock()
	defer t.mu.Unlock()

	f.setTimer(t)

	if !f.GetTime().After(time.Now()) {
		t.run(f)
		return
	}

	var prev TimerFunc
	next := t.next

	for {
		if next == nil || next.GetTime().After(f.GetTime()) {
			f.setNext(next)
			if prev != nil {
				prev.setNext(f)
			} else {
				t.next = f
			}
			break
		}
		prev = next
		next = next.getNext()
	}

	for _, g := range f.GetGroups() {
		if timers, ok := t.group[g]; ok {
			t.group[g] = append(timers, f)
		} else {
			t.group[g] = []TimerFunc{f}
		}
	}

}
func (t *timer) delUnixTimerFunc(f TimerFunc) {

	var prev TimerFunc
	next := t.next
	for {
		if next == f {
			if prev == nil {
				t.next = next.getNext()
			} else {
				prev.setNext(next.getNext())
			}
			break
		}
		prev = next
		next = next.getNext()
	}
}

func (t *timer) delGroupTimerFunc(f TimerFunc) {

	for _, g := range f.GetGroups() {
		fns, ok := t.group[g]
		if !ok {
			continue
		}

		for k, t := range fns {
			if t == f {
				fns = append(fns[:k], fns[k+1:]...)
				break
			}
		}

		t.group[g] = fns
		if len(fns) == 0 {
			delete(t.group, g)
		}
	}
}

func (t *timer) delTimerFunc(f TimerFunc) {
	t.delUnixTimerFunc(f)
	t.delGroupTimerFunc(f)
}

func (t *timer) DelGroup(groupId string) []TimerFunc {
	t.mu.Lock()
	defer t.mu.Unlock()
	fns, ok := t.group[groupId]
	if !ok {
		return nil
	}

	for _, f := range fns {
		t.delTimerFunc(f)
	}
	return fns
}

func (t *timer) Start() {

	t.mu.Lock()
	defer t.mu.Unlock()

	ticker := time.NewTicker(time.Second)
	ctx, cancel := context.WithCancel(context.TODO())
	t.state = OPEN

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go t.exec()
			}
		}
	}()

	t.ctxCancel = cancel
}

func (t *timer) SetErrhandler(i func(interface{})) {
	t.errHandler = i
}

func (t *timer) Stop() {

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.state != OPEN {
		t.state = CLOSED
		t.ctxCancel()
	}

}

func (t *timer) exec() {

	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Println("exec", time.Now().Format("2006-01-02 15:04:05.999999"))

	if t.state == CLOSED {
		return
	}

	unix := time.Now()
	f := t.next
	for {
		if f == nil {
			break
		}

		if !f.GetTime().After(unix) {
			t.run(f)
			t.delTimerFunc(f)
		} else {
			break
		}

		f = f.getNext()
	}

}
