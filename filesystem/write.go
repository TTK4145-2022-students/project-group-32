package filesystem

import (
	// "fmt"
	"encoding/json"
	"io/ioutil"
)

func SaveElevatorState(elevatorState interface{}) {
	write("filesystem/elevator_state.json", elevatorState)
}

func SaveOrders(orders interface{}) {
	write("filesystem/orders.json", orders)
}

func write(filepath string, elevatorState interface{}) {
	// fmt.Println("Filesystem/write.go")

	file, _ := json.MarshalIndent(elevatorState, "", " ")
	_ = ioutil.WriteFile(filepath, file, 0644)
}
