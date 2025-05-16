import socket
import threading
import time

def start_tcp_listener(port):
    # Create an IPv4 TCP socket
    ipv4_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    ipv4_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    ipv4_socket.bind(('', port))
    ipv4_socket.listen()

    # Create an IPv6 TCP socket
    ipv6_socket = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
    ipv6_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    ipv6_socket.bind(('', port))
    ipv6_socket.listen()

    print(f"Listening for TCP connections on port {port} (IPv4 and IPv6)...")

    def accept_connections(sock):
        while True:
            client_socket, client_address = sock.accept()
            with client_socket:
                print(f"TCP Connection from {client_address}")
                data = client_socket.recv(1024)
                if data:
                    print(f"Received data: {data.decode()}")
                    client_socket.sendall(b"Data received")

    # Start a thread for IPv4 and IPv6
    threading.Thread(target=accept_connections, args=(ipv4_socket,), daemon=True).start()
    threading.Thread(target=accept_connections, args=(ipv6_socket,), daemon=True).start()

def start_udp_listener(port):
    # Create an IPv4 UDP socket
    ipv4_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    ipv4_socket.bind(('', port))

    # Create an IPv6 UDP socket
    ipv6_socket = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
    ipv6_socket.bind(('', port))

    print(f"Listening for UDP connections on port {port} (IPv4 and IPv6)...")

    def receive_data(sock):
        while True:
            data, client_address = sock.recvfrom(1024)
            print(f"UDP Connection from {client_address}")
            print(f"Received data: {data.decode()}")
            sock.sendto(b"Data received", client_address)

    # Start a thread for IPv4 and IPv6
    threading.Thread(target=receive_data, args=(ipv4_socket,), daemon=True).start()
    threading.Thread(target=receive_data, args=(ipv6_socket,), daemon=True).start()

if __name__ == "__main__":
    print("Luna Port Listener 1.0")
    try:
        port = int(input("Enter the port number to listen on (1-65535): "))
        if 1 <= port <= 65535:
            protocol = input("Choose protocol (tcp/udp): ").strip().lower()
            if protocol == 'tcp':
                start_tcp_listener(port)
            elif protocol == 'udp':
                start_udp_listener(port)
            else:
                print("Invalid protocol. Please choose 'tcp' or 'udp'.")
        else:
            print("Port number must be between 1 and 65535.")
        
        # Keep the main thread alive
        while True:
            time.sleep(1)  # Sleep to prevent busy waiting
    except ValueError:
        print("Invalid input. Please enter a valid port number.")
    except KeyboardInterrupt:
        print("Shutting down the listener.")
