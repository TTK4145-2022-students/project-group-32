package network

import (
	//"fmt"
	"elevators/filesystem"
	"encoding/json"
	"net"
)



func InitUDPReceivingSocket(port int) (net.UDPAddr, *net.UDPConn){
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

func ReceiveOrderState(conn *net.UDPConn) (filesystem.OrderState) {
	var orderState filesystem.OrderState
	buf := ReceiveUDPMessage(conn)
	json.Unmarshal(buf, &orderState) // convert from json/[]byte to struct/OrderState
	return orderState
}

func ReceiveUDPMessage(conn *net.UDPConn) ([]byte) {
	var buf []byte
	_, _, err := conn.ReadFromUDP(buf[:])

	if err != nil {
		panic(err)
	} 
	return buf
}
