package network

import (
	//"elevators/controlunit"
	"elevators/controlunit/cabstate"
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

	wconn, err := net.DialUDP("udp", nil, &sendaddr) // code does not block here

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
	// fmt.Println("You sent: msg: ", message)
}

func TestSendAndReceive() {

	// var state filesystem.OrderState
	// state.Dir = "up"
	// state.Floor = 3
	// state.Name = "Elevator"
	var state = orderstate.GetOrders()
	jsState, _ := json.Marshal(state)
	json.Unmarshal(jsState, &state)
	fmt.Println(string(jsState))
	// fmt.Println(string(state))

	//Initialize sockets
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	_, conn := InitUDPReceivingSocket(UDPPort)

	//Close sockets when program terminates
	defer conn.Close()
	// defer wconn.Close()

	//Send and receive message
	for {
		// BroadcastMessage(json.RawMessage(`{"precomputed": true}`), wconn)
		//BroadcastMessage(string(jsState), wconn)
		state := orderstate.GetOrders()
		BroadcastOrderState(state, wconn)
		time.Sleep(time.Millisecond * 1000)

		// state = ReceiveOrderState(conn)
		msg := ReceiveUDPMessage(conn)
		json.Unmarshal(msg, &state)

		// msg := ReceiveCopy(conn)
		// print msg
		fmt.Println(string(msg))
		s := string(msg)
		fmt.Println(s)
		json.Unmarshal(msg, &state)

		cabstate := cabstate.Cab

		fmt.Println("Your state:", cabstate)

		time.Sleep(time.Millisecond * 1000)
	}
}
