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

	drv_buttons := make(chan hardware.ButtonEvent)
	drv_floor_arrival := make(chan int)
	drv_floor_leave := make(chan bool)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timer := make(chan bool)
	drv_recieve := make(chan orderstate.AllOrders)

	go hardware.PollButtons(drv_buttons)
	go hardware.PollFloorSensor(drv_floor_arrival, drv_floor_leave)
	go hardware.PollObstructionSwitch(drv_obstr)
	go hardware.PollStopButton(drv_stop)

	go timer.PollTimer(drv_timer)

	go network.PollReceiveOrderState(drv_recieve)

	go network.SendOrderStatePeriodically()

	go filesystem.SaveStatePeriodically()

	for {
		select {
		case a := <-drv_buttons:
			// fmt.Printf("%+v\n", a)
			orderstate.AcceptNewOrder(a.Button, a.Floor)
			orders := orderstate.GetOrders()
			cabstate.FSMNewOrder(a.Floor, orders)

		case a := <-drv_floor_arrival:
			// fmt.Printf("%+v\n", a)
			hardware.SetFloorIndicator(a)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(a, orders)

		case a := <-drv_floor_leave:
			// fmt.Printf("%+v\n", a)
			_ = a
			cabstate.FSMFloorLeave()

		case a := <-drv_obstr:
			// fmt.Printf("%+v\n", a)
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(a, orders)

		case a := <-drv_timer:
			// fmt.Printf("%+v\n", a)
			if a {
				orders := orderstate.GetOrders()
				cabstate.FSMDoorTimeout(orders)
			}

		case a := <-drv_stop:
			_ = a
			// fmt.Printf("%+v\n", a)
			// for f := 0; f < hardware.FloorCount; f++ {
			// 	for b := hardware.ButtonType(0); b < 3; b++ {
			// 		hardware.SetButtonLamp(b, f, false)
			// 	}
			// }
		case recievedOrderState := <-drv_recieve:
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
