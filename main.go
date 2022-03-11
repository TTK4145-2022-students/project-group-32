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

	hardware.Init("localhost:15657", hardware.FloorCount)
	controlunit.Init()

	go controlunit.RunElevatorLoop()
	network.TestSend()
	// go network.TestSendAndReceive()


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
