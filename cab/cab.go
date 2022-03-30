package cab

import (
	"elevators/eta"
	"elevators/hardware"
	"elevators/orders"
	"elevators/prioritize"
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
	allOrders orders.AllOrders) ElevatorBehaviour {

	switch Cab.Behaviour {
	case Idle:
		eta.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false,
			allOrders)

		if cabInFloor(orderFloor) {
			FSMFloorStop(
				orderFloor,
				allOrders)
		} else {
			timer.DecisionDeadlineTimer.TimerStart()
		}

	case Moving:
		if cabInFloor(orderFloor) {
			FSMFloorArrival(
				Cab.AboveOrAtFloor,
				allOrders)
		}

	case DoorOpen:
		eta.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			true,
			allOrders)

		orderSummary := orders.GetOrderSummary(
			allOrders,
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
	allOrders orders.AllOrders) ElevatorBehaviour {

	Cab.AboveOrAtFloor = floor
	Cab.BetweenFloors = false

	orderSummary := orders.GetOrderSummary(
		allOrders,
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
				allOrders)
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

func FSMDecisionDeadline(allOrders orders.AllOrders) ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		allOrders, internalETAs := eta.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false,
			allOrders)

		currentOrderSummary := orders.GetOrderSummary(
			allOrders,
			Cab.AboveOrAtFloor)

		prioritizedDirection := eta.PrioritizedDirection(
			Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			allOrders,
			internalETAs)

		motorAction := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
			currentOrderSummary)

		setMotorAndCabState(motorAction)
	}
	timer.DecisionDeadlineTimer.TimerStop()
	return Cab.Behaviour
}

func FSMPoke(allOrders orders.AllOrders) ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		eta.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			false,
			allOrders)

		timer.DecisionDeadlineTimer.TimerStart()

	case DoorOpen:
		eta.UpdateOrderAndInternalETAs(
			Cab.RecentDirection,
			Cab.AboveOrAtFloor,
			true,
			allOrders)

		timer.DecisionDeadlineTimer.TimerStart()
	}
	return Cab.Behaviour
}
