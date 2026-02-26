package adapter

import (
	"fmt"
	"strings"
)

// Factory creates ClawAdapter instances based on type
type Factory struct {
	adapters map[AdapterType]func() ClawAdapter
}

// NewFactory creates a new adapter factory with registered adapters
func NewFactory() *Factory {
	f := &Factory{
		adapters: make(map[AdapterType]func() ClawAdapter),
	}

	// Register built-in adapters
	f.Register(AdapterTypeOpenClaw, func() ClawAdapter {
		return NewOpenClawAdapter()
	})

	// TODO: Register NanoClaw adapter when implemented
	// f.Register(AdapterTypeNanoClaw, func() ClawAdapter {
	// 	return NewNanoClawAdapter()
	// })

	return f
}

// Register registers a new adapter type
func (f *Factory) Register(adapterType AdapterType, constructor func() ClawAdapter) {
	f.adapters[adapterType] = constructor
}

// Create creates a ClawAdapter for the given type
func (f *Factory) Create(adapterType AdapterType) (ClawAdapter, error) {
	constructor, ok := f.adapters[adapterType]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrAdapterNotFound, adapterType)
	}

	return constructor(), nil
}

// CreateByString creates a ClawAdapter using a string type name
func (f *Factory) CreateByString(typeName string) (ClawAdapter, error) {
	adapterType := AdapterType(strings.TrimSpace(typeName))
	return f.Create(adapterType)
}

// IsSupported checks if an adapter type is supported
func (f *Factory) IsSupported(adapterType AdapterType) bool {
	_, ok := f.adapters[adapterType]
	return ok
}

// GetSupportedTypes returns a list of supported adapter types
func (f *Factory) GetSupportedTypes() []AdapterType {
	types := make([]AdapterType, 0, len(f.adapters))
	for t := range f.adapters {
		types = append(types, t)
	}
	return types
}

// DefaultFactory is the global default factory instance
var DefaultFactory = NewFactory()

// Create creates a ClawAdapter using the default factory
func Create(adapterType AdapterType) (ClawAdapter, error) {
	return DefaultFactory.Create(adapterType)
}

// CreateByString creates a ClawAdapter using the default factory
func CreateByString(typeName string) (ClawAdapter, error) {
	return DefaultFactory.CreateByString(typeName)
}

// IsSupported checks if an adapter type is supported by the default factory
func IsSupported(adapterType AdapterType) bool {
	return DefaultFactory.IsSupported(adapterType)
}
