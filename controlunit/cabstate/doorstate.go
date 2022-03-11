package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
)

// type DoorState struct {
// 	Obstructed bool
// 	DoorBehaviour hardware.DoorState
// }

// var Door DoorState

// func InitDoorState() {
// 	Door := new(DoorState)
// 	_ = Door
// }

// func setDoorState(state hardware.DoorState) {
// 	hardware.SetDoorOpenLamp(bool(state))
// 	switch state {
// 	case hardware.DS_Open:
// 		.Open = true
// 		timer.TimerStart(3)
// 	case hardware.DS_Close:
// 		Door.Open = false
// 		timer.TimerStop()
// 	default:
// 		panic("door state not implemented")
// 	}
// }

// func FSMCloseDoor() DoorState {
// 	switch Door.Obstructed {
// 	case true:
// 		break
// 	case false:
// 		setDoorState(hardware.DS_Open)
// 	default:
// 		panic("door obstruction not boolean on input")
// 	}
// 	return Door
// }

func setDoorAndCabState(state hardware.DoorState) {
	switch state {
	case hardware.DS_Open_Cab:
		openDoor()
		orderstate.CompleteOrderCab(Cab.aboveOrAtFloor)
	case hardware.DS_Open_Up:
		openDoor()
		orderstate.CompleteOrderCabAndUp(Cab.aboveOrAtFloor)
		Cab.recentDirection = hardware.MD_Up
	case hardware.DS_Open_Down:
		openDoor()
		orderstate.CompleteOrderCabAndDown(Cab.aboveOrAtFloor)
		Cab.recentDirection = hardware.MD_Down
	case hardware.DS_Close:
		closeDoor()
	default:
		panic("door state not implemented")
	}
}

func openDoor() {
	hardware.SetDoorOpenLamp(true)
	Cab.behaviour = DoorOpen
	timer.TimerStart(3)
}

func closeDoor() {
	hardware.SetDoorOpenLamp(false)
	Cab.behaviour = Idle
	timer.TimerStop()
}

func FSMObstructionChange(obstructed bool, orders orderstate.AllOrders) {
	Cab.doorObstructed = obstructed
	switch obstructed {
	case true:
		switch Cab.behaviour {
		case DoorOpen:
			Cab.behaviour = CabObstructed
		}
	case false:
		switch Cab.behaviour {
		case CabObstructed:
			Cab.behaviour = DoorOpen
			if timer.TimedOut() {
				FSMDoorTimeout(orders)
			}
		}
	}
}

func FSMDoorTimeout(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.aboveOrAtFloor)
	switch Cab.behaviour {
	case DoorOpen:
		doorAction := prioritize.DoorActionOnDoorTimeout(
			Cab.recentDirection,
			Cab.doorObstructed,
			currentOrderStatus)
		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Close {
			return FSMDoorClose(orders)
		}
	case CabObstructed:
		break
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}

func FSMFloorStop(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.aboveOrAtFloor)
	switch Cab.behaviour {
	case Idle:
		doorAction := prioritize.DoorActionOnFloorStop(
			Cab.recentDirection,
			currentOrderStatus)
		setDoorAndCabState(doorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}
