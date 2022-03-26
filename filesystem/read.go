package filesystem

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadCabState() cabstate.CabState {
	var cabState cabstate.CabState
	json.Unmarshal(read(cabFile), &cabState)
	return cabState
}

func ReadOrders() orderstate.AllOrders {
	var orderState orderstate.AllOrders
	json.Unmarshal(read(orderFile), &orderState)
	return orderState
}

func read(filepath string) []byte {
	jsonFile, _ := os.Open(filepath)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	return byteValue
}
