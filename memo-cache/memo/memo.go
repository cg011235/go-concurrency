package memo

type Func func (key string) (interface{}, error)

type result struct {
	value	interface{}
	err		error
}

type Memo struct {
	f		Func				// function that is to be cached
	cache	map[string]result	// Maps all url to result so we lookup them fast
}

// TODO(chaitanya): use regular expression to validate url
func isValidHttpURL(key string) bool {
	return true
}

func (m *Memo) Get(key string) (interface{}, error) {
	res, ok := m.cache[key]
	if !ok {
		res.value, res.err = m.f(key)
		m.cache[key] = res
	}
	return res.value, res.err
}

func New(f Func) *Memo {
	return &Memo{f, make(map[string]result)}
}