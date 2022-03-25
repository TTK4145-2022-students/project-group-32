package filesystem

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"encoding/json"
	"io/ioutil"
	"time"
)

func SaveStatesPeriodically() {
	for {
		SaveCabState(cabstate.Cab)
		SaveOrders(orderstate.GetOrders())
		time.Sleep(time.Millisecond * 50)
	}
}

func SaveCabState(cabState cabstate.CabState) {
	write("filesystem/cabState.json", cabState)
}

func SaveOrders(orders orderstate.AllOrders) {
	write("filesystem/orderState.json", orders)
}

func write(filepath string, elevatorState interface{}) {
	file, _ := json.MarshalIndent(elevatorState, "", " ")
	_ = ioutil.WriteFile(filepath, file, 0644)
}
