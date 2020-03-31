package main

import (
	"log"

	spectator "github.com/ugo-framework/ugo-spectator/lib"
)

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")
	ch := make(chan string)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	// to hold the program from exiting
	<-ch
	defer watcher.Close() // handle error
}
