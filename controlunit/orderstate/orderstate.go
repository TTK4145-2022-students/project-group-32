package orderstate

import (
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
var cabOrders [hardware.FloorCount]bool

func InitOrderState() {
	upOrders := new([hardware.FloorCount]OrderState)
	downOrders := new([hardware.FloorCount]OrderState)
	cabOrders := new([hardware.FloorCount]bool)

	_ = upOrders
	_ = downOrders
	_ = cabOrders
}

func AcceptNewOrder(orderType hardware.ButtonType, floor int) {
	switch orderType {
	case hardware.BT_HallUp:
		upOrders[floor].lastOrderTime = time.Now()
		upOrders[floor].isOrder = true
	case hardware.BT_HallDown:
		downOrders[floor].lastOrderTime = time.Now()
		downOrders[floor].isOrder = true
	case hardware.BT_Cab:
		cabOrders[floor] = true
	default:
		panic("order type not implemented " + string(rune(orderType)))
	}
	hardware.SetButtonLamp(orderType, floor, true)
}

func CompleteOrder(floor int) {
	// switch orderType {
	// case hardware.BT_HallUp:
	// 	upOrders[floor].lastCompleteTime = time.Now()
	// 	upOrders[floor].isOrder = false
	// case hardware.BT_HallDown:
	// 	downOrders[floor].lastCompleteTime = time.Now()
	// 	downOrders[floor].isOrder = false
	// case hardware.BT_Cab:
	// 	cabOrders[floor] = false
	// default:
	// 	panic("order type not implemented " + string(rune(orderType)))
	// }
	hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
	hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
	hardware.SetButtonLamp(hardware.BT_Cab, floor, false)
	upOrders[floor].lastCompleteTime = time.Now()
	upOrders[floor].isOrder = false
	downOrders[floor].lastCompleteTime = time.Now()
	downOrders[floor].isOrder = false
	cabOrders[floor] = false
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
			if OrderInFloor(floor, hardware.MD_Up) {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if OrderInFloor(floor, hardware.MD_Down) {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func OrderInFloor(floor int, direction hardware.MotorDirection) bool {
	switch direction {
	case hardware.MD_Up:
		return upOrders[floor].isOrder || cabOrders[floor]
	case hardware.MD_Down:
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
		if OrderInFloor(floor, hardware.MD_Up) {
			return true
		}
	}
	return false
}

func OrdersAtOrBelow(currentFloor int) bool {
	for floor := currentFloor; floor >= 0; floor-- {
		if OrderInFloor(floor, hardware.MD_Down) {
			return true
		}
	}
	return false
}
