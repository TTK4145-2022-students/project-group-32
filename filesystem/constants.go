package filesystem

import (
	"elevators/orders"
	"os"
	"time"
)

var cabFile = "filesystem/cabState.json"
var orderFile = "filesystem/orderState.json"
var saveToFileRate = time.Millisecond * orders.WaitBeforeGuaranteeTime / 2

func Init() {
	if len(os.Args) > 1 {
		cabFile = "filesystem/cabState" + os.Args[1] + ".json"
		orderFile = "filesystem/orderState" + os.Args[1] + ".json"
	}
}
