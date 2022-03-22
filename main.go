package main

import (

	// "elevators/filesystem"
	"elevators/controlunit"
	//"elevators/filesystem"
	"elevators/hardware"
	// "fmt"
	// "os"
	"time"
	// "io/ioutil"
	// "elevators/phoenix"
)

func main() {
	// phoenix.Init()
	// go phoenix.Phoenix()
	hardware.Init("localhost:15657", hardware.FloorCount)
	controlunit.Init()

	go controlunit.RunElevatorLoop()
	// for{
	// 	time.Sleep(time.Second)
	// 	jsonFile, _ := os.Open("filesystem/orderState.json")
	// 	defer jsonFile.Close()
	// 	byteValue, _ := ioutil.ReadAll(jsonFile)
	// 	fmt.Println(string(byteValue))
	// }

	for {
		time.Sleep(1 * time.Hour)
	}
}
