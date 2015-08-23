package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
)

var (
	L = log.New(os.Stdout, "forward: ", log.Lshortfile|log.LstdFlags)

	FNetwork = flag.String("network", "tcp", "network protocol")
	FListen  = flag.String("listen", ":20022", "listen on this port")
	FForward = flag.String("forward", ":80", "destination address and port")
)

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

		go func(myid int) {
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
				wait <- 1
			}()
			go func() {
				n, err := io.Copy(conn, c)
				if err != nil {
					L.Printf("%d: %s\n", myid, err)
				}
				L.Printf("%d: %s -> %s %d bytes\n", myid, c.RemoteAddr(), conn.RemoteAddr(), n)
				wait <- 1
			}()
			<-wait
			<-wait
		}(id)
	}
}
