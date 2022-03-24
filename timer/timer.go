package timer

import (
	// "fmt"
	"time"
)

const _pollRate = 20 * time.Millisecond
const _doorOpenTime = 3 * time.Second
const _waitForDecisionTime = 250 * time.Millisecond
const _forceActionTime = 1 * time.Second

type Timer struct {
	isActive      bool
	endTime       time.Time
	timerDuration time.Duration
}

var DoorTimer = Timer{timerDuration: _doorOpenTime}
var DecisionTimer = Timer{timerDuration: _waitForDecisionTime}

// var NewOrderDecisionTimer = Timer{timerDuration: _waitForDecisionTime}
var ForceActionTimer = Timer{timerDuration: _forceActionTime, isActive: true}

func (timer *Timer) PollTimerOut(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := timer.TimedOut()
		if v != prev && v {
			receiver <- v
		}
		prev = v
		// fmt.Println(timer.timerDuration)
	}
}

func (timer *Timer) TimerStart() {
	timer.endTime = time.Now().Add(timer.timerDuration)
	timer.isActive = true
}

func (timer *Timer) TimerStop() {
	timer.isActive = false
}

func (timer *Timer) TimedOut() bool {
	return timer.isActive && time.Now().After(timer.endTime)
}
