package main 

import (
	"log"
	"net"
	"fmt"
)

func SocketClient(m []byte) {
	conn, err := net.Dial("tcp", "your_tcp_endpoint:your_port")

	defer conn.Close()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("sending logs to syslog-ng:", string(m))
	conn.Write(m)

	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])

}


func main(){
	b := []byte("i'am testing")
	SocketClient(b)
}