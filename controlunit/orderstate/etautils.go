package orderstate

import (
	"elevators/hardware"
	"time"
)

func InternalETABest(
	order OrderState,
	internalETA time.Time) bool {

	return order.BestETA.Equal(internalETA) &&
		!internalETA.IsZero()
}

func InternalETABestAndNotExpired(
	order OrderState,
	internalETA time.Time,
	currentTime time.Time) bool {

	return InternalETABest(
		order,
		internalETA) &&
		order.BestETA.Before(currentTime)
}

func hasOrderAndBestETABetweenTimes(
	order OrderState,
	startTime time.Time,
	endTime time.Time) bool {

	return order.hasOrder() &&
		startTime.Before(order.BestETA) &&
		order.BestETA.Before(endTime)
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
	// Max time to find minimum eta
	// the actual max time of 64-bit Time can't compare because of overflow
	return time.Unix(1<<62, 0)
}
