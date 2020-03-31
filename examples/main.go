package main

import (
	"fmt"
	"log"

	spectator "github.com/ugo-framework/ugo-spectator/lib"
)

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer watcher.Close() // handle error
	// event to catch for file change
	for {
		select {
		case res := <-watcher.Ch:
			fmt.Println("RES", res)
		}

	}
}
