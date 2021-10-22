package timer_test

import (
	"fmt"
	"runtime/debug"
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
		fmt.Println(i)
		fmt.Println(string(debug.Stack()))
	})
	ti.Start()

	ti.AddTimerFunc(timer.NewTimerFunc(
		time.Now().Add(-1*time.Second),
		func() {
			var w *int
			fmt.Println(w)
		},
		"111", "222",
	))

	for {
		time.Sleep(time.Second)
	}

}
