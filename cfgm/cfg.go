package cfgm

import (
	"context"
	"reflect"

	"github.com/caarlos0/env/v11"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/injectz"
)

type contextKey int

const (
	cfgContextKey contextKey = iota
)

// Common sentinel values.
const (
	DisabledValue = "<disabled>"
	DefaultValue  = "<default>"
)

// Config describes a set of configuration values.
type Config interface {
	Config()
}

// ConfigLoader describes a function that can load a [Config], for example from environment variables or JSON file.
type ConfigLoader[T Config] func(ctx context.Context) (T, error)

// EnvConfigLoaderOptions describes the options for [MustNewEnvConfigLoader].
type EnvConfigLoaderOptions = env.Options

// MustNewEnvConfigLoader returns a [ConfigLoader] that loads the config from environment variables.
// Under the hood it uses "github.com/caarlos0/env/v11". T must be a struct pointer.
func MustNewEnvConfigLoader[T Config](options *EnvConfigLoaderOptions) ConfigLoader[T] {
	{
		var cfg T
		t := reflect.TypeOf(cfg)
		errorz.Assertf(t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct, "T must be a pointer to a struct")
	}

	if options == nil {
		options = &EnvConfigLoaderOptions{}
	}

	return func(ctx context.Context) (T, error) {
		var cfg T
		reflect.ValueOf(&cfg).Elem().Set(reflect.New(reflect.TypeOf(cfg).Elem()))
		return cfg, errorz.MaybeWrap(env.ParseWithOptions(cfg, *options))
	}
}

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
