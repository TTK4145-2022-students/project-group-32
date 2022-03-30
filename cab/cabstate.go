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
	}
}

func FSMInitBetweenFloors() ElevatorBehaviour {
	Cab.BetweenFloors = true
	setMotorAndCabState(hardware.MD_Down)
	return Cab.Behaviour
}

func FSMNewOrder(
	orderFloor int,
	orders orderstate.AllOrders) ElevatorBehaviour {

	switch Cab.Behaviour {
	case Idle:
		orderstate.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false)

		if cabInFloor(orderFloor) {
			FSMFloorStop(
				orderFloor,
				orders)
		} else {
			timer.DecisionDeadlineTimer.TimerStart()
		}

	case Moving:
		if cabInFloor(orderFloor) {
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

func FSMFloorArrival(
	floor int,
	orders orderstate.AllOrders) ElevatorBehaviour {

	Cab.AboveOrAtFloor = floor
	Cab.BetweenFloors = false

	orderSummary := orderstate.GetOrderSummary(
		orders,
		floor)

	switch Cab.Behaviour {
	case Moving:
		motorAction := prioritize.MotorActionOnFloorArrival(
			Cab.RecentDirection,
			orderSummary)

		setMotorAndCabState(motorAction)

		if motorAction == hardware.MD_Stop {
			return FSMFloorStop(
				floor,
				orders)
		}
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
		panic("Invalid cab state on floor leave")
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

		currentOrderSummary := orderstate.GetOrderSummary(
			orders,
			Cab.AboveOrAtFloor)

		prioritizedDirection := orderstate.PrioritizedDirection(
			Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			orders,
			internalETAs)

		motorAction := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
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
