package controlunit

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/hardware"
	"elevators/network"
	"elevators/timer"
	"fmt"
)

func Init() {
	cabstate.FSMInitBetweenFloors()
}

func RunCommunicationLoop(receiver chan<- [hardware.FloorCount]bool) {

	drv_recieve := make(chan orderstate.AllOrders)

	go network.Send()
	go network.Receive(drv_recieve)

	for {
		select {
		case orders := <-drv_recieve:
			// fmt.Printf("%+v\n", a)
			updatedOrders := orderstate.UpdateOrders(orders)
			for _, newOrder := range updatedOrders {
				if newOrder {
					receiver <- updatedOrders
				}
			}

		}
	}
}

func RunElevatorLoop() {

	drv_buttons := make(chan hardware.ButtonEvent)
	drv_floor_arrival := make(chan int)
	drv_floor_leave := make(chan bool)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timer := make(chan bool)
	drv_order_update := make(chan [hardware.FloorCount]bool)

	go hardware.PollButtons(drv_buttons)
	go hardware.PollFloorSensor(drv_floor_arrival, drv_floor_leave)
	go hardware.PollObstructionSwitch(drv_obstr)
	go hardware.PollStopButton(drv_stop)
	go timer.PollTimer(drv_timer)
	// go PollOrderUpdate(drv_order_update)
	go RunCommunicationLoop(drv_order_update)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			orderstate.AcceptNewOrder(a.Button, a.Floor)
			orders := orderstate.GetOrders()
			cabstate.FSMNewOrder(a.Floor, orders)

		case a := <-drv_floor_arrival:
			fmt.Printf("%+v\n", a)
			hardware.SetFloorIndicator(a)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(a, orders)

		case a := <-drv_floor_leave:
			fmt.Printf("%+v\n", a)
			cabstate.FSMFloorLeave()

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(a, orders)

		case a := <-drv_timer:
			fmt.Printf("%+v\n", a)
			if a {
				orders := orderstate.GetOrders()
				cabstate.FSMDoorTimeout(orders)
			}

		case a := <-drv_stop:
			fmt.Printf("%+v\n", a)
			// for f := 0; f < hardware.FloorCount; f++ {
			// 	for b := hardware.ButtonType(0); b < 3; b++ {
			// 		hardware.SetButtonLamp(b, f, false)
			// 	}
			// }
		case a := <-drv_order_update:
			fmt.Printf("%+v\n", a)
			fmt.Println("updating orders")
			// todo better handling of bunch update of new orders
			orders := orderstate.GetOrders()
			for floor, newOrder := range a {
				if newOrder {
					cabstate.FSMNewOrder(floor, orders)
				}
			}
		}
	}
}

// // temp order update
// func PollOrderUpdate(receiver chan<- [hardware.FloorCount]bool) {
// 	for {
// 		time.Sleep(7 * time.Second)
// 		var newOrders orderstate.AllOrders
// 		newOrders.Up[0] = orderstate.OrderState{time.Now(), time.Now().Add(-5 * time.Second), time.Now().Add(-5 * time.Second)}
// 		updatedOrders := orderstate.UpdateOrders(newOrders)
// 		receiver <- updatedOrders
// 	}
// }
