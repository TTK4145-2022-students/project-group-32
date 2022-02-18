package main

import (
	"elevators/controlunit"
	"elevators/network"
	"time"
)

func main() {

	controlunit.Init()
	UDPPort := 20014
	_, wconn := network.InitUDPSendingSocket(UDPPort, "255.255.255.255")
	_, conn := network.InitUDPReceivingSocket(UDPPort)
	for {
		network.BroadcastMessage("Kan dette virke mon tro", wconn)
		time.Sleep(time.Millisecond * 4000)
		network.ReceiveUDPMessage(conn)
	}
}
