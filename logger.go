package mallory

import (
	"log"
	"os"
)

// global logger
var L = log.New(os.Stdout, "mallory: ", log.Lshortfile|log.LstdFlags)
