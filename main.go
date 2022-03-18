package main

import (

	// "elevators/filesystem"
	"elevators/controlunit"
	"elevators/hardware"

	// "elevators/phoenix"
	//"time"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()

	hardware.Init("localhost:15657", hardware.FloorCount)
	controlunit.Init()

	go controlunit.RunElevatorLoop()

	for {
	}
}
