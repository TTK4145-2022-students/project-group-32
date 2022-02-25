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

func FSMInitBetweenFloors() ElevatorBehaviour {
	Cab.motorDirection = cab.Down
	Cab.motorRunning = true
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
		} else if orderstate.OrderInFloor(Cab.aboveOrAtFloor, cab.Up) {
			Cab.doorOpen = true
			Cab.motorDirection = cab.Up
			Cab.behaviour = DoorOpen
		} else if orderstate.OrderInFloor(Cab.aboveOrAtFloor, cab.Down) {
			Cab.doorOpen = true
			Cab.motorDirection = cab.Down
			Cab.behaviour = DoorOpen
		} else if (Cab.motorDirection == cab.Up && orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor)) ||
			(Cab.motorDirection == cab.Down && orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor)) {
			Cab.motorRunning = true
			Cab.behaviour = Moving
		} else if Cab.motorDirection == cab.Up && orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor) {
			Cab.motorDirection = cab.Down
			Cab.motorRunning = true
			Cab.behaviour = Moving
		} else if Cab.motorDirection == cab.Down && orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor) {
			Cab.motorDirection = cab.Up
			Cab.motorRunning = true
			Cab.behaviour = Moving
		}
	}
	return Cab.behaviour
}
