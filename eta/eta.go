package eta

import (
	"elevators/hardware"
	"elevators/orders"
	"elevators/prioritize"
	"elevators/timer"
	"time"
)

type InternalETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const travelDuration = 4 * time.Second
const offsetDuration = 2 * timer.PokeRate
const orderDuration = timer.DoorOpenTime + offsetDuration
const doorDuration = timer.DoorOpenTime

var internalETAs InternalETAs

func GetInternalETAs() InternalETAs {
	return internalETAs
}

func UpdateOrderAndInternalETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int,
	doorOpen bool,
	allOrders orders.AllOrders) (
	orders.AllOrders,
	InternalETAs) {

	prioritizedDirection := PrioritizedDirection(
		currentFloor,
		recentDirection,
		allOrders,
		internalETAs)

	newETAs := ComputeETAs(
		currentFloor,
		prioritizedDirection,
		recentDirection,
		doorOpen,
		allOrders)

	for floor := 0; floor < hardware.FloorCount; floor++ {
		if !newETAs.Up[floor].IsZero() &&
			(newETAs.Up[floor].Before(allOrders.Up[floor].BestETA) ||
				allOrders.Up[floor].BestETA.Before(time.Now())) {
			allOrders.Up[floor].BestETA = newETAs.Up[floor]
		} else if internalETAs.Up[floor].Equal(allOrders.Up[floor].BestETA) &&
			!allOrders.Up[floor].BestETA.IsZero() {
			newETAs.Up[floor] = allOrders.Up[floor].BestETA
		}

		if !newETAs.Down[floor].IsZero() &&
			(newETAs.Down[floor].Before(allOrders.Down[floor].BestETA) ||
				allOrders.Down[floor].BestETA.Before(time.Now())) {
			allOrders.Down[floor].BestETA = newETAs.Down[floor]
		} else if internalETAs.Down[floor].Equal(allOrders.Down[floor].BestETA) &&
			!allOrders.Down[floor].BestETA.IsZero() {
			newETAs.Down[floor] = allOrders.Down[floor].BestETA
		}
		allOrders.Up[floor].LocalETA = newETAs.Up[floor]
		allOrders.Down[floor].LocalETA = newETAs.Down[floor]
		allOrders.Up[floor].Now = time.Now()
		allOrders.Down[floor].Now = time.Now()
	}
	orders.SetOrders(allOrders)
	internalETAs = newETAs
	return allOrders, internalETAs
}

func ComputeETAs(
	currentFloor int,
	prioritizedDirection hardware.MotorDirection,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	allOrders orders.AllOrders) InternalETAs {

	if prioritizedDirection != hardware.MD_Stop {
		return SimulateETAs(
			prioritizedDirection,
			currentFloor,
			recentDirection,
			doorOpen,
			allOrders)
	} else {
		ETAsBelow := calculateETAforDirection(
			currentFloor,
			hardware.MD_Down,
			doorOpen,
			allOrders)
		ETAsAbove := calculateETAforDirection(
			currentFloor,
			hardware.MD_Up,
			doorOpen,
			allOrders)
		return bestETA(
			currentFloor,
			allOrders,
			ETAsBelow,
			ETAsAbove)
	}
}

func SimulateETAs(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	allOrders orders.AllOrders) InternalETAs {

	now := time.Now()
	simulationFloor := currentFloor
	simulationDirection := recentDirection
	simulationOrders := allOrders
	simulationTime := time.Now().Add(offsetDuration)
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
			now)
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
	now time.Time) hardware.MotorDirection {

	if etas.Cab[*floor].Equal(time.Time{}) {
		etas.Cab[*floor] = *simulationTime
	}

	doorAction := prioritize.DoorActionOnDoorTimeout(
		prioritizedDirection,
		false,
		orders.GetOrderSummary(
			*allOrders,
			*floor))

	switch doorAction {
	case hardware.DS_Close:
		newDirection := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
			orders.GetOrderSummary(
				*allOrders,
				*floor))

		if newDirection != prioritizedDirection {
			return hardware.MD_Stop
		}
		*floor += int(newDirection)
		*simulationTime = simulationTime.Add(travelDuration)

	case hardware.DS_Open_Down:
		etas.Down[*floor] = *simulationTime
		allOrders.Down[*floor].LastCompleteTime = now
		allOrders.Cab[*floor] = false

	case hardware.DS_Open_Up:
		etas.Up[*floor] = *simulationTime
		allOrders.Up[*floor].LastCompleteTime = now
		allOrders.Cab[*floor] = false

	case hardware.DS_Open_Cab:
		allOrders.Cab[*floor] = false
	default:
		panic("Invalid door action in eta simulation")
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
	allOrders orders.AllOrders) InternalETAs {

	var calculatedETAs InternalETAs

	simulationFloor := currentFloor
	simulationDirection := direction
	currentTime := time.Now().Add(offsetDuration)
	if doorOpen {
		currentTime = currentTime.Add(doorDuration)
	}
	for {
		calculatedETAs.setETA(
			simulationDirection,
			simulationFloor,
			currentTime)

		simulationFloor += int(simulationDirection)
		currentTime = currentTime.Add(travelDuration)

		if !hardware.ValidFloor(simulationFloor) {
			simulationDirection = -simulationDirection
			simulationFloor += int(simulationDirection)
			currentTime = currentTime.Add(-travelDuration)
		}

		if simulationFloor == currentFloor {
			return calculatedETAs
		}

		if orderToServe(
			simulationDirection,
			allOrders.Up[simulationFloor],
			allOrders.Down[simulationFloor],
			allOrders.Cab[simulationFloor]) {

			currentTime = currentTime.Add(orderDuration)

		}
	}

}

func bestETA(
	startFloor int,
	allOrders orders.AllOrders,
	ETAsBelow InternalETAs,
	ETAsAbove InternalETAs) InternalETAs {

	ETAsBelowFloor := startFloor
	ETAsAboveFloor := startFloor

	ETAsBelowDirection := hardware.MD_Down
	ETAsAboveDirection := hardware.MD_Up

	now := time.Now()
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
			now) &&

			!newETABetterOrBestETAExpiredWithOrder(
				ETAsAboveFloorOrderState,
				ETAsAboveFloorETA,
				now) {

			return ETAsBelow

		}

		if newETABetterOrBestETAExpiredWithOrder(
			ETAsAboveFloorOrderState,
			ETAsAboveFloorETA,
			now) &&

			!newETABetterOrBestETAExpiredWithOrder(
				ETAsBelowFloorOrderState,
				ETAsBelowFloorETA,
				now) {

			return ETAsAbove

		}
	}

	if hardware.FloorBelowMiddleFloor(startFloor) {
		return ETAsAbove
	} else {
		return ETAsBelow
	}
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
		if (allOrders.Up[floor].HasOrder() &&
			InternalETABest(
				allOrders.Up[floor],
				allETAs.Up[floor])) ||
			(allOrders.Down[floor].HasOrder() &&
				InternalETABest(
					allOrders.Down[floor],
					allETAs.Down[floor])) ||
			allOrders.Cab[floor] {
			return true
		}
	}
	return false
}

func internalETABest(
	orderState orders.OrderState,
	internalETA time.Time) bool {

	return orderState.BestETA.Equal(internalETA) && internalETA.After(time.Now())
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

func AllInternalETAsBest(allOrders orders.AllOrders) bool {
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if !internalETABest(
			allOrders.Down[floor],
			internalETAs.Down[floor]) ||
			!internalETABest(
				allOrders.Up[floor],
				internalETAs.Up[floor]) {

			return false
		}
	}
	return true
}

func FirstInternalETAExpiration(currentETAs InternalETAs) time.Time {
	now := time.Now()
	etaExpiration := maxTime()
	for _, floor := range hardware.ValidFloors() {
		if currentETAs.ETABetweenTimes(
			hardware.MD_Up,
			floor,
			now,
			etaExpiration) {

			etaExpiration = currentETAs.getETA(
				hardware.MD_Up,
				floor)

		}
		if internalETAs.ETABetweenTimes(
			hardware.MD_Down,
			floor,
			now,
			etaExpiration) {

			etaExpiration = internalETAs.getETA(
				hardware.MD_Down,
				floor)

		}
	}
	if etaExpiration.Equal(maxTime()) {
		etaExpiration = now
	}
	return etaExpiration
}
