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
	orderstate.InitOrderState()
	cabstate.InitCabState()
	newETA := eta.ComputeETA(2, 3)
	// eta.ComputeETA(2, 3)
	// fmt.Println(fmt.Println(newETA.String()))
	fmt.Println(newETA.String())
}

func RunElevatorLoop() {

	drv_buttons := make(chan hardware.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)
	drv_timer := make(chan bool)

	go hardware.PollButtons(drv_buttons)
	go hardware.PollFloorSensor(drv_floors)
	go hardware.PollObstructionSwitch(drv_obstr)
	go hardware.PollStopButton(drv_stop)
	go timer.PollTimer(drv_timer)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			orderstate.AcceptNewOrder(a.Button, a.Floor)
			cabstate.FSMNewOrder(a.Floor)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			orderstate.CompleteOrder(a)
			cabstate.FSMFloorArrival(a)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				hardware.SetMotorDirection(hardware.MD_Stop)
			}

		case a := <-drv_timer:
			fmt.Printf("%+v\n", a)
			if a {
				cabstate.FSMDoorTimeout()
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
