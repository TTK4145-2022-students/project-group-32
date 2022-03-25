package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"fmt"
	"time"
)

type AllDurations struct {
	Up   [hardware.FloorCount]time.Duration
	Down [hardware.FloorCount]time.Duration
	Cab  [hardware.FloorCount]time.Duration
}

type InternalETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const travelDuration = 3 * time.Second
const orderDuration = 4 * time.Second
const offsetDuration = 1 * time.Second

var internalETAs InternalETAs

func GetInternalETAs() InternalETAs {
	return internalETAs
}

func UpdateETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int) {
	allOrdersMtx.Lock()
	defer allOrdersMtx.Unlock()

	// newDurations := ComputeDurations(currentFloor, recentDirection, allOrders, internalETAs)
	// newETAs := ComputeInternalETAs(newDurations)
	newETAs := ComputeETAs(currentFloor, recentDirection, allOrders, internalETAs)

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
}

func ComputeInternalETAs(durations AllDurations) InternalETAs {
	var newETAs InternalETAs
	var now = time.Now()
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if durations.Cab[floor] != time.Duration(0) {
			newETAs.Cab[floor] = now.Add(durations.Cab[floor])
		}
		if durations.Up[floor] != time.Duration(0) {
			newETAs.Up[floor] = now.Add(durations.Up[floor])
		}
		if durations.Down[floor] != time.Duration(0) {
			newETAs.Down[floor] = now.Add(durations.Down[floor])
		}
	}
	return newETAs
}

func ComputeDurations(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs InternalETAs) AllDurations {

	prioritizedDirection := ETADirection(
		currentFloor,
		recentDirection,
		orders,
		allETAs)
	if prioritizedDirection != hardware.MD_Stop {
		return SimulateDurations(
			prioritizedDirection,
			currentFloor,
			recentDirection,
			orders)
	} else {
		durationsBelow := calculateDurationforDirection(
			currentFloor,
			hardware.MD_Down,
			orders)
		durationsAbove := calculateDurationforDirection(
			currentFloor,
			hardware.MD_Up,
			orders)
		ETAs := []InternalETAs{ComputeInternalETAs(durationsBelow), ComputeInternalETAs(durationsAbove)}
		ETAindex := bestDurations(
			currentFloor,
			orders,
			ETAs)
		switch ETAindex {
		case 0:
			return durationsBelow
		case 1:
			return durationsAbove
		default:
			panic("ugly code failde")
		}
	}
}

func SimulateDurations(
	prioritizedDirection hardware.MotorDirection,
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders) AllDurations {

	simulationFloor := currentFloor
	simulationDirection := recentDirection
	simulationOrders := orders
	simulationTime := offsetDuration
	var simulatedDurations AllDurations
	for prioritizedDirection != hardware.MD_Stop {
		prioritizedDirection = simulateDurationStep(
			prioritizedDirection,
			&simulationFloor,
			&simulationDirection,
			&simulationOrders,
			&simulationTime,
			&simulatedDurations)
	}
	return simulatedDurations
}

func simulateDurationStep(
	prioritizedDirection hardware.MotorDirection,
	floor *int,
	direction *hardware.MotorDirection,
	orders *AllOrders,
	simTime *time.Duration,
	durations *AllDurations) hardware.MotorDirection {

	if durations.Cab[*floor] == time.Duration(0) {
		durations.Cab[*floor] = *simTime
	}

	doorAction := prioritize.DoorActionOnDoorTimeout(
		prioritizedDirection,
		false,
		GetOrderStatus(*orders, *floor))

	switch doorAction {
	case hardware.DS_Close:
		newDirection := prioritize.MotorActionOnDoorClose(
			prioritizedDirection,
			GetOrderStatus(*orders, *floor))
		if newDirection != prioritizedDirection {
			return hardware.MD_Stop
		}
		*floor += int(newDirection)
		*simTime += travelDuration

	case hardware.DS_Open_Down:
		durations.Down[*floor] = *simTime
		orders.Down[*floor].LastCompleteTime = time.Now()
		orders.Cab[*floor] = false
	case hardware.DS_Open_Up:
		durations.Up[*floor] = *simTime
		orders.Up[*floor].LastCompleteTime = time.Now()
		orders.Cab[*floor] = false
	case hardware.DS_Open_Cab:
		orders.Cab[*floor] = false
	default:
		panic("Invalid door action in eta simulation")
	}

	if doorAction != hardware.DS_Close {
		*simTime += orderDuration
	}
	return prioritizedDirection
}

func ComputeETAs(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	internalETAs InternalETAs) InternalETAs {

	prioritizedDirection := ETADirection(
		currentFloor,
		recentDirection,
		orders,
		internalETAs)
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
	simulationTime := time.Now().Add(offsetDuration)
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
		newDirection := prioritize.MotorActionOnDoorClose(
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

func calculateDurationforDirection(
	currentFloor int,
	direction hardware.MotorDirection,
	orders AllOrders) AllDurations {

	var computedDurations AllDurations

	simCabFloor := currentFloor
	simCabDirection := int(direction)
	currentTime := offsetDuration
	for {
		if simCabDirection == int(hardware.MD_Down) {
			computedDurations.Down[simCabFloor] = currentTime
		}
		if simCabDirection == int(hardware.MD_Up) {
			computedDurations.Up[simCabFloor] = currentTime
		}

		simCabFloor += simCabDirection
		currentTime += travelDuration
		if simCabFloor < 0 || simCabFloor >= hardware.FloorCount {
			simCabDirection = -simCabDirection
			simCabFloor += simCabDirection
			currentTime -= travelDuration
		}

		if simCabFloor == currentFloor {
			return computedDurations
		}

		if simCabDirection == int(hardware.MD_Down) &&
			hasOrder(orders.Down[simCabFloor]) ||
			(simCabDirection == int(hardware.MD_Up) &&
				hasOrder(orders.Up[simCabFloor])) {
			currentTime += orderDuration
		}
	}

}

func calculateETAforDirection(
	currentFloor int,
	direction hardware.MotorDirection,
	orders AllOrders) InternalETAs {

	var calculatedETAs InternalETAs

	simulationFloor := currentFloor
	simulationDirection := direction
	currentTime := time.Now().Add(offsetDuration)
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
			currentTime = currentTime.Add(-travelDuration)
		}

		if simulationFloor == currentFloor {
			return calculatedETAs
		}

		if simulationDirection == hardware.MD_Down &&
			hasOrder(orders.Down[simulationFloor]) ||
			(simulationDirection == hardware.MD_Up &&
				hasOrder(orders.Up[simulationFloor])) {
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
		ETAsBelowFloorOrder := hasOrder(orders.Down[ETAsBelowFloor])
		if ETAsBelowDirection == hardware.MD_Up {
			ETAsBelowFloorETA = ETAsBelow.Up[ETAsBelowFloor]
			ETAsBelowFloorBestETA = orders.Up[ETAsBelowFloor].BestETA
			ETAsBelowFloorOrder = hasOrder(orders.Up[ETAsBelowFloor])
		}
		ETAsAboveFloorETA := ETAsAbove.Up[ETAsAboveFloor]
		ETAsAboveFloorBestETA := orders.Up[ETAsAboveFloor].BestETA
		ETAsAboveFloorOrder := hasOrder(orders.Up[ETAsAboveFloor])
		if ETAsAboveDirection == hardware.MD_Down {
			ETAsAboveFloorETA = ETAsAbove.Down[ETAsAboveFloor]
			ETAsAboveFloorBestETA = orders.Down[ETAsAboveFloor].BestETA
			ETAsAboveFloorOrder = hasOrder(orders.Down[ETAsAboveFloor])
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

func bestDurations(
	floor int,
	orders AllOrders,
	ETAs []InternalETAs) int {

	ETAsBelow := ETAs[0]
	ETAsAbove := ETAs[1]

	belowFloor := floor
	aboveFloor := floor

	belowDir := -1
	aboveDir := 1

	now := time.Now()
	for {
		belowFloor += belowDir
		aboveFloor += aboveDir
		if belowFloor <= 0 {
			belowFloor = 0
			belowDir = 1
		}
		if aboveFloor >= hardware.FloorCount-1 {
			aboveFloor = hardware.FloorCount - 1
			aboveDir = -1
		}

		if aboveFloor == floor || belowFloor == floor {
			break
		}

		floorETABelow := ETAsBelow.Down[belowFloor]
		floorOrderETABelow := hasOrder(orders.Down[belowFloor])
		orderETABelow := orders.Down[belowFloor].BestETA
		if belowDir == 1 {
			floorETABelow = ETAsBelow.Up[belowFloor]
			floorOrderETABelow = hasOrder(orders.Up[belowFloor])
			orderETABelow = orders.Up[belowFloor].BestETA
		}
		floorETAAbove := ETAsAbove.Up[aboveFloor]
		floorOrderETAAbove := hasOrder(orders.Up[aboveFloor])
		orderETAAbove := orders.Up[aboveFloor].BestETA
		if aboveDir == -1 {
			floorETAAbove = ETAsAbove.Down[aboveFloor]
			floorOrderETAAbove = hasOrder(orders.Down[aboveFloor])
			orderETAAbove = orders.Down[aboveFloor].BestETA
		}
		if (floorETABelow.Before(orderETABelow) ||
			(floorOrderETABelow &&
				now.After(orderETABelow))) &&

			!(floorETAAbove.Before(orderETAAbove) ||
				(floorOrderETAAbove &&
					now.After(orderETAAbove))) {
			return 0
		}
		if (floorETAAbove.Before(orderETAAbove) ||
			(floorOrderETAAbove &&
				now.After(orderETAAbove))) &&
			!(floorETABelow.Before(orderETABelow) ||
				(floorOrderETABelow &&
					now.After(orderETABelow))) {
			return 1
		}
	}

	if 2*floor < hardware.FloorCount {
		return 1
	} else {
		return 0
	}
}

// func min(a, b int) int {
// 	if a <= b {
// 		return a
// 	}
// 	return b
// }

func internalETABest(orderState OrderState, internalETA time.Time) bool {
	return orderState.BestETA.Equal(internalETA) && internalETA.After(time.Now()) // && !internalETA.IsZero()
}

func orderAndInternalETABest(
	direction hardware.MotorDirection,
	currentFloor int,
	orders AllOrders,
	allETAs InternalETAs) bool {
	switch direction {
	case hardware.MD_Up:
		if hasOrder(orders.Up[currentFloor]) {
			return true
		}
	case hardware.MD_Down:
		if hasOrder(orders.Down[currentFloor]) {
			return true
		}
	}
	for floor := currentFloor + int(direction); 0 <= floor && floor < hardware.FloorCount; floor += int(direction) {
		if (hasOrder(orders.Up[floor]) &&
			internalETABest(orders.Up[floor], allETAs.Up[floor])) ||
			(hasOrder(orders.Down[floor]) &&
				internalETABest(orders.Down[floor], allETAs.Down[floor])) ||
			orders.Cab[floor] {
			return true
		}
	}
	return false
}

func ETADirection(
	floor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs InternalETAs) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if orderAndInternalETABest(hardware.MD_Up, floor, orders, allETAs) {
			fmt.Println("going up, best up")
			return hardware.MD_Up
		}
		if orderAndInternalETABest(hardware.MD_Down, floor, orders, allETAs) {
			fmt.Println("going up, best down")
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if orderAndInternalETABest(hardware.MD_Down, floor, orders, allETAs) {
			fmt.Println("going deon, best down")
			return hardware.MD_Down
		}
		if orderAndInternalETABest(hardware.MD_Up, floor, orders, allETAs) {
			fmt.Println("going deon, best up")
			return hardware.MD_Up
		}
	}
	// if !AnyOrders(orders) && !AllInternalETAsBest(orders) {
	// 	fmt.Println("prioritizing to prepare")
	// 	if 0 < floor && floor < hardware.FloorCount-1 &&
	// 		internalETABest(orders.Up[floor-1], allETAs.Up[floor-1]) {
	// 		return hardware.MD_Down
	// 	} else if 0 < floor && floor < hardware.FloorCount-1 &&
	// 		internalETABest(orders.Down[floor+1], allETAs.Down[floor+1]) {
	// 		return hardware.MD_Up
	// 	}
	// }
	return hardware.MD_Stop
}

func PrioritizedDirection(currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs InternalETAs) hardware.MotorDirection {

	fmt.Println("Prioritizing direction")
	fmt.Print(recentDirection)
	fmt.Println(" is recent direction")
	fmt.Println("")

	etaDirection := ETADirection(currentFloor, recentDirection, orders, allETAs)
	// if etaDirection == hardware.MD_Stop {
	// 	return recentDirection
	// } else {
	return etaDirection
	// }
}

func AllInternalETAsBest(orders AllOrders) bool {
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if !internalETABest(orders.Down[floor], internalETAs.Down[floor]) ||
			!internalETABest(orders.Up[floor], internalETAs.Up[floor]) {
			return false
		}
	}
	return true
}
