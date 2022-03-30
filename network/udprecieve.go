package network

import (
	"elevators/orders"
	"encoding/json"
	"net"
	"time"
)

const _pollRate = _sendRate / 3

const bufferSize = 2 * 2048
const listenAddr = "224.0.0.251"

func InitUDPReceivingSocket(port int) (
	net.UDPAddr,
	*net.UDPConn) {

	addr := net.UDPAddr{
		IP: net.ParseIP(
			listenAddr),
		Port: port,
	}

	conn, err := net.ListenUDP(
		"udp",
		&addr)
	if err != nil {
		panic(err)
	}

	return addr,
		conn
}

func ReceiveOrders(conn *net.UDPConn) orders.AllOrders {
	var allOrders orders.AllOrders
	buf := receiveUDPMessage(conn)
	json.Unmarshal(
		buf,
		&allOrders)
	return allOrders
}

func receiveUDPMessage(conn *net.UDPConn) []byte {
	var buf [bufferSize]byte
	rlen,
		_, err := conn.ReadFromUDP(buf[:])
	if err != nil {
		panic(err)
	}
	return buf[:rlen]
}

func PollReceiveOrders(receiver chan<- orders.AllOrders) {
	_, conn := InitUDPReceivingSocket(UDPPort)
	defer conn.Close()

	for {
		time.Sleep(_pollRate)
		allOrders := ReceiveOrders(conn)
		receiver <- allOrders
	}
}
