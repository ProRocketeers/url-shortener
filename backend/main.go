package main

import (
	"fmt"
	"os"

	_ "github.com/ProRocketeers/url-shortener/docs"
	"github.com/ProRocketeers/url-shortener/infrastructure"
)

var Version = "dev"

func main() {
	if err := infrastructure.RunServerGracefully(Version); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
