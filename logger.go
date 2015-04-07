package mallory

import (
	"log"
	"os"
)

var logger = log.NewLogger(os.Stdout, "mallory: ", log.Lshortfile|log.LstdFlags)
