package vm

import "sync"

type f64pool struct {
	items [][]float64
	max   int
	cap   int
	mux   sync.Mutex
}

func (p *f64pool) Put(item []float64) {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.items) == p.max {
		return
	}

	p.items = append(p.items, item)
}

func (p *f64pool) Get() []float64 {
	p.mux.Lock()
	defer p.mux.Unlock()

	if len(p.items) == 0 {
		item := make([]float64, 0, p.cap)
		return item
	}

	item := p.items[len(p.items)-1]
	p.items = p.items[:len(p.items)-1]
	return item
}
