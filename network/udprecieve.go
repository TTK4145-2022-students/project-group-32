package network

import (
	"elevators/controlunit/orderstate"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

const bufferSize = 2048

func InitUDPReceivingSocket(port int) (net.UDPAddr, *net.UDPConn) {
	addr := net.UDPAddr{
		Port: port,
		//IP:   net.ParseIP("10.100.23.240:39205"),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		panic(err)
	}

	return addr, conn
}

func ReceiveOrderState(conn *net.UDPConn) orderstate.AllOrders {
	var allOrders orderstate.AllOrders
	buf := ReceiveUDPMessage(conn)
	json.Unmarshal(buf, &allOrders)
	fmt.Println("msg: ", buf)
	return allOrders
}

func ReceiveUDPMessage(conn *net.UDPConn) []byte {
	var buf [bufferSize]byte
	rlen, _, err := conn.ReadFromUDP(buf[:])

	if err != nil {
		panic(err)
	}
	return buf[:rlen]
}

func TestReceive() {
	_, conn := InitUDPReceivingSocket(UDPPort)
	defer conn.Close()

	for {
		// var state orderstate.OrderStatus
		// msg := ReceiveUDPMessage(conn)
		// json.Unmarshal(msg, &state)
		state := ReceiveOrderState(conn)

		// fmt.Println(time.Now())
		fmt.Println("state:", state, "\n\n", time.Now())
	}
}
