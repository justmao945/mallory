package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var (
	L = log.New(os.Stdout, "forward: ", log.Lshortfile|log.LstdFlags)

	FNetwork = flag.String("network", "tcp", "network protocol")
	FListen  = flag.String("listen", ":20022", "listen on this port")
	FForward = flag.String("forward", ":80", "destination address and port")
)

func isChanClose(ch chan int) bool {
	select {
	case _, received := <-ch:
		return !received
	default:
	}
	return false
}
func main() {
	flag.Parse()

	L.Printf("Listening on %s for %s...\n", *FListen, *FNetwork)
	ln, err := net.Listen(*FNetwork, *FListen)
	if err != nil {
		L.Fatal(err)
	}

	for id := 0; ; id++ {
		conn, err := ln.Accept()
		if err != nil {
			L.Printf("%d: %s\n", id, err)
			continue
		}
		L.Printf("%d: new %s\n", id, conn.RemoteAddr())

		if tcpConn := conn.(*net.TCPConn); tcpConn != nil {
			L.Printf("%d: setup keepalive for TCP connection\n", id)
			tcpConn.SetKeepAlive(true)
			tcpConn.SetKeepAlivePeriod(30 * time.Second)
		}

		go func(myid int, conn net.Conn) {
			defer conn.Close()
			c, err := net.Dial(*FNetwork, *FForward)
			if err != nil {
				L.Printf("%d: %s\n", myid, err)
				return
			}
			L.Printf("%d: new %s <-> %s\n", myid, c.RemoteAddr(), conn.RemoteAddr())
			defer c.Close()
			wait := make(chan int)
			go func() {
				n, err := io.Copy(c, conn)
				if err != nil {
					L.Printf("%d: %s\n", myid, err)
				}
				L.Printf("%d: %s -> %s %d bytes\n", myid, conn.RemoteAddr(), c.RemoteAddr(), n)
				if !isChanClose(wait) {
					close(wait)
				}
			}()
			go func() {
				n, err := io.Copy(conn, c)
				if err != nil {
					L.Printf("%d: %s\n", myid, err)
				}
				L.Printf("%d: %s -> %s %d bytes\n", myid, c.RemoteAddr(), conn.RemoteAddr(), n)
				if !isChanClose(wait) {
					close(wait)
				}

			}()
			<-wait
			L.Printf("%d: connection closed\n", myid)
		}(id, conn)
	}
}
