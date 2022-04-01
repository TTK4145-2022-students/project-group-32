package eta

import (
	"elevators/hardware"
	"elevators/orders"
	"time"
)

func internalETABest(
	order orders.OrderState,
	internalETA time.Time) bool {

	return order.BestETA.Equal(internalETA) &&
		!internalETA.IsZero()
}

func orderToServe(
	direction hardware.MotorDirection,
	orderUp orders.OrderState,
	orderDown orders.OrderState,
	orderCab bool) bool {

	return (direction == hardware.MD_Down &&
		orderDown.HasOrder()) ||
		(direction == hardware.MD_Up &&
			orderUp.HasOrder())
}

func (internalETAs *InternalETAs) getETA(
	direction hardware.MotorDirection,
	floor int) time.Time {

	switch direction {
	case hardware.MD_Down:
		return internalETAs.down[floor]

	case hardware.MD_Up:
		return internalETAs.up[floor]

	case hardware.MD_Stop:
		return internalETAs.cab[floor]

	default:
		panic("Invalid direction to get eta")
	}
}

func (internalETAs *InternalETAs) setETA(
	direction hardware.MotorDirection,
	floor int,
	eta time.Time) {

	switch direction {
	case hardware.MD_Down:
		internalETAs.down[floor] = eta

	case hardware.MD_Up:
		internalETAs.up[floor] = eta

	case hardware.MD_Stop:
		internalETAs.cab[floor] = eta

	default:
		panic("Invalid direction to get eta")
	}
}

func newETABetterOrBestETAExpired(
	order orders.OrderState,
	newETA time.Time,
	currentTime time.Time) bool {

	return !newETA.IsZero() &&
		(newETA.Before(order.BestETA) ||
			order.BestETA.Before(currentTime))
}

func newETABetterOrBestETAExpiredWithOrder(
	order orders.OrderState,
	newETA time.Time,
	currentTime time.Time) bool {

	return (newETA.Before(order.BestETA) ||
		(order.HasOrder() &&
			currentTime.After(order.BestETA)))
}

func noETA(eta time.Time) bool {
	return eta.Equal(time.Time{})
}

func orderAndInternalETABest(
	direction hardware.MotorDirection,
	currentFloor int,
	allOrders orders.AllOrders,
	allETAs InternalETAs) bool {

	switch direction {
	case hardware.MD_Up:
		if allOrders.Up[currentFloor].HasOrder() {
			return true
		}

	case hardware.MD_Down:
		if allOrders.Down[currentFloor].HasOrder() {
			return true
		}
	}
	for floor := currentFloor + int(direction); hardware.ValidFloor(floor); floor += int(direction) {
		if internalETABestWithOrder(
			allOrders.Up[floor],
			allETAs.up[floor]) ||
			internalETABestWithOrder(
				allOrders.Down[floor],
				allETAs.down[floor]) ||
			allOrders.Cab[floor] {
			return true
		}
	}
	return false
}

func internalETABestWithOrder(
	order orders.OrderState,
	eta time.Time) bool {

	return order.HasOrder() &&
		internalETABest(
			order,
			eta)
}
