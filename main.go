package main

import (
	"elevators/network"
	"time"
)

func main() {
	_, wconn := network.InitUDPSendingSocket(20014, "255.255.255.255")
	_, conn := network.InitUDPReceivingSocket(20014)
	network.BroadcastMessage("Kan dette virke mon tro", wconn)
	time.Sleep(time.Millisecond * 400)
	network.ReceiveUDPMessage(conn)
}
