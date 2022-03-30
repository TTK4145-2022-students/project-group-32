package main

import (
	"elevators/cab"
	"elevators/filesystem"
	"elevators/hardware"
	"elevators/network"
	"elevators/orders"
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
	orders.Init(filesystem.ReadOrders())
	cab.Init()

	buttonPress := make(chan hardware.ButtonEvent)
	floorArrival := make(chan int)
	floorLeft := make(chan bool)
	obstructionChange := make(chan bool)
	stopChange := make(chan bool)

	doorTimedOut := make(chan bool)
	decisionDeadlineTimedOut := make(chan bool)
	etaExpiredAlarmRinging := make(chan bool)

	pokeCab := make(chan bool)

	ordersRecieved := make(chan orders.AllOrders)

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
			orders.AcceptNewOrder(
				buttonEvent.Button,
				buttonEvent.Floor)
			allOrders := orders.GetOrders()
			cab.FSMNewOrder(
				buttonEvent.Floor,
				allOrders)

		case floor := <-floorArrival:
			hardware.SetFloorIndicator(floor)
			allOrders := orders.GetOrders()
			cab.FSMFloorArrival(
				floor,
				allOrders)
		case <-floorLeft:
			cab.FSMFloorLeave()

		case obstruction := <-obstructionChange:
			allOrders := orders.GetOrders()
			cab.FSMObstructionChange(
				obstruction,
				allOrders)
		case <-doorTimedOut:
			allOrders := orders.GetOrders()
			cab.FSMDoorTimeout(allOrders)

		case <-decisionDeadlineTimedOut:
			allOrders := orders.GetOrders()
			cab.FSMDecisionDeadline(allOrders)

		case <-pokeCab:
			allOrders := orders.GetOrders()
			cab.FSMPoke(allOrders)
			timer.PokeCabTimer.TimerStart()

		case <-stopChange:
			orders.ResetOrders()
			for f := 0; f < hardware.FloorCount; f++ {
				for b := hardware.ButtonType(0); b < 3; b++ {
					hardware.SetButtonLamp(
						b,
						f,
						false)
				}
			}

		case recievedOrderState := <-ordersRecieved:
			newOrdersInFloors := orders.UpdateOrders(recievedOrderState)
			allOrders := orders.GetOrders()
			for floor, newOrder := range newOrdersInFloors {
				if newOrder {
					cab.FSMNewOrder(
						floor,
						allOrders)
				}
			}
		}
	}
}
