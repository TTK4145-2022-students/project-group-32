package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	Moving
)

type CabState struct {
	aboveOrAtFloor int
	betweenFloors  bool
	doorOpen       bool
	motorDirection hardware.MotorDirection
	behaviour      ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
}

func FSMInitBetweenFloors() ElevatorBehaviour {
	Cab.motorDirection = hardware.MD_Down
	Cab.behaviour = Moving
	return Cab.behaviour
}

func FSMNewOrder(orderFloor int) ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			Cab.doorOpen = true
			Cab.behaviour = DoorOpen
			break
		}

		if Cab.aboveOrAtFloor < orderFloor {
			Cab.motorDirection = hardware.MD_Up
		} else {
			Cab.motorDirection = hardware.MD_Down
		}
		Cab.motorRunning = true
		Cab.behaviour = Moving
	case Moving:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			Cab.motorRunning = false
			Cab.doorOpen = true
			Cab.behaviour = DoorOpen
		}
	case DoorOpen:
		break
	}
	return Cab.behaviour
}

func FSMFloorArrival(floor int) ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		break
	case Moving:
		if orderstate.OrderInFloor(floor, Cab.motorDirection) {
			Cab.motorRunning = false
			Cab.doorOpen = true
			Cab.behaviour = DoorOpen
		}
	case DoorOpen:
		break
	}
	Cab.aboveOrAtFloor = floor
	return Cab.behaviour
}

func FSMDoorTimeout() ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		break
	case Moving:
		break
	case DoorOpen:
		//todo check orders
		Cab.doorOpen = false
		if orderstate.OrderInFloor(Cab.aboveOrAtFloor, Cab.motorDirection) {
			Cab.doorOpen = true
			Cab.behaviour = DoorOpen
		} else if orderstate.OrderInFloor(Cab.aboveOrAtFloor, hardware.MD_Up) {
			Cab.doorOpen = true
			Cab.motorDirection = hardware.MD_Up
			Cab.behaviour = DoorOpen
		} else if orderstate.OrderInFloor(Cab.aboveOrAtFloor, hardware.MD_Down) {
			Cab.doorOpen = true
			Cab.motorDirection = hardware.MD_Down
			Cab.behaviour = DoorOpen
		} else if (Cab.motorDirection == hardware.MD_Up && orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor)) ||
			(Cab.motorDirection == hardware.MD_Down && orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor)) {
			Cab.motorRunning = true
			Cab.behaviour = Moving
		} else if Cab.motorDirection == hardware.MD_Up && orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor) {
			Cab.motorDirection = hardware.MD_Down
			Cab.motorRunning = true
			Cab.behaviour = Moving
		} else if Cab.motorDirection == hardware.MD_Down && orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor) {
			Cab.motorDirection = hardware.MD_Up
			Cab.motorRunning = true
			Cab.behaviour = Moving
		}
	}
	return Cab.behaviour
}
