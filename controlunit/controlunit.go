package controlunit

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/eta"
	"elevators/controlunit/orderstate"
	"fmt"
)

func Init() {
	orderstate.InitOrderState()
	cabstate.InitCabState()
	newETA := eta.ComputeETA(2, 3)
	// eta.ComputeETA(2, 3)
	// fmt.Println(fmt.Println(newETA.String()))
	fmt.Println(newETA.String())
}
