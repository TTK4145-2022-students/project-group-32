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
	timer.DecisionDeadlineTimer.TimerStart()
}

func FSMObstructionChange(obstructed bool, orders orderstate.AllOrders) {
	Cab.DoorObstructed = obstructed
	switch obstructed {
	case true:
		if Cab.Behaviour == DoorOpen {
			Cab.Behaviour = CabObstructed
		}
	case false:
		if Cab.Behaviour == CabObstructed {
			Cab.Behaviour = DoorOpen
			if timer.DoorTimer.TimedOut() {
				FSMDoorTimeout(orders)
			}
		}
	}
}

func FSMDoorTimeout(orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case DoorOpen:
		doorAction := prioritize.DoorActionOnDoorTimeout(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			Cab.DoorObstructed,
			currentOrderStatus)

		setDoorAndCabState(doorAction)

		if doorAction == hardware.DS_Close {
			orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
			timer.DecisionDeadlineTimer.TimerStart()
		}
	case CabObstructed:
		break
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}

func FSMFloorStop(floor int, orders orderstate.AllOrders) ElevatorBehaviour {
	currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		doorAction := prioritize.DoorActionOnFloorStop(
			orderstate.PrioritizedDirection(
				Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			currentOrderStatus)

		setDoorAndCabState(doorAction)
	default:
		panic("Invalid cab state on door timeout")
	}
	return Cab.Behaviour
}
