package cab

func cabInFloor(floor int) bool {
	return Cab.AboveOrAtFloor == floor &&
		!Cab.BetweenFloors
}
