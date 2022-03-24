package timer

import (
	"time"
)

const _pollRate = 20 * time.Millisecond
const _doorOpenTime = 3 * time.Second
const _waitForDecisionTime = 100 * time.Millisecond

type Timer struct {
	isActive      bool
	endTime       time.Time
	timerDuration time.Duration
}

var DoorTimer = Timer{timerDuration: _doorOpenTime}
var DecisionTimer = Timer{timerDuration: _waitForDecisionTime}

func (timer *Timer) PollTimerOut(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := timer.TimedOut()
		if v != prev && v {
			receiver <- v
		}
		prev = v
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
