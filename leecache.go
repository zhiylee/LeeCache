package leecache

// Getter : Data Getter use a key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc : A function type implements  Getter
type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var Localhost string
