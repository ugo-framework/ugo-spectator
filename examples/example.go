package examples

import (
	ugo "github.com/ugo-framework/ugo-spectator"
	"log"
)

func example() {
	watcher, err := ugo.Init("..")
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
