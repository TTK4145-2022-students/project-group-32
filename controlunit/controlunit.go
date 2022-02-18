package controlunit

import (
	"elevators/controlunit/orderstate"
	"fmt"
)

func Init() {
	fmt.Println("Hello, world")
	orderstate.InitOrderState()
}
