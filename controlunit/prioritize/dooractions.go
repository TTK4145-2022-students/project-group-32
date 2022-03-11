package prioritize

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func DoorActionOnDoorTimeout(
	recentDirection hardware.MotorDirection,
	doorObstructed bool,
	currentOrders orderstate.OrderStatus) hardware.DoorState {
	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor ||
			(currentOrders.CabAtFloor && currentOrders.AboveFloor) {
			return hardware.DS_Open_Up
		} else if currentOrders.DownAtFloor && !currentOrders.AboveFloor {
			return hardware.DS_Open_Down
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor ||
			(currentOrders.CabAtFloor && currentOrders.AboveFloor) {
			return hardware.DS_Open_Down
		} else if currentOrders.UpAtFloor && !currentOrders.BelowFloor {
			return hardware.DS_Open_Up
		}
	default:
		panic("Invalid recent direction on door close")
	}
	if currentOrders.CabAtFloor {
		return hardware.DS_Open_Cab
	} else if doorObstructed {
		return hardware.DS_Open_Cab
	}
	return hardware.DS_Close
}

func DoorActionOnFloorStop(
	recentDirection hardware.MotorDirection,
	currentOrders orderstate.OrderStatus) hardware.DoorState {

	switch recentDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor ||
			(currentOrders.CabAtFloor && currentOrders.AboveFloor) {
			return hardware.DS_Open_Up
		} else if currentOrders.DownAtFloor {
			return hardware.DS_Open_Down
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor ||
			(currentOrders.CabAtFloor && currentOrders.AboveFloor) {
			return hardware.DS_Open_Down
		} else if currentOrders.UpAtFloor {
			return hardware.DS_Open_Up
		}
	default:
		panic("Invalid recent direction on floor stop")
	}
	if currentOrders.CabAtFloor {
		return hardware.DS_Open_Cab
	}
	return hardware.DS_Close
}
