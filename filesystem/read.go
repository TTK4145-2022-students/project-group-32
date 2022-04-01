package filesystem

import (
	"elevators/orders"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadOrders() orders.AllOrders {
	var orderState orders.AllOrders
	json.Unmarshal(
		read(orderFile),
		&orderState)
	return orderState
}

func read(filepath string) []byte {
	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
