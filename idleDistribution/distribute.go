package idledistribution

import (
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func MotorActionOnDistribute(
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
