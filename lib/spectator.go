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
	"context"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

// UgoSpectator struct with Watcher  and all methods
type UgoSpectator struct {
	Watcher       *fsnotify.Watcher  // *fsnotify watcher instance
	dirname       string             // dirname to watch
	fileToRestart string             // Function to restart after watching
	osV           string             // osV ttake the current os
	CancelCtx     context.CancelFunc // context to cancel
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
	//wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	ugoWatcher.CancelCtx = cancel
	ch := make(chan string)

	cPath, err := os.Getwd()
	if err != nil {
		return &UgoSpectator{}, err
	}
	pathToWatch := path.Join(cPath, dirname)
	ugoWatcher.dirname = pathToWatch
	fmt.Printf("\033[1;33m%s%s\n\033[0m", "\nat ", pathToWatch)
	// calling fsNotifiyFunc in a goroutine
	go fsNotifiyFunc(ctx, ch, ugoWatcher.osV, ugoWatcher)

	err = watcher.Add(pathToWatch)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("Done: ", <-done)
	return ugoWatcher, nil
}

// Close return error if an error occurs during the closing of the
// fsnotify watcher instance
func (u *UgoSpectator) Close() error {
	u.CancelCtx()
	fmt.Printf("\033[1;31m%s\033[0m", "Ugo Spector Closing")
	return u.Watcher.Close()
}

// Clear screen function
func clear(osV string) {
	if osV == "linux" {
		cmd := exec.Command("go run main.go")
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

func fsNotifiyFunc(ctx context.Context, ch chan string, osV string, u *UgoSpectator) {
	for {
		select {
		case event, ok := <-u.Watcher.Events:
			if !ok {
				return
			}
			fileSplit := strings.SplitN(event.Name, "/", -1)
			if event.Op == fsnotify.Create {
				fmt.Printf("modified file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
			}
			if event.Op == fsnotify.Write {
				fmt.Printf("modified file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
			}
			if event.Op&fsnotify.Remove == fsnotify.Remove {
				fmt.Printf("Removed file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
				fmt.Printf("\033[1;33m%s%s\n\033[0m", "\nat ", "Reloading...")

			}

			fmt.Printf("\033[1;36m%s\033[0m", "Reloading...")
			time.Sleep(1 * time.Second)
			clear(osV)
			fmt.Printf("\033[1;36m%s\033[0m", "Ugo Spectator is watching your files")
			fmt.Printf("\033[1;33m%s%s\n\033[0m", "\nat ", u.dirname)
			ctx.Done()
			if event.Op&fsnotify.Rename == fsnotify.Rename {
				fmt.Printf("Removed file: %s/%s\n", fileSplit[len(fileSplit)-2], fileSplit[len(fileSplit)-1])
			}

		case err, ok := <-u.Watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("error:", err)
		}
	}

}
