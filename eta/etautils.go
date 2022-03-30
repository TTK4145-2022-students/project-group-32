package eta

import (
	"elevators/hardware"
	"elevators/orders"
	"time"
)

func InternalETABest(
	order orders.OrderState,
	internalETA time.Time) bool {

	return order.BestETA.Equal(internalETA) &&
		!internalETA.IsZero()
}

func InternalETABestAndNotExpired(
	order orders.OrderState,
	internalETA time.Time,
	currentTime time.Time) bool {

	return InternalETABest(
		order,
		internalETA) &&
		order.BestETA.Before(currentTime)
}

func HasOrderAndBestETABetweenTimes(
	order orders.OrderState,
	startTime time.Time,
	endTime time.Time) bool {

	return order.HasOrder() &&
		startTime.Before(order.BestETA) &&
		order.BestETA.Before(endTime)
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
		return internalETAs.Down[floor]

	case hardware.MD_Up:
		return internalETAs.Up[floor]

	case hardware.MD_Stop:
		return internalETAs.Cab[floor]

	default:
		panic("Invalid direction to get eta")
	}
}

func (internalETAs *InternalETAs) ETABetweenTimes(
	direction hardware.MotorDirection,
	floor int,
	startTime time.Time,
	endTime time.Time) bool {

	return startTime.Before(internalETAs.getETA(
		direction,
		floor)) &&
		internalETAs.getETA(
			direction,
			floor).Before(endTime)
}

func (internalETAs *InternalETAs) setETA(
	direction hardware.MotorDirection,
	floor int,
	eta time.Time) {

	switch direction {
	case hardware.MD_Down:
		internalETAs.Down[floor] = eta

	case hardware.MD_Up:
		internalETAs.Up[floor] = eta

	case hardware.MD_Stop:
		internalETAs.Cab[floor] = eta

	default:
		panic("Invalid direction to get eta")
	}
}

func maxTime() time.Time {
	return time.Unix(
		1<<62,
		0)
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
		if InternalETABestWithOrder(
			allOrders.Up[floor],
			allETAs.Up[floor]) ||
			InternalETABestWithOrder(
				allOrders.Down[floor],
				allETAs.Down[floor]) ||
			allOrders.Cab[floor] {
			return true
		}
	}
	return false
}

func InternalETABestWithOrder(
	order orders.OrderState,
	eta time.Time) bool {

	return order.HasOrder() &&
		InternalETABest(
			order,
			eta)
}

func inputBestETABetterOrBestETAExpired(
	inputOrder orders.OrderState,
	currentOrder orders.OrderState,
	currentTime time.Time) bool {

	return newETABetterOrBestETAExpired(
		currentOrder,
		inputOrder.BestETA,
		currentTime)
}

func inputBestETAExpired(
	inputOrder orders.OrderState,
	currentTime time.Time) bool {

	return inputOrder.BestETA.Before(currentTime)
}

func noETA(eta time.Time) bool {
	return eta.Equal(time.Time{})
}
