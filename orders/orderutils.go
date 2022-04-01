package orders

import (
	"elevators/hardware"
	"time"
)

func (orderState *OrderState) HasOrder() bool {
	return orderState.LastOrderTime.After(orderState.LastCompleteTime)
}

func (orders *AllOrders) GetOrderState(
	direction hardware.MotorDirection,
	floor int) OrderState {

	switch direction {
	case hardware.MD_Down:
		return orders.Down[floor]

	case hardware.MD_Up:
		return orders.Up[floor]

	default:
		panic("Invalid direction to get order")
	}
}

func inputETABetterOrCurrentETAExpired(
	inputOrder OrderState,
	currentOrder OrderState,
	currentTime time.Time) bool {

	return (inputOrder.BestETA.Before(currentOrder.BestETA) ||
		currentOrder.BestETA.Before(currentTime)) &&
		inputOrder.BestETA.After(currentTime)
}
