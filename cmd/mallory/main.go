package main

import (
	"flag"
	"net/http"

	m "github.com/justmao945/mallory"
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

	L.Printf("Connecting remote SSH server: %s\n", c.File.RemoteServer)
	smart, err := m.NewServer(m.SmartSrv, c)
	if err != nil {
		L.Fatalln(err)
	}
	normal, err := m.NewServer(m.NormalSrv, c)
	if err != nil {
		L.Fatalln(err)
	}

	go func() {
		L.Printf("Local normal HTTP proxy: %s\n", c.File.LocalNormalServer)
		L.Fatalln(http.ListenAndServe(c.File.LocalNormalServer, normal))
	}()

	L.Printf("Local smart HTTP proxy: %s\n", c.File.LocalSmartServer)
	L.Fatalln(http.ListenAndServe(c.File.LocalSmartServer, smart))
}
