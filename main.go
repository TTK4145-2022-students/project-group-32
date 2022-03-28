package main

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/filesystem"
	"elevators/hardware"
	"elevators/network"
	"elevators/timer"
	"os"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()
	if len(os.Args) > 1 {
		hardware.Init("localhost:"+os.Args[1], hardware.FloorCount)
	} else {
		hardware.Init("localhost:15657", hardware.FloorCount)
	}
	orderstate.Init(filesystem.ReadOrders())
	cabstate.Init(filesystem.ReadCabState())

	buttonPress := make(chan hardware.ButtonEvent)
	floorArrival := make(chan int)
	floorLeft := make(chan bool)
	obstructionChange := make(chan bool)
	stopChange := make(chan bool)

	doorTimedOut := make(chan bool)
	decisionTimedOut := make(chan bool)
	pokeElevator := make(chan bool)

	ordersRecieved := make(chan orderstate.AllOrders)

	go hardware.PollButtons(buttonPress)
	go hardware.PollFloorSensor(floorArrival, floorLeft)
	go hardware.PollObstructionSwitch(obstructionChange)
	go hardware.PollStopButton(stopChange)

	go timer.DoorTimer.PollTimerOut(doorTimedOut)
	go timer.DecisionTimer.PollTimerOut(decisionTimedOut)
	go timer.PokeElevatorTimer.PollTimerOut(pokeElevator)

	go network.PollReceiveOrderState(ordersRecieved)

	go network.SendOrderStatePeriodically()

	go filesystem.SaveStatePeriodically()

	//Glorious loop
	for {
		select {
		case buttonEvent := <-buttonPress:
			// fmt.Println("Button pressed")
			orderstate.AcceptNewOrder(buttonEvent.Button, buttonEvent.Floor)
			orders := orderstate.GetOrders()
			timer.DecisionTimer.TimerStart()
			cabstate.FSMNewOrder(buttonEvent.Floor, orders)

		case floor := <-floorArrival:
			// fmt.Println("Ariived at floor")
			hardware.SetFloorIndicator(floor)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(floor, orders)

		case <-floorLeft:
			cabstate.FSMFloorLeave()

		case obstruction := <-obstructionChange:
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(obstruction, orders)

		case <-doorTimedOut:
			// fmt.Println("Door timed out")
			orders := orderstate.GetOrders()
			cabstate.FSMDoorTimeout(orders)

		case <-decisionTimedOut:
			// fmt.Println("Decision timed out")
			orders := orderstate.GetOrders()
			cabstate.FSMDecisionTimeout(orders)

		case <-pokeElevator:
			orders := orderstate.GetOrders()
			timer.DecisionTimer.TimerStart()
			cabstate.FSMDecisionTimeout(orders)
			timer.PokeElevatorTimer.TimerStart()

		case <-stopChange:
			orderstate.ResetOrders()
			for f := 0; f < hardware.FloorCount; f++ {
				for b := hardware.ButtonType(0); b < 3; b++ {
					hardware.SetButtonLamp(b, f, false)
				}
			}
		case recievedOrderState := <-ordersRecieved:
			newOrdersInFloors := orderstate.UpdateOrders(recievedOrderState)
			orders := orderstate.GetOrders()
			for floor, newOrder := range newOrdersInFloors {
				if newOrder {
					// fmt.Println("recieved new order")
					cabstate.FSMNewOrder(floor, orders)
				}
			}
		}
	}
}
