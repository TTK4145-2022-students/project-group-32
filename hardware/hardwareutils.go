package hardware

const FloorCount = 4

const BottomFloor = 0

const TopFloor = FloorCount - 1

func ValidFloor(floor int) bool {
	return BottomFloor <= floor && floor <= TopFloor
}

func FloorBelowMiddleFloor(floor int) bool {
	return floor < FloorCount/2
}

func ValidFloors() [FloorCount]int {
	var floorIndices [FloorCount]int
	for floor := BottomFloor; ValidFloor(floor); floor++ {
		floorIndices[floor] = floor
	}
	return floorIndices
}
