package filesystem

import (
	// "fmt"
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"encoding/json"
	"io/ioutil"
	"sync"
	"time"
)

var fileMtx = new(sync.RWMutex)

func SaveStatePeriodically() {
	for {
		SaveCabState(cabstate.Cab)
		SaveOrders(orderstate.GetOrders())
		time.Sleep(time.Millisecond * 50)
	}
}

func SaveCabState(cabState cabstate.CabState) {
	fileMtx.Lock()
	defer fileMtx.Unlock()
	write("filesystem/cabState.json", cabState)
}

func SaveOrders(orderState orderstate.AllOrders) {
	fileMtx.Lock()
	defer fileMtx.Unlock()
	write("filesystem/orderState.json", orderState)
}

func write(filepath string, elevatorState interface{}) {
	// // fmt.Println("Filesystem/write.go")

	file, _ := json.MarshalIndent(elevatorState, "", " ")
	_ = ioutil.WriteFile(filepath, file, 0644)
}
