package tcfgm

import (
	"context"

	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/injectz"
	"github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/cfgm"
)

var (
	_ fixturez.BeforeSuite = (*Helper[cfgm.Config])(nil)
	_ fixturez.AfterSuite  = (*Helper[cfgm.Config])(nil)
)

// Helper is a test helper.
type Helper[T cfgm.Config] struct {
	ConfigLoader cfgm.ConfigLoader[T]
	releaser     injectz.Releaser
}

// BeforeSuite implements [fixturez.BeforeSuite].
func (h *Helper[T]) BeforeSuite(ctx context.Context, _ *gomega.WithT) context.Context {
	injector, releaser := cfgm.NewInitializer[T](h.ConfigLoader)(ctx)
	h.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements [fixturez.AfterSuite].
func (h *Helper[T]) AfterSuite(_ context.Context, _ *gomega.WithT) {
	h.releaser()
	h.releaser = nil
}
