package main

import (
	"elevators/cab"
	"elevators/eta"
	"elevators/filesystem"
	"elevators/hardware"
	"elevators/orders"
	"fmt"
	"time"
)

// {
// filename := "test"
// filestate := orders.GetInternalETAs()
// if len(os.Args) > 1 {
// 	file,_ := json.MarshalIndent(filestate,"", " ")
// 	_ = ioutil.WriteFile("testresults/" + filename+os.Args[1]+".json",file, 0644)
// } else {
// 	file,_ := json.MarshalIndent(filestate,"", " ")
// 	_ = ioutil.WriteFile("testresults/" + filename+".json",file,0644)
// }
// }

func testAnyOrders() {
	filesystem.Init()
	allOrders := filesystem.ReadOrders()
	fmt.Print("Any orders: ")
	fmt.Println(orders.AnyOrders(allOrders))
}

func testFirstInternalETA() {
	filesystem.Init()
	allOrders := filesystem.ReadOrders()
	cabstate := filesystem.ReadCabState()
	etas := eta.ComputeETAs(cabstate.AboveOrAtFloor,
		hardware.MD_Stop,
		hardware.MD_Down,
		cabstate.Behaviour == cab.DoorOpen,
		allOrders)
	fmt.Println("now: ")
	fmt.Println(time.Now())
	fmt.Println("first internal eta expire: ")
	fmt.Println(eta.FirstInternalETAExpiration(etas))
}

func testComputeETAs() {
	filesystem.Init()
	allOrders := filesystem.ReadOrders()
	cabstate := filesystem.ReadCabState()
	etas := eta.ComputeETAs(cabstate.AboveOrAtFloor,
		hardware.MD_Stop,
		cabstate.RecentDirection,
		cabstate.Behaviour == cab.DoorOpen,
		allOrders)
	filesystem.Write("testresults/"+"computeETAs.json",
		etas)
	etas = eta.ComputeETAs(cabstate.AboveOrAtFloor,
		hardware.MD_Stop,
		hardware.MD_Down,
		cabstate.Behaviour == cab.DoorOpen,
		allOrders)
	filesystem.Write("testresults/"+"computeETAsDown.json",
		etas)
	etas = eta.ComputeETAs(cabstate.AboveOrAtFloor,
		hardware.MD_Stop,
		hardware.MD_Up,
		cabstate.Behaviour == cab.DoorOpen,
		allOrders)
	filesystem.Write("testresults/"+"computeETAsUp.json",
		etas)
}

func main() {
	// testHasOrder()
	// testAnyOrders()
	// testFirstExternalETA()
	// testComputeETAs()
	fmt.Println(hardware.ValidFloors())
	fmt.Println(time.Now().Before(time.Unix(1<<62, 0)))
	testFirstInternalETA()
}
