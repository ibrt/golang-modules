// Package clk implements a clock module.
package clkm

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/ibrt/golang-utils/injectz"
)

type contextKey int

const (
	clkContextKey contextKey = iota
)

var (
	_ injectz.Initializer = Initializer
)

// Clock describes a clock.
type Clock clock.Clock

// Initializer initializes.
func Initializer(_ context.Context) (injectz.Injector, injectz.Releaser) {
	return NewSingletonInjector(clock.New()), injectz.NewNoopReleaser()
}

// NewSingletonInjector injects.
func NewSingletonInjector(clk Clock) injectz.Injector {
	return injectz.NewSingletonInjector(clkContextKey, clk)
}

// MustGet extracts, panics if not found.
func MustGet(ctx context.Context) Clock {
	return ctx.Value(clkContextKey).(Clock)
}
