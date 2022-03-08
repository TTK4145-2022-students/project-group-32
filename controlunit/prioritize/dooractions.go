package prioritize

import (
	"elevators/hardware"
)

func DoorActionOnDoorTimeout(
	recentDirection hardware.MotorDirection,
	upOrdersInFloor bool,
	downOrdersInFloor bool,
	cabOrdersInFloor bool) hardware.DoorState {

	switch recentDirection {
	case hardware.MD_Up:
		if upOrdersInFloor {
			return hardware.DS_Open
		}
	case hardware.MD_Down:
		if downOrdersInFloor {
			return hardware.DS_Open
		}
	default:
		panic("Invalid recent direction on door close")
	}
	if cabOrdersInFloor {
		return hardware.DS_Open
	}
	return hardware.DS_Close
}

func DoorActionOnFloorStop(
	recentDirection hardware.MotorDirection,
	upOrdersInFloor bool,
	downOrdersInFloor bool,
	cabOrdersInFloor bool,
	ordersAbove bool,
	ordersBelow bool) hardware.DoorState {

	if cabOrdersInFloor {
		return hardware.DS_Open
	}
	switch recentDirection {
	case hardware.MD_Up:
		if upOrdersInFloor {
			return hardware.DS_Open
		} else if downOrdersInFloor && !ordersAbove {
			return hardware.DS_Open
		}
	case hardware.MD_Down:
		if downOrdersInFloor {
			return hardware.DS_Open
		} else if upOrdersInFloor && !ordersBelow {
			return hardware.DS_Open
		}
	default:
		panic("Invalid recent direction on floor stop")
	}
	return hardware.DS_Close
}
