package main

import (
	"fmt"
	"log"

	spectator "github.com/ugo-framework/ugo-spectator/lib"
)

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")
	//ch := make(chan string)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer watcher.Close() // handle error
	for {
		select {
		case res := <-watcher.Ch:
			fmt.Println("RES", res)
		}

	}
	// to hold the program from exiting
	//<-ch
}
