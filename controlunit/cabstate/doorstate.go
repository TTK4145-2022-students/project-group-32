package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
)

func setDoorAndCabState(state hardware.DoorAction) {
	switch state {
	case hardware.DS_Open_Cab:
		openDoor()
		orderstate.CompleteOrderCab(Cab.AboveOrAtFloor)

	case hardware.DS_Open_Up:
		openDoor()
		orderstate.CompleteOrderCabAndUp(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Up

	case hardware.DS_Open_Down:
		openDoor()
		orderstate.CompleteOrderCabAndDown(Cab.AboveOrAtFloor)
		Cab.RecentDirection = hardware.MD_Down

	case hardware.DS_Close:
		closeDoor()
		timer.DecisionDeadlineTimer.TimerStart()

	case hardware.DS_Do_Nothing:
		break
	default:
		panic("door state not implemented")
	}
}

func openDoor() {
	hardware.SetDoorOpenLamp(true)
	Cab.Behaviour = DoorOpen
	timer.DoorTimer.TimerStart()
}

func closeDoor() {
	hardware.SetDoorOpenLamp(false)
	Cab.Behaviour = Idle
	timer.DoorTimer.TimerStop()
}

func FSMObstructionChange(
	obstructed bool,
	orders orderstate.AllOrders) {

	Cab.DoorObstructed = obstructed
	switch obstructed {
	case true:
		switch Cab.Behaviour {
		case DoorOpen:
			Cab.Behaviour = CabObstructed
		}

	case false:
		switch Cab.Behaviour {
		case CabObstructed:
			Cab.Behaviour = DoorOpen
			if timer.DoorTimer.TimedOut() {
				FSMDoorTimeout(orders)
			}
		}
	}
}

func FSMDoorTimeout(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderSummary := orderstate.GetOrderSummary(
		orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case DoorOpen:
		doorAction := prioritize.DoorActionOnDoorTimeout(
			orderstate.PrioritizedDirection(
				Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			Cab.DoorObstructed,
			currentOrderSummary)
		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Close {
			orderstate.UpdateOrderAndInternalETAs(
				Cab.RecentDirection,
				Cab.AboveOrAtFloor,
				false)
			timer.DecisionDeadlineTimer.TimerStart()
		}

	case CabObstructed:
		break
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}

func FSMFloorStop(
	floor int,
	orders orderstate.AllOrders) ElevatorBehaviour {

	currentOrderSummary := orderstate.GetOrderSummary(
		orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		doorAction := prioritize.DoorActionOnFloorStop(
			orderstate.PrioritizedDirection(
				Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			currentOrderSummary)
		setDoorAndCabState(doorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}
