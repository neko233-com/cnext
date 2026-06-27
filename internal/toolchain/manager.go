package toolchain

// Manager discovers and configures compiler toolchains.
type Manager struct{}

// New creates a new Manager.
func New() *Manager {
	return &Manager{}
}
