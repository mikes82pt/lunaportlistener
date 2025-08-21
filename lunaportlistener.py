import socket
import logging
import threading
import os
import platform
from typing import Tuple
from colorama import init, Fore, Style

VERSION = "Luna Port Listener v2.1"

# Initialize colorama
init(autoreset=True)

# Set Windows console title
if platform.system() == "Windows":
    os.system(f"title {VERSION}")

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s\n%(message)s\n",
    datefmt="%Y-%m-%d %H:%M:%S"
)

BUFFER_SIZE = 1024


def handle_tcp_client(client_socket: socket.socket, client_address: Tuple, buffer_size: int) -> None:
    """
    Handle communication with a single TCP client.
    """
    with client_socket:
        logging.info(Fore.CYAN + f"[TCP] New connection from {client_address}")
        try:
            while True:
                data = client_socket.recv(buffer_size)
                if not data:
                    break
                logging.info(
                    Fore.CYAN + f"[TCP] Data from {client_address}:\n    {data.decode(errors='replace')}"
                )
                client_socket.sendall(b"Data received")
        except Exception as e:
            logging.error(Fore.RED + f"[TCP] Error with {client_address}: {e}")
        finally:
            logging.info(Fore.CYAN + f"[TCP] Connection closed: {client_address}")


def tcp_listener(port: int, buffer_size: int = BUFFER_SIZE) -> None:
    """
    Start multi-threaded TCP listeners on both IPv4 and IPv6 (Windows).
    """
    def listen(sock: socket.socket, proto: str):
        with sock:
            sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            sock.listen()
            logging.info(Fore.CYAN + f"[TCP] Listening on {proto} all interfaces, port {port}...")
            while True:
                client_socket, client_address = sock.accept()
                threading.Thread(
                    target=handle_tcp_client,
                    args=(client_socket, client_address, buffer_size),
                    daemon=True
                ).start()

    if platform.system() == "Windows":
        # IPv4 socket
        ipv4_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        ipv4_sock.bind(("0.0.0.0", port))

        # IPv6 socket
        ipv6_sock = socket.socket(socket.AF_INET6, socket.SOCK_STREAM)
        ipv6_sock.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 1)
        ipv6_sock.bind(("::", port))

        threading.Thread(target=listen, args=(ipv4_sock, "IPv4"), daemon=True).start()
        threading.Thread(target=listen, args=(ipv6_sock, "IPv6"), daemon=True).start()

    else:
        # Fallback dual-stack for non-Windows
        with socket.socket(socket.AF_INET6, socket.SOCK_STREAM) as server_socket:
            server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            server_socket.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 0)
            server_socket.bind(('', port))
            server_socket.listen()
            logging.info(Fore.CYAN + f"[TCP] Listening (dual-stack) on port {port}...")
            while True:
                client_socket, client_address = server_socket.accept()
                threading.Thread(
                    target=handle_tcp_client,
                    args=(client_socket, client_address, buffer_size),
                    daemon=True
                ).start()


def udp_listener(port: int, buffer_size: int = BUFFER_SIZE) -> None:
    """
    Start UDP listeners on both IPv4 and IPv6 (Windows).
    """
    def listen(sock: socket.socket, proto: str):
        with sock:
            sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            logging.info(Fore.GREEN + f"[UDP] Listening on {proto} all interfaces, port {port}...")
            while True:
                data, client_address = sock.recvfrom(buffer_size)
                logging.info(
                    Fore.GREEN + f"[UDP] Datagram from {client_address}:\n    {data.decode(errors='replace')}"
                )
                sock.sendto(b"Data received", client_address)

    if platform.system() == "Windows":
        # IPv4 socket
        ipv4_sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        ipv4_sock.bind(("0.0.0.0", port))

        # IPv6 socket
        ipv6_sock = socket.socket(socket.AF_INET6, socket.SOCK_DGRAM)
        ipv6_sock.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 1)
        ipv6_sock.bind(("::", port))

        threading.Thread(target=listen, args=(ipv4_sock, "IPv4"), daemon=True).start()
        threading.Thread(target=listen, args=(ipv6_sock, "IPv6"), daemon=True).start()

    else:
        # Fallback dual-stack for non-Windows
        with socket.socket(socket.AF_INET6, socket.SOCK_DGRAM) as server_socket:
            server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
            server_socket.setsockopt(socket.IPPROTO_IPV6, socket.IPV6_V6ONLY, 0)
            server_socket.bind(('', port))
            logging.info(Fore.GREEN + f"[UDP] Listening (dual-stack) on port {port}...")
            while True:
                data, client_address = server_socket.recvfrom(buffer_size)
                logging.info(
                    Fore.GREEN + f"[UDP] Datagram from {client_address}:\n    {data.decode(errors='replace')}"
                )
                server_socket.sendto(b"Data received", client_address)


if __name__ == "__main__":
    print("\n" + Fore.YELLOW + "=" * 50)
    print(Fore.YELLOW + f"  {VERSION}")
    print(Fore.YELLOW + "=" * 50 + "\n")

    try:
        port = int(input("Enter the port number to listen on (1-65535): "))
        if not (1 <= port <= 65535):
            raise ValueError("Port number must be between 1 and 65535.")

        # Start TCP and UDP listeners in separate threads
        threading.Thread(target=tcp_listener, args=(port,), daemon=True).start()
        threading.Thread(target=udp_listener, args=(port,), daemon=True).start()

        logging.info(Fore.YELLOW + f"TCP and UDP listeners started on port {port}.\nPress Ctrl+C to stop.")

        # Keep main thread alive
        while True:
            threading.Event().wait(1)

    except ValueError as ve:
        logging.error(Fore.RED + str(ve))
    except KeyboardInterrupt:
        logging.info(Fore.YELLOW + "Shutting down listeners...")
