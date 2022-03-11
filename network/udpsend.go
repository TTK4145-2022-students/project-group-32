package network

import (
	//"elevators/controlunit"

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
func SendHello(){
	//Initialize sockets
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	// _, conn := InitUDPReceivingSocket(UDPPort)

	// Close sockets when program terminates
	// defer conn.Close()
	//defer wconn.Close()
	
	for{
		message := []byte("hello43")
		// broadcastMessage(message, wconn)
		wconn.Write(message)
		time.Sleep(time.Millisecond * 1000)
		fmt.Println(message)
	}
	
}

func TestSend() {
	_, wconn := InitUDPSendingSocket(UDPPort, broadcastAddr)
	defer wconn.Close()

	// var state orderstate.OrderStatus
	// state.UpAtFloor = true
	// state.DownAtFloor = false
	// state.CabAtFloor = true
	// state.AboveFloor = false
	// state.BelowFloor = false
	for{
		var state orderstate.AllOrders // = orderstate.GetOrders() 
		state
		fmt.Println(state, "\n\n")
		jsState, _ := json.Marshal(state)
	
		broadcastMessage(jsState, wconn)
		time.Sleep(time.Millisecond * 1000)
	}
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

	// Close sockets when program terminates
	defer conn.Close()
	//defer wconn.Close()

	//Send and receive message
	for {
		// BroadcastMessage(json.RawMessage(`{"precomputed": true}`), wconn)
		//BroadcastMessage(string(jsState), wconn)
		state := orderstate.GetOrders()
		BroadcastOrderState(state, wconn)
		time.Sleep(time.Millisecond * 1000)

		//state = ReceiveOrderState(conn)
		msg := ReceiveUDPMessage(conn)
		json.Unmarshal(msg, &state)

		fmt.Println("I am 1,   ", time.Now(), "orders:", msg,"\n\n")

		time.Sleep(time.Millisecond * 1000)
	}
}
