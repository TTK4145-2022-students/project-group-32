package timer

import (
	"time"
)

const _pollRate = 20 * time.Millisecond
const _doorOpenTime = 3 * time.Second
const _decisionDeadline = 250 * time.Millisecond
const _etaExpirationMargin = 50 * time.Millisecond

const _pokeRate = 500 * time.Millisecond

type Timer struct {
	isActive      bool
	endTime       time.Time
	timerDuration time.Duration
}

type Alarm struct {
	isActive    bool
	alarmTime   time.Time
	alarmOffset time.Duration
}

var DoorTimer = Timer{timerDuration: _doorOpenTime}
var DecisionDeadlineTimer = Timer{timerDuration: _decisionDeadline}

var PokeCabTimer = Timer{timerDuration: _pokeRate}

var ETAExpiredAlarm = Alarm{alarmOffset: _etaExpirationMargin}

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

func (alarm *Alarm) PollAlarm(receiver chan<- bool) {
	prev := false
	for {
		time.Sleep(_pollRate)
		v := alarm.AlarmRinging()
		if v != prev && v {
			receiver <- v
		}
		prev = v
	}
}

func (alarm *Alarm) SetAlarm(alarmTime time.Time) {
	alarm.alarmTime = alarmTime.Add(alarm.alarmOffset)
	alarm.isActive = true
}

func (alarm *Alarm) DismissAlarm() {
	alarm.isActive = false
}

func (alarm *Alarm) AlarmRinging() bool {
	return alarm.isActive && time.Now().After(alarm.alarmTime)
}
