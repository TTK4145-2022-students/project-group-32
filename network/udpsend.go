package network

import (
	"elevators/orders"
	"encoding/json"
	"net"
	"time"
)

const UDPPort = 20014
const broadcastAddr = "255.255.255.255"
const sendRate = orders.WaitBeforeGuaranteeTime / 5

func InitUDPSendingSocket(
	port int,
	sendAddr string) (
	net.UDPAddr,
	*net.UDPConn) {

	sendaddr := net.UDPAddr{
		Port: port,
		IP: net.ParseIP(
			sendAddr),
	}

	wconn, _ := net.DialUDP(
		"udp",
		nil,
		&sendaddr)

	return sendaddr,
		wconn
}

func BroadcastOrders(
	allOrders orders.AllOrders,
	wconn *net.UDPConn) {

	message, _ := json.Marshal(allOrders)

	broadcastMessage(
		message,
		wconn)
}

func broadcastMessage(
	message []byte,
	wconn *net.UDPConn) {

	wconn.Write(message)
}

func SendOrdersPeriodically() {
	_, wconn := InitUDPSendingSocket(
		UDPPort,
		broadcastAddr)

	defer wconn.Close()

	for {
		allOrders := orders.GetOrders()

		BroadcastOrders(
			allOrders,
			wconn)

		time.Sleep(sendRate)
	}
}
