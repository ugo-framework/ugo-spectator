// Copyright 2020 The Go Authors and UGO Authors
/*
Package ugoSpectator implements a simple library to watch files on the
directory specified.

Methods:
		Init() (*UgoSpectator, error)
		Close() error

*/
package spectator

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
)

// UgoSpectator struct with Watcher  and all methods
type UgoSpectator struct {
	Watcher *fsnotify.Watcher // *fsnotify watcher instance
	Cb      string            // Function to restart after watching
	osV     string
}

// Init initializes the fsnotify NewWatcher and
//  a *fsnotify watcher instance and an error
func Init(dirname string) (*UgoSpectator, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return &UgoSpectator{}, err
	}
	ugoWatcher := &UgoSpectator{Watcher: watcher, osV: runtime.GOOS}
	clear(ugoWatcher.osV)
	fmt.Printf("\033[1;36m%s\033[0m", "Ugo Spectator is watching your files")
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event.Op)
				fileSplit := strings.SplitN(event.Name, "/", -1)
				if event.Op == fsnotify.Create {
					fmt.Printf("modified file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
				}
				if event.Op == fsnotify.Write {
					fmt.Printf("modified file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
				}
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Printf("Removed file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
				}
				if event.Op&fsnotify.Rename == fsnotify.Rename {
					fmt.Printf("Removed file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
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
	fmt.Printf("\033[1;33m%s%s\n\033[0m", "\nat ", pathToWatch)

	err = watcher.Add(pathToWatch)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Done: ", <-done)
	return ugoWatcher, nil
}

func (u *UgoSpectator) Close() error {
	fmt.Printf("\033[1;31m%s\033[0m", "Ugo Spector Closing")
	return u.Watcher.Close()
}

func clear(osV string) {
	if osV == "linux" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

	}
	if osV == "darwin" {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

	}
	if osV == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
