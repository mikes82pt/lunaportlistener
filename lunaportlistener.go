package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

const VERSION = "Luna Port Listener v3.4 (Go Rewrite)"

func printHelp() {
	fmt.Println(VERSION)
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  --port string        Ports to listen on (single, list, or range)")
	fmt.Println("  --ipv4               Bind to all IPv4 addresses")
	fmt.Println("  --ipv6               Bind to all IPv6 addresses")
	fmt.Println("  --autoclose int      Auto close after N minutes in non-interactive mode")
	fmt.Println("                       Default: 15, 0 = never close")
	fmt.Println("  --help               Show this help message")
	fmt.Println()
	fmt.Println("If no parameters are provided, the listener runs in interactive mode.")
	fmt.Println("CTRL + C to close")
}

func parsePorts(input string) ([]string, error) {
	var ports []string
	parts := strings.Split(input, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			r := strings.Split(part, "-")
			if len(r) != 2 {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			start, err1 := strconv.Atoi(r[0])
			end, err2 := strconv.Atoi(r[1])
			if err1 != nil || err2 != nil || start > end {
				return nil, fmt.Errorf("invalid port range: %s", part)
			}

			for p := start; p <= end; p++ {
				ports = append(ports, strconv.Itoa(p))
			}
		} else {
			p, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid port: %s", part)
			}
			ports = append(ports, strconv.Itoa(p))
		}
	}

	return ports, nil
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

func tcpListener(network, port string) {
	address := net.JoinHostPort("", port)
	l, err := net.Listen(network, address)
	if err != nil {
		log.Printf("[TCP] Failed to listen on %s %s: %v", network, address, err)
		return
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

func udpListener(network, port string) {
	address := net.JoinHostPort("", port)
	addr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		log.Printf("[UDP] Resolve failed on %s: %v", network, err)
		return
	}

	conn, err := net.ListenUDP(network, addr)
	if err != nil {
		log.Printf("[UDP] Listen failed on %s: %v", network, err)
		return
	}

	log.Printf("[UDP] Listening on %s %s", network, address)
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
}

func main() {
	var portInput string
	var ipv4Only bool
	var ipv6Only bool
	var showHelp bool
	var autoClose int

	flag.StringVar(&portInput, "port", "", "Ports to listen on (single, list, or range)")
	flag.BoolVar(&ipv4Only, "ipv4", false, "Bind to all IPv4 addresses")
	flag.BoolVar(&ipv6Only, "ipv6", false, "Bind to all IPv6 addresses")
	flag.IntVar(&autoClose, "autoclose", 15, "Auto close after N minutes (non-interactive mode)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.Parse()

	if showHelp {
		printHelp()
		return
	}

	interactive := portInput == ""

	if interactive {
		fmt.Println("===================================")
		fmt.Println(" ", VERSION)
		fmt.Println("===================================\n")
		fmt.Print("Enter ports: ")
		fmt.Scanln(&portInput)
	}

	ports, err := parsePorts(portInput)
	if err != nil {
		log.Fatal(err)
	}

	// Default: both stacks enabled
	if !ipv4Only && !ipv6Only {
		ipv4Only = true
		ipv6Only = true
	}

	for _, port := range ports {
		if ipv4Only {
			go tcpListener("tcp4", port)
			go udpListener("udp4", port)
		}
		if ipv6Only {
			go tcpListener("tcp6", port)
			go udpListener("udp6", port)
		}
	}

	// Auto-close logic (non-interactive only)
	if !interactive && autoClose > 0 {
		log.Printf("Auto-close enabled: exiting after %d minutes", autoClose)
		time.AfterFunc(time.Duration(autoClose)*time.Minute, func() {
			log.Println("Auto-close timer reached, exiting")
			os.Exit(0)
		})
	}

	select {}
}

