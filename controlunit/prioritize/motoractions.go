package prioritize

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func MotorActionOnDoorClose(
	recentDirection hardware.MotorDirection,
	currentOrders orderstate.OrderStatus) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.AboveFloor {
			return hardware.MD_Up
		}
	case hardware.MD_Down:
		if currentOrders.BelowFloor {
			return hardware.MD_Down
		}
	default:
		panic("Invalid recent direction on door close")
	}

	if currentOrders.AboveFloor {
		return hardware.MD_Up
	} else if currentOrders.BelowFloor {
		return hardware.MD_Down
	} else {
		return hardware.MD_Stop
	}
}

func MotorActionOnFloorArrival(
	recentDirection hardware.MotorDirection,
	currentOrders orderstate.OrderStatus) hardware.MotorDirection {

	if currentOrders.CabAtFloor {
		return hardware.MD_Stop
	}
	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor {
			return hardware.MD_Stop
		} else if currentOrders.AboveFloor {
			return hardware.MD_Up
		} else if currentOrders.DownAtFloor {
			return hardware.MD_Stop
		} else if currentOrders.BelowFloor {
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor {
			return hardware.MD_Stop
		} else if currentOrders.BelowFloor {
			return hardware.MD_Down
		} else if currentOrders.UpAtFloor {
			return hardware.MD_Stop
		} else if currentOrders.AboveFloor {
			return hardware.MD_Up
		}
	default:
		panic("Invalid recent direction on floor arrival")
	}
	return hardware.MD_Stop
}
