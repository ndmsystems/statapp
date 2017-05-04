/*Package statapp implements monitoring application status using http
Example:
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"github.com/dzen-it/statapp"
	"syscall"
)

func main() {
	// Runs a implementation of the server in background
	statapp.Start(":9999")

	// sets custom parameters
	statapp.Set("param-1", 42)
	statapp.Set("param-2", 1970)
	statapp.Set("param-3", 2000)

	// Get the value of the parameter
	val := statapp.Get("param-2")
	fmt.Println(val)

	// Deletes parameter #2
	statapp.Delete("param-2")

	// Increment the value of the parameter
	statapp.Inc("param-3", 33) // Now the value is 2033

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
}
*/
package statapp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var (
	globalServer *Server
)

// Server implement a http server and contains storage of parameters
type Server struct {
	server    http.Server
	isRunning bool
	params    *params
}

// New returns new implementation of Server
func New() *Server {
	return new(Server)
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	params := s.params.GetAll()
	params["mem"], params["goroute"] = m.HeapSys, uint64(runtime.NumGoroutine())
	json.NewEncoder(w).Encode(params)

	w.Header().Set("Content-Type", "application/json")
}

// Start creates new implementation of a parameters storage
// and starts a new implementation of a server
func (s *Server) Start(addr string) (err error) {
	if s.isRunning {
		err = fmt.Errorf("Server already started")
		return
	}

	s.params = newParams()
	s.server = http.Server{
		Addr:    addr,
		Handler: http.HandlerFunc(s.rootHandler),
	}

	go start(s)
	go signalWaitingForShutdown(s)

	return
}

func start(s *Server) {
	s.isRunning = true
	defer func() {
		s.isRunning = false
	}()

	if err := s.server.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

// The method is gracefull shutdown
func signalWaitingForShutdown(s *Server) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)
	<-sigCh
	s.Stop()
}

// Stop sends a stop signal to the server.
// The parameters storage will also be cleaned.
func (s *Server) Stop() {
	defer func() { s.isRunning = false }()
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Println(err)
	}
}

// Set sets the parameter associated with integer value into a storage.
func (s *Server) Set(param string, val uint64) {
	s.params.Set(param, val)
}

// Get gets the value of the parameter.
func (s *Server) Get(param string) uint64 {
	return s.params.Get(param)
}

// Inc increments the value of the parameter
func (s *Server) Inc(param string, inc int) {
	s.params.Set(param, s.params.Get(param)+uint64(inc))
	return
}

// Delete deletes the parameter from a storage
func (s *Server) Delete(param string) {
	s.params.Delete(param)
}

func Start(addr string) (err error) {
	globalServer = New()
	err = globalServer.Start(addr)
	return
}

func Get(param string) (val uint64) {
	return globalServer.Get(param)
}

func Set(param string, val uint64) {
	globalServer.Set(param, val)
}

func Delete(param string) {
	globalServer.Delete(param)
}

func Inc(param string, inc int) {
	globalServer.Inc(param, inc)
}
