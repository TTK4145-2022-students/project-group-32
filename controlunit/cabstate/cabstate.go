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
	motorDirection hardware.MotorDirection
	behaviour      ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
}

func setMotorAndCabState(state hardware.MotorDirection) {
	hardware.SetMotorDirection(state)
	Cab.motorDirection = state
	if state != hardware.MD_Stop {
		Cab.behaviour = Moving
	} else {
		Cab.behaviour = Idle
	}
}

func setDoorAndCabState(state hardware.DoorState) {
	hardware.SetDoorOpenLamp(state)
	if state == hardware.DS_Open {
		Cab.behaviour = DoorOpen
	} else {
		Cab.behaviour = Idle
	}
}

func FSMInitBetweenFloors() ElevatorBehaviour {
	setMotorAndCabState(hardware.MD_Down)
	Cab.behaviour = Moving
	return Cab.behaviour
}

func FSMNewOrder(orderFloor int) ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			setDoorAndCabState(hardware.DS_Open)
		} else if Cab.aboveOrAtFloor < orderFloor {
			setMotorAndCabState(hardware.MD_Up)
		} else {
			setMotorAndCabState(hardware.MD_Down)
		}
	case Moving:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			setMotorAndCabState(hardware.MD_Stop)
			setDoorAndCabState(hardware.DS_Open)
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
			setMotorAndCabState(hardware.MD_Stop)
			setDoorAndCabState(hardware.DS_Open)
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
