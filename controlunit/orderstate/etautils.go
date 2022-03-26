package orderstate

import (
	"elevators/hardware"
	"time"
)

func newETABetter(
	newETA time.Time,
	order OrderState,
	currentTime time.Time) bool {

	return !newETA.IsZero() &&
		(newETA.Before(order.BestETA) ||
			order.BestETA.Before(currentTime))
}

func InternalETABest(
	orderState OrderState,
	internalETA time.Time) bool {

	return orderState.BestETA.Equal(internalETA) &&
		!internalETA.IsZero()
}

func InternalETABestAndNotExpired(
	newETA time.Time,
	order OrderState,
	currentTime time.Time) bool {

	return !newETA.Equal(order.BestETA) &&
		order.BestETA.Before(currentTime)
}

func orderToServe(
	direction hardware.MotorDirection,
	orderUp OrderState,
	orderDown OrderState,
	orderCab bool) bool {

	return (direction == hardware.MD_Down &&
		orderDown.hasOrder()) ||
		(direction == hardware.MD_Up &&
			orderUp.hasOrder()) ||
		orderCab
}

func invalidFloor(floor int) bool {
	return floor < 0 || floor >= hardware.FloorCount

}
