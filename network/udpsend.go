package network

import (
	"elevators/controlunit/orderstate"
	"encoding/json"
	"net"
	"time"
)

const UDPPort = 20014
const broadcastAddr = "255.255.255.255"
const _sendRate = orderstate.WaitBeforeGuaranteeTime / 5

func InitUDPSendingSocket(port int, sendAddr string) (net.UDPAddr, *net.UDPConn) {
	sendaddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(sendAddr),
	}

	wconn, err := net.DialUDP("udp", nil, &sendaddr)

	if err != nil {
		panic(err)
	}

	return sendaddr, wconn
}

func BroadcastOrders(orders orderstate.AllOrders, wconn *net.UDPConn) {
	message, _ := json.Marshal(orders)
	broadcastMessage(message, wconn)
}

func broadcastMessage(message []byte, wconn *net.UDPConn) {
	wconn.Write(message)
}

func SendOrdersPeriodically() {
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	defer wconn.Close()

	for {
		orders := orderstate.GetOrders()
		BroadcastOrders(orders, wconn)
		time.Sleep(_sendRate)
	}
}
