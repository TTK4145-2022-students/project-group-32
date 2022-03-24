package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
	// "time"
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
// 		timer.TimerStart(doorOpenSecs)
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
		orderstate.CompleteOrderCab(Cab.AboveOrAtFloor)
	case hardware.DS_Open_Up:
		openDoor()
		orderstate.CompleteOrderCabAndUp(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Up
	case hardware.DS_Open_Down:
		openDoor()
		orderstate.CompleteOrderCabAndDown(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Down
	case hardware.DS_Close:
		closeDoor()
	default:
		panic("door state not implemented")
	}
}

func openDoor() {
	hardware.SetDoorOpenLamp(true)
	Cab.Behaviour = DoorOpen
	timer.DoorTimer.TimerStart()
}

func closeDoor() {
	hardware.SetDoorOpenLamp(false)
	Cab.Behaviour = Idle
	timer.DoorTimer.TimerStop()
}

func FSMObstructionChange(obstructed bool, orders orderstate.AllOrders) {
	Cab.DoorObstructed = obstructed
	switch obstructed {
	case true:
		switch Cab.Behaviour {
		case DoorOpen:
			Cab.Behaviour = CabObstructed
		}
	case false:
		switch Cab.Behaviour {
		case CabObstructed:
			Cab.Behaviour = DoorOpen
			if timer.DoorTimer.TimedOut() {
				FSMDoorTimeout(orders)
			}
		}
	}
}

func FSMDoorTimeout(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case DoorOpen:
		doorAction := prioritize.DoorActionOnDoorTimeout(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			Cab.DoorObstructed,
			currentOrderStatus)
		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Close {
			orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
			timer.DoorCloseDecisionTimer.TimerStart()
		}
	case CabObstructed:
		break
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}

func FSMFloorStop(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		doorAction := prioritize.DoorActionOnFloorStop(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			currentOrderStatus)
		setDoorAndCabState(doorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}
