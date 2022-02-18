package network

import (
	"net"
)

const broadcastAddr = "255.255.255.255"
const port = 20014

func InitUDPSendingSocket(port int, sendAddr string) (net.UDPAddr, *net.UDPConn) {
	sendaddr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(broadcastAddr),
	}
	wconn, err := net.DialUDP("udp", nil, &sendaddr) // code does not block here
	if err != nil {
		panic(err)
	}
	//defer wconn.Close() //Close at the end of program

	return sendaddr, wconn
}

func BroadcastMessage(message string, wconn *net.UDPConn) {
	sendMessage := []byte(message)
	// var buf [1024]byte
	_, err := wconn.Write(sendMessage)
	if err != nil {
		panic(err)
	}
}
