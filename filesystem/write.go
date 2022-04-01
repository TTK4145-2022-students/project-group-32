package filesystem

import (
	"elevators/orders"
	"encoding/json"
	"io/ioutil"
	"time"
)

func SaveOrdersPeriodically() {
	for {
		saveOrders(orders.GetOrders())
		time.Sleep(saveToFileRate)
	}
}

func saveOrders(allOrders orders.AllOrders) {
	write(
		orderFile,
		allOrders)
}

func write(
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
