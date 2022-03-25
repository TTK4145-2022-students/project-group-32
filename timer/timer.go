package timer

import (
	// "fmt"
	"time"
)

const _pollRate = 20 * time.Millisecond
const _doorOpenTime = 3 * time.Second
const _decisionDeadline = 250 * time.Millisecond
const _pokeRate = 1 * time.Second

type Timer struct {
	isActive      bool
	endTime       time.Time
	timerDuration time.Duration
}

var DoorTimer = Timer{timerDuration: _doorOpenTime}
var DecisionDeadlineTimer = Timer{timerDuration: _decisionDeadline}
var PokeCabTimer = Timer{timerDuration: _pokeRate}

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
