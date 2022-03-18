package filesystem

import (
	// "fmt"
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"encoding/json"
	"io/ioutil"
)

func SaveState(cabState cabstate.CabState, orderState orderstate.AllOrders) {
	write("filesystem/cabState.json", cabState)
	write("filesystem/orderState.json", orderState)
}

func write(filepath string, elevatorState interface{}) {
	// fmt.Println("Filesystem/write.go")

	file, _ := json.MarshalIndent(elevatorState, "", " ")
	_ = ioutil.WriteFile(filepath, file, 0644)
}
