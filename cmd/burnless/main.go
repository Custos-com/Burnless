package main

import (
	"fmt"
	"os"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println("burnless", version)
		return
	}
	fmt.Println("burnless — SRE config as code")
	fmt.Println("Run 'burnless --help' for usage")
}
