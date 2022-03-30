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

func (orders *AllOrders) setOrderETA(
	direction hardware.MotorDirection,
	floor int,
	eta time.Time) {

	switch direction {
	case hardware.MD_Down:
		orders.Down[floor].BestETA = eta

	case hardware.MD_Up:
		orders.Up[floor].BestETA = eta

	default:
		panic("Invalid direction to set order")
	}
}
