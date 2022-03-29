package filesystem

import (
	"os"
	"time"
)

var cabFile = "filesystem/cabState.json"
var orderFile = "filesystem/orderState.json"
var _saveToFileRate = time.Millisecond * 500

func Init() {
	if len(os.Args) > 1 {
		cabFile = "filesystem/cabState" + os.Args[1] + ".json"
		orderFile = "filesystem/orderState" + os.Args[1] + ".json"
	}
}
