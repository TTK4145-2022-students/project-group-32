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

type InternalETAs struct {
	Up   [hardware.FloorCount]time.Time
	Down [hardware.FloorCount]time.Time
	Cab  [hardware.FloorCount]time.Time
}

const travelDuration = 5 * time.Second
const orderDuration = 4 * time.Second
const offsetDuration = 2 * time.Second

// const directionChangeCost = 2*travelDuration + orderDuration

// var internalDurations AllDurations

var internalETAs InternalETAs

func GetInternalETAs() InternalETAs {
	return internalETAs
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
// 	return time.Duration(durationSecs) * time.Secders(orders) && !AllInternalETAsBest(orders) {
// 	fmt.Println("prioritizing to prepare")
// 	if 0 < floor && 2*floor < hardware.FloorCount &&
// 		internalETABest(orders.Up[floor-1], allETAs.Up[floor-1]) {
// 		return hardware.MD_Down
// 	} else if floor < hardware.FloorCount-1 &&
// 		internalETABest(orders.Down[floor+1], allETAs.Down[floor+1]) {
// 		return hardware.MD_Up
// 	}
// }
// func stopDuration(caborder bool) time.Duration {
// 	if caborder {
// 		return time.Duration(secsPerOrder) * time.Second
// 	} else {
// 		return time.Duration(0) * time.Second
// 	}
// }

func UpdateETAs(
	recentDirection hardware.MotorDirection,
	currentFloor int) {

	// fmt.Println("updateing eats")

	newDurations := ComputeDurations(currentFloor, recentDirection, allOrders, internalETAs)
	newETAs := ComputeInternalETAs(newDurations)

	for floor := 0; floor < hardware.FloorCount; floor++ {
		// if newDurations.Up[floor] < allDurations.Up[floor] &&
		if !newETAs.Up[floor].IsZero() &&
			(newETAs.Up[floor].Before(allOrders.Up[floor].BestETA) ||
				allOrders.Up[floor].BestETA.Before(time.Now())) {
			// fmt.Println("setting best up")
			allOrders.Up[floor].BestETA = newETAs.Up[floor]
		} else if internalETAs.Up[floor].Equal(allOrders.Up[floor].BestETA) &&
			!allOrders.Up[floor].BestETA.IsZero() { /*&&
			internalETAs.Down[floor].After(time.Now())*/
			// Make sure to keep ownership
			// fmt.Println("has a best up")
			newETAs.Up[floor] = allOrders.Up[floor].BestETA
		}

		// if newDurations.Down[floor] < allDurations.Down[floor] &&
		if !newETAs.Down[floor].IsZero() &&
			(newETAs.Down[floor].Before(allOrders.Down[floor].BestETA) ||
				allOrders.Down[floor].BestETA.Before(time.Now())) {
			// fmt.Println("setting best down")
			allOrders.Down[floor].BestETA = newETAs.Down[floor]
		} else if internalETAs.Down[floor].Equal(allOrders.Down[floor].BestETA) &&
			!allOrders.Down[floor].BestETA.IsZero() { /*&&
			internalETAs.Down[floor].After(time.Now())*/
			// Make sure to keep ownership
			// fmt.Println("has a best down")
			newETAs.Down[floor] = allOrders.Down[floor].BestETA
		}
		allOrders.Up[floor].LocalETA = newETAs.Up[floor]
		allOrders.Down[floor].LocalETA = newETAs.Down[floor]
		allOrders.Up[floor].Now = time.Now()
		allOrders.Down[floor].Now = time.Now()
		// fmt.Print(newETAs)
	}
	// updateInternalETAs(newDurations, newETAs)
	// internalDurations = newDurations
	internalETAs = newETAs
}

// func updateInternalETAs(
// 	simulatedDurations AllDurations,
// 	simulatedETAs AllETAs) {

// 	for floor := 0; floor < hardware.FloorCount; floor++ {
// 		if allDurations.Up[floor] == time.Duration(0) ||
// 			(simulatedDurations.Up[floor] < allDurations.Up[floor] &&
// 				simulatedDurations.Up[floor] != time.Duration(0)) {
// 			allETAs.Up[floor] = simulatedETAs.Up[floor]
// 			allDurations.Up[floor] = simulatedDurations.Up[floor]
// 		}

// 		if allDurations.Up[floor] == time.Duration(0) ||
// 			(simulatedDurations.Down[floor] < allDurations.Down[floor] &&
// 				simulatedDurations.Down[floor] != time.Duration(0)) {
// 			allETAs.Down[floor] = simulatedETAs.Down[floor]
// 			allDurations.Down[floor] = simulatedDurations.Down[floor]
// 		}
// 	}
// }

func ComputeInternalETAs(durations AllDurations) InternalETAs {
	var newETAs InternalETAs
	var now = time.Now()
	// fmt.Println("compute internal etas")
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if durations.Cab[floor] != time.Duration(0) {
			newETAs.Cab[floor] = now.Add(durations.Cab[floor])
			// fmt.Print("uc,")
		}
		if durations.Up[floor] != time.Duration(0) {
			newETAs.Up[floor] = now.Add(durations.Up[floor])
			// fmt.Print("uu,")
		}
		if durations.Down[floor] != time.Duration(0) {
			newETAs.Down[floor] = now.Add(durations.Down[floor])
			// fmt.Print("ud,")
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
		// fmt.Print("simulate durations")
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
			// fmt.Print("split durations below")
			return durationsBelow
		case 1:
			// fmt.Print("split durations above")
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

// func TestHallDurations(
// 	currentFloor int,
// 	direction hardware.MotorDirection,
// 	orders AllOrders){
// 		var computedDurations AllDurations
// 		testTime := offsetDuration
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
	// fmt.Println("Comparing durations")
	// fmt.Println(" ")
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

		// if floorOrderETAAbove {
		// 	// fmt.Print("sim has order in")
		// 	fmt.Println(aboveFloor)
		// }
		// if floorOrderETABelow {
		// 	// fmt.Print("sim has order up in")
		// 	fmt.Println(belowFloor)
		// }
		// if floorETABelow.Before(orderETABelow) {
		// 	// fmt.Println("below before order")
		// }
		// if floorOrderETABelow &&
		// 	now.After(orderETABelow) {
		// 	fmt.Println("below order and eta expired")
		// }
		// if floorETAAbove.Before(orderETAAbove) {
		// 	fmt.Println("above before order")
		// }
		// if floorOrderETAAbove &&
		// 	now.After(orderETAAbove) {
		// 	fmt.Println("above order and eta expired")
		// }
		if (floorETABelow.Before(orderETABelow) ||
			(floorOrderETABelow &&
				now.After(orderETABelow))) &&

			!(floorETAAbove.Before(orderETAAbove) ||
				(floorOrderETAAbove &&
					now.After(orderETAAbove))) {
			// fmt.Println("splitting down")
			return 0
		}
		if (floorETAAbove.Before(orderETAAbove) ||
			(floorOrderETAAbove &&
				now.After(orderETAAbove))) &&
			!(floorETABelow.Before(orderETABelow) ||
				(floorOrderETABelow &&
					now.After(orderETABelow))) {
			// fmt.Println("splitting  up")
			return 1
		}
	}

	if 2*floor < hardware.FloorCount {
		// fmt.Println("splitting down because of position and end of simulation")
		return 1
	} else {
		// fmt.Println("splitting up because of position and end of simulation")
		return 0
	}
}

// The function will return the best ETA based on which ETA table it first finds an improved ETA in from the current floor.
func bestETA(
	floor int,
	orders AllOrders,
	ETAsBelow InternalETAs,
	ETAsAbove InternalETAs) InternalETAs {

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

func internalETABest(orderState OrderState, internalETA time.Time) bool {
	return orderState.BestETA.Equal(internalETA) && !internalETA.IsZero()
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
			// fmt.Println("best above")
			return hardware.MD_Up
		}
		if orderAndInternalETABest(hardware.MD_Down, floor, orders, allETAs) {
			// fmt.Println("best below")
			return hardware.MD_Down
		}
	case hardware.MD_Down:
		if orderAndInternalETABest(hardware.MD_Down, floor, orders, allETAs) {
			// fmt.Println("best below")
			return hardware.MD_Down
		}
		if orderAndInternalETABest(hardware.MD_Up, floor, orders, allETAs) {
			// fmt.Println("best above")
			return hardware.MD_Up
		}
	}
	// if !AnyOrders(orders) && !AllInternalETAsBest(orders) {
	// 	fmt.Println("prioritizing to prepare")
	// 	if 0 < floor && 2*floor < hardware.FloorCount &&
	// 		internalETABest(orders.Up[floor-1], allETAs.Up[floor-1]) {
	// 		return hardware.MD_Down
	// 	} else if floor < hardware.FloorCount-1 &&
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
