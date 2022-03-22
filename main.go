package main

import (

	// "elevators/filesystem"
	"elevators/controlunit"
	"elevators/hardware"
	"time"
	// "elevators/phoenix"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()
	hardware.Init("localhost:15657", hardware.FloorCount)
	controlunit.Init()

	go controlunit.RunElevatorLoop()

	for {
		time.Sleep(1 * time.Hour)
	}
}
