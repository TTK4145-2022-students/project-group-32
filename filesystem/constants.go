package filesystem

import "os"

var cabFile = "filesystem/cabState.json"
var orderFile = "filesystem/orderState.json"

func Init() {
	if len(os.Args) > 1 {
		cabFile = "filesystem/cabState" + os.Args[1] + ".json"
		orderFile = "filesystem/orderState" + os.Args[1] + ".json"
	}
}
