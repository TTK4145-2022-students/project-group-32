package orderstate

import (
	"elevators/hardware"
	"time"
)

func (orders *AllOrders) getOrderState(
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
		panic("Invalid direction to get order")
	}
}

func newETABetterOrBestETAExpired(
	order OrderState,
	newETA time.Time,
	currentTime time.Time) bool {

	return !newETA.IsZero() &&
		(newETA.Before(order.BestETA) ||
			order.BestETA.Before(currentTime))
}

func newETABetterOrBestETAExpiredWithOrder(
	order OrderState,
	newETA time.Time,
	currentTime time.Time) bool {

	return !newETA.IsZero() &&
		(newETA.Before(order.BestETA) ||
			(order.hasOrder() &&
				order.BestETA.Before(currentTime)))
}

// func inputBestETABetterOrBestETAExpired(
// 	inputOrder OrderState,
// 	currentOrder OrderState,
// 	currentTime time.Time) bool {

// 	return newETABetterOrBestETAExpired(
// 		inputOrder.BestETA,
// 		currentOrder,
// 		currentTime)
// }

func inputBestETAExpired(
	inputOrder OrderState,
	currentTime time.Time) bool {

	return inputOrder.BestETA.Before(currentTime)
}
