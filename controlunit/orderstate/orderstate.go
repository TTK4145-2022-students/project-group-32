package orderstate

import (
	"elevators/hardware"
	"elevators/hardware/cab"
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
var cabOrders [hardware.FloorCount]bool

func InitOrderState() {
	upOrders := new([hardware.FloorCount]OrderState)
	downOrders := new([hardware.FloorCount]OrderState)
	cabOrders := new([hardware.FloorCount]bool)

	_ = upOrders
	_ = downOrders
	_ = cabOrders
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
			if OrderInFloor(floor, cab.Up) {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if OrderInFloor(floor, cab.Down) {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func OrderInFloor(floor int, direction cab.Direction) bool {
	switch direction {
	case cab.Up:
		return upOrders[floor].isOrder || cabOrders[floor]
	case cab.Down:
		return downOrders[floor].isOrder || cabOrders[floor]
	default:
		panic("direction not implemented " + string(rune(direction)))
	}
}

func AnyOrders() bool {
	for _, order := range upOrders {
		if order.isOrder {
			return true
		}
	}
	for _, order := range downOrders {
		if order.isOrder {
			return true
		}
	}
	for _, order := range cabOrders {
		if order {
			return true
		}
	}
	return false
}

func OrdersAtOrAbove(currentFloor int) bool {
	for floor := currentFloor; floor < hardware.FloorCount; floor++ {
		if OrderInFloor(floor, cab.Up) {
			return true
		}
	}
	return false
}

func OrdersAtOrBelow(currentFloor int) bool {
	for floor := currentFloor; floor >= 0; floor-- {
		if OrderInFloor(floor, cab.Down) {
			return true
		}
	}
	return false
}
