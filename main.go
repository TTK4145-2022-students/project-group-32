package main

import (

	// "elevators/filesystem"

	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/filesystem"
	"elevators/network"
	"elevators/timer"
	"os"

	//"elevators/filesystem"
	"elevators/hardware"
	// "fmt"
	// "io/ioutil"
	// "elevators/phoenix"
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
	// stopChange := make(chan bool)
	doorTimedOut := make(chan bool)
	decisionTimedOut := make(chan bool)
	ordersRecieved := make(chan orderstate.AllOrders)

	go hardware.PollButtons(buttonPress)
	go hardware.PollFloorSensor(floorArrival, floorLeft)
	go hardware.PollObstructionSwitch(obstructionChange)
	// go hardware.PollStopButton(stopChange)

	go timer.DoorTimer.PollTimerOut(doorTimedOut)

	go network.PollReceiveOrderState(ordersRecieved)

	go network.SendOrderStatePeriodically()

	go filesystem.SaveStatePeriodically()

	for {
		select {
		case buttonEvent := <-buttonPress:
			// fmt.Printf("%+v\n", a)
			orderstate.AcceptNewOrder(buttonEvent.Button, buttonEvent.Floor)
			orders := orderstate.GetOrders()
			cabstate.FSMNewOrder(buttonEvent.Floor, orders)

		case floor := <-floorArrival:
			// fmt.Printf("%+v\n", a)
			hardware.SetFloorIndicator(floor)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(floor, orders)

		case <-floorLeft:
			// fmt.Printf("%+v\n", a)
			cabstate.FSMFloorLeave()

		case obstruction := <-obstructionChange:
			// fmt.Printf("%+v\n", a)
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(obstruction, orders)

		case <-doorTimedOut:
			// fmt.Printf("%+v\n", a)
			orders := orderstate.GetOrders()
			cabstate.FSMDoorTimeout(orders)

		case <-decisionTimedOut:
			// fmt.Printf("%+v\n", a)
			orders := orderstate.GetOrders()
			cabstate.FSMDoorClose(orders)

		// case a := <-stopChange:
		// 	_ = a
		// fmt.Printf("%+v\n", a)
		// for f := 0; f < hardware.FloorCount; f++ {
		// 	for b := hardware.ButtonType(0); b < 3; b++ {
		// 		hardware.SetButtonLamp(b, f, false)
		// 	}
		// }
		case recievedOrderState := <-ordersRecieved:
			// fmt.Printf("%+v\n", a)
			// fmt.Println("updating orders")
			// todo better handling of bunch update of new orders
			newOrdersInFloors := orderstate.UpdateOrders(recievedOrderState)
			orders := orderstate.GetOrders()
			for floor, newOrder := range newOrdersInFloors {
				if newOrder {
					cabstate.FSMNewOrder(floor, orders)
				}
			}
		}
	}
}
