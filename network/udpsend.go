package network

import (
	"elevators/controlunit/orderstate"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const UDPPort = 20014
const broadcastAddr = "255.255.255.255"

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

func TestSend() {
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	defer wconn.Close()

	for {
		var state orderstate.AllOrders
		state.Up[0] = orderstate.OrderState{time.Now(), time.Now().Add(-5 * time.Second), time.Now().Add(-5 * time.Second)}
		state.Cab[1] = true

		fmt.Println("state send:", state, "\n\n")
		BroadcastOrderState(state, wconn)
		time.Sleep(time.Second)
	}
}