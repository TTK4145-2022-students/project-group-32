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

type Direction int

const (
	Up Direction = iota
	Down
)

type CabState struct {
	aboveOrAtFloor  int
	betweenFloors   bool
	recentDirection Direction
	motorDirection  hardware.MotorDirection
	behaviour       ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
}

func setMotorAndCabState(state hardware.MotorDirection) {
	hardware.SetMotorDirection(state)
	Cab.motorDirection = state
	switch state {
	case hardware.MD_Up:
		Cab.behaviour = Moving
		Cab.recentDirection = Up
	case hardware.MD_Down:
		Cab.behaviour = Moving
		Cab.recentDirection = Down
	case hardware.MD_Stop:
		Cab.behaviour = Idle
	default:
		panic("motor state not implemented " + string(rune(state)))
	}
}

func setDoorAndCabState(state hardware.DoorState) {
	hardware.SetDoorOpenLamp(bool(state))
	switch state {
	case hardware.DS_Open:
		Cab.behaviour = DoorOpen
	case hardware.DS_Closed:
		Cab.behaviour = Idle
	default:
		panic("door state not implemented")
	}
}

func FSMInitBetweenFloors() ElevatorBehaviour {
	setMotorAndCabState(hardware.MD_Down)
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
		setDoorAndCabState(hardware.DS_Closed)
		if orderstate.OrderInFloor(Cab.aboveOrAtFloor, Cab.motorDirection) {
			setDoorAndCabState(hardware.DS_Open)
		} else if Cab.recentDirection == Up && orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor) {
			setMotorAndCabState(hardware.MD_Up)
		} else if Cab.recentDirection == Down && orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor) {
			setMotorAndCabState(hardware.MD_Down)
		} else if orderstate.OrdersAtOrBelow(Cab.aboveOrAtFloor) {
			setMotorAndCabState(hardware.MD_Down)
		} else if orderstate.OrdersAtOrAbove(Cab.aboveOrAtFloor) {
			setMotorAndCabState(hardware.MD_Up)
		}
	}
	return Cab.behaviour
}
