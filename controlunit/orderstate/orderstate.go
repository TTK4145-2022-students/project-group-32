package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"

	// "fmt"
	"sync"
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
	LocalETA	time.Time
}

type AllOrders struct {
	Up   [hardware.FloorCount]OrderState
	Down [hardware.FloorCount]OrderState
	Cab  [hardware.FloorCount]bool
}

var allOrders AllOrders
var allOrdersMtx = new(sync.RWMutex)

var MaxTime = time.Unix(1<<63-1, 999999999)

func ResetOrders(){
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	allOrders = AllOrders{}
}


func Init(orderState AllOrders) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	allOrders = orderState
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if allOrders.Cab[floor] {
			hardware.SetButtonLamp(hardware.BT_Cab, floor, true)
		}
		if floor < hardware.FloorCount-1 {
			if hasOrder(allOrders.Up[floor]) {
				hardware.SetButtonLamp(hardware.BT_HallUp, floor, true)
			}
			if hasOrder(allOrders.Down[floor]) {
				hardware.SetButtonLamp(hardware.BT_HallDown, floor, true)
			}
		}
	}
}

func GetOrders() AllOrders {
	allOrdersMtx.RLock()
	defer allOrdersMtx.RUnlock()
	return allOrders
}

func AcceptNewOrder(orderType hardware.ButtonType, floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
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
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(hardware.BT_Cab, floor, false)
	allOrders.Cab[floor] = false
}

func clearUpOrder(floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(hardware.BT_HallUp, floor, false)
	allOrders.Up[floor].LastCompleteTime = time.Now()
}

func clearDownOrder(floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(hardware.BT_HallDown, floor, false)
	allOrders.Down[floor].LastCompleteTime = time.Now()
}

func updateFloorOrderState(inputState OrderState, currentState *OrderState) OrderChange {
	currentOrder := hasOrder(*currentState)
	if inputState.LastOrderTime.After(currentState.LastOrderTime) {
		currentState.LastOrderTime = inputState.LastOrderTime
	}
	// fmt.Println(currentState.LastOrderTime.String())
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
