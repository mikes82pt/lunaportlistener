package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

const VERSION = "Luna Port Listener v3.0 (Go Rewrite)"

func setWindowsTitle(title string) {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "title", title)
		_ = cmd.Run()
	}
}

func handleTCP(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("[TCP] New connection from %s", addr)

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[TCP] Connection closed: %s", addr)
			return
		}

		log.Printf("[TCP] Data from %s:\n    %s", addr, strings.TrimSpace(data))
		conn.Write([]byte("Data received\n"))
	}
}

func tcpListener(port string) {
	listeners := []string{"tcp4", "tcp6"}

	for _, network := range listeners {
		go func(netw string) {
			l, err := net.Listen(netw, ":"+port)
			if err != nil {
				log.Printf("[TCP] Failed to listen on %s: %v", netw, err)
				return
			}
			log.Printf("[TCP] Listening on %s, port %s...", netw, port)

			for {
				conn, err := l.Accept()
				if err != nil {
					log.Printf("[TCP] Accept error: %v", err)
					continue
				}
				go handleTCP(conn)
			}
		}(network)
	}
}

func udpListener(port string) {
	listeners := []string{"udp4", "udp6"}

	for _, network := range listeners {
		go func(netw string) {
			addr, err := net.ResolveUDPAddr(netw, ":"+port)
			if err != nil {
				log.Printf("[UDP] Failed to resolve %s: %v", netw, err)
				return
			}

			conn, err := net.ListenUDP(netw, addr)
			if err != nil {
				log.Printf("[UDP] Failed to listen on %s: %v", netw, err)
				return
			}

			log.Printf("[UDP] Listening on %s, port %s...", netw, port)
			buffer := make([]byte, 1024)

			for {
				n, remote, err := conn.ReadFromUDP(buffer)
				if err != nil {
					log.Printf("[UDP] Error: %v", err)
					continue
				}

				msg := string(buffer[:n])
				log.Printf("[UDP] Datagram from %s:\n    %s", remote, strings.TrimSpace(msg))

				conn.WriteToUDP([]byte("Data received"), remote)
			}
		}(network)
	}
}

func main() {
	setWindowsTitle(VERSION)

	fmt.Println("===================================")
	fmt.Println(" ", VERSION)
	fmt.Println("===================================\n")

	fmt.Print("Enter port to listen on (1-65535): ")
	var port string
	fmt.Scanln(&port)

	log.Printf("Starting TCP and UDP listeners on port %s...", port)

	tcpListener(port)
	udpListener(port)

	select {} // keep running forever
}

