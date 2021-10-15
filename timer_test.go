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

	ti := timer.New()

	ti.SetErrhandler(func(i interface{}) {
		fmt.Print(i)
	})
	ti.Start()

	ti.AddTimerFunc(timer.TimerFunc{
		GroupId: 123,
		Unix:    time.Now().Add(1 * time.Second).Unix(),
		Run: func() {
			var w *int
			fmt.Println(*w)
		},
	})

	for {
		time.Sleep(time.Second)
	}

}
