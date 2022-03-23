package timer

import (
	"time"
)

const _pollRate = 20 * time.Millisecond

var timerActive = false
var timerEndTime time.Time

func PollTimerOut(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := TimedOut()
		if v != prev && v {
			receiver <- v
		}
		prev = v
	}
}

func TimerStart(secondsDuration int) {
	timerEndTime = time.Now().Add(time.Second * time.Duration(secondsDuration))
	timerActive = true
}

func TimerStop() {
	timerActive = false
}

func TimedOut() bool {
	return timerActive && time.Now().After(timerEndTime)
}
