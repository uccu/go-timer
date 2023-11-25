package timer_test

import (
	"fmt"
	"log"
	"runtime/debug"
	"testing"
	"time"

	"github.com/uccu/go-timer"
)

func TestA(t *testing.T) {

	ti := timer.New()
	ti.SetErrhandler(func(i interface{}) {
		log.Printf("[ERROR] error: %s, stack: %s", fmt.Sprint(i), string(debug.Stack()))
	})
	ti.Start()

	now := time.Now()

	for i := 0; i < 100; i++ {

		ti.AddTimerFunc(timer.NewTimerFunc(now.Add(2*time.Second), func() {
			fmt.Println("timer2", time.Since(now))
		}, "activity", "activity:2"))
		ti.AddTimerFunc(timer.NewTimerFunc(now.Add(1*time.Second), func() {
			fmt.Println("timer1", time.Since(now))
		}, "activity", "activity:1"))
		ti.AddTimerFunc(timer.NewTimerFunc(now.Add(-1*time.Second), func() {
			fmt.Println("timer-1", time.Since(now))
		}, "activity", "activity:-1"))
		ti.AddTimerFunc(timer.NewTimerFunc(now, func() {
			fmt.Println("timer0", time.Since(now))
		}, "activity", "activity:0"))
		ti.AddTimerFunc(timer.NewTimerFunc(now.Add(3*time.Second), func() {
			fmt.Println("timer3", time.Since(now))
		}, "activity", "activity:3"))
		ti.AddTimerFunc(timer.NewTimerFunc(now.Add(4*time.Second), func() {
			fmt.Println("timer4", time.Since(now))
		}, "activity", "activity:4"))
	}

	for {
		fmt.Println(*ti.GetCountData())
		time.Sleep(time.Second * 1)
	}

}
