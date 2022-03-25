package main

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/filesystem"
	"elevators/hardware"
	"elevators/network"
	"elevators/phoenix"
	"elevators/timer"
	"fmt"
	"os"
)

func main() {
	phoenix.Init()
	go phoenix.Phoenix()
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
	decisionDeadlineTimedOut := make(chan bool)
	PokeCabTimedOut := make(chan bool)

	ordersRecieved := make(chan orderstate.AllOrders)

	go hardware.PollButtons(buttonPress)
	go hardware.PollFloorSensor(floorArrival, floorLeft)
	go hardware.PollObstructionSwitch(obstructionChange)
	go hardware.PollStopButton(stopChange)

	go timer.DoorTimer.PollTimerOut(doorTimedOut)
	go timer.DecisionDeadlineTimer.PollTimerOut(decisionDeadlineTimedOut)
	go timer.PokeCabTimer.PollTimerOut(PokeCabTimedOut)

	go network.PollReceiveOrders(ordersRecieved)
	go network.SendOrdersPeriodically()

	go filesystem.SaveStatesPeriodically()

	for {
		select {
		case buttonEvent := <-buttonPress:
			orderstate.AcceptNewOrder(buttonEvent.Button, buttonEvent.Floor)
			orders := orderstate.GetOrders()
			timer.DecisionDeadlineTimer.TimerStart()
			cabstate.FSMNewOrder(buttonEvent.Floor, orders)

		case floor := <-floorArrival:
			hardware.SetFloorIndicator(floor)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(floor, orders)

		case <-floorLeft:
			cabstate.FSMFloorLeave()

		case obstruction := <-obstructionChange:
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(obstruction, orders)

		case <-doorTimedOut:
			orders := orderstate.GetOrders()
			cabstate.FSMDoorTimeout(orders)

		case <-decisionDeadlineTimedOut:
			orders := orderstate.GetOrders()
			cabstate.FSMDecisionTimeout(orders)

		case <-PokeCabTimedOut:
			orders := orderstate.GetOrders()
			timer.DecisionDeadlineTimer.TimerStart()
			cabstate.FSMDecisionTimeout(orders)
			timer.PokeCabTimer.TimerStart()

		case a := <-stopChange:
			_ = a
			fmt.Printf("%+v\n", a)
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
					fmt.Println("recieved new order")
					cabstate.FSMNewOrder(floor, orders)
				}
			}
		}
	}
}
