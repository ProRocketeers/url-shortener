package main

import (
	"fmt"
	"os"

	"github.com/ProRocketeers/url-shortener/infrastructure"
)

func main() {
	if err := infrastructure.RunServerGracefully(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
