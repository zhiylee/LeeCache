package singleflight

import "sync"

// package singleflight make the duplicate function call at some time into one

// call : record the flight and bring back the value
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group : duplicate function call will be assigned to the same flight
type Group struct {
	mux sync.Mutex       //protects var m
	m   map[string]*call //lazily initialized
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mux.Lock()
	//lazily initialized
	if g.m == nil {
		g.m = make(map[string]*call)
	}

	//if already hava a flight for given key at a time
	if c, ok := g.m[key]; ok {
		g.mux.Unlock()
		c.wg.Wait()

		return c.val, c.err
	}

	//new a flight
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mux.Unlock()

	c.val, c.err = fn()
	c.wg.Done()

	g.mux.Lock()
	delete(g.m, key)
	g.mux.Unlock()

	return c.val, c.err
}
