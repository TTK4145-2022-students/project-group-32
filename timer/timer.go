package timer

import (
	"time"
)

const _pollRate = 20 * time.Millisecond

var timerActive = false
var timerEndTime time.Time

func PollTimer(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := TimedOut()
		if v != prev {
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
