package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
)

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	Moving
)

type CabState struct {
	aboveOrAtFloor  int
	betweenFloors   bool
	recentDirection hardware.MotorDirection
	motorDirection  hardware.MotorDirection
	behaviour       ElevatorBehaviour
}

var Cab CabState

func InitCabState() {
	Cab := new(CabState)
	_ = Cab
}

func GetCabDirection() hardware.MotorDirection {
	return Cab.motorDirection
}

func setMotorAndCabState(state hardware.MotorDirection) {
	hardware.SetMotorDirection(state)
	Cab.motorDirection = state
	switch state {
	case hardware.MD_Up:
		Cab.behaviour = Moving
		Cab.recentDirection = state
	case hardware.MD_Down:
		Cab.behaviour = Moving
		Cab.recentDirection = state
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
		timer.TimerStart(3)
	case hardware.DS_Close:
		Cab.behaviour = Idle
		timer.TimerStop()
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
		} else if floor == 0 || floor == hardware.FloorCount-1 {
			setMotorAndCabState(hardware.MD_Stop)
		}
	case DoorOpen:
		break
	}
	Cab.aboveOrAtFloor = floor
	return Cab.behaviour
}

func FSMDoorTimeout() ElevatorBehaviour {
	switch Cab.behaviour {
	case DoorOpen:
		//todo check orders
		doorAction := prioritize.ActionOnDoorTimeout(
			Cab.recentDirection,
			orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.CabOrdersInFloor(Cab.aboveOrAtFloor))
		setDoorAndCabState(doorAction)

		motorAction := prioritize.ActionOnDoorClose(
			Cab.recentDirection,
			orderstate.OrdersAbove(Cab.aboveOrAtFloor),
			orderstate.OrdersBelow(Cab.aboveOrAtFloor))
		setMotorAndCabState(motorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}
