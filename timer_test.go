package timer_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/uccu/go-timer"
)

type Gr struct {
	groupId int64
	uinx    int64
	run     func()
}

func (t *Gr) GetGroupId() interface{} {
	return t.groupId
}

func (t *Gr) GetUnix() int64 {
	return t.uinx
}

func (t *Gr) Run() {
	t.run()
}

func TestA(t *testing.T) {

	timer := timer.New()
	timer.Start()

	timerFunc := &Gr{}
	timerFunc.groupId = 123
	timerFunc.uinx = time.Now().Add(2 * time.Second).Unix()
	timerFunc.run = func() {
		fmt.Println(111)
	}
	timer.AddTimerFunc(timerFunc)
	timer.Stop()

}
