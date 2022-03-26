package prioritize

import (
	"elevators/hardware"
)

func DoorActionOnDoorTimeout(
	prioritizedDirection hardware.MotorDirection,
	doorObstructed bool,
	currentOrders OrderSummary) hardware.DoorAction {

	switch prioritizedDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor ||
			(currentOrders.CabAtFloor &&
				currentOrders.AboveFloor) {
			return hardware.DS_Open_Up
		} else if currentOrders.DownAtFloor &&
			!currentOrders.AboveFloor {
			return hardware.DS_Open_Down
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor ||
			(currentOrders.CabAtFloor &&
				currentOrders.AboveFloor) {
			return hardware.DS_Open_Down
		} else if currentOrders.UpAtFloor &&
			!currentOrders.BelowFloor {
			return hardware.DS_Open_Up
		}
	case hardware.MD_Stop:
		break
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
	prioritizedDirection hardware.MotorDirection,
	currentOrders OrderSummary) hardware.DoorAction {

	switch prioritizedDirection {
	case hardware.MD_Up:
		if currentOrders.UpAtFloor ||
			(currentOrders.CabAtFloor &&
				currentOrders.AboveFloor) {
			return hardware.DS_Open_Up
		} else if currentOrders.DownAtFloor {
			return hardware.DS_Open_Down
		}
	case hardware.MD_Down:
		if currentOrders.DownAtFloor ||
			(currentOrders.CabAtFloor &&
				currentOrders.AboveFloor) {
			return hardware.DS_Open_Down
		} else if currentOrders.UpAtFloor {
			return hardware.DS_Open_Up
		}
	case hardware.MD_Stop:
		break
	default:
		panic("Invalid recent direction on floor stop")
	}
	if currentOrders.CabAtFloor {
		return hardware.DS_Open_Cab
	}
	return hardware.DS_Close
}

func DoorActionOnNewOrder(
	prioritizedDirection hardware.MotorDirection,
	currentOrders OrderSummary) hardware.DoorAction {

	switch prioritizedDirection {
	case hardware.MD_Up:
		if currentOrders.DownAtFloor {
			return hardware.DS_Open_Down
		}
	case hardware.MD_Down:
		if currentOrders.UpAtFloor {
			return hardware.DS_Open_Up
		}
	case hardware.MD_Stop:
		break
	default:
		panic("Invalid recent direction on floor stop")
	}
	if currentOrders.CabAtFloor {
		return hardware.DS_Open_Cab
	}
	return hardware.DS_Do_Nothing
}
