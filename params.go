package statapp

// params implements a paramsIface

type params struct {
	storage map[string]uint64
}

func newParams() (p *params) {
	p = new(params)
	p.storage = make(map[string]uint64)
	return
}

func (p *params) Set(param string, val uint64) {
	p.storage[param] = val
}

func (p *params) Get(param string) (val uint64) {
	return p.storage[param]
}
func (p *params) GetAll() map[string]uint64 {
	return p.storage
}

func (p *params) Delete(param string) {
	delete(p.storage, param)
}
