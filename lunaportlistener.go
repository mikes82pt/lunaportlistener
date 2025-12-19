package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

const VERSION = "Luna Port Listener v3.2 (Go Rewrite)"

func printHelp() {
	fmt.Println(VERSION)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  --port <number>     Port number to listen on")
	fmt.Println("  --ipv6-only         Listen on IPv6 only")
	fmt.Println("  --ipv4-only         Listen on IPv4 only")
	fmt.Println("  --help              Show this help message")
	fmt.Println()
	fmt.Println("Default behavior:")
	fmt.Println("  IPv6-first dual-stack listener (Windows)")
	fmt.Println("CTRL + C to exit")
}

func handleTCP(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("[TCP] Connection from %s", addr)

	reader := bufio.NewReader(conn)

	for {
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[TCP] Connection closed: %s", addr)
			return
		}

		log.Printf("[TCP] Data from %s:\n    %s", addr, strings.TrimSpace(data))
		_, _ = conn.Write([]byte("Data received\n"))
	}
}

func tcpListener(network, address string) {
	l, err := net.Listen(network, address)
	if err != nil {
		log.Fatalf("[TCP] Listen failed (%s %s): %v", network, address, err)
	}

	log.Printf("[TCP] Listening on %s %s", network, address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Printf("[TCP] Accept error: %v", err)
			continue
		}
		go handleTCP(conn)
	}
}

func udpListener(network, address string) {
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		log.Fatalf("[UDP] Resolve failed (%s %s): %v", network, address, err)
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		log.Fatalf("[UDP] Listen failed (%s %s): %v", network, address, err)
	}

	log.Printf("[UDP] Listening on %s %s", network, address)

	buffer := make([]byte, 1024)

	for {
		n, remote, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("[UDP] Read error: %v", err)
			continue
		}

		msg := string(buffer[:n])
		log.Printf("[UDP] Datagram from %s:\n    %s", remote, strings.TrimSpace(msg))
		_, _ = conn.WriteToUDP([]byte("Data received"), remote)
	}
}

func main() {
	var port string
	var ipv6Only bool
	var ipv4Only bool
	var showHelp bool

	flag.StringVar(&port, "port", "", "Port number to listen on")
	flag.BoolVar(&ipv6Only, "ipv6-only", false, "IPv6 only")
	flag.BoolVar(&ipv4Only, "ipv4-only", false, "IPv4 only")
	flag.BoolVar(&showHelp, "help", false, "Show help")
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	if ipv6Only && ipv4Only {
		log.Fatal("Cannot use --ipv6-only and --ipv4-only together")
	}

	if port == "" {
		fmt.Println("===================================")
		fmt.Println(" ", VERSION)
		fmt.Println("===================================\n")
		fmt.Print("Enter port to listen on (1-65535): ")
		fmt.Scanln(&port)
	}

	var tcpNet, udpNet, address string

	switch {
	case ipv4Only:
		tcpNet = "tcp4"
		udpNet = "udp4"
		address = "0.0.0.0:" + port
		log.Println("Mode: IPv4-only")

	case ipv6Only:
		tcpNet = "tcp6"
		udpNet = "udp6"
		address = "[::]:" + port
		log.Println("Mode: IPv6-only")

	default:
		tcpNet = "tcp"
		udpNet = "udp"
		address = "[::]:" + port
		log.Println("Mode: IPv6 dual-stack)")
	}

	go tcpListener(tcpNet, address)
	go udpListener(udpNet, address)

	select {}
}

