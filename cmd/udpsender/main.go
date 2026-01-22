package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)


const addr = "localhost:42069"

func main() {
	updAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalln("error while resolving an address:", err)
	}
	conn, err := net.DialUDP(updAddr.Network(), nil, updAddr)
	if err != nil {
		log.Fatalln("error while dialing:", err)
	}
	defer conn.Close()
	fmt.Println("listening", addr)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		_, err = conn.Write([]byte (line))
		if err != nil {
			log.Println(err)
		}

		fmt.Print("Message sent: ", line)
	}
}