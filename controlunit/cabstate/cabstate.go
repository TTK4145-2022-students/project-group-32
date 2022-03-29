package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
	// "fmt"
	// "time"
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

// func Init(cabState CabState) {
// 	Cab = cabState
// 	FSMInitBetweenFloors()
// }

func Init() {
	FSMInitBetweenFloors()
}

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
		orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false)
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorStop(orderFloor, orders)
		} else {
			timer.DecisionDeadlineTimer.TimerStart()
		}
	case Moving:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorArrival(
				Cab.AboveOrAtFloor,
				orders)
		}

	case DoorOpen:
		orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			true)
		orderSummary := orderstate.GetOrderSummary(
			orders,
			Cab.AboveOrAtFloor)
		doorAction := prioritize.DoorActionOnNewOrder(
			Cab.RecentDirection,
			orderSummary)

		setDoorAndCabState(doorAction)
	}
	return Cab.Behaviour
}

func FSMFloorArrival(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	Cab.AboveOrAtFloor = floor
	Cab.BetweenFloors = false
	// orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
	OrderSummary := orderstate.GetOrderSummary(orders, floor)
	switch Cab.Behaviour {
	case Moving:
		motorAction := prioritize.MotorActionOnFloorArrival(
			Cab.RecentDirection,
			OrderSummary)
		setMotorAndCabState(motorAction)

		if motorAction == hardware.MD_Stop {
			return FSMFloorStop(floor, orders)
		}
	default:
		// panic("Invalid cab state on floor arrival")
		// fmt.Println("nomoarrive")
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
		// fmt.Println("nomoleave")
	}
	return Cab.Behaviour
}

func FSMDecisionDeadline() ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		orders, internalETAs := orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false)
		currentOrderSummary := orderstate.GetOrderSummary(orders, Cab.AboveOrAtFloor)
		motorAction := prioritize.MotorActionOnDecisionDeadline(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				internalETAs),
			currentOrderSummary)
		setMotorAndCabState(motorAction)
	}
	timer.DecisionDeadlineTimer.TimerStop()
	return Cab.Behaviour
}

func FSMPoke() ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false)
		timer.DecisionDeadlineTimer.TimerStart()
	case DoorOpen:
		orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			true)
		timer.DecisionDeadlineTimer.TimerStart()
	}
	return Cab.Behaviour
}
