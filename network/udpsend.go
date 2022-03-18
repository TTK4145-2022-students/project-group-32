package network

import (
	"elevators/controlunit/orderstate"
	"encoding/json"
	"net"
	"time"
)

const UDPPort = 20014
const broadcastAddr = "255.255.255.255"
const _sendRate = 40 * time.Millisecond

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

func BroadcastOrderState(orderState orderstate.AllOrders, wconn *net.UDPConn) {
	message, _ := json.Marshal(orderState)
	broadcastMessage(message, wconn)
}

func broadcastMessage(message []byte, wconn *net.UDPConn) {
	_, err := wconn.Write(message)
	if err != nil {
		panic(err)
	}
}

func Send() {
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	defer wconn.Close()

	for {
		state := orderstate.GetOrders()

		// fmt.Println("state send:", state, "\n\n")
		BroadcastOrderState(state, wconn)
		time.Sleep(_sendRate)
	}
}
