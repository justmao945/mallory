package main

import (
	"fmt"
	"golang.org/x/net/publicsuffix"
	"os"
)

func main() {
	for _, host := range os.Args[1:] {
		fmt.Printf("Host: %s\n", host)
		tld, _ := publicsuffix.EffectiveTLDPlusOne(host)
		fmt.Printf("\tEffectiveTLDPlusOne: %s\n", tld)
		suffix, _ := publicsuffix.PublicSuffix(host)
		fmt.Printf("\tPublicSuffix: %s\n", suffix)
	}
}
