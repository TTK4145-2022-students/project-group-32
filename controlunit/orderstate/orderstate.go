package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"fmt"
	"time"
)

type OrderChange int

const (
	NoChange OrderChange = iota
	OrderCleared
	OrderPlaced
)

type OrderState struct {
	LastOrderTime    time.Time
	LastCompleteTime time.Time
	BestETA          time.Time
}

type AllOrders struct {
	Up   [hardware.FloorCount]OrderState
	Down [hardware.FloorCount]OrderState
	Cab  [hardware.FloorCount]bool
}

var allOrders AllOrders

func Init(orderState AllOrders) {
	allOrders = orderState
}

func GetOrders() AllOrders {
	return allOrders
}

func AcceptNewOrder(orderType hardware.ButtonType, floor int) {
	switch orderType {
	case hardware.BT_HallUp:
		allOrders.Up[floor].LastOrderTime = time.Now()
	case hardware.BT_HallDown:
		allOrders.Down[floor].LastOrderTime = time.Now()
	case hardware.BT_Cab:
		allOrders.Cab[floor] = true
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
	allOrders.Cab[floor] = false
}

func clearUpOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
	allOrders.Up[floor].LastCompleteTime = time.Now()
}

func clearDownOrder(floor int) {
	hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
	allOrders.Down[floor].LastCompleteTime = time.Now()
}

func updateFloorOrderState(inputState OrderState, currentState *OrderState) OrderChange {
	currentOrder := hasOrder(*currentState)
	if inputState.LastOrderTime.After(currentState.LastOrderTime) {
		currentState.LastOrderTime = inputState.LastOrderTime
	}
	fmt.Println(currentState.LastOrderTime.String())
	if inputState.LastCompleteTime.After(currentState.LastCompleteTime) {
		currentState.LastCompleteTime = inputState.LastCompleteTime
	}
	if inputState.BestETA.Before(currentState.BestETA) && inputState.BestETA.After(time.Now()) {
		currentState.BestETA = inputState.BestETA
	}

	newCurrentOrder := hasOrder(*currentState)
	orderChange := NoChange
	if newCurrentOrder && !currentOrder {
		orderChange = OrderPlaced
	} else if !newCurrentOrder && currentOrder {
		orderChange = OrderCleared
	}
	return orderChange
}

func hasOrder(inputState OrderState) bool {
	return inputState.LastOrderTime.After(inputState.LastCompleteTime)
}

func UpdateOrders(inputOrders AllOrders) [hardware.FloorCount]bool {
	var newOrders [hardware.FloorCount]bool
	for floor := 0; floor < hardware.FloorCount; floor++ {
		switch updateFloorOrderState(inputOrders.Down[floor], &allOrders.Down[floor]) {
		case OrderCleared:
			hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
		case OrderPlaced:
			hardware.SetButtonLamp(hardware.BT_HallDown, floor, true)
			newOrders[floor] = true
		}
		switch updateFloorOrderState(inputOrders.Up[floor], &allOrders.Up[floor]) {
		case OrderCleared:
			hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
		case OrderPlaced:
			hardware.SetButtonLamp(hardware.BT_HallUp, floor, true)
			newOrders[floor] = true
		}
	}
	return newOrders
}

func UpdateETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int) {
	newDurations := ComputeDurations(currentFloor, recentDirection, allOrders, allETAs)
	newETAs := ComputeInternalETAs(newDurations)
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if newDurations.Up[floor] < allDurations.Up[floor] && newETAs.Up[floor].Before(allOrders.Up[floor].BestETA) {
			allOrders.Up[floor].BestETA = newETAs.Up[floor]
		} else if allETAs.Up[floor].Equal(allOrders.Up[floor].BestETA) {
			// Make sure to keep ownership
			newETAs.Up[floor] = allETAs.Up[floor]
		}

		if newDurations.Down[floor] < allDurations.Down[floor] && newETAs.Down[floor].Before(allOrders.Down[floor].BestETA) {
			allOrders.Down[floor].BestETA = newETAs.Down[floor]
		} else if allETAs.Down[floor].Equal(allOrders.Down[floor].BestETA) {
			// Make sure to keep ownership
			newETAs.Down[floor] = allETAs.Down[floor]
		}
	}
	allDurations = newDurations
	allETAs = newETAs
}

func OrdersBetween(orders AllOrders, startFloor int, destinationFloor int) int {
	if startFloor == destinationFloor {
		return 0
	}
	ordersBetweenCount := 0
	if startFloor < destinationFloor {
		for floor := startFloor; floor < destinationFloor; floor++ {
			if hasOrder(orders.Up[floor]) || orders.Cab[floor] {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if hasOrder(orders.Down[floor]) || orders.Cab[floor] {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func OrdersAbove(orders AllOrders, currentFloor int) bool {
	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if hasOrder(orders.Up[floor]) || hasOrder(orders.Down[floor]) || orders.Cab[floor] {
			return true
		}
	}
	return false
}

func CabOrdersAbove(cabOrders [hardware.FloorCount]bool, currentFloor int) bool {
	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if cabOrders[floor] {
			return true
		}
	}
	return false
}

func OrdersBelow(orders AllOrders, currentFloor int) bool {
	for floor := currentFloor - 1; floor >= 0; floor-- {
		if hasOrder(orders.Up[floor]) || hasOrder(orders.Down[floor]) || orders.Cab[floor] {
			return true
		}
	}
	return false
}

func CabOrdersBelow(cabOrders [hardware.FloorCount]bool, currentFloor int) bool {
	for floor := currentFloor - 1; floor >= 0; floor-- {
		if cabOrders[floor] {
			return true
		}
	}
	return false
}

func GetOrderStatus(orders AllOrders, floor int) prioritize.OrderStatus {
	var orderStatus prioritize.OrderStatus
	orderStatus.UpAtFloor = hasOrder(orders.Up[floor])
	orderStatus.DownAtFloor = hasOrder(orders.Down[floor])
	orderStatus.CabAtFloor = orders.Cab[floor]
	orderStatus.AboveFloor = OrdersAbove(orders, floor)
	orderStatus.BelowFloor = OrdersBelow(orders, floor)
	return orderStatus
}
