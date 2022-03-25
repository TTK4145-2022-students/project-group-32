package idledistribution

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/hardware"
)

func Distribute() {
	orders := orderstate.GetOrders()
	actualFloor := cabstate.Cab.AboveOrAtFloor

	if !orderstate.AnyOrders(orders) && cabstate.Cab.MotorDirection == hardware.MD_Stop {
		highestAttractiveness := 0
		var bestFloor int
		for simFloor := 0; simFloor < hardware.FloorCount; simFloor++ {
			floorAttractiveness := getNumberOfBestETAsForFloor(simFloor)
			if highestAttractiveness <= floorAttractiveness {
				highestAttractiveness = floorAttractiveness

				if (abs(actualFloor-simFloor) < abs(actualFloor-bestFloor)) {
					bestFloor = simFloor
				}
			}
		}
		if bestFloor < actualFloor {
			// gå ned
		} else if bestFloor > actualFloor {
			// gå opp
		}
	}
}

func getNumberOfBestETAsForFloor(floor int) int {
	orders := orderstate.GetOrders()
	internalETAs := orderstate.GetInternalETAs()

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