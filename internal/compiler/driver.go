package compiler

// Driver wraps the underlying C/C++ compiler (gcc, clang, or msvc).
type Driver struct{}

// New creates a new compiler Driver.
func New() *Driver {
	return &Driver{}
}
