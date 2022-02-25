package network

import (
	"elevators/filesystem"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

// const broadcastAddr = "255.255.255.255"
// const port = 20014

func InitUDPSendingSocket(port int, sendAddr string) (net.UDPAddr, *net.UDPConn) {
	sendaddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(sendAddr),
	}
	wconn, err := net.DialUDP("udp", nil, &sendaddr) // code does not block here
	if err != nil {
		panic(err)
	}
	//defer wconn.Close() //Close at the end of program

	return sendaddr, wconn
}

func BroadcastMessage(message string, wconn *net.UDPConn) {  //Sending string
	sendMessage := []byte(message)
	// var buf [1024]byte
	_, err := wconn.Write(sendMessage)
	if err != nil {
		panic(err)
	}
}

// func BroadcastMessage(message json.RawMessage, wconn *net.UDPConn) {  //Sending json.RawMessage
// 	// sendMessage := []byte(message)
// 	_, err := wconn.Write(message)
// 	if err != nil {
// 		panic(err)
// 	}
// }





func TestSendAndReceive(){
	UDPPort := 20014

	var state  filesystem.OrderState
	state.Dir = "up"
	state.Floor = 3
	state.Name = "Elevator"

	jsState, _ := json.Marshal(state)
	json.Unmarshal(jsState, &state)
	fmt.Println(string(jsState))
	// fmt.Println(string(state))

	//Initialize sockets
	_, wconn := InitUDPSendingSocket(UDPPort, "255.255.255.255")
	_, conn := InitUDPReceivingSocket(UDPPort)

	//Close sockets when program terminates
	defer conn.Close()
	defer wconn.Close()

	//Send and receive message
	for {
		// BroadcastMessage(json.RawMessage(`{"precomputed": true}`), wconn)
		BroadcastMessage(string(jsState), wconn)
		time.Sleep(time.Millisecond * 2000)
		ReceiveUDPMessage(conn)
	}
}

