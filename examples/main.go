package main

import (
	spectator "github.com/ugo-framework/ugo-spectator/lib"
	"log"
)

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")
	ch := make(chan string)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	<-ch
	defer watcher.Close() // handle error
}
