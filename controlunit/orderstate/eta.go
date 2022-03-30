package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
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
	doorOpen bool) (
	AllOrders,
	InternalETAs) {

	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()

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
	internalETAs = newETAs
	return allOrders, internalETAs
}

func ComputeETAs(
	currentFloor int,
	prioritizedDirection hardware.MotorDirection,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	orders AllOrders) InternalETAs {

	if prioritizedDirection != hardware.MD_Stop {
		return SimulateETAs(
			prioritizedDirection,
			currentFloor,
			recentDirection,
			doorOpen,
			orders)
	} else {
		ETAsBelow := calculateETAforDirection(
			currentFloor,
			hardware.MD_Down,
			doorOpen,
			orders)
		ETAsAbove := calculateETAforDirection(
			currentFloor,
			hardware.MD_Up,
			doorOpen,
			orders)
		return bestETA(
			currentFloor,
			orders,
			ETAsBelow,
			ETAsAbove)
	}
}

func SimulateETAs(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	recentDirection hardware.MotorDirection,
	doorOpen bool,
	orders AllOrders) InternalETAs {

	now := time.Now()
	simulationFloor := currentFloor
	simulationDirection := recentDirection
	simulationOrders := orders
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
	orders *AllOrders,
	simulationTime *time.Time,
	etas *InternalETAs,
	now time.Time) hardware.MotorDirection {

	if etas.Cab[*floor].Equal(time.Time{}) {
		etas.Cab[*floor] = *simulationTime
	}

	doorAction := prioritize.DoorActionOnDoorTimeout(
		prioritizedDirection,
		false,
		GetOrderSummary(
			*orders,
			*floor))

	switch doorAction {
	case hardware.DS_Close:
		newDirection := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
			GetOrderSummary(
				*orders,
				*floor))

		if newDirection != prioritizedDirection {
			return hardware.MD_Stop
		}
		*floor += int(newDirection)
		*simulationTime = simulationTime.Add(travelDuration)

	case hardware.DS_Open_Down:
		etas.Down[*floor] = *simulationTime
		orders.Down[*floor].LastCompleteTime = now
		orders.Cab[*floor] = false

	case hardware.DS_Open_Up:
		etas.Up[*floor] = *simulationTime
		orders.Up[*floor].LastCompleteTime = now
		orders.Cab[*floor] = false

	case hardware.DS_Open_Cab:
		orders.Cab[*floor] = false
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
	orders AllOrders) InternalETAs {

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
			orders.Up[simulationFloor],
			orders.Down[simulationFloor],
			orders.Cab[simulationFloor]) {

			currentTime = currentTime.Add(orderDuration)

		}
	}

}

func bestETA(
	startFloor int,
	orders AllOrders,
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
		ETAsBelowFloorOrderState := orders.getOrderState(
			ETAsBelowDirection,
			ETAsBelowFloor)

		ETAsAboveFloorETA := ETAsAbove.getETA(
			ETAsAboveDirection,
			ETAsAboveFloor)
		ETAsAboveFloorOrderState := orders.getOrderState(
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
	orders AllOrders,
	allETAs InternalETAs) bool {

	switch direction {
	case hardware.MD_Up:
		if orders.Up[currentFloor].hasOrder() {
			return true
		}

	case hardware.MD_Down:
		if orders.Down[currentFloor].hasOrder() {
			return true
		}
	}
	for floor := currentFloor + int(direction); hardware.ValidFloor(floor); floor += int(direction) {
		if (orders.Up[floor].hasOrder() &&
			InternalETABest(
				orders.Up[floor],
				allETAs.Up[floor])) ||
			(orders.Down[floor].hasOrder() &&
				InternalETABest(
					orders.Down[floor],
					allETAs.Down[floor])) ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func internalETABest(
	orderState OrderState,
	internalETA time.Time) bool {

	return orderState.BestETA.Equal(internalETA) && internalETA.After(time.Now())
}

func PrioritizedDirection(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs InternalETAs) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:

		if orderAndInternalETABest(
			hardware.MD_Up,
			currentFloor,
			orders,
			allETAs) {

			return hardware.MD_Up

		}
		if orderAndInternalETABest(
			hardware.MD_Down,
			currentFloor,
			orders,
			allETAs) {

			return hardware.MD_Down

		}

	case hardware.MD_Down:

		if orderAndInternalETABest(
			hardware.MD_Down,
			currentFloor,
			orders,
			allETAs) {

			return hardware.MD_Down

		}
		if orderAndInternalETABest(
			hardware.MD_Up,
			currentFloor,
			orders,
			allETAs) {

			return hardware.MD_Up

		}
	}
	return hardware.MD_Stop
}

func AllInternalETAsBest(orders AllOrders) bool {
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if !internalETABest(
			orders.Down[floor], internalETAs.Down[floor]) ||
			!internalETABest(
				orders.Up[floor],
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
