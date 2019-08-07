package dyn

// PaginationOption is a basic interface for passing optional parameters to a function
type PaginationOption func(WithPagination)

// WithPagination defines an interface with limit and offset options
type WithPagination interface {
	setLimit(limit int)
	setOffset(offset int)
}

// Limit provides a limit option value
func Limit(limit int) PaginationOption {
	return func(w WithPagination) {
		w.setLimit(limit)
	}
}

// Offset provides an offset option value
func Offset(offset int) PaginationOption {
	return func(w WithPagination) {
		w.setOffset(offset)
	}
}
