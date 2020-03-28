// Copyright 2020 The Go Authors and UGO Authors
/*
Package ugoSpectator implements a simple library to watch files on the
directory specified.

Methods:
		Init() (*UGO, error)
		Close() error

*/
package ugoSpectator

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
)

// UgoSpectator struct with Watcher  and all methods
type UgoSpectator struct {
	Watcher *fsnotify.Watcher
}

// Init initializes the fsnotify NewWatcher and
//  a *fsnotify watcher instance and an error
func Init(dirname string) (*UgoSpectator, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return &UgoSpectator{}, err
	}
	fmt.Printf("\033[1;36m%s\033[0m", "Ugo Spectator is watching your files")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	cPath, err := os.Getwd()
	if err != nil {
		return &UgoSpectator{}, err
	}
	pathToWatch := path.Join(cPath, dirname)
	fmt.Printf("\033[1;36m%s%s\033[0m", "\n at ", pathToWatch)

	err = watcher.Add(pathToWatch)
	if err != nil {
		fmt.Println(err)
	}
	<-done
	return &UgoSpectator{Watcher: watcher}, nil
}

func (u *UgoSpectator) Close() error {
	fmt.Printf("\033[1;31m%s\033[0m", "Ugo Spector Closing")
	return u.Watcher.Close()
}
