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
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// UgoSpectator struct with Watcher  and all methods
type UgoSpectator struct {
	Watcher   *fsnotify.Watcher  // *fsnotify watcher instance
	osV       string             // osV take the current os
	dirname   string             // dirname to watch
	Ch        chan bool          // ch to trigger on file change
	CancelCtx context.CancelFunc // context to cancel
}

// Init initializes the fsnotify NewWatcher and
//  a *fsnotify watcher instance and an error
func Init(dirname string) (*UgoSpectator, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return &UgoSpectator{}, err
	}
	ugoWatcher := &UgoSpectator{Watcher: watcher, osV: runtime.GOOS, Ch: make(chan bool)}
	clear(ugoWatcher.osV)
	fmt.Printf("\033[1;36m%s\033[0m", "Ugo Spectator is watching your files")
	ctx, cancel := context.WithCancel(context.Background())
	ugoWatcher.CancelCtx = cancel
	cPath, err := os.Getwd()
	if err != nil {
		return &UgoSpectator{}, err
	}
	pathToWatch := path.Join(cPath, dirname)
	ugoWatcher.dirname = pathToWatch
	fmt.Printf("\033[1;33m%s%s\n\033[0m", "\nat ", pathToWatch)
	dirsToWatch := []string{pathToWatch}
	files, err := ioutil.ReadDir(pathToWatch)
	if err != nil {
		fmt.Println(err)
	}
	for _, f := range files {
		if f.IsDir() {
			err := filepath.Walk(f.Name(),
				func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() && info.Name()[:1] != "." {
						dirsToWatch = append(dirsToWatch, path)
					}
					return nil
				})
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	for _, p := range dirsToWatch {
		err = watcher.Add(p)
		if err != nil {
			fmt.Println(err)
		}
	}
	// calling fsNotifiyFunc in a goroutine
	go fsNotifiyFunc(ctx, ugoWatcher.osV, ugoWatcher)
	return ugoWatcher, nil
}

// Close return error if an error occurs during the closing of the
// fsnotify watcher instance
func (u *UgoSpectator) Close() error {
	u.CancelCtx()
	fmt.Printf("\033[1;31m%s\033[0m", "\nUgo Spector Closing\n")
	return u.Watcher.Close()
}

// Clear screen function
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

func fsNotifiyFunc(ctx context.Context, osV string, u *UgoSpectator) {
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
			u.Ch <- true
			fmt.Printf("\033[1;36m%s\033[0m", "Reloading...")
			time.Sleep(time.Second * 1)
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
