package main

import (
	"fmt"
	// "elevators/filesystem"
	"elevators/controlunit"
	"elevators/hardware"
	"elevators/network"
	// "elevators/phoenix"
	//"time"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()

	controlunit.Init()

	network.TestSendAndReceive()
	hardware.Init("localhost:15657", hardware.FloorCount)

	var d hardware.MotorDirection = hardware.MD_Up
	hardware.SetMotorDirection(d)

	drv_buttons := make(chan hardware.ButtonEvent)
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_stop := make(chan bool)

	go hardware.PollButtons(drv_buttons)
	go hardware.PollFloorSensor(drv_floors)
	go hardware.PollObstructionSwitch(drv_obstr)
	go hardware.PollStopButton(drv_stop)

	for {
		select {
		case a := <-drv_buttons:
			fmt.Printf("%+v\n", a)
			hardware.SetButtonLamp(a.Button, a.Floor, true)

		case a := <-drv_floors:
			fmt.Printf("%+v\n", a)
			if a == hardware.FloorCount-1 {
				d = hardware.MD_Down
			} else if a == 0 {
				d = hardware.MD_Up
			}
			hardware.SetMotorDirection(d)

		case a := <-drv_obstr:
			fmt.Printf("%+v\n", a)
			if a {
				hardware.SetMotorDirection(hardware.MD_Stop)
			} else {
				hardware.SetMotorDirection(d)
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

	// elevatorState := filesystem.ElevatorState {
	// 	Name:  "Elevator 6",
	// 	Floor: 1,
	// 	Dir:   "up",
	// }

	// orderState := filesystem.OrderState {
	// 	Name:  "Order 1",
	// 	Floor: 1,
	// 	Dir:   "up",
	// }

	// filesystem.SaveElevatorState(elevatorState)
	// data := filesystem.ReadElevatorState()
	// fmt.Println(data)

	// filesystem.SaveOrders(orderState)
	// data_order := filesystem.ReadOrders()
	// fmt.Println(data_order)
}
