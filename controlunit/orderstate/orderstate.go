package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"sync"
	"time"
)

const WaitBeforeGuaranteeTime = 200 * time.Millisecond

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
	LocalETA         time.Time
	Now              time.Time
}

type AllOrders struct {
	Up   [hardware.FloorCount]OrderState
	Down [hardware.FloorCount]OrderState
	Cab  [hardware.FloorCount]bool
}

var allOrders AllOrders

var allOrdersMtx = new(sync.RWMutex)

func Init(initOrders AllOrders) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	allOrders = initOrders
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if allOrders.Cab[floor] {
			hardware.SetButtonLamp(
				hardware.BT_Cab,
				floor,
				true)
		}
		if allOrders.Up[floor].hasOrder() {
			hardware.SetButtonLamp(
				hardware.BT_HallUp,
				floor,
				true)
		}
		if allOrders.Down[floor].hasOrder() {
			hardware.SetButtonLamp(
				hardware.BT_HallDown,
				floor,
				true)
		}
	}
}

func ResetOrders() {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	allOrders = AllOrders{}
}

func GetOrders() AllOrders {
	allOrdersMtx.RLock()
	defer allOrdersMtx.RUnlock()
	return allOrders
}

func AcceptNewOrder(
	orderType hardware.ButtonType,
	floor int) {

	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	switch orderType {
	case hardware.BT_HallUp:
		allOrders.Up[floor].LastOrderTime = time.Now()
	case hardware.BT_HallDown:
		allOrders.Down[floor].LastOrderTime = time.Now()
	case hardware.BT_Cab:
		allOrders.Cab[floor] = true
	}
	go waitForOrderGuarantee(
		orderType,
		floor)
}

func waitForOrderGuarantee(
	orderType hardware.ButtonType,
	floor int) {

	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	time.Sleep(WaitBeforeGuaranteeTime)
	switch orderType {
	case hardware.BT_HallUp:
		if allOrders.Up[floor].hasOrder() {
			hardware.SetButtonLamp(
				orderType,
				floor,
				true)
		}
	case hardware.BT_HallDown:
		if allOrders.Down[floor].hasOrder() {
			hardware.SetButtonLamp(
				orderType,
				floor,
				true)
		}
	case hardware.BT_Cab:
		if allOrders.Cab[floor] {
			hardware.SetButtonLamp(
				orderType,
				floor,
				true)
		}
	}
}

func CompleteCabAndUpOrder(floor int) {
	clearCabOrder(floor)
	clearUpOrder(floor)
}

func CompleteCabAndDownOrder(floor int) {
	clearCabOrder(floor)
	clearDownOrder(floor)
}

func CompleteCabOrder(floor int) {
	clearCabOrder(floor)
}

func clearCabOrder(floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(
		hardware.BT_Cab,
		floor,
		false)
	allOrders.Cab[floor] = false
}

func clearUpOrder(floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(
		hardware.BT_HallUp,
		floor,
		false)
	allOrders.Up[floor].LastCompleteTime = time.Now()
}

func clearDownOrder(floor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	hardware.SetButtonLamp(
		hardware.BT_HallDown,
		floor,
		false)
	allOrders.Down[floor].LastCompleteTime = time.Now()
}

func (orderState *OrderState) hasOrder() bool {
	return orderState.LastOrderTime.After(orderState.LastCompleteTime)
}

func AnyOrders(orders AllOrders) bool {
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if orders.Up[floor].hasOrder() ||
			orders.Down[floor].hasOrder() ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func UpdateOrders(inputOrders AllOrders) [hardware.FloorCount]bool {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()
	var newOrders [hardware.FloorCount]bool
	for floor := 0; floor < hardware.FloorCount; floor++ {
		switch updateFloorOrderState(
			inputOrders.Down[floor],
			&allOrders.Down[floor]) {

		case OrderCleared:
			hardware.SetButtonLamp(
				hardware.BT_HallDown,
				floor,
				false)
		case OrderPlaced:
			hardware.SetButtonLamp(
				hardware.BT_HallDown,
				floor,
				true)
			newOrders[floor] = true
		}

		switch updateFloorOrderState(
			inputOrders.Up[floor],
			&allOrders.Up[floor]) {

		case OrderCleared:
			hardware.SetButtonLamp(
				hardware.BT_HallUp,
				floor,
				false)
		case OrderPlaced:
			hardware.SetButtonLamp(
				hardware.BT_HallUp,
				floor,
				true)
			newOrders[floor] = true
		}
	}
	return newOrders
}

func updateFloorOrderState(
	inputState OrderState,
	currentState *OrderState) OrderChange {

	currentOrder := currentState.hasOrder()
	now := time.Now()
	if inputState.LastOrderTime.After(currentState.LastOrderTime) {
		currentState.LastOrderTime = inputState.LastOrderTime
	}
	if inputState.LastCompleteTime.After(currentState.LastCompleteTime) {
		currentState.LastCompleteTime = inputState.LastCompleteTime
	}
	if inputState.BestETA.After(now) &&
		(inputState.BestETA.Before(currentState.BestETA) ||
			currentState.BestETA.Before(now)) {
		currentState.BestETA = inputState.BestETA
	}

	newCurrentOrder := currentState.hasOrder()
	orderChange := NoChange
	if newCurrentOrder &&
		!currentOrder {
		orderChange = OrderPlaced
	} else if !newCurrentOrder &&
		currentOrder {
		orderChange = OrderCleared
	}
	return orderChange
}

func OrdersBetween(
	orders AllOrders,
	startFloor int,
	destinationFloor int) int {

	if startFloor == destinationFloor {
		return 0
	}
	ordersBetweenCount := 0
	if startFloor < destinationFloor {
		for floor := startFloor; floor < destinationFloor; floor++ {
			if orders.Up[floor].hasOrder() ||
				orders.Cab[floor] {
				ordersBetweenCount++
			}
		}
	} else {
		for floor := startFloor; floor > destinationFloor; floor-- {
			if orders.Down[floor].hasOrder() ||
				orders.Cab[floor] {
				ordersBetweenCount++
			}
		}
	}
	return ordersBetweenCount
}

func OrdersAbove(
	orders AllOrders,
	currentFloor int) bool {

	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if orders.Up[floor].hasOrder() ||
			orders.Down[floor].hasOrder() ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func OrdersBelow(
	orders AllOrders,
	currentFloor int) bool {

	for floor := currentFloor - 1; floor >= 0; floor-- {
		if orders.Up[floor].hasOrder() ||
			orders.Down[floor].hasOrder() ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func CabOrdersAbove(
	cabOrders [hardware.FloorCount]bool,
	currentFloor int) bool {

	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if cabOrders[floor] {
			return true
		}
	}
	return false
}

func CabOrdersBelow(
	cabOrders [hardware.FloorCount]bool,
	currentFloor int) bool {

	for floor := currentFloor - 1; floor >= 0; floor-- {
		if cabOrders[floor] {
			return true
		}
	}
	return false
}

func GetOrderSummary(
	orders AllOrders,
	floor int) prioritize.OrderSummary {

	var orderSummary prioritize.OrderSummary
	orderSummary.UpAtFloor = orders.Up[floor].hasOrder()
	orderSummary.DownAtFloor = orders.Down[floor].hasOrder()
	orderSummary.CabAtFloor = orders.Cab[floor]
	orderSummary.AboveFloor = OrdersAbove(
		orders,
		floor)
	orderSummary.BelowFloor = OrdersBelow(
		orders,
		floor)
	return orderSummary
}
