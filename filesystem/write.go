package filesystem

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"encoding/json"
	"io/ioutil"
	"time"
)

func SaveOrdersPeriodically() {
	for {
		SaveOrders(orderstate.GetOrders())
		time.Sleep(_saveToFileRate)
	}
}

func SaveCabState(cabState cabstate.CabState) {
	Write(
		cabFile,
		cabState)
}

func SaveOrders(orders orderstate.AllOrders) {
	Write(
		orderFile,
		orders)
}

func Write(
	filepath string,
	state interface{}) {

	file, _ := json.MarshalIndent(
		state,
		"",
		" ")
	_ = ioutil.WriteFile(
		filepath,
		file,
		0644)
}
