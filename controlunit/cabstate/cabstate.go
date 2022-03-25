package cabstate

import (
	"elevators/controlunit/orderstate"
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"elevators/timer"
	"fmt"
	// "time"
	// "fmt"
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

// func ForceActivationLoop(){
// 	for{
// 		if timer.ForceActionTimer.TimedOut(){
// 			fmt.Println("Force action timed out ")
// 			timer.ForceActionTimer.TimerStart()
// 		}else{
// 			timer.ForceActionTimer.TimerStop()
// 		}
// 		time.Sleep(time.Second)
// 	}
// }

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
	orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
	switch Cab.Behaviour {
	case Idle:
		if (Cab.AboveOrAtFloor == orderFloor) && !Cab.BetweenFloors {
			FSMFloorStop(orderFloor, orders)
		}
		//  else if Cab.AboveOrAtFloor < orderFloor {
		// 	setMotorAndCabState(hardware.MD_Up)
		// } else {
		// 	setMotorAndCabState(hardware.MD_Down)
		// }
		// timer.DecisionTimer.TimerStart()
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
	orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
	orderStatus := orderstate.GetOrderStatus(orders, floor)
	fmt.Println("")
	fmt.Println("FSMFloorArrival")
	fmt.Println("")
	fmt.Println("")
	switch Cab.Behaviour {
	case Moving:
		var pridir = orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
			Cab.RecentDirection,
			orders,
			orderstate.GetInternalETAs())
		fmt.Print("Recent direction:")
		fmt.Println(Cab.RecentDirection)
		fmt.Print("Prioritized direction:")
		fmt.Print(pridir)
		fmt.Println("")
		motorAction := prioritize.MotorActionOnFloorArrival(
			// orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
			// 	Cab.RecentDirection,
			// 	orders,
			// 	orderstate.GetInternalETAs()),
			Cab.RecentDirection,
			orderStatus)
		// fmt.Println("Moving")
		setMotorAndCabState(motorAction)

		if motorAction == hardware.MD_Stop {
			// fmt.Println("Should stop now")
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

func FSMDecisionTimeout(orders orderstate.AllOrders) ElevatorBehaviour {
	switch Cab.Behaviour {
	case Idle:
		orderstate.UpdateETAs(Cab.RecentDirection, Cab.AboveOrAtFloor)
		currentOrderStatus := orderstate.GetOrderStatus(orders, Cab.AboveOrAtFloor)
		motorAction := prioritize.MotorActionOnDoorClose(
			orderstate.PrioritizedDirection(Cab.AboveOrAtFloor,
				Cab.RecentDirection,
				orders,
				orderstate.GetInternalETAs()),
			currentOrderStatus)
		setMotorAndCabState(motorAction)
	}
	timer.DecisionTimer.TimerStop()
	// if orderstate.AnyOrders(orders) {
	// timer.DecisionTimer.TimerStart()
	// }
	return Cab.Behaviour
}
