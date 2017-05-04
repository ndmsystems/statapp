/*Package statapp implements monitoring application status using http
Example:
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

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
}
*/
package statapp

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime"
)

var (
	server http.Server
	params map[string]uint64
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	
	params["mem"], params["goroute"] = mem.HeapSys, uint64(runtime.NumGoroutine())
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}

// Start creates new implementation of a parameters storage
// and starts a new implementation of a server
func Start(addr string) {
	params = make(map[string]uint64)
	server = http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(rootHandler),
	}
	go start()
}

func start() {
	if err := server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// Stop sends a stop signal to the server.
// The parameters storage will also be cleaned.
func Stop() {
	if err := server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}

// Set sets the parameter associated with integer value into a storage.
func Set(param string, val uint64) {
	params[param] = val
}

// Inc increments the value of the parameter
func Inc(param string, inc int) {
	params[param] = params[param] + uint64(inc)
}
