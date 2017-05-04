package statapp

import (
	"fmt"
	"sync"
)

type paramsIface interface {
	Set(string, interface{})
	Get(string) (interface{}, error)
	GetAll() map[string]interface{}
	Delete(string)
}

// params implements a paramsIface
type params struct {
	sync.Mutex
	storage map[string]interface{}
}

func newParams() (p *params) {
	p = new(params)
	p.storage = make(map[string]interface{})
	return
}

func (p *params) Set(param string, val interface{}) {
	p.Lock()
	p.storage[param] = val
	p.Unlock()
}

func (p *params) Get(param string) (val interface{}, err error) {
	if val = p.storage[param]; val == nil {
		err = fmt.Errorf("param is not exist")

	}
	return
}
func (p *params) GetAll() map[string]interface{} {
	return p.storage
}

func (p *params) Delete(param string) {
	p.Lock()
	delete(p.storage, param)
	p.Unlock()
}
