package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
)

// type DoorState struct {
// 	Obstructed bool
// 	Open   bool
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
	hardware.SetDoorOpenLamp(bool(state))
	switch state {
	case hardware.DS_Open:
		Cab.behaviour = DoorOpen
		timer.TimerStart(3)
	case hardware.DS_Close:
		Cab.behaviour = Idle
		timer.TimerStop()
	default:
		panic("door state not implemented")
	}
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

		if doorAction == hardware.DS_Open {
			orderstate.CompleteOrder(Cab.aboveOrAtFloor,
				Cab.recentDirection,
				currentOrderStatus)
		} else {
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

		if doorAction == hardware.DS_Open {
			orderstate.CompleteOrder(Cab.aboveOrAtFloor,
				Cab.recentDirection,
				currentOrderStatus)
		}
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}
