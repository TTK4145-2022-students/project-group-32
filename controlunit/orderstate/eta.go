package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"time"
)

type InternalETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const travelDuration = 5 * time.Second
const orderDuration = 4 * time.Second
const OffsetDuration = 1 * time.Second
const directionChangeDuration = 200 * time.Millisecond

var internalETAs InternalETAs

func GetInternalETAs() InternalETAs {
	return internalETAs
}

func UpdateOrderAndInternalETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int) (AllOrders, InternalETAs) {
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
		allOrders)

	for floor := 0; floor < hardware.FloorCount; floor++ {
		// if allOrders.Up[floor].BestETA.Before(time.Now()) {
		// 	allOrders.Up[floor].BestETA = time.Time{}
		// }
		if !newETAs.Up[floor].IsZero() &&
			(newETAs.Up[floor].Before(allOrders.Up[floor].BestETA) ||
				// allOrders.Up[floor].BestETA.IsZero()) {
				allOrders.Up[floor].BestETA.Before(time.Now())) {
			allOrders.Up[floor].BestETA = newETAs.Up[floor]
		} else if internalETAs.Up[floor].Equal(allOrders.Up[floor].BestETA) &&
			// !allOrders.Up[floor].BestETA.IsZero() {
			!allOrders.Up[floor].BestETA.Before(time.Now()) {
			newETAs.Up[floor] = allOrders.Up[floor].BestETA
		}

		// if allOrders.Down[floor].BestETA.Before(time.Now()) {
		// 	allOrders.Down[floor].BestETA = time.Time{}
		// }
		if !newETAs.Down[floor].IsZero() &&
			(newETAs.Down[floor].Before(allOrders.Down[floor].BestETA) ||
				// allOrders.Down[floor].BestETA.IsZero()) {
				allOrders.Down[floor].BestETA.Before(time.Now())) {
			allOrders.Down[floor].BestETA = newETAs.Down[floor]
		} else if internalETAs.Down[floor].Equal(allOrders.Down[floor].BestETA) &&
			// !allOrders.Down[floor].BestETA.IsZero() {
			!allOrders.Down[floor].BestETA.Before(time.Now()) {
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
	orders AllOrders) InternalETAs {

	if prioritizedDirection != hardware.MD_Stop {
		return SimulateETAs(
			prioritizedDirection,
			currentFloor,
			recentDirection,
			orders)
	} else {
		ETAsBelow := calculateETAforDirection(
			currentFloor,
			hardware.MD_Down,
			orders)
		ETAsAbove := calculateETAforDirection(
			currentFloor,
			hardware.MD_Up,
			orders)
		return bestETA(currentFloor, orders, ETAsBelow, ETAsAbove)
	}
}

func SimulateETAs(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders) InternalETAs {

	simulationFloor := currentFloor
	simulationDirection := recentDirection
	simulationOrders := orders
	simulationTime := time.Now().Add(OffsetDuration)
	var simulatedDurations InternalETAs
	for prioritizedDirection != hardware.MD_Stop {
		prioritizedDirection = simulateETAStep(
			prioritizedDirection,
			&simulationFloor,
			&simulationDirection,
			&simulationOrders,
			&simulationTime,
			&simulatedDurations)
	}
	return simulatedDurations
}

func simulateETAStep(
	prioritizedDirection hardware.MotorDirection,
	floor *int,
	direction *hardware.MotorDirection,
	orders *AllOrders,
	simulationTime *time.Time,
	etas *InternalETAs) hardware.MotorDirection {

	if etas.Cab[*floor].Equal(time.Time{}) {
		etas.Cab[*floor] = *simulationTime
	}

	doorAction := prioritize.DoorActionOnDoorTimeout(
		prioritizedDirection,
		false,
		GetOrderStatus(*orders, *floor))

	switch doorAction {
	case hardware.DS_Close:
		newDirection := prioritize.MotorActionOnDecisionDeadline(
			prioritizedDirection,
			GetOrderStatus(*orders, *floor))
		if newDirection != prioritizedDirection {
			return hardware.MD_Stop
		}
		*floor += int(newDirection)
		*simulationTime = simulationTime.Add(travelDuration)
	case hardware.DS_Open_Down:
		etas.Down[*floor] = *simulationTime
		orders.Down[*floor].LastCompleteTime = time.Now()
		orders.Cab[*floor] = false
	case hardware.DS_Open_Up:
		etas.Up[*floor] = *simulationTime
		orders.Up[*floor].LastCompleteTime = time.Now()
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
	orders AllOrders) InternalETAs {

	var calculatedETAs InternalETAs

	simulationFloor := currentFloor
	simulationDirection := direction
	currentTime := time.Now().Add(OffsetDuration)
	for {
		if simulationDirection == hardware.MD_Down {
			calculatedETAs.Down[simulationFloor] = currentTime
		}
		if simulationDirection == hardware.MD_Up {
			calculatedETAs.Up[simulationFloor] = currentTime
		}

		simulationFloor += int(simulationDirection)
		currentTime = currentTime.Add(travelDuration)
		if simulationFloor < 0 || simulationFloor >= hardware.FloorCount {
			simulationDirection = -simulationDirection
			simulationFloor += int(simulationDirection)
			currentTime = currentTime.Add(-travelDuration).Add(directionChangeDuration)
		}

		if simulationFloor == currentFloor {
			return calculatedETAs
		}

		if simulationDirection == hardware.MD_Down &&
			orders.Down[simulationFloor].hasOrder() ||
			(simulationDirection == hardware.MD_Up &&
				orders.Up[simulationFloor].hasOrder()) {
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
		if ETAsBelowFloor < 0 {
			ETAsBelowFloor = 0
			ETAsBelowDirection = hardware.MD_Up
		}
		if ETAsAboveFloor >= hardware.FloorCount {
			ETAsAboveFloor = hardware.FloorCount - 1
			ETAsAboveDirection = hardware.MD_Down
		}

		if ETAsAboveFloor == startFloor || ETAsBelowFloor == startFloor {
			break
		}

		ETAsBelowFloorETA := ETAsBelow.Down[ETAsBelowFloor]
		ETAsBelowFloorBestETA := orders.Down[ETAsBelowFloor].BestETA
		ETAsBelowFloorOrder := orders.Down[ETAsBelowFloor].hasOrder()
		if ETAsBelowDirection == hardware.MD_Up {
			ETAsBelowFloorETA = ETAsBelow.Up[ETAsBelowFloor]
			ETAsBelowFloorBestETA = orders.Up[ETAsBelowFloor].BestETA
			ETAsBelowFloorOrder = orders.Up[ETAsBelowFloor].hasOrder()
		}
		ETAsAboveFloorETA := ETAsAbove.Up[ETAsAboveFloor]
		ETAsAboveFloorBestETA := orders.Up[ETAsAboveFloor].BestETA
		ETAsAboveFloorOrder := orders.Up[ETAsAboveFloor].hasOrder()
		if ETAsAboveDirection == hardware.MD_Down {
			ETAsAboveFloorETA = ETAsAbove.Down[ETAsAboveFloor]
			ETAsAboveFloorBestETA = orders.Down[ETAsAboveFloor].BestETA
			ETAsAboveFloorOrder = orders.Down[ETAsAboveFloor].hasOrder()
		}

		if (ETAsBelowFloorETA.Before(ETAsBelowFloorBestETA) ||
			(ETAsBelowFloorOrder && now.After(ETAsBelowFloorBestETA))) &&

			!(ETAsAboveFloorETA.Before(ETAsAboveFloorBestETA) ||
				(ETAsAboveFloorOrder && now.After(ETAsAboveFloorBestETA))) {
			return ETAsBelow
		}
		if (ETAsAboveFloorETA.Before(ETAsAboveFloorBestETA) ||
			(ETAsAboveFloorOrder && now.After(ETAsAboveFloorBestETA))) &&

			!(ETAsBelowFloorETA.Before(ETAsBelowFloorBestETA) ||
				(ETAsBelowFloorOrder && now.After(ETAsBelowFloorBestETA))) {
			return ETAsAbove
		}
	}

	if 2*startFloor < hardware.FloorCount {
		return ETAsAbove
	} else {
		return ETAsBelow
	}
}

func InternalETABest(orderState OrderState, internalETA time.Time) bool {
	return orderState.BestETA.Equal(internalETA) && !internalETA.IsZero() // && internalETA.After(time.Now())
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
	for floor := currentFloor + int(direction); 0 <= floor && floor < hardware.FloorCount; floor += int(direction) {
		if (orders.Up[floor].hasOrder() &&
			InternalETABest(orders.Up[floor], allETAs.Up[floor])) ||
			(orders.Down[floor].hasOrder() &&
				InternalETABest(orders.Down[floor], allETAs.Down[floor])) ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func PrioritizedDirection(currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs InternalETAs) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if orderAndInternalETABest(hardware.MD_Up, currentFloor, orders, allETAs) {
			return hardware.MD_Up
		}
		if orderAndInternalETABest(hardware.MD_Down, currentFloor, orders, allETAs) {
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if orderAndInternalETABest(hardware.MD_Down, currentFloor, orders, allETAs) {
			return hardware.MD_Down
		}
		if orderAndInternalETABest(hardware.MD_Up, currentFloor, orders, allETAs) {
			return hardware.MD_Up
		}
	}
	return hardware.MD_Stop
}

func AllInternalETAsBest(orders AllOrders) bool {
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if !InternalETABest(orders.Down[floor], internalETAs.Down[floor]) ||
			!InternalETABest(orders.Up[floor], internalETAs.Up[floor]) {
			return false
		}
	}
	return true
}

func FirstBestETAexpirationWithOrder(orders AllOrders) time.Time {
	now := time.Now()
	var etaExpiration time.Time
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if orders.Down[floor].hasOrder() &&
			now.Before(orders.Down[floor].BestETA) &&
			(orders.Down[floor].BestETA.Before(etaExpiration) ||
				etaExpiration.IsZero()) {
			etaExpiration = orders.Down[floor].BestETA
		}

		if orders.Up[floor].hasOrder() &&
			now.Before(orders.Up[floor].BestETA) &&
			(orders.Up[floor].BestETA.Before(etaExpiration) ||
				etaExpiration.IsZero()) {
			etaExpiration = orders.Up[floor].BestETA
		}
	}
	if etaExpiration.IsZero() {
		etaExpiration = now
	}
	return etaExpiration
}
