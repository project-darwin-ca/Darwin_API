package main

import (
	"github.com/xuxant/Darwin_API/cmd/data"
	"log"
)

func main() {
	cmd := data.NewRootCommand()
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
