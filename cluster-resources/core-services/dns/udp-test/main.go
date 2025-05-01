package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func server() {
	addr, _ := net.ResolveUDPAddr("udp", ":1234")
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on port 1234: %s", err)
	}
	defer conn.Close()

	for {
		b := make([]byte, 2048)
		n, paddr, err := conn.ReadFromUDP(b)
		if err != nil {
			log.Fatalf("Error reading from connection: %s", err)
		}
		log.Printf("Received message: %s", b[:n])

		_, err = conn.WriteTo(append([]byte("1: "), b[:n]...), paddr)
		if err != nil {
			log.Fatalf("Failed to write to client: %s", err)
		}
		conn.WriteTo(append([]byte("2: "), b[:n]...), paddr)
		if err != nil {
			log.Fatalf("Failed to write to client(2): %s", err)
		}
		conn.WriteTo(append([]byte("3: "), b[:n]...), paddr)
		if err != nil {
			log.Fatalf("Failed to write to client(3): %s", err)
		}
	}
}

func client(endpoint, message string) {
	addr, _ := net.ResolveUDPAddr("udp", endpoint)
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatalf("Failed to dial endpoint: %s", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Fatalf("Failed to write message: %s", err)
	}

	deadline := time.Now().Add(time.Second * 15)
	ctx, _ := context.WithDeadline(context.Background(), deadline)
	for {
		select {
		case <-ctx.Done():
			log.Println("Timeout after 15 seconds.")
			return
		default:
			b := make([]byte, 2048)
			if err = conn.SetDeadline(deadline); err != nil {
				log.Fatalf("Error setting deadline: %s", err)
			}
			n, err := conn.Read(b)
			if err != nil {
				if err == io.EOF {
					log.Println("Got EOF")
					break
				}
				log.Printf("Error: %s", err)
				break
			}
			log.Printf("Msg: %s", b[:n])

			if string(b[:n]) == "3: "+message {
				return
			}
		}
	}

}

func main() {
	switch os.Args[1] {
	case "server":
		server()
	case "client":
		client(os.Args[2], os.Args[3])
	}

}
