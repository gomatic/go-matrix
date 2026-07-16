package matrix

import (
	errs "github.com/gomatic/go-error"
)

// Every error the package can emit is a const of the ecosystem's [errs.Const]
// sentinel type, so each is matchable with errors.Is rather than by string.
const (
	// ErrNonPositiveWidth is returned by [New] when width is not greater than zero.
	ErrNonPositiveWidth errs.Const = "width must be positive"
	// ErrNonPositiveHeight is returned by [New] when height is not greater than zero.
	ErrNonPositiveHeight errs.Const = "height must be positive"
)
