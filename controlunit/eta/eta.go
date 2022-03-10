package eta

import (
	"elevators/controlunit/orderstate"
	"time"
)

const timeEstimateBetweenFloors = 2

const timeEstimatePerOrder = 4

func ComputeETA(
	orders orderstate.AllOrders, aboveOrAtFloor int, destinationFloor int) time.Time {
	// Todo: get more realistic newETA, take direction into consideration
	var newETA = time.Now()
	newETA = newETA.Add(time.Second * timeEstimateBetweenFloors * time.Duration(destinationFloor-aboveOrAtFloor))
	newETA = newETA.Add(time.Second * timeEstimatePerOrder * time.Duration(orderstate.OrdersBetween(orders, aboveOrAtFloor, destinationFloor)))
	return newETA
}
