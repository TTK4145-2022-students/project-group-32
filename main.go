package main

import (

	// "elevators/filesystem"
	"elevators/controlunit"
	"elevators/hardware"

	// "elevators/phoenix"
	//"time"
	"os"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()

	if len(os.Args) > 1 {
		hardware.Init("localhost:"+os.Args[1], hardware.FloorCount)
	} else {
		hardware.Init("localhost:15657", hardware.FloorCount)
	}
	controlunit.Init()

	go controlunit.RunElevatorLoop()

	for {
	}
}
