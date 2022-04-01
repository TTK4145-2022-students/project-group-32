package filesystem

import (
	"elevators/orders"
	"encoding/json"
	"io/ioutil"
	"time"
)

func SaveOrdersPeriodically() {
	for {
		SaveOrders(orders.GetOrders())
		time.Sleep(saveToFileRate)
	}
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
