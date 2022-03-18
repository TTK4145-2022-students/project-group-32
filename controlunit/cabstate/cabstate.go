package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
)

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	CabObstructed
	Moving
)

type CabState struct {
	aboveOrAtFloor  int
	betweenFloors   bool
	recentDirection hardware.MotorDirection
	motorDirection  hardware.MotorDirection
	doorObstructed  bool
	behaviour       ElevatorBehaviour
}

var Cab CabState

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

func FSMInitBetweenFloors() ElevatorBehaviour {
	Cab.betweenFloors = true
	setMotorAndCabState(hardware.MD_Down)
	return Cab.behaviour
}

func FSMNewOrder(orderFloor int, orders orderstate.AllOrders) ElevatorBehaviour {
	switch Cab.behaviour {
	case Idle:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			FSMFloorStop(orderFloor, orders)
		} else if Cab.aboveOrAtFloor < orderFloor {
			setMotorAndCabState(hardware.MD_Up)
		} else {
			setMotorAndCabState(hardware.MD_Down)
		}
	case Moving:
		if (Cab.aboveOrAtFloor == orderFloor) && !Cab.betweenFloors {
			setMotorAndCabState(hardware.MD_Stop)
			FSMFloorStop(orderFloor, orders)
		}
	}
	return Cab.behaviour
}

func FSMFloorArrival(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	Cab.aboveOrAtFloor = floor
	Cab.betweenFloors = false
	orderStatus := orderstate.GetOrderStatus(orders, floor)
	switch Cab.behaviour {
	case Moving:
		motorAction := prioritize.MotorActionOnFloorArrival(
			Cab.recentDirection,
			orderStatus)
		setMotorAndCabState(motorAction)

		if motorAction == hardware.MD_Stop {
			return FSMFloorStop(floor, orders)
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

func FSMDoorClose(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.aboveOrAtFloor)
	switch Cab.behaviour {
	case Idle:
		motorAction := prioritize.MotorActionOnDoorClose(
			Cab.recentDirection,
			currentOrderStatus)
		setMotorAndCabState(motorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.behaviour
}
