package mallory

import (
	"log"
	"os"
)

// global logger
var L = log.NewLogger(os.Stdout, "mallory: ", log.Lshortfile|log.LstdFlags)
