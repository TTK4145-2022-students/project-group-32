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
	Write(cabFile, cabState)
}

func SaveOrders(orders orderstate.AllOrders) {
	Write(orderFile, orders)
}

func Write(filepath string, state interface{}) {
	file, _ := json.MarshalIndent(state, "", " ")
	_ = ioutil.WriteFile(filepath, file, 0644)
}
