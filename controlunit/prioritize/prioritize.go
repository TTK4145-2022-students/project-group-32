package prioritize

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func PrioritizedDirection(currentFloor int,
	recentDirection hardware.MotorDirection,
	orders orderstate.AllOrders,
	allETAs orderstate.AllETAs) hardware.MotorDirection {

	etaDirection := orderstate.ETADirection(currentFloor, recentDirection, orders, allETAs)
	if etaDirection == hardware.MD_Stop {
		return recentDirection
	} else {
		return etaDirection
	}
}
