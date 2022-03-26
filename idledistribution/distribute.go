package idledistribution

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
	"fmt"
	"time"
)

// const maxDurationDiffFromFloor = orderstate.OffsetDuration

func MotorActionOnDistributeTrygve(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders orderstate.AllOrders) hardware.MotorDirection {

	if !orderstate.AnyOrders(orders) {
		highestAttractiveness := 0
		var bestFloor int
		for simFloor := 0; simFloor < hardware.FloorCount; simFloor++ {

			direction := recentDirection
			if simFloor < currentFloor {
				direction = hardware.MD_Up
			} else if simFloor > currentFloor {
				direction = hardware.MD_Down
			}
			floorAttractiveness := getNumberOfBestETAsForFloor(
				simFloor,
				direction,
				orders)
			if highestAttractiveness <= floorAttractiveness {
				highestAttractiveness = floorAttractiveness

				if abs(currentFloor-simFloor) < abs(currentFloor-bestFloor) {
					bestFloor = simFloor
				}
			}
		}
		if bestFloor < currentFloor {
			return hardware.MD_Down
		} else if bestFloor > currentFloor {
			return hardware.MD_Up
		}
	}
	return hardware.MD_Stop
}

func MotorActionOnDistribute(
	currentFloor int,
	recentDirection hardware.MotorDirection,
	orders orderstate.AllOrders,
	internalETAs orderstate.InternalETAs) hardware.MotorDirection {

	if !orderstate.AnyOrders(orders) &&
		!orderstate.AllInternalETAsBest(orders) {
		if currentFloor == 0 &&
			orderstate.InternalETABest(
				orders.Up[currentFloor+1],
				internalETAs.Up[currentFloor+1]) {
			fmt.Println("preparing bottom up")
			return hardware.MD_Up
		} else if currentFloor == hardware.FloorCount-1 &&
			orderstate.InternalETABest(
				orders.Down[currentFloor-1],
				internalETAs.Down[currentFloor-1]) {
			fmt.Println("preparing top down")
			return hardware.MD_Down
		} else if 0 < currentFloor &&
			currentFloor < hardware.FloorCount-1 &&
			(!orderstate.InternalETABest(orders.Up[currentFloor],
				internalETAs.Up[currentFloor]) &&
				time.Until(orders.Up[currentFloor].BestETA) < orderstate.OffsetDuration) &&
			orderstate.InternalETABest(
				orders.Up[currentFloor-1],
				internalETAs.Up[currentFloor-1]) {
			fmt.Println("preparing middle down")
			return hardware.MD_Down
		} else if 0 < currentFloor &&
			currentFloor < hardware.FloorCount-1 &&
			(!orderstate.InternalETABest(orders.Down[currentFloor],
				internalETAs.Down[currentFloor]) &&
				time.Until(orders.Up[currentFloor].BestETA) < orderstate.OffsetDuration) &&
			orderstate.InternalETABest(
				orders.Down[currentFloor+1],
				internalETAs.Down[currentFloor+1]) {
			fmt.Println("preparing middle up")
			return hardware.MD_Up
		}
	}
	return hardware.MD_Stop
}

func getNumberOfBestETAsForFloor(
	floor int,
	recentDirection hardware.MotorDirection,
	orders orderstate.AllOrders) int {

	internalETAs := orderstate.ComputeETAs(
		floor,
		hardware.MD_Stop,
		recentDirection,
		orders)

	numBestETAs := 0
	for floor := 0; floor < hardware.FloorCount; floor++ {
		if orders.Down[floor].BestETA.Before(internalETAs.Down[floor]) {
			numBestETAs++
		}
		if orders.Up[floor].BestETA.Before(internalETAs.Up[floor]) {
			numBestETAs++
		}
	}

	return numBestETAs
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// func AssumeCabPositionsFromETAs(
// 	orders orderstate.AllOrders,
// 	internalETAs orderstate.InternalETAs) [hardware.FloorCount]bool {
// 	var cabPositions [hardware.FloorCount]bool
// 	now := time.Now()
// 	for floor := 0; floor < hardware.FloorCount; floor++ {
// 		if (
// !orderstate.InternalETABest(orders.Up[floor],
// internalETAs.Up[floor]) &&
// 			orders.Up[floor].BestETA.Sub(now) < maxDurationDiffFromFloor &&
// 			orders.Up[floor].BestETA.Sub(now) > -maxDurationDiffFromFloor) ||
// 			(
// !orderstate.InternalETABest(orders.Down[floor],
// internalETAs.Down[floor]) &&
// 				orders.Down[floor].BestETA.Sub(now) < maxDurationDiffFromFloor &&
// 				orders.Down[floor].BestETA.Sub(now) > -maxDurationDiffFromFloor) {
// 			cabPositions[floor] = true
// 		}
// 	}
// 	return cabPositions
// }
