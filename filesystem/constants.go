package filesystem

import (
	"elevators/orders"
	"os"
)

var (
	orderFile      = "filesystem/orderState.json"
	saveToFileRate = orders.WaitBeforeGuaranteeTime / 2
)

func Init() {
	if len(os.Args) > 1 {
		orderFile = "filesystem/orderState" + os.Args[1] + ".json"
	}
}
