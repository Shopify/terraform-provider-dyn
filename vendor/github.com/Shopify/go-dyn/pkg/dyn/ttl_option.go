package dyn

// TTLOption is a basic interface for passing optional parameters for a record
type TTLOption func(WithTTL)

// WithTTL defines an interface with a ttl option
type WithTTL interface {
	setTTL(ttl int)
}

// TTL provides a ttl option value
func TTL(ttl int) TTLOption {
	return func(w WithTTL) {
		w.setTTL(ttl)
	}
}
