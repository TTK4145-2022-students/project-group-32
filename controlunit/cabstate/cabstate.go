package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"fmt"
)

type ElevatorBehaviour int

const (
	Idle ElevatorBehaviour = iota
	DoorOpen
	CabObstructed
	Moving
)

type CabState struct {
	AboveOrAtFloor  int
	BetweenFloors   bool
	RecentDirection hardware.MotorDirection
	MotorDirection  hardware.MotorDirection
	DoorObstructed  bool
	Behaviour       ElevatorBehaviour
}

var Cab CabState

func setMotorAndCabState(state hardware.MotorDirection) {
	hardware.SetMotorDirection(state)
	Cab.MotorDirection = state
	switch state {
	case hardware.MD_Up:
		Cab.Behaviour = Moving
		Cab.RecentDirection = state
	case hardware.MD_Down:
		Cab.Behaviour = Moving
		Cab.RecentDirection = state
	case hardware.MD_Stop:
		Cab.Behaviour = Idle
	default:
		panic("motor state not implemented " + string(rune(state)))
	}
}

func FSMInitBetweenFloors() ElevatorBehaviour {
	Cab.BetweenFloors = true
	setMotorAndCabState(hardware.MD_Down)
	return Cab.Behaviour
}

func FSMNewOrder(orderFloor int, orders orderstate.AllOrders) ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorStop(orderFloor, orders)
		} else if Cab.AboveOrAtFloor < orderFloor {
			setMotorAndCabState(hardware.MD_Up)
		} else {
			setMotorAndCabState(hardware.MD_Down)
		}
	case Moving:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			setMotorAndCabState(hardware.MD_Stop)
			FSMFloorStop(orderFloor, orders)
		}
	}
	return Cab.Behaviour
}

func FSMFloorArrival(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	Cab.AboveOrAtFloor = floor
	Cab.BetweenFloors = false
	orderStatus := orderstate.GetOrderStatus(orders, floor)
	switch Cab.Behaviour {
	case Moving:
		motorAction := prioritize.MotorActionOnFloorArrival(
			Cab.RecentDirection,
			orderStatus)
		setMotorAndCabState(motorAction)

		if motorAction == hardware.MD_Stop {
			return FSMFloorStop(floor, orders)
		}
	default:
		// panic("Invalid cab state on floor arrival")
		fmt.Println("nomoarrive")
	}
	return Cab.Behaviour
}

func FSMFloorLeave() ElevatorBehaviour {
	Cab.BetweenFloors = true
	switch Cab.Behaviour {
	case Moving:
		switch Cab.MotorDirection {
		case hardware.MD_Up:
			break
		case hardware.MD_Down:
			Cab.AboveOrAtFloor = Cab.AboveOrAtFloor - 1
		default:
			panic("Invalid motor direction on floor leave")
		}
	default:
		// panic("Invalid cab state on floor leave")
		fmt.Println("nomoleave")
	}
	return Cab.Behaviour
}

func FSMDoorClose(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		motorAction := prioritize.MotorActionOnDoorClose(
			Cab.RecentDirection,
			currentOrderStatus)
		setMotorAndCabState(motorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}
