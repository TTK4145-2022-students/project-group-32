package prioritize

import (
	"elevators/hardware"
)

func MotorActionOnDoorClose(
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

func MotorActionOnFloorArrival(
	recentDirection hardware.MotorDirection,
	upOrdersInFloor bool,
	downOrdersInFloor bool,
	cabOrdersInFloor bool,
	ordersAbove bool,
	ordersBelow bool) hardware.MotorDirection {

	if cabOrdersInFloor {
		return hardware.MD_Stop
	}
	switch recentDirection {
	case hardware.MD_Up:
		if upOrdersInFloor {
			return hardware.MD_Stop
		} else if ordersAbove {
			return hardware.MD_Up
		} else if downOrdersInFloor {
			return hardware.MD_Stop
		} else if ordersBelow {
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if downOrdersInFloor {
			return hardware.MD_Stop
		} else if ordersBelow {
			return hardware.MD_Down
		} else if upOrdersInFloor {
			return hardware.MD_Stop
		} else if ordersAbove {
			return hardware.MD_Up
		}
	default:
		panic("Invalid recent direction on floor arrival")
	}
	return hardware.MD_Stop
}
