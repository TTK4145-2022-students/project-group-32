package main

import (

	// "elevators/filesystem"
	"elevators/controlunit"
	"elevators/hardware"

	// "elevators/phoenix"
	//"time"
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
	controlunit.Init()

	go controlunit.RunElevatorLoop()

	// go network.TestSend()
	// go network.TestReceive()

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

	for {
	}
}
