package network

import (
	"fmt"
	"net"
)

//const broadcastAddr = "255.255.255.255"
//const port = 20014 //TODO: Set to appropiate port number

func InitUDPReceivingSocket(port int) (net.UDPAddr, *net.UDPConn){
	addr := net.UDPAddr{
		Port: port,
		//IP:   net.ParseIP("10.100.23.240:39205"),
	}

	conn, err := net.ListenUDP("udp", &addr) // code does not block here
	if err != nil {
		panic(err)
	}
	// defer conn.Close()

	return addr, conn
}

func ReceiveUDPMessage(conn *net.UDPConn) ([1024]byte){  //Perhaps unneccesary function. Enough with readfromudp?
	var buf [1024]byte
	rlen, addr, err := conn.ReadFromUDP(buf[:])
	fmt.Println(addr, "sent:", string(buf[0:rlen]))
	if err != nil {
		panic(err)
	}
	return buf
}

