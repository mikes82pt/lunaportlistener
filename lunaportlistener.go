package main

import (
    "bufio"
    "flag"
    "fmt"
    "log"
    "net"
    "os/exec"
    "runtime"
    "strings"
)

const VERSION = "Luna Port Listener v3.1 (Go Rewrite)"

func setWindowsTitle(title string) {
    if runtime.GOOS == "windows" {
        cmd := exec.Command("cmd", "/c", "title", title)
        _ = cmd.Run()
    }
}

func printHelp() {
	fmt.Println(VERSION)
	fmt.Println("Usage:")
	fmt.Println("  --port <number>     Port number to listen on")
	fmt.Println("  --bind <IP>         Bind to a specific IPv4 or IPv6 address")
	fmt.Println("  --help              Show this help message")
	fmt.Println()
	fmt.Println("If no parameters are provided, the listener waits for user input and binds to all IPv4 and IPv6 addresses.")
	fmt.Println("CTRL + C to close")
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

func tcpListener(bindAddr, port string) {
    networks := []string{"tcp4", "tcp6"}

    for _, netw := range networks {
        go func(network string) {
            address := net.JoinHostPort(bindAddr, port)
            l, err := net.Listen(network, address)
            if err != nil {
                log.Printf("[TCP] Failed to listen on %s (%s): %v", network, address, err)
                return
            }
            log.Printf("[TCP] Listening on %s, %s...", network, address)

            for {
                conn, err := l.Accept()
                if err != nil {
                    log.Printf("[TCP] Accept error: %v", err)
                    continue
                }
                go handleTCP(conn)
            }
        }(netw)
    }
}

func udpListener(bindAddr, port string) {
    networks := []string{"udp4", "udp6"}

    for _, netw := range networks {
        go func(network string) {
            address := net.JoinHostPort(bindAddr, port)
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

            log.Printf("[UDP] Listening on %s, %s...", network, address)
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
        }(netw)
    }
}

func main() {
    setWindowsTitle(VERSION)

    var port string
    var bindAddr string
    var showHelp bool

    flag.StringVar(&port, "port", "", "Port number to listen on")
    flag.StringVar(&bindAddr, "bind", "", "Bind to a specific IPv4/IPv6 address")
    flag.BoolVar(&showHelp, "help", false, "Show help message")
    flag.Parse()

    if showHelp {
        printHelp()
        return
    }

    if port == "" {
        fmt.Println("===================================")
        fmt.Println(" ", VERSION)
        fmt.Println("===================================\n")

        fmt.Print("Enter port to listen on (1-65535): ")
        fmt.Scanln(&port)

        bindAddr = "" // all addresses
    }

    if bindAddr == "" {
        bindAddr = "" // all IPv4 + IPv6
    }

    log.Printf("Starting listeners on port %s (bind: '%s')...", port, bindAddr)
    tcpListener(bindAddr, port)
    udpListener(bindAddr, port)

    select {} // keep running
}
