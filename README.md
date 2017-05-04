StatApp
=======
**Monitoring application status using http**
---

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/hyperium/hyper/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/dzen-it/statapp?status.svg)](https://godoc.org/github.com/dzen-it/statapp)

Download and install it: `go get github.com/dzen-it/statapp` 

Returning parameters to JSON is easy:
``` go
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"statapp"
	"syscall"
)

func main() {
	s := statapp.New()
	// Runs a implementation of the server in background
	s.Start(":9999")

	// sets custom parameters
	s.Set("param-1", 42)
	s.Set("param-2", 1970)
	s.Set("param-3", 2000)

	// Get the value of the parameter
	val, err := s.Get("param-2")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(val)

	// Deletes parameter #2
	s.Delete("param-2")

	// Increment the value of the parameter
	err = s.Inc("param-3", 33) // Now the value is 2033
	if err != nil {
		log.Fatalln(err)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
}
```
Make a request:
``` bash
    $ curl localhost:9999
    {"goroute":8,"mem":1703936,"param-1":42,,"param-3":2033}
```
You could notice `goroute` and `mem`. Yes, the parameters are present by default.