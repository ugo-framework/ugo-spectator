// Copyright 2020 The Go Authors and UGO Authors
/*
Package ugo implements a simple library to watch files on the
directory specified.

Methods:
		Init() (*UGO, error)
		Close() error

*/
package ugo

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path"
)

// UGO struct with Watcher  and all methods
type UGO struct {
	Watcher *fsnotify.Watcher
}

// Init initializes the fsnotify NewWatcher and
//  a *fsnotify watcher instance and an error
func Init(dirname string) (*UGO, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return &UGO{}, err
	}
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	cPath, err := os.Getwd()
	if err != nil {
		return &UGO{}, err
	}
	err = watcher.Add(path.Join(cPath, dirname))
	if err != nil {
		log.Fatal(err)
	}
	<-done
	return &UGO{Watcher: watcher}, nil
}

func (u *UGO) Close() error {
	return u.Watcher.Close()
}
