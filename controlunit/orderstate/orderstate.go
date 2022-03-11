package orderstate

import (
	"elevators/hardware"
	"fmt"
	"time"
)

type OrderState struct {
	lastOrderTime    time.Time
	lastCompleteTime time.Time
	bestETA          time.Time
}

type AllOrders struct {
	up   [hardware.FloorCount]OrderState
	down [hardware.FloorCount]OrderState
	cab  [hardware.FloorCount]bool
}

type OrderStatus struct {
	UpAtFloor   bool
	DownAtFloor bool
	CabAtFloor  bool
	AboveFloor  bool
	BelowFloor  bool
}

var allOrders AllOrders

func InitOrders() {
	allOrders := new(AllOrders)

	_ = allOrders
}

func GetOrders() AllOrders {
	return allOrders
}

func AcceptNewOrder(orderType hardware.ButtonType, floor int) {
	switch orderType {
	case hardware.BT_HallUp:
		allOrders.up[floor].lastOrderTime = time.Now()
	case hardware.BT_HallDown:
		allOrders.down[floor].lastOrderTime = time.Now()
	case hardware.BT_Cab:
		allOrders.cab[floor] = true
	default:
		panic("order type not implemented " + string(rune(orderType)))
	}
	hardware.SetButtonLamp(orderType, floor, true)
}

func CompleteOrderCabAndUp(floor int) {
	clearCabOrder(floor)
	clearUpOrder(floor)
}

func CompleteOrderCabAndDown(floor int) {
	clearCabOrder(floor)
	clearDownOrder(floor)
}

func CompleteOrderCab(floor int) {
	clearCabOrder(floor)
}

func clearCabOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_Cab, floor, false)
	allOrders.cab[floor] = false
}

func clearUpOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
	allOrders.up[floor].lastCompleteTime = time.Now()
}

func clearDownOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
	allOrders.down[floor].lastCompleteTime = time.Now()
}

func updateFloorOrderState(inputState OrderState, currentState *OrderState) {
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
}

func hasOrder(inputState OrderState) bool {
	return inputState.lastOrderTime.After(inputState.lastCompleteTime)
}

func OrdersBetween(orders AllOrders, startFloor int, destinationFloor int) int {
	if startFloor == destinationFloor {
		return 0
	}
	ordersBetweenCount := 0
	if startFloor < destinationFloor {
		for floor := startFloor; floor < destinationFloor; floor++ {
			if hasOrder(orders.up[floor]) || orders.cab[floor] {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if hasOrder(orders.down[floor]) || orders.cab[floor] {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func OrdersAbove(orders AllOrders, currentFloor int) bool {
	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if hasOrder(orders.up[floor]) || hasOrder(orders.down[floor]) || orders.cab[floor] {
			return true
		}
	}
	return false
}

func OrdersBelow(orders AllOrders, currentFloor int) bool {
	for floor := currentFloor - 1; floor >= 0; floor-- {
		if hasOrder(orders.up[floor]) || hasOrder(orders.down[floor]) || orders.cab[floor] {
			return true
		}
	}
	return false
}

func GetOrderStatus(orders AllOrders, floor int) OrderStatus {
	var orderStatus OrderStatus
	orderStatus.UpAtFloor = hasOrder(orders.up[floor])
	orderStatus.DownAtFloor = hasOrder(orders.down[floor])
	orderStatus.CabAtFloor = orders.cab[floor]
	orderStatus.AboveFloor = OrdersAbove(orders, floor)
	orderStatus.BelowFloor = OrdersBelow(orders, floor)
	return orderStatus
}
