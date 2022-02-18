package main

import (
	"fmt"
	"elevators/filesystem"
)


func main() {
	elevatorState := filesystem.ElevatorState {
		Name:  "Elevator 6",
		Floor: 1,
		Dir:   "up",
	}

	orderState := filesystem.OrderState {
		Name:  "Order 1",
		Floor: 1,
		Dir:   "up",
	}


	filesystem.SaveElevatorState(elevatorState)
	data := filesystem.ReadElevatorState()
	fmt.Println(data)

	filesystem.SaveOrders(orderState)
	data_order := filesystem.ReadOrders()
	fmt.Println(data_order)
}
