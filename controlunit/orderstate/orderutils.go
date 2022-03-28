package orderstate

import "elevators/hardware"

func (orders AllOrders) getOrderState(
	direction hardware.MotorDirection,
	floor int) OrderState {

	switch direction {
	case hardware.MD_Down:
		return orders.Down[floor]

	case hardware.MD_Up:
		return orders.Up[floor]

	default:
		panic("Invalid direction to get order")
	}
}
