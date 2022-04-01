package eta

import (
	"elevators/hardware"
	"elevators/orders"
	"elevators/prioritize"
	"elevators/timer"
	"time"
)

type InternalETAs struct {
	up   [hardware.FloorCount]time.Time
	down [hardware.FloorCount]time.Time
	cab  [hardware.FloorCount]time.Time
}

const (
	travelDuration = 4 * time.Second
	offsetDuration = 2 * timer.PokeRate
	orderDuration  = timer.DoorOpenTime + offsetDuration
	doorDuration   = timer.DoorOpenTime
)

var internalETAs InternalETAs

func GetInternalETAs() InternalETAs {
	return internalETAs
}

func UpdateOrderAndInternalETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int,
	doorOpen bool,
	allOrders orders.AllOrders) (orders.AllOrders, InternalETAs) {

	currentTime := time.Now()

	prioritizedDirection := PrioritizedDirection(
		currentFloor,
		recentDirection,
		allOrders,
		internalETAs)

	newETAs := computeETAs(
		currentFloor,
		prioritizedDirection,
		recentDirection,
		doorOpen,
		allOrders,
		currentTime)

	for _, floor := range hardware.ValidFloors() {
		if newETABetterOrBestETAExpired(
			allOrders.Up[floor],
			newETAs.up[floor],
			currentTime) {

			allOrders.Up[floor].BestETA = newETAs.up[floor]
		} else if internalETABest(
			allOrders.Up[floor],
			internalETAs.up[floor]) {

			newETAs.up[floor] = allOrders.Up[floor].BestETA
		}
		if newETABetterOrBestETAExpired(
			allOrders.Down[floor],
			newETAs.down[floor],
			currentTime) {

			allOrders.Down[floor].BestETA = newETAs.down[floor]

		} else if internalETABest(
			allOrders.Down[floor],
			internalETAs.down[floor]) {

			newETAs.down[floor] = allOrders.Down[floor].BestETA
		}
	}

	orders.SetOrders(allOrders)
	internalETAs = newETAs

	return allOrders, internalETAs
}

func computeETAs(
	currentFloor int,
	prioritizedDirection hardware.MotorDirection,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	allOrders orders.AllOrders,
	currentTime time.Time) InternalETAs {

	if prioritizedDirection != hardware.MD_Stop {
		return simulateETAs(
			prioritizedDirection,
			currentFloor,
			recentDirection,
			doorOpen,
			allOrders,
			currentTime)

	} else {
		ETAsBelow := calculateETAforDirection(
			currentFloor,
			hardware.MD_Down,
			doorOpen,
			allOrders,
			currentTime)

		ETAsAbove := calculateETAforDirection(
			currentFloor,
			hardware.MD_Up,
			doorOpen,
			allOrders,
			currentTime)

		return bestETA(
			currentFloor,
			allOrders,
			ETAsBelow,
			ETAsAbove,
			currentTime)
	}
}

func simulateETAs(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	allOrders orders.AllOrders,
	currentTime time.Time) InternalETAs {

	simulationFloor := currentFloor
	simulationDirection := recentDirection
	simulationOrders := allOrders
	simulationTime := currentTime.Add(offsetDuration)

	if doorOpen {
		simulationTime = simulationTime.Add(doorDuration)
	}

	var simulatedDurations InternalETAs
	for prioritizedDirection != hardware.MD_Stop {
		prioritizedDirection = simulateETAStep(
			prioritizedDirection,
			&simulationFloor,
			&simulationDirection,
			&simulationOrders,
			&simulationTime,
			&simulatedDurations,
			currentTime)
	}

	return simulatedDurations
}

func simulateETAStep(
	prioritizedDirection hardware.MotorDirection,
	floor *int,
	direction *hardware.MotorDirection,
	allOrders *orders.AllOrders,
	simulationTime *time.Time,
	etas *InternalETAs,
	currentTime time.Time) hardware.MotorDirection {

	if noETA(etas.cab[*floor]) {
		etas.cab[*floor] = *simulationTime
	}

	orderSummary := orders.GetOrderSummary(
		*allOrders,
		*floor)

	doorAction := prioritize.DoorActionOnDoorTimeout(
		prioritizedDirection,
		false,
		orderSummary)

	switch doorAction {
	case hardware.DS_Close:
		newDirection := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
			orderSummary)

		if newDirection != prioritizedDirection {
			return hardware.MD_Stop
		}

		*floor += int(newDirection)
		*simulationTime = simulationTime.Add(travelDuration)

	case hardware.DS_Open_Down:
		etas.down[*floor] = *simulationTime
		allOrders.Down[*floor].LastCompleteTime = currentTime
		allOrders.Cab[*floor] = false

	case hardware.DS_Open_Up:
		etas.up[*floor] = *simulationTime
		allOrders.Up[*floor].LastCompleteTime = currentTime
		allOrders.Cab[*floor] = false

	case hardware.DS_Open_Cab:
		allOrders.Cab[*floor] = false
	}

	if doorAction != hardware.DS_Close {
		*simulationTime = simulationTime.Add(orderDuration)
	}

	return prioritizedDirection
}

func calculateETAforDirection(
	currentFloor int,
	direction hardware.MotorDirection,
	doorOpen bool,
	allOrders orders.AllOrders,
	currentTime time.Time) InternalETAs {

	var calculatedETAs InternalETAs

	simulationFloor := currentFloor
	simulationDirection := direction
	simulationTime := currentTime.Add(offsetDuration)
	if doorOpen {
		simulationTime = simulationTime.Add(doorDuration)
	}

	for {
		calculatedETAs.setETA(
			simulationDirection,
			simulationFloor,
			simulationTime)

		simulationFloor += int(simulationDirection)
		simulationTime = simulationTime.Add(travelDuration)

		if !hardware.ValidFloor(simulationFloor) {
			simulationDirection = -simulationDirection
			simulationFloor += int(simulationDirection)
			simulationTime = simulationTime.Add(-travelDuration)
		}
		if simulationFloor == currentFloor {
			return calculatedETAs
		}
		if orderToServe(
			simulationDirection,
			allOrders.Up[simulationFloor],
			allOrders.Down[simulationFloor],
			allOrders.Cab[simulationFloor]) {

			simulationTime = simulationTime.Add(orderDuration)
		}
	}
}

func bestETA(
	startFloor int,
	allOrders orders.AllOrders,
	ETAsBelow InternalETAs,
	ETAsAbove InternalETAs,
	currentTime time.Time) InternalETAs {

	ETAsBelowFloor := startFloor
	ETAsAboveFloor := startFloor

	ETAsBelowDirection := hardware.MD_Down
	ETAsAboveDirection := hardware.MD_Up

	for {
		ETAsBelowFloor += int(ETAsBelowDirection)
		ETAsAboveFloor += int(ETAsAboveDirection)

		if ETAsBelowFloor <= hardware.BottomFloor {
			ETAsBelowFloor = hardware.BottomFloor
			ETAsBelowDirection = hardware.MD_Up
		}
		if ETAsAboveFloor >= hardware.TopFloor {
			ETAsAboveFloor = hardware.TopFloor
			ETAsAboveDirection = hardware.MD_Down
		}
		if ETAsAboveFloor == startFloor ||
			ETAsBelowFloor == startFloor {
			break
		}

		ETAsBelowFloorETA := ETAsBelow.getETA(
			ETAsBelowDirection,
			ETAsBelowFloor)
		ETAsBelowFloorOrderState := allOrders.GetOrderState(
			ETAsBelowDirection,
			ETAsBelowFloor)

		ETAsAboveFloorETA := ETAsAbove.getETA(
			ETAsAboveDirection,
			ETAsAboveFloor)
		ETAsAboveFloorOrderState := allOrders.GetOrderState(
			ETAsAboveDirection,
			ETAsAboveFloor)

		if newETABetterOrBestETAExpiredWithOrder(
			ETAsBelowFloorOrderState,
			ETAsBelowFloorETA,
			currentTime) &&

			!newETABetterOrBestETAExpiredWithOrder(
				ETAsAboveFloorOrderState,
				ETAsAboveFloorETA,
				currentTime) {

			return ETAsBelow
		}
		if newETABetterOrBestETAExpiredWithOrder(
			ETAsAboveFloorOrderState,
			ETAsAboveFloorETA,
			currentTime) &&

			!newETABetterOrBestETAExpiredWithOrder(
				ETAsBelowFloorOrderState,
				ETAsBelowFloorETA,
				currentTime) {

			return ETAsAbove
		}
	}

	if hardware.FloorBelowMiddleFloor(startFloor) {
		return ETAsAbove
	} else {
		return ETAsBelow
	}
}

func PrioritizedDirection(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	allOrders orders.AllOrders,
	allETAs InternalETAs) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if orderAndInternalETABest(
			hardware.MD_Up,
			currentFloor,
			allOrders,
			allETAs) {

			return hardware.MD_Up
		}
		if orderAndInternalETABest(
			hardware.MD_Down,
			currentFloor,
			allOrders,
			allETAs) {

			return hardware.MD_Down
		}

	case hardware.MD_Down:
		if orderAndInternalETABest(
			hardware.MD_Down,
			currentFloor,
			allOrders,
			allETAs) {

			return hardware.MD_Down
		}
		if orderAndInternalETABest(
			hardware.MD_Up,
			currentFloor,
			allOrders,
			allETAs) {

			return hardware.MD_Up
		}
	}
	return hardware.MD_Stop
}
