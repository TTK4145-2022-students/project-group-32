package main

import (
	"elevators/controlunit/cabstate"
	"elevators/controlunit/orderstate"
	"elevators/filesystem"
	"elevators/hardware"
	"elevators/network"
	"elevators/phoenix"
	"elevators/timer"
	"os"
)

func main() {
	phoenix.Init()
	go phoenix.Phoenix()
	if len(os.Args) > 1 {
		hardware.Init(
			"localhost:"+os.Args[1],
			hardware.FloorCount)
	} else {
		hardware.Init(
			"localhost:15657",
			hardware.FloorCount)
	}
	filesystem.Init()
	orderstate.Init(filesystem.ReadOrders())
	// cabstate.Init(filesystem.ReadCabState())
	cabstate.Init()

	buttonPress := make(chan hardware.ButtonEvent)
	floorArrival := make(chan int)
	floorLeft := make(chan bool)
	obstructionChange := make(chan bool)
	stopChange := make(chan bool)

	doorTimedOut := make(chan bool)
	decisionDeadlineTimedOut := make(chan bool)
	etaExpiredAlarmRinging := make(chan bool)
	// internalETAExpiringAlarmRinging := make(chan bool)

	pokeCab := make(chan bool)

	ordersRecieved := make(chan orderstate.AllOrders)

	go hardware.PollButtons(buttonPress)
	go hardware.PollFloorSensor(
		floorArrival,
		floorLeft)
	go hardware.PollObstructionSwitch(obstructionChange)
	go hardware.PollStopButton(stopChange)

	go timer.DoorTimer.PollTimerOut(doorTimedOut)
	go timer.DecisionDeadlineTimer.PollTimerOut(decisionDeadlineTimedOut)
	go timer.ETAExpiredAlarm.PollAlarm(etaExpiredAlarmRinging)
	go timer.InternalETAExpiringAlarm.PollAlarm(etaExpiredAlarmRinging)
	go timer.PokeCabTimer.PollTimerOut(pokeCab)

	go network.PollReceiveOrders(ordersRecieved)
	go network.SendOrdersPeriodically()

	go filesystem.SaveOrdersPeriodically()

	timer.PokeCabTimer.TimerStart()
	// All hail The Loop!
	// that grants us non-concurrency and fastest poll rate
	for {
		select {
		case buttonEvent := <-buttonPress:
			// fmt.Println("Button pressed")
			orderstate.AcceptNewOrder(
				buttonEvent.Button,
				buttonEvent.Floor)
			orders := orderstate.GetOrders()
			cabstate.FSMNewOrder(
				buttonEvent.Floor,
				orders)

		case floor := <-floorArrival:
			// fmt.Println("Ariived at floor")
			hardware.SetFloorIndicator(floor)
			orders := orderstate.GetOrders()
			cabstate.FSMFloorArrival(floor, orders)

		case <-floorLeft:
			cabstate.FSMFloorLeave()

		case obstruction := <-obstructionChange:
			orders := orderstate.GetOrders()
			cabstate.FSMObstructionChange(obstruction, orders)

		case <-doorTimedOut:
			// fmt.Println("Door timed out")
			orders := orderstate.GetOrders()
			cabstate.FSMDoorTimeout(orders)

		case <-decisionDeadlineTimedOut:
			// fmt.Println("Decision timed out")
			cabstate.FSMDecisionDeadline()

		case <-pokeCab:
			cabstate.FSMPoke()
			timer.PokeCabTimer.TimerStart()

		case <-stopChange:
			orderstate.ResetOrders()
			for f := 0; f < hardware.FloorCount; f++ {
				for b := hardware.ButtonType(0); b < 3; b++ {
					hardware.SetButtonLamp(b, f, false)
				}
			}
		case recievedOrderState := <-ordersRecieved:
			newOrdersInFloors := orderstate.UpdateOrders(recievedOrderState)
			orders := orderstate.GetOrders()
			for floor, newOrder := range newOrdersInFloors {
				if newOrder {
					// fmt.Println("recieved new order")
					cabstate.FSMNewOrder(floor, orders)
				}
			}
		}
	}
}
