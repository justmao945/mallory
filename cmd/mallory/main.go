package main

import (
	"flag"
	m "github.com/justmao945/mallory"
	"net/http"
)

var L = m.L

func main() {
	f := flag.String("config", "$HOME/.config/mallory.json", "config file")
	flag.Parse()

	L.Printf("Starting...\n")
	c, err := m.NewConfig(*f)
	if err != nil {
		L.Fatalln(err)
	}
	L.Println("Connecting...")
	srv, err := m.NewServer(c)
	if err != nil {
		L.Fatalln(err)
	}

	L.Printf("Listen and serve HTTP proxy on %s\n", c.File.LocalServer)
	L.Printf("\tRemote SSH server: %s\n", c.File.RemoteServer)
	L.Fatalln(http.ListenAndServe(c.File.LocalServer, srv))
}
