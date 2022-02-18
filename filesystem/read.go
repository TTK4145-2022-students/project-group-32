package filesystem

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadElevatorState() ElevatorState {
	var elevatorState ElevatorState
	json.Unmarshal(read("filesystem/elevator_state.json"), &elevatorState)
	return elevatorState
}

func ReadOrders() OrderState {
	var orderState OrderState
	json.Unmarshal(read("filesystem/orders.json"), &orderState)
	return orderState
}

func read(filepath string) []byte {
	// fmt.Println("Filesystem/read.go")

	jsonFile, err := os.Open(filepath)

	if err != nil {
		fmt.Println(err)
	}
	
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}