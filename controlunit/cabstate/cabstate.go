package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	idledistribution "elevators/idleDistribution"
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

func Init(cabState CabState) {
	Cab = cabState
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
	orderstate.UpdateOrderAndInternalETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorStop(orderFloor, orders)
		} else {
			timer.DecisionDeadlineTimer.TimerStart()
		}
	case Moving:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorArrival(Cab.AboveOrAtFloor, orders)
		}
	case DoorOpen:
		orderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
		doorAction := prioritize.DoorActionOnNewOrder(Cab.RecentDirection, orderStatus)
		setDoorAndCabState(doorAction)
	}
	return Cab.Behaviour
}

func FSMFloorArrival(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	Cab.AboveOrAtFloor = floor
	Cab.BetweenFloors = false
	// orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
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
	}
	return Cab.Behaviour
}

func FSMDecisionDeadline() ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		orders, internalETAs := orderstate.UpdateOrderAndInternalETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
		currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
		motorAction := prioritize.MotorActionOnDecisionDeadline(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				internalETAs),
			// Cab.RecentDirection,
			currentOrderStatus)

		setMotorAndCabState(motorAction)
		if motorAction == hardware.MD_Stop && orderstate.AnyOrders(orders) {
			timer.ETAExpiredAlarm.SetAlarm(orderstate.FirstBestETAexpirationWithOrder(orders))
		}
	}
	timer.DecisionDeadlineTimer.TimerStop()
	return Cab.Behaviour
}

func FSMDistribute(
	orders orderstate.AllOrders,
	internalETAs orderstate.InternalETAs) ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		// orders, internalETAs = orderstate.UpdateOrderAndInternalETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
		motorAction := idledistribution.MotorActionOnDistribute(
			Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			orders,
			internalETAs)

		setMotorAndCabState(motorAction)

		orderstate.UpdateOrderAndInternalETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)

	}
	return Cab.Behaviour
}
