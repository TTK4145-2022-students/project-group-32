package orderstate

import (
	"elevators/hardware"
	"time"
)

type AllDurations struct {
	Up   [hardware.FloorCount]time.Duration
	Down [hardware.FloorCount]time.Duration
	Cab  [hardware.FloorCount]time.Duration
}

type AllETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const secsPerFloor = 2

const travelDuration = time.Duration(secsPerFloor) * time.Second

const secsPerOrder = 4

const secsPerDirectionChange = 3

var allDurations AllDurations

var allETAs AllETAs

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

// func travelDuration() time.Duration{
// 	return time.Duration(secsPerFloor) * time.Second
// }

func stopDuration(caborder bool) time.Duration {
	if caborder {
		return time.Duration(secsPerOrder) * time.Second
	} else {
		return time.Duration(0) * time.Second
	}
}

func ComputeAllETAs(durations AllDurations) AllETAs {
	var newETAs AllETAs
	var now = time.Now()
	for floor := 0; floor < hardware.FloorCount; floor++ {
		newETAs.Cab[floor] = now.Add(durations.Cab[floor])
		newETAs.Up[floor] = now.Add(durations.Up[floor])
		newETAs.Down[floor] = now.Add(durations.Down[floor])
	}
	return newETAs
}

func ComputeAllDurations(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	orders AllOrders) AllDurations {
	var newDurations AllDurations
	newDurations.Cab[currentFloor] = time.Duration(0)
	switch prioritizedDirection {
	// Todo: get more realistic newDuration, take orders in floors into consideration, not only cab above/below
	case hardware.MD_Up:
		floor := currentFloor
		for {
			if floor > currentFloor {
				newDurations.Cab[floor] = newDurations.Cab[floor-1] + travelDuration + stopDuration(orders.Cab[floor-1])
			}
			if CabOrdersAbove(orders.Cab, floor) {
				newDurations.Up[floor] = newDurations.Cab[floor]
				floor++
			} else if CabOrdersBelow(orders.Cab, currentFloor) {
				break // loop
			} else {
				newDurations.Up[floor] = newDurations.Cab[floor]
			}

			floor++
			if floor == hardware.FloorCount {
				floor--
				break // loop
			}
		}

		highestReachedFloor := floor
		newDurations.Down[floor] = newDurations.Cab[floor]

		for {
			if floor > 0 {
				floor--
			} else {
				break // loop
			}

			newDurations.Down[floor] = newDurations.Down[floor+1] + travelDuration
			if floor < currentFloor {
				newDurations.Cab[floor] = newDurations.Down[floor]
			}
		}

		if currentFloor > floor {
			newDurations.Up[floor] = newDurations.Cab[floor]
			for floor := floor + 1; floor < hardware.FloorCount; floor++ {
				if floor < currentFloor || floor >= highestReachedFloor {
					newDurations.Up[floor] = newDurations.Up[floor-1] + travelDuration
					if floor > highestReachedFloor {
						newDurations.Cab[floor] = newDurations.Up[floor]
					}
				}
			}
		}
	}
	return newDurations
}
