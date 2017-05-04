StatApp
=======
**Monitoring application status using http**
---

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/dzen-it/statapp/blob/master/LICENSE)
[![GoDoc](https://godoc.org/github.com/dzen-it/statapp?status.svg)](https://godoc.org/github.com/dzen-it/statapp)

Download and install it: `go get github.com/dzen-it/statapp` 

Returning parameters to JSON is easy:
``` go
package main

import (
	"time"

	"github.com/dzen-it/statapp"
)

func main() {
	// Runs a implementation of the server in background
	statapp.Start(":9999")

	// sets custom parameters
	statapp.Set("param-1", 42)
	statapp.Set("param-2", 2030)

	// Increment the value of the parameter
	statapp.Inc("param-2", 3) // Now the value is 2033
	
	sigCh:=make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
}
```
Make a request:
``` bash
    $ curl localhost:9999
    {"goroute":8,"mem":1703936,"param-1":42,,"param-2":2033}
```
You could notice `goroute` and `mem`. Yes, the parameters are present by default.
