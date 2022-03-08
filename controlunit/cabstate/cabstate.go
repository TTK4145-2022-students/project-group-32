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
	FSMInitBetweenFloors()
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
	Cab.betweenFloors = true
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
	Cab.aboveOrAtFloor = floor
	Cab.betweenFloors = false
	switch Cab.behaviour {
	case Moving:
		motorAction := prioritize.MotorActionOnFloorArrival(
			Cab.recentDirection,
			orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.CabOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.OrdersAbove(Cab.aboveOrAtFloor),
			orderstate.OrdersBelow(Cab.aboveOrAtFloor))
		setMotorAndCabState(motorAction)

		if motorAction != hardware.MD_Stop {
			return Cab.behaviour
		}
		doorAction := prioritize.DoorActionOnFloorStop(
			Cab.recentDirection,
			orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.CabOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.OrdersAbove(Cab.aboveOrAtFloor),
			orderstate.OrdersBelow(Cab.aboveOrAtFloor))
		setDoorAndCabState(doorAction)
		if doorAction == hardware.DS_Open {
			orderstate.CompleteOrder(floor,
				Cab.recentDirection,
				orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
				orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
				orderstate.OrdersAbove(Cab.aboveOrAtFloor),
				orderstate.OrdersBelow(Cab.aboveOrAtFloor))
		}
	default:
		panic("Invalid cab state on floor arrival")
	}
	return Cab.behaviour
}

func FSMFloorLeave() ElevatorBehaviour {
	Cab.betweenFloors = true
	switch Cab.behaviour {
	case Moving:
		switch Cab.motorDirection {
		case hardware.MD_Up:
			break
		case hardware.MD_Down:
			Cab.aboveOrAtFloor = Cab.aboveOrAtFloor - 1
		default:
			panic("Invalid motor direction on floor leave")
		}
	default:
		panic("Invalid cab state on floor leave")
	}
	return Cab.behaviour
}

func FSMDoorTimeout() ElevatorBehaviour {
	switch Cab.behaviour {
	case DoorOpen:
		//todo check orders
		doorAction := prioritize.DoorActionOnDoorTimeout(
			Cab.recentDirection,
			orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
			orderstate.CabOrdersInFloor(Cab.aboveOrAtFloor))
		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Open {
			orderstate.CompleteOrder(Cab.aboveOrAtFloor,
				Cab.recentDirection,
				orderstate.UpOrdersInFloor(Cab.aboveOrAtFloor),
				orderstate.DownOrdersInFloor(Cab.aboveOrAtFloor),
				orderstate.OrdersAbove(Cab.aboveOrAtFloor),
				orderstate.OrdersBelow(Cab.aboveOrAtFloor))
			return Cab.behaviour
		}
		motorAction := prioritize.MotorActionOnDoorClose(
			Cab.recentDirection,
			orderstate.OrdersAbove(Cab.aboveOrAtFloor),
			orderstate.OrdersBelow(Cab.aboveOrAtFloor))
		setMotorAndCabState(motorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}
