package network

import (
	//"fmt"
	// "encoding/json"
	"net"
	// "elevators/filesystem"
)



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

func ReceiveUDPMessage(conn *net.UDPConn) ([]byte){  //return in string format
	var buf []byte
	//rlen, addr, err := conn.ReadFromUDP(buf[:])
	_, _, err := conn.ReadFromUDP(buf[:])
	//fmt.Println(addr, "You received:", string(buf[0:rlen]))
	if err != nil {
		panic(err)
	} 
	return buf
}

// func ReceiveUDPMessage(conn *net.UDPConn) filesystem.OrderState {  //return in struct format
// 	var buf []byte
// 	rlen, addr, err := conn.ReadFromUDP(buf[:])
// 	// _, addr, err := conn.ReadFromUDP(buf[:])
// 	fmt.Println(string(buf[:]))
// 	if err != nil {
// 		panic(err)
// 	}
// 	var orderState filesystem.OrderState 
// 	//message, _ := json.Marshal(string(buf[:]))
// 	message := (buf[:rlen])
// 	fmt.Println(addr, "sent:", string(message))
// 	json.Unmarshal(message, &orderState)
// 	fmt.Println(addr, "sent:", string(message))
// 	return orderState
// }
