# Luna Port Listener

**Luna Port Listener** is a lightweight TCP and UDP listener designed for network testing, debugging, and packet inspection.
This release is a full Go rewrite of the original Python implementation and now runs natively on Windows with no external dependencies.

---

## **Features**

* Listens on both **TCP** and **UDP**
* Supports **IPv4 and IPv6** simultaneously
* Fully concurrent using Go goroutines
* Windows-native executable
* No installation or runtime required
* Static build — no DLLs, no Python dependency
* Automatically updates the Windows console title

Two Windows builds are available:

* **64-bit:** `lunaportlistener.exe`
* **32-bit:** `lunaportlistener-x86.exe`

---

## **What it does**

Luna Port Listener allows you to:

* Accept incoming TCP connections and display/log the received data
* Receive UDP datagrams and show their contents
* Test local or remote network connectivity
* Troubleshoot routing, firewalls, NAT, and port forwarding
* Verify IPv6 connectivity and behavior
* Simulate a minimal echo-style server

Useful for:

* Network engineers
* Sysadmins
* Pen-testing labs
* Developers working with sockets
* Home lab / port-forwarding tests

---

## **Usage**

Run the executable and enter a port number when prompted, or use command-line options:

```
Usage:
  --port <number>     Port number to listen on
  --bind <IP>         Bind to a specific IPv4 or IPv6 address
  --help              Show this help message
```

Example:

```
lunaportlistener.exe --port 8080 --bind 0.0.0.0
```

---

## **System Requirements**

* **Windows 8.1** (32-bit or 64-bit)
* Compatible with Windows 10 (32-bit or 64-bit) and Windows 11 as well

---
