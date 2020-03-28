package main

import (
	ugoSpectator "github.com/ugo-framework/ugo-spectator/lib"
	"log"
)

func main() {
	watcher, err := ugoSpectator.Init("")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = watcher.Close()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
