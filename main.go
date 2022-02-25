package main

import (

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

	go network.TestSendAndReceive()
	hardware.Init("localhost:15657", hardware.FloorCount)

	controlunit.RunElevatorLoop()

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
