package main

import (
	"log"

	spectator "github.com/ugo-framework/ugo-spectator/lib"
)

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")

	select {
	case _ = <-watcher.Ch:
		watcher.Close() // Closes the watcher
	}

	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	//defer watcher.Close() // handle error
	// event to catch for file change
}
