package cfgm

import (
	"context"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/injectz"
)

type contextKey int

const (
	cfgContextKey contextKey = iota
)

// Config describes a set of configuration values.
type Config interface {
	Config()
}

// ConfigLoader describes a function that can load a [Config], for example from environment variables or JSON file.
type ConfigLoader[T Config] func(ctx context.Context) (T, error)

// NewInitializer returns a [injectz.Initializer] that uses the given [ConfigLoader].
func NewInitializer[T Config](cfgLoader ConfigLoader[T]) injectz.Initializer {
	return func(ctx context.Context) (injectz.Injector, injectz.Releaser) {
		cfg, err := cfgLoader(ctx)
		errorz.MaybeMustWrap(err)
		return NewSingletonInjector(cfg), injectz.NewNoopReleaser()
	}
}

// NewSingletonInjector injects.
func NewSingletonInjector[T Config](cfg T) injectz.Injector {
	return injectz.NewSingletonInjector(cfgContextKey, cfg)
}

// MustGet extracts, panics if not found or T is wrong.
func MustGet[T Config](ctx context.Context) T {
	return ctx.Value(cfgContextKey).(T)
}
