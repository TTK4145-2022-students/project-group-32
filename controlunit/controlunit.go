package controlunit

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/eta"
	"elevators/controlunit/orderstate"
	"elevators/hardware"
	"elevators/timer"
	"fmt"
)

func Init() {
	orderstate.InitOrders()
	cabstate.InitCabState()
	orders := orderstate.GetOrders()
	newETA := eta.ComputeETA(orders, 2, 3)
	fmt.Println(newETA.String())
}

func RunElevatorLoop() {

	drv_buttons := make(chan hardware.ButtonEvent)
	drv_floor_arrival := make(chan int)
	drv_floor_leave := make(chan bool)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timer := make(chan bool)

	go hardware.PollButtons(drv_buttons)
	go hardware.PollFloorSensor(drv_floor_arrival, drv_floor_leave)
	go hardware.PollObstructionSwitch(drv_obstr)
	go hardware.PollStopButton(drv_stop)
	go timer.PollTimer(drv_timer)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			orderstate.AcceptNewOrder(a.Button, a.Floor)
			cabstate.FSMNewOrder(a.Floor)

		case a := <-drv_floor_arrival:
			fmt.Printf("%+v\n", a)
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
			for f := 0; f < hardware.FloorCount; f++ {
				for b := hardware.ButtonType(0); b < 3; b++ {
					hardware.SetButtonLamp(b, f, false)
				}
			}
		}
	}
}
