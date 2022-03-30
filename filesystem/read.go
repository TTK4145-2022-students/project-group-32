package filesystem

import (
	"elevators/cab"
	"elevators/orders"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadCabState() cab.CabState {
	var cabState cab.CabState
	json.Unmarshal(
		read(
			cabFile),
		&cabState)
	return cabState
}

func ReadOrders() orders.AllOrders {
	var orderState orders.AllOrders
	json.Unmarshal(
		read(
			orderFile),
		&orderState)
	return orderState
}

func read(filepath string) []byte {
	jsonFile, _ := os.Open(filepath)
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
