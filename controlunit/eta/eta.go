package eta

import (
	"elevators/hardware"
	"time"
)

const secsPerFloor = 2

const secsPerOrder = 4

func ComputeETA(
	direction hardware.MotorDirection, aboveOrAtFloor int, destinationFloor int) time.Time {
	return time.Now().Add(ComputeDurationToFloor(direction, aboveOrAtFloor, destinationFloor))

}

func ComputeDurationToFloor(
	direction hardware.MotorDirection, aboveOrAtFloor int, destinationFloor int) time.Duration {
	// Todo: get more realistic newETA, take orders into consideration
	var durationSecs = 0
	for floor := aboveOrAtFloor; (floor < hardware.FloorCount) && (floor >= 0) && (floor != destinationFloor); floor += int(direction) {
		durationSecs += secsPerFloor
		if floor == 0 {
			durationSecs += destinationFloor * secsPerFloor
		} else if floor == hardware.FloorCount-1 {
			durationSecs += (floor - destinationFloor) * secsPerFloor
		}
	}
	return time.Duration(durationSecs) * time.Second
}
