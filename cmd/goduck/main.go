package main

import (
	"log"

	"github.com/shishir1290/goduck/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}