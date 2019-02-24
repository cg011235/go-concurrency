package memo

import "sync"

type Func func (key string) (interface{}, error)

type result struct {
	value	interface{}
	err		error
}

type entry struct {
	res		result
	ready	chan struct{}
}
type Memo struct {
	f		Func				// function that is to be cached
	cache	map[string]*entry	// Maps all url to result so we lookup them fast
	mu		sync.RWMutex		// mutex to control access
}

// TODO(chaitanya): use regular expression to validate url
func isValidHttpURL(key string) bool {
	return true
}

func (m *Memo) Get(key string) (interface{}, error) {
	m.mu.Lock()
	e := m.cache[key]
	if e == nil {  // This is first time request
		e = &entry{ready:make(chan struct{})}
		m.cache[key] = e
		m.mu.Unlock()  // Complete critical section fast

		e.res.value, e.res.err = m.f(key)

		close(e.ready)  // Broadcast that result is ready
	} else {  // This is repeat request
		m.mu.Unlock()
		<-e.ready
	}
	return e.res.value, e.res.err
}

func New(f Func) *Memo {
	return &Memo{f, make(map[string]*entry), sync.RWMutex{}}
}