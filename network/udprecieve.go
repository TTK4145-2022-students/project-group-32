package network

import (
	"elevators/controlunit/orderstate"
	"encoding/json"
	"net"
)

const bufferSize = 2048

func InitUDPReceivingSocket(port int) (net.UDPAddr, *net.UDPConn) {
	addr := net.UDPAddr{
		Port: port,
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	return addr, conn
}

func ReceiveOrderState(conn *net.UDPConn) orderstate.AllOrders {
	var allOrders orderstate.AllOrders
	buf := receiveUDPMessage(conn)
	json.Unmarshal(buf, &allOrders)
	return allOrders
}

func receiveUDPMessage(conn *net.UDPConn) []byte {
	var buf [bufferSize]byte
	rlen, _, err := conn.ReadFromUDP(buf[:])

	if err != nil {
		panic(err)
	}
	return buf[:rlen]
}

func Receive(receiver chan<- orderstate.AllOrders) {
	_, conn := InitUDPReceivingSocket(UDPPort)
	defer conn.Close()

	for {
		state := ReceiveOrderState(conn)
		// // fmt.Println("state recieve:", state, "\n\n", time.Now())
		receiver <- state
	}
}
