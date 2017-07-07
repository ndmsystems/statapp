/*Package statapp implements monitoring application status using http
Example:
package main

import (
	"github.com/dzen-it/statapp"
)

func main() {
	// Runs a implementation of the server in background
	statapp.Start(":9999")   // http://localhost:9999

	// sets custom parameters
	statapp.Set("param-1", 42)

	// Increment the value of the parameter
	statapp.Inc("param-2", 3) // Now the value is 2033

}
*/
package statapp

import (
	"encoding/json"
	"net/http"
	"runtime"
	"sync"
)

var (
	server http.Server
	params = struct{
    sync.RWMutex
    m map[string]uint64
  }{m: make(map[string]uint64)}
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	params.RLock()
	defer params.RUnlock()
	params.m["mem"], params.m["goroute"] = mem.HeapSys, uint64(runtime.NumGoroutine())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(params)
}

// Start creates new implementation of a parameters storage
// and starts a new implementation of a server
func Start(addr string) {
	server = http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(rootHandler),
	}
	go start()
}

func start() {
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Set sets the parameter associated with integer value into a storage.
func Set(param string, val uint64) {
	params.RLock()
	defer params.RUnlock()
	params.m[param] = val
}

// Inc increments the value of the parameter
func Inc(param string, inc int) {
	params.RLock()
	defer params.RUnlock()
	params.m[param] += uint64(inc)
}
