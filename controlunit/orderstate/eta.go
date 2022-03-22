package orderstate

import (
	"elevators/controlunit/prioritize"
	"elevators/hardware"
	"time"
)

type AllDurations struct {
	Up   [hardware.FloorCount]time.Duration
	Down [hardware.FloorCount]time.Duration
	Cab  [hardware.FloorCount]time.Duration
}

type AllETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const travelDuration = 2 * time.Second
const orderDuration = 4 * time.Second
const directionChangeCost = 2*travelDuration + orderDuration

var allDurations AllDurations

var allETAs AllETAs

func GetInternalETAs() AllETAs {
	return allETAs
}

// func ComputeETA(
// 	direction hardware.MotorDirection, aboveOrAtFloor int, destinationFloor int) time.Time {
// 	return time.Now().Add(ComputeDurationToFloor(direction, aboveOrAtFloor, destinationFloor))

// }

// func ComputeDurationToFloor(
// 	direction hardware.MotorDirection, aboveOrAtFloor int, destinationFloor int) time.Duration {
// 	// Todo: get more realistic newETA, take orders into consideration
// 	var durationSecs = 0
// 	for floor := aboveOrAtFloor; (floor < hardware.FloorCount) && (floor >= 0) && (floor != destinationFloor); floor += int(direction) {
// 		durationSecs += secsPerFloor
// 		if floor == 0 {
// 			durationSecs += destinationFloor * secsPerFloor
// 		} else if floor == hardware.FloorCount-1 {
// 			durationSecs += (floor - destinationFloor) * secsPerFloor
// 		}
// 	}
// 	return time.Duration(durationSecs) * time.Second
// }

// func travelDuration() time.Duration{
// 	return time.Duration(secsPerFloor) * time.Second
// }

// func stopDuration(caborder bool) time.Duration {
// 	if caborder {
// 		return time.Duration(secsPerOrder) * time.Second
// 	} else {
// 		return time.Duration(0) * time.Second
// 	}
// }

func UpdateInternalETAs(
	simulatedDurations AllDurations,
	simulatedETAs AllETAs) {

	for floor := 0; floor < hardware.FloorCount; floor++ {
		if simulatedDurations.Up[floor] < allDurations.Up[floor] && simulatedDurations.Up[floor] != time.Duration(0) {
			allETAs.Up[floor] = simulatedETAs.Up[floor]
			allDurations.Up[floor] = simulatedDurations.Up[floor]
		}

		if simulatedDurations.Down[floor] < simulatedDurations.Down[floor] && simulatedDurations.Down[floor] != time.Duration(0) {
			allETAs.Down[floor] = simulatedETAs.Down[floor]
			allDurations.Down[floor] = simulatedDurations.Down[floor]
		}
	}
}

func ComputeInternalETAs(durations AllDurations) AllETAs {
	var newETAs AllETAs
	var now = time.Now()
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if durations.Cab[floor] != time.Duration(0) {
			newETAs.Cab[floor] = now.Add(durations.Cab[floor])
		}
		if durations.Up[floor] != time.Duration(0) {
			newETAs.Up[floor] = now.Add(durations.Up[floor])
		}
		if durations.Cab[floor] != time.Duration(0) {
			newETAs.Down[floor] = now.Add(durations.Down[floor])
		}
	}
	return newETAs
}

// func ComputeAllDurations(
// 	prioritizedDirection hardware.MotorDirection,
// 	currentFloor int,
// 	orders AllOrders) AllDurations {
// 	var newDurations AllDurations
// 	newDurations.Cab[currentFloor] = time.Duration(0)
// 	switch prioritizedDirection {
// 	// Todo: get more realistic newDuration, take orders in floors into consideration, not only cab above/below
// 	case hardware.MD_Up:
// 		floor := currentFloor
// 		for {
// 			if floor > currentFloor {
// 				newDurations.Cab[floor] = newDurations.Cab[floor-1] + travelDuration + stopDuration(orders.Cab[floor-1])
// 			}
// 			if CabOrdersAbove(orders.Cab, floor) {
// 				newDurations.Up[floor] = newDurations.Cab[floor]
// 			} else if CabOrdersBelow(orders.Cab, currentFloor) {
// 				break // loop
// 			} else {
// 				newDurations.Up[floor] = newDurations.Cab[floor]
// 			}

// 			floor++
// 			if floor == hardware.FloorCount {
// 				floor--
// 				break // loop
// 			}
// 		}

// 		highestReachedFloor := floor
// 		newDurations.Down[floor] = newDurations.Cab[floor]

// 		for {
// 			if floor > 0 {
// 				floor--
// 			} else {
// 				break // loop
// 			}

// 			newDurations.Down[floor] = newDurations.Down[floor+1] + travelDuration
// 			if floor < currentFloor {
// 				newDurations.Cab[floor] = newDurations.Down[floor]
// 			}
// 		}

// 		if currentFloor > floor {
// 			newDurations.Up[floor] = newDurations.Cab[floor]
// 			for floor := floor + 1; floor < hardware.FloorCount; floor++ {
// 				if floor < currentFloor || floor >= highestReachedFloor {
// 					newDurations.Up[floor] = newDurations.Up[floor-1] + travelDuration
// 					if floor > highestReachedFloor {
// 						newDurations.Cab[floor] = newDurations.Up[floor]
// 					}
// 				}
// 			}
// 		}
// 	case hardware.MD_Down:
// 		floor := currentFloor
// 		for {
// 			if floor < currentFloor {
// 				newDurations.Cab[floor] = newDurations.Cab[floor+1] + travelDuration + stopDuration(orders.Cab[floor+1])
// 			}
// 			if CabOrdersAbove(orders.Cab, floor) {
// 				newDurations.Down[floor] = newDurations.Cab[floor]
// 			} else if CabOrdersBelow(orders.Cab, currentFloor) {
// 				break // loop
// 			} else {
// 				newDurations.Down[floor] = newDurations.Cab[floor]
// 			}

// 			floor--
// 			if floor == -1 {
// 				floor++
// 				break // loop
// 			}
// 		}

// 		lowestReachedFloor := floor
// 		newDurations.Up[floor] = newDurations.Cab[floor]

// 		for {
// 			if floor < hardware.FloorCount-1 {
// 				floor++
// 			} else {
// 				break // loop
// 			}

// 			newDurations.Up[floor] = newDurations.Up[floor-1] + travelDuration
// 			if floor > currentFloor {
// 				newDurations.Cab[floor] = newDurations.Up[floor]
// 			}
// 		}

// 		if currentFloor < floor {
// 			newDurations.Down[floor] = newDurations.Cab[floor]
// 			for floor := floor - 1; floor >= 0; floor-- {
// 				if floor > currentFloor || floor <= lowestReachedFloor {
// 					newDurations.Down[floor] = newDurations.Down[floor+1] + travelDuration
// 					if floor < lowestReachedFloor {
// 						newDurations.Cab[floor] = newDurations.Down[floor]
// 					}
// 				}
// 			}
// 		}
// 	}
// 	return newDurations
// }

func ComputeDurations(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs AllETAs) AllDurations {

	prioritizedDirection := PrioritizedDirection(
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
		durationsBelow := calculateETAforDirection(
			currentFloor,
			hardware.MD_Down,
			orders)
		durationsAbove := calculateETAforDirection(
			currentFloor,
			hardware.MD_Up,
			orders)
		ETAs := []AllETAs{ComputeInternalETAs(durationsBelow), ComputeInternalETAs(durationsAbove)}
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
	simulationTime := time.Duration(1)
	var simulatedDurations AllDurations
	for prioritizedDirection != hardware.MD_Stop {
		prioritizedDirection = simulateStep(
			prioritizedDirection,
			&simulationFloor,
			&simulationDirection,
			&simulationOrders,
			&simulationTime,
			&simulatedDurations)
	}
	return simulatedDurations
}

// func TestHallDurations(
// 	currentFloor int,
// 	direction hardware.MotorDirection,
// 	orders AllOrders){
// 		var computedDurations AllDurations
// 		testTime := time.Duration(1)
// 		floor := currentFloor
// 		for 0 < floor && floor < hardware.FloorCount - 1 {
// 			switch direction {
// 			case hardware.MD_Up:
// 				computedDurations.Up[floor] = testTime
// 				if hasOrder(orders.Up[floor]){
// 					testTime += orderDuration
// 				}
// 			case hardware.MD_Down:
// 				computedDurations.Down[floor] = testTime
// 				if hasOrder(orders.Down[floor]){
// 					testTime += orderDuration
// 				}
// 			}
// 			floor += int(direction)
// 			testTime += travelDuration
// 		}

// 		for floor != currentFloor{
// 			switch direction {
// 			case hardware.MD_Down:
// 				computedDurations.Up[floor] = testTime
// 				if hasOrder(orders.Up[floor]){
// 					testTime += orderDuration
// 				}
// 			case hardware.MD_Down:
// 				computedDurations.Down[floor] = testTime
// 				if hasOrder(orders.Down[floor]){
// 					testTime += orderDuration
// 				}
// 			}
// 			floor -= int(direction)
// 			testTime += travelDuration
// 		}
// 	}

/*
	var bestETA
	if elevator has direction
		bestETA = get ETA (Directions) // give eta in direction to the furthest order in direction
	else
		bestETA = max(get ETA(up), get ETA(down)) // get best eta in both directions to the furthest order in direction
	end
	updateETA(bestETA)

	o x o o x x o x o o o

	o > o o > o o o > o o
	0 0 < 0 0 < 0 0 < 0 0

	    #>
	o > o o > o o o o o o o
	0 0 0 0 0 0 < 0 0 < 0 0


	# ETAs
	simCabFloor := currentFloor
	simCabDirection := direction
	currentTime := getTime()
	for {
		simCabFloor += simCabDirection
		currentTime += travelDuration
		if simCabFloor < 0 || simCabFloor >= hardware.FloorCount {
			simCabDirection = !simCabDirection
			simCabFloor += simCabDirection
			currentTime -= travelDuration
		}

		if (simCabFloor == currentFloor) {
			return ETAs
		}

		if simCabDirection == -1 && orderDown[simCabFloor] || simCabDirection == 1 && orderUp[simCabFloor] {
			currentTime += stopDuration
		}
		if currentTime.Before(orders.Down[simCabFloor].bestETA) && simCabDirection == -1 {
			ETAs.Down[simCabFloor] = currentTime
		}
		if currentTime.Before(orders.Up[simCabFloor].bestETA) && simCabDirection == 1 {
			ETAs.Up[simCabFloor] = currentTime
		}
	}
*/

func calculateETAforDirection(
	currentFloor int,
	direction hardware.MotorDirection,
	orders AllOrders) AllDurations {

	var computedDurations AllDurations

	simCabFloor := currentFloor
	simCabDirection := int(direction)
	currentTime := time.Duration(1)
	for {
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

		if simCabDirection == -1 && hasOrder(orders.Down[simCabFloor]) || (simCabDirection == 1 && hasOrder(orders.Up[simCabFloor])) {
			currentTime += orderDuration
		}
		if simCabDirection == -1 {
			computedDurations.Down[simCabFloor] = currentTime
		}
		if simCabDirection == 1 {
			computedDurations.Up[simCabFloor] = currentTime
		}
	}

}

func bestDurations(
	floor int,
	orders AllOrders,
	ETAs []AllETAs) int {

	ETAsBelow := ETAs[0]
	ETAsAbove := ETAs[1]

	belowFloor := floor
	aboveFloor := floor

	belowDir := -1
	aboveDir := 1

	for {
		belowFloor += belowDir
		aboveFloor += aboveDir
		if belowFloor < 0 {
			belowFloor = 0
			belowDir = 1
		}
		if aboveFloor >= hardware.FloorCount {
			aboveFloor = hardware.FloorCount - 1
			aboveDir = -1
		}

		if aboveFloor == floor || belowFloor == floor {
			break
		}

		floorETABelow := ETAsBelow.Down[belowFloor]
		orderETABelow := orders.Down[belowFloor].BestETA
		if belowDir == 1 {
			floorETABelow = ETAsBelow.Up[belowFloor]
			orderETABelow = orders.Up[belowFloor].BestETA
		}
		floorETAAbove := ETAsAbove.Up[aboveFloor]
		orderETAAbove := orders.Up[aboveFloor].BestETA
		if aboveDir == -1 {
			floorETAAbove = ETAsAbove.Down[aboveFloor]
			orderETAAbove = orders.Down[aboveFloor].BestETA
		}

		if floorETABelow.Before(orderETABelow) && !floorETAAbove.Before(orderETAAbove) {
			return 0
		}
		if floorETAAbove.Before(orderETAAbove) && !floorETABelow.Before(orderETABelow) {
			return 1
		}
	}

	if 2*floor < hardware.FloorCount {
		return 0
	} else {
		return 1
	}
}

// The function will return the best ETA based on which ETA table it first finds an improved ETA in from the current floor.
func bestETA(
	floor int,
	orders AllOrders,
	ETAsBelow AllETAs,
	ETAsAbove AllETAs) AllETAs {

	belowFloor := floor
	aboveFloor := floor

	belowDir := -1
	aboveDir := 1

	for {
		belowFloor += belowDir
		aboveFloor += aboveDir
		if belowFloor < 0 {
			belowFloor = 0
			belowDir = 1
		}
		if aboveFloor >= hardware.FloorCount {
			aboveFloor = hardware.FloorCount - 1
			aboveDir = -1
		}

		if aboveFloor == floor || belowFloor == floor {
			break
		}

		floorETABelow := ETAsBelow.Down[belowFloor]
		orderETABelow := orders.Down[belowFloor].BestETA
		if belowDir == 1 {
			floorETABelow = ETAsBelow.Up[belowFloor]
			orderETABelow = orders.Up[belowFloor].BestETA
		}
		floorETAAbove := ETAsAbove.Up[aboveFloor]
		orderETAAbove := orders.Up[aboveFloor].BestETA
		if aboveDir == -1 {
			floorETAAbove = ETAsAbove.Down[aboveFloor]
			orderETAAbove = orders.Down[aboveFloor].BestETA
		}

		if floorETABelow.Before(orderETABelow) && !floorETAAbove.Before(orderETAAbove) {
			return ETAsBelow
		}
		if floorETAAbove.Before(orderETAAbove) && !floorETABelow.Before(orderETABelow) {
			return ETAsAbove
		}
	}

	if 2*floor < hardware.FloorCount {
		return ETAsAbove
	} else {
		return ETAsBelow
	}

	// belowConcatETAs := append(ETAsBelow.Down[:floor], ETAsBelow.Up[:floor]...)
	// aboveConcatETAs := append(ETAsAbove.Up[floor+1:], ETAsAbove.Down[floor+1:]...)

	// belowGlobalConcatETAs := append(globalETAs.Down[:floor], globalETAs.Up[:floor]...)
	// aboveGlobalConcatETAs := append(globalETAs.Up[floor+1:], globalETAs.Down[floor+1:]...)

	// for i := 0; i < min(len(belowConcatETAs), len(aboveConcatETAs)); i++ {
	// 	if belowConcatETAs[i] < glo
	// }

	// if len(belowConcatETAs) > len(aboveConcatETAs) {
	// 	return ETAsBelow
	// } else {
	// 	return ETAsAbove
	// }
}

// func min(a, b int) int {
// 	if a <= b {
// 		return a
// 	}
// 	return b
// }

func simulateStep(
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
			return newDirection
		}
		*floor += int(newDirection)
		*simTime += travelDuration

	case hardware.DS_Open_Down:
		durations.Down[*floor] = *simTime
	case hardware.DS_Open_Up:
		durations.Up[*floor] = *simTime
	case hardware.DS_Open_Cab:
		break
	default:
		panic("Invalid door action in eta simulation")
	}

	if doorAction != hardware.DS_Close {
		*simTime += orderDuration
		//Todo handle rest of simulation step
	}
	return prioritizedDirection
}

func internalETABest(orderState OrderState, internalETA time.Time) bool {
	return orderState.BestETA.Equal(internalETA) || !internalETA.IsZero()
}

func orderAndInternalETABestAbove(
	currentFloor int,
	orders AllOrders,
	allETAs AllETAs) bool {
	for floor := currentFloor + 1; floor < hardware.FloorCount; floor++ {
		if (hasOrder(orders.Up[floor]) && internalETABest(orders.Up[floor], allETAs.Up[floor])) || (hasOrder(orders.Down[floor]) && internalETABest(orders.Down[floor], allETAs.Up[floor])) || orders.Cab[floor] {
			return true
		}
	}
	return false
}

func orderAndInternalETABestBelow(
	currentFloor int,
	orders AllOrders,
	allETAs AllETAs) bool {
	for floor := currentFloor - 1; floor >= 0; floor-- {
		if (hasOrder(orders.Up[floor]) && internalETABest(orders.Up[floor], allETAs.Up[floor])) || (hasOrder(orders.Down[floor]) && internalETABest(orders.Down[floor], allETAs.Up[floor])) || orders.Cab[floor] {
			return true
		}
	}
	return false
}

func ETADirection(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs AllETAs) hardware.MotorDirection {

	switch recentDirection {
	case hardware.MD_Up:
		if orderAndInternalETABestAbove(currentFloor, orders, allETAs) {
			return hardware.MD_Up
		}
		if orderAndInternalETABestBelow(currentFloor, orders, allETAs) {
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if orderAndInternalETABestBelow(currentFloor, orders, allETAs) {
			return hardware.MD_Down
		}
		if orderAndInternalETABestAbove(currentFloor, orders, allETAs) {
			return hardware.MD_Up
		}
	}
	return hardware.MD_Stop
}

func PrioritizedDirection(currentFloor int,
	recentDirection hardware.MotorDirection,
	orders AllOrders,
	allETAs AllETAs) hardware.MotorDirection {

	etaDirection := ETADirection(currentFloor, recentDirection, orders, allETAs)
	if etaDirection == hardware.MD_Stop {
		return recentDirection
	} else {
		return etaDirection
	}
}
