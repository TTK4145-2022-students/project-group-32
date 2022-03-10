package prioritize

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func DoorActionOnDoorTimeout(
	recentDirection hardware.MotorDirection,
	currentOrders orderstate.OrderStatus) hardware.DoorState {

	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor {
			return hardware.DS_Open
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor {
			return hardware.DS_Open
		}
	default:
		panic("Invalid recent direction on door close")
	}
	if currentOrders.CabAtFloor {
		return hardware.DS_Open
	}
	return hardware.DS_Close
}

func DoorActionOnFloorStop(
	recentDirection hardware.MotorDirection,
	currentOrders orderstate.OrderStatus) hardware.DoorState {

	if currentOrders.CabAtFloor {
		return hardware.DS_Open
	}
	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor {
			return hardware.DS_Open
		} else if currentOrders.DownAtFloor && !currentOrders.AboveFloor {
			return hardware.DS_Open
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor {
			return hardware.DS_Open
		} else if currentOrders.UpAtFloor && !currentOrders.BelowFloor {
			return hardware.DS_Open
		}
	default:
		panic("Invalid recent direction on floor stop")
	}
	return hardware.DS_Close
}
