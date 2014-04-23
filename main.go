package main

import (
	"github.com/justmao945/mallory/mallory"
	"log"
	"net/http"
)

func main() {
	var env mallory.Env
	if err := env.Parse(); err != nil {
		log.Fatal(err)
	}

	hint := env.Fallback()
	for _, v := range hint {
		log.Printf("Warning: %s\n", v)
	}
	if len(hint) != 0 {
		log.Println("Fallback to default")
	}

	log.Printf("Starting...\n")
	srv, err := mallory.CreateServer(&env)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Listen and serve on %s\n", env.Addr)
	log.Printf("\tEngine: %s\n", env.Engine)
	if env.Engine == "gae" {
		log.Printf("\tAppSpot: %s\n", env.AppSpot)
	}

	if env.PAC != "" {
		log.Printf("\tService: PAC file at http://%s/pac\n", env.Addr)
	}
	log.Fatal(http.ListenAndServe(env.Addr, srv))
}
