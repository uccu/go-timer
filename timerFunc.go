package timer

import "time"

type TimerFunc interface {
	GetGroups() []string               // 获取所有分组
	Set(string, interface{}) TimerFunc // 设置数据
	Get(string) interface{}            // 获取数据
	Run()                              // 运行函数
	GetTime() time.Time                // 获取运行的时间
	IsExpired() bool                   // 时间是否已经过期

	setTimer(Timer)
	getTimer() Timer

	getNext() TimerFunc
	setNext(TimerFunc)
}

type timerFunc struct {
	data     map[string]interface{}
	run      func()
	groupIds []string
	time     time.Time
	timer    Timer
	next     TimerFunc
}

func (f *timerFunc) GetGroups() []string {
	return f.groupIds
}

func (f *timerFunc) Set(key string, value interface{}) TimerFunc {
	if f.data == nil {
		f.data = map[string]interface{}{}
	}
	f.data[key] = value
	return f
}

func (f *timerFunc) setTimer(t Timer) {
	f.timer = t
}

func (f *timerFunc) getTimer() Timer {
	return f.timer
}

func (f *timerFunc) GetTime() time.Time {
	return f.time
}

func (f *timerFunc) getNext() TimerFunc {
	return f.next
}

func (f *timerFunc) setNext(next TimerFunc) {
	f.next = next
}

func (f *timerFunc) Get(key string) interface{} {
	if f.data == nil {
		return nil
	}
	value, ok := f.data[key]
	if !ok {
		return nil
	}
	return value
}

func (f *timerFunc) Run() {
	if f.run != nil {
		f.run()
	}
}

func (f *timerFunc) IsExpired() bool {
	return f.time.Before(time.Now())
}
