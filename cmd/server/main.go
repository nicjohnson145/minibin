package main

import (
	"fmt"
	"os"

	"github.com/nicjohnson145/minibin/server"
)

func main() {
	if err := server.Run(); err != nil {
		fmt.Println("error running server: %w", err)
		os.Exit(1)
	}
}
