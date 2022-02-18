package orderstate

import (
	"fmt"
	"time"
)

const floorCount = 4

type OrderState struct {
	isOrder          bool
	lastOrderTime    time.Time
	lastCompleteTime time.Time
	bestETA          time.Time
}

func InitOrderState() {
	upOrders := new([floorCount]OrderState)
	downOrders := new([floorCount]OrderState)

	// fmt.Println(upOrders[0].lastOrderTime.String())
	fmt.Println(downOrders[0].bestETA.String())

	updateUpFloorOrderState(OrderState{true, time.Now(), time.Now(), time.Now()}, &upOrders[0])

	fmt.Println(upOrders[0].lastOrderTime.String())
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
