package orderstate

import (
	"elevators/controlunit/cabstate"
	"elevators/hardware"
	"fmt"
	"time"
)

type OrderState struct {
	isOrder          bool
	lastOrderTime    time.Time
	lastCompleteTime time.Time
	bestETA          time.Time
}

var upOrders [hardware.FloorCount]OrderState
var downOrders [hardware.FloorCount]OrderState

func InitOrderState() {
	upOrders := new([hardware.FloorCount]OrderState)
	downOrders := new([hardware.FloorCount]OrderState)

	_ = upOrders
	_ = downOrders

	// fmt.Println(upOrders[0].lastOrderTime.String())
	// fmt.Println(downOrders[0].bestETA.String())

	// updateUpFloorOrderState(OrderState{true, time.Now(), time.Now(), time.Now()}, &upOrders[0])

	// fmt.Println(upOrders[0].lastOrderTime.String())
}

func updateUpFloorOrderState(inputState OrderState, currentState *OrderState) {
	if inputState.lastOrderTime.After(currentState.lastOrderTime) {
		currentState.lastOrderTime = inputState.lastOrderTime
	}
	fmt.Println(currentState.lastOrderTime.String())
	if inputState.lastCompleteTime.After(currentState.lastCompleteTime) {
		currentState.lastCompleteTime = inputState.lastCompleteTime
	}
	if inputState.bestETA.Before(currentState.bestETA) && inputState.bestETA.After(time.Now()) {
		currentState.bestETA = inputState.bestETA
	}
	if currentState.lastOrderTime.After(currentState.lastCompleteTime) {
		currentState.isOrder = true
	} else {
		currentState.isOrder = false
	}
}

func OrdersBetween(startFloor int, destinationFloor int) int {
	if startFloor == destinationFloor {
		return 0
	}
	ordersBetweenCount := 0
	if startFloor < destinationFloor {
		for floor := startFloor; floor < destinationFloor; floor++ {
			if orderInFloor(floor, cabstate.Up) {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if orderInFloor(floor, cabstate.Down) {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func orderInFloor(floor int, direction cabstate.Direction) bool {
	switch direction {
	case cabstate.Up:
		return upOrders[floor].isOrder
	case cabstate.Down:
		return downOrders[floor].isOrder
	default:
		panic("direction not implemented " + string(rune(direction)))
	}
}
