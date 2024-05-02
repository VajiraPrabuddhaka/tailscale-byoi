package main

import (
	"io"
	"log"
	"net"

	"golang.org/x/net/proxy"
)

const (
	listenAddr      = "0.0.0.0:8080"
	proxyAddr       = "localhost:1055"
	destinationAddr = "100.95.126.2:8090" // Replace with actual IP and port
)

func handleConnection(conn net.Conn, dialer proxy.Dialer) {
	defer conn.Close()

	// Dial the destination address through the SOCKS5 proxy
	proxyConn, err := dialer.Dial("tcp", destinationAddr)
	if err != nil {
		log.Printf("Failed to connect to destination through proxy: %v", err)
		return
	}
	defer proxyConn.Close()

	// Use io.Copy to forward the traffic
	go func() {
		if _, err := io.Copy(proxyConn, conn); err != nil {
			log.Printf("Failed to copy data to proxy connection: %v", err)
		}
	}()
	if _, err := io.Copy(conn, proxyConn); err != nil {
		log.Printf("Failed to copy data to client connection: %v", err)
	}
}

func main() {
	// Create a proxy dialer
	dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
	if err != nil {
		log.Fatalf("Failed to connect to SOCKS5 proxy: %v", err)
	}

	// Listen on the specified port
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}
	defer listener.Close()
	log.Printf("Listening on %s and forwarding to %s via proxy %s", listenAddr, destinationAddr, proxyAddr)

	// Accept connections and handle them
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}
		go handleConnection(conn, dialer)
	}
}
