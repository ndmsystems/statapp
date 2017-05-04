/*Package statapp implements monitoring application status using http
Example:
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

// Server implement a http server and contains storage of parameters
type Server struct {
	server    http.Server
	isRunning bool
	params    paramsIface
}

// New returns new implementation of Server
func New() *Server {
	return new(Server)
}

func (s *Server) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	params := s.params.GetAll()
	params["mem"] = m.HeapSys
	params["goroute"] = runtime.NumGoroutine()

	json.NewEncoder(w).Encode(params)
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
func (s *Server) Get(param string) (uint64, error) {
	val, err := s.params.Get(param)
	return val.(uint64), err
}

// Inc increments the value of the parameter
func (s *Server) Inc(param string, inc int) (err error) {
	val, err := s.params.Get(param)
	if err != nil {
		return
	}

	s.params.Set(param, val.(uint64)+uint64(inc))
	return
}

// Delete deletes the parameter from a storage
func (s *Server) Delete(param string) {
	s.params.Delete(param)
}
