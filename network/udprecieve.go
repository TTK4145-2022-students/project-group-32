package network

import (
	"elevators/hardware"
	"elevators/orders"
	"encoding/json"
	"net"
	"time"
)

const (
	pollRate   = sendRate / 3
	bufferSize = hardware.FloorCount * 512
	listenAddr = "224.0.0.251"
)

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

func receiveOrders(conn *net.UDPConn) orders.AllOrders {
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
	_, conn := InitUDPReceivingSocket(udpPort)
	defer conn.Close()

	for {
		time.Sleep(pollRate)
		allOrders := receiveOrders(conn)
		receiver <- allOrders
	}
}
