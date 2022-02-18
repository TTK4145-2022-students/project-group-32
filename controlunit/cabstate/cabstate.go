package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware/cab"
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
	motorRunning   bool
	motorDirection cab.Direction
	behaviour      ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
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
			Cab.motorDirection = cab.Up
		} else {
			Cab.motorDirection = cab.Down
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

func FSMDoorTimeout(floor int) ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		break
	case Moving:
		break
	case DoorOpen:
		//todo check orders
		Cab.doorOpen = false
	}
	return Cab.behaviour
}
