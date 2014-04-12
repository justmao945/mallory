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

	log.Printf("Listen and serve on %s\n", env.Addr)
	log.Printf("\tEngine: %s\n", env.Engine)
	if env.Engine == "gae" {
		log.Printf("\tAppSpot: %s\n", env.AppSpot)
	}

	srv := mallory.NewServer(&env)
	if err := srv.Init(); err != nil {
		log.Fatal(err)
	}
	log.Fatal(http.ListenAndServe(env.Addr, srv))
}
