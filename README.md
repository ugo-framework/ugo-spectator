<div align="center">
    <img width="300" src="./assets/ugo-spectator.png" alt="ugo spectator logo">
</div>

# Ugo Spectator
 _The only file watcher you will ever need written in go (Golang)._

## Installation
```shell script
go get github.com/ugo-framework/ugo-spectator
```

Then use it in your code use the following sysntax
```go
package examples

import spectator "github.com/ugo-framework/ugo-spectator/lib"

func main() {
	// initialise the spectator with the dirname
	watcher, err := spectator.Init(".")
	if err != nil {
		// Handle Error
	}
	defer watcher.Close() // handle error
	// event to catch for file change
	for {
		select {
		case res := <-watcher.Ch:
			// Do Something
		}
	}
}

```
The full example can be viewed [here](./examples/main.go)
This documentation is not Complete and still in its early stages
