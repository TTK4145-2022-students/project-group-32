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

func CompleteOrder(floor int,
	recentDirection hardware.MotorDirection,
	upOrdersInFloor bool,
	downOrdersInFloor bool,
	ordersAbove bool,
	ordersBelow bool) {

	clearCabOrder(floor)
	switch recentDirection {
	case hardware.MD_Up:
		clearUpOrder(floor)
		if downOrdersInFloor && (!ordersAbove && !upOrdersInFloor) {
			clearDownOrder(floor)
		}
	case hardware.MD_Down:
		clearDownOrder(floor)
		if upOrdersInFloor && (!ordersBelow && !downOrdersInFloor) {
			clearUpOrder(floor)
		}
	default:
		panic("Invalid recent direction on floor stop")
	}
}

func clearCabOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_Cab, floor, false)
	cabOrders[floor] = false
}

func clearUpOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
	upOrders[floor].lastCompleteTime = time.Now()
	upOrders[floor].isOrder = false
}

func clearDownOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
	downOrders[floor].lastCompleteTime = time.Now()
	downOrders[floor].isOrder = false
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
		return UpOrdersInFloor(floor) || CabOrdersInFloor(floor)
	case hardware.MD_Down:
		return DownOrdersInFloor(floor) || CabOrdersInFloor(floor)
	default:
		panic("direction not implemented " + string(rune(direction)))
	}
}

func OrdersInFloor(floor int) bool {
	return UpOrdersInFloor(floor) || DownOrdersInFloor(floor) || CabOrdersInFloor(floor)
}

func DownOrdersInFloor(floor int) bool {
	return downOrders[floor].isOrder
}

func UpOrdersInFloor(floor int) bool {
	return upOrders[floor].isOrder
}

func CabOrdersInFloor(floor int) bool {
	return cabOrders[floor]
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

func OrdersAbove(currentFloor int) bool {
	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if OrdersInFloor(floor) {
			return true
		}
	}
	return false
}

func OrdersBelow(currentFloor int) bool {
	for floor := currentFloor - 1; floor >= 0; floor-- {
		if OrdersInFloor(floor) {
			return true
		}
	}
	return false
}
