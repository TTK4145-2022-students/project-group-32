package main

// Copy and paste the following code snippet into a function
// to get a snapshot of the orders in json file format:

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

// Write unit test functions below when debugging,
// remember to make functions temporarily global:

// func testComputeETAs() {
// 	filesystem.Init()
// 	allOrders := filesystem.ReadOrders()
// 	cabstate := filesystem.ReadCabState()
// 	etas := eta.ComputeETAs(cabstate.AboveOrAtFloor,
// 		hardware.MD_Stop,
// 		cabstate.RecentDirection,
// 		cabstate.Behaviour == cab.DoorOpen,
// 		allOrders)
// 	filesystem.Write("testresults/"+"computeETAs.json",
// 		etas)
// 	etas = eta.ComputeETAs(cabstate.AboveOrAtFloor,
// 		hardware.MD_Stop,
// 		hardware.MD_Down,
// 		cabstate.Behaviour == cab.DoorOpen,
// 		allOrders)
// 	filesystem.Write("testresults/"+"computeETAsDown.json",
// 		etas)
// 	etas = eta.ComputeETAs(cabstate.AboveOrAtFloor,
// 		hardware.MD_Stop,
// 		hardware.MD_Up,
// 		cabstate.Behaviour == cab.DoorOpen,
// 		allOrders)
// 	filesystem.Write("testresults/"+"computeETAsUp.json",
// 		etas)
// }

func main() {
	// testComputeETAs()
}
