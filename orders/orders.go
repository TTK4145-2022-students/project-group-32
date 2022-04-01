package orders

import (
	"elevators/hardware"
	"elevators/prioritize"
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
}

type AllOrders struct {
	Up   [hardware.FloorCount]OrderState
	Down [hardware.FloorCount]OrderState
	Cab  [hardware.FloorCount]bool
}

var allOrders AllOrders
var allOrdersMtx = new(sync.RWMutex)

func Init(orderState AllOrders) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()

	allOrders = orderState

	for _, floor := range hardware.ValidFloors() {
		if allOrders.Cab[floor] {
			hardware.SetButtonLamp(
				hardware.BT_Cab,
				floor,
				true)
		}
		if allOrders.Up[floor].HasOrder() {
			hardware.SetButtonLamp(
				hardware.BT_HallUp,
				floor,
				true)
		}
		if allOrders.Down[floor].HasOrder() {
			hardware.SetButtonLamp(
				hardware.BT_HallDown,
				floor,
				true)
		}
	}
}

func GetOrders() AllOrders {
	allOrdersMtx.RLock()
	defer allOrdersMtx.RUnlock()

	return allOrders
}

func SetOrders(orders AllOrders) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()

	allOrders = orders
}

func ResetOrders() {
	SetOrders(AllOrders{})
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

	default:
		panic("order type not implemented " + string(rune(orderType)))
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
		if allOrders.Up[floor].HasOrder() {
			hardware.SetButtonLamp(
				orderType,
				floor,
				true)
		}

	case hardware.BT_HallDown:
		if allOrders.Down[floor].HasOrder() {
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

func updateFloorOrderState(
	inputState OrderState,
	currentState *OrderState) OrderChange {

	currentTime := time.Now()

	currentOrder := currentState.HasOrder()
	if inputState.LastOrderTime.After(currentState.LastOrderTime) {
		currentState.LastOrderTime = inputState.LastOrderTime
	}
	if inputState.LastCompleteTime.After(currentState.LastCompleteTime) {
		currentState.LastCompleteTime = inputState.LastCompleteTime
	}
	if inputETABetterOrCurrentETAExpired(
		inputState,
		*currentState,
		currentTime) {

		currentState.BestETA = inputState.BestETA
	}

	newCurrentOrder := currentState.HasOrder()
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

func AnyOrders(orders AllOrders) bool {
	for _, floor := range hardware.ValidFloors() {

		if orders.Up[floor].HasOrder() ||
			orders.Down[floor].HasOrder() ||
			orders.Cab[floor] {

			return true

		}
	}
	return false
}

func UpdateOrders(inputOrders AllOrders) [hardware.FloorCount]bool {

	var newOrders [hardware.FloorCount]bool

	for _, floor := range hardware.ValidFloors() {
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

func OrdersAbove(
	orders AllOrders,
	currentFloor int) bool {

	for floor := currentFloor + 1; hardware.ValidFloor(floor); floor++ {

		if orders.Up[floor].HasOrder() ||
			orders.Down[floor].HasOrder() ||
			orders.Cab[floor] {

			return true

		}
	}
	return false
}

func CabOrdersAbove(
	cabOrders [hardware.FloorCount]bool,
	currentFloor int) bool {

	for floor := currentFloor + 1; hardware.ValidFloor(floor); floor++ {

		if cabOrders[floor] {

			return true

		}
	}
	return false
}

func OrdersBelow(
	orders AllOrders,
	currentFloor int) bool {

	for floor := currentFloor - 1; hardware.ValidFloor(floor); floor-- {

		if orders.Up[floor].HasOrder() ||
			orders.Down[floor].HasOrder() ||
			orders.Cab[floor] {

			return true

		}
	}
	return false
}

func CabOrdersBelow(
	cabOrders [hardware.FloorCount]bool,
	currentFloor int) bool {

	for floor := currentFloor - 1; hardware.ValidFloor(floor); floor-- {
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

	orderSummary.UpAtFloor = orders.Up[floor].HasOrder()
	orderSummary.DownAtFloor = orders.Down[floor].HasOrder()
	orderSummary.CabAtFloor = orders.Cab[floor]

	orderSummary.AboveFloor = OrdersAbove(
		orders,
		floor)
	orderSummary.BelowFloor = OrdersBelow(
		orders,
		floor)

	return orderSummary
}
