package tclkm

import (
	"context"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/injectz"
	"github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	"github.com/ibrt/golang-modules/clkm"
)

var (
	_ fixturez.BeforeSuite = (*RealHelper)(nil)
	_ fixturez.AfterSuite  = (*RealHelper)(nil)
	_ fixturez.BeforeSuite = (*MockHelper)(nil)
	_ fixturez.AfterSuite  = (*MockHelper)(nil)
	_ fixturez.BeforeTest  = (*MockHelper)(nil)
)

// RealHelper is a test helper.
type RealHelper struct {
	releaser injectz.Releaser
}

// BeforeSuite implements [fixturez.BeforeSuite].
func (h *RealHelper) BeforeSuite(ctx context.Context, _ *gomega.WithT) context.Context {
	injector, releaser := clkm.Initializer(ctx)
	h.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements [fixturez.AfterSuite].
func (h *RealHelper) AfterSuite(_ context.Context, _ *gomega.WithT) {
	h.releaser()
	h.releaser = nil
}

// MockHelper is a test helper.
type MockHelper struct {
	mock *clock.Mock
}

// BeforeSuite implements [fixturez.BeforeSuite].
func (h *MockHelper) BeforeSuite(ctx context.Context, _ *gomega.WithT) context.Context {
	h.mock = clock.NewMock()
	h.mock.Set(time.Now())
	return clkm.NewSingletonInjector(h.mock)(ctx)
}

// AfterSuite implements [fixturez.AfterSuite].
func (h *MockHelper) AfterSuite(_ context.Context, _ *gomega.WithT) {
	h.mock = nil
}

// BeforeTest implements [fixturez.BeforeTest].
func (h *MockHelper) BeforeTest(ctx context.Context, _ *gomega.WithT, _ *gomock.Controller) context.Context {
	h.mock.Set(time.Now().UTC())
	return ctx
}

// GetMock returns the mock.
func (h *MockHelper) GetMock() *clock.Mock {
	return h.mock
}
