package prioritize

import (
	"elevators/hardware"
)

func ActionOnDoorClose(
	recentDirection hardware.MotorDirection,
	ordersAbove bool,
	ordersBelow bool) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if ordersAbove {
			return hardware.MD_Up
		}
	case hardware.MD_Down:
		if ordersBelow {
			return hardware.MD_Down
		}
	default:
		panic("Invalid recent direction on door close")
	}

	if ordersAbove {
		return hardware.MD_Up
	} else if ordersBelow {
		return hardware.MD_Down
	} else {
		return hardware.MD_Stop
	}
}

func ActionOnDoorTimeout(
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
