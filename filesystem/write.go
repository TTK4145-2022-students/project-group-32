package filesystem

import (
	"elevators/cab"
	"elevators/orders"
	"encoding/json"
	"io/ioutil"
	"time"
)

func SaveOrdersPeriodically() {
	for {
		SaveOrders(orders.GetOrders())
		time.Sleep(_saveToFileRate)
	}
}

func SaveCabState(cabState cab.CabState) {
	Write(
		cabFile,
		cabState)
}

func SaveOrders(allOrders orders.AllOrders) {
	Write(
		orderFile,
		allOrders)
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
