package pkg

// Resolver handles package dependency resolution.
type Resolver struct{}

// New creates a new Resolver.
func New() *Resolver {
	return &Resolver{}
}
