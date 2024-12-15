package tlogm

import (
	"context"
	"runtime"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/outz"
	"github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"go.uber.org/mock/gomock"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/logm"
)

var (
	_ fixturez.BeforeSuite = (*MockHelper)(nil)
	_ fixturez.AfterSuite  = (*MockHelper)(nil)
	_ fixturez.BeforeTest  = (*MockHelper)(nil)
)

// MockHelper is a test helper.
type MockHelper struct {
	EnableLogrusOutput bool
	mock               *MockSender
	client             *libhoney.Client
}

// BeforeSuite implements [fixturez.BeforeSuite].
func (h *MockHelper) BeforeSuite(ctx context.Context, g *gomega.WithT) context.Context {
	mock := NewMockSender()
	var logger *logrus.Logger

	if h.EnableLogrusOutput {
		logger = outz.NewLogger()
		logger.SetLevel(logrus.TraceLevel)
		logger.SetFormatter(outz.NewHumanLogFormatter().SetInitTime(clkm.MustGet(ctx).Now()))
	}

	client, err := libhoney.NewClient(libhoney.ClientConfig{
		APIKey:       "test-honeycomb-api-key",
		Dataset:      "test-dataset",
		SampleRate:   1,
		Transmission: logm.NewSink(logger, mock),
	})
	g.Expect(err).To(gomega.Succeed())

	h.mock = mock
	h.client = client

	return logm.NewSingletonInjector(logm.NewRawLogFromClient(client))(ctx)
}

// AfterSuite implements [fixturez.AfterSuite].
func (h *MockHelper) AfterSuite(_ context.Context, _ *gomega.WithT) {
	h.client.Flush()
	h.client.Close()

	// Yield to other goroutines to ensure logs are fully flushed.
	runtime.Gosched()

	h.mock = nil
	h.client = nil
}

// BeforeTest implements [fixturez.BeforeTest].
func (h *MockHelper) BeforeTest(ctx context.Context, _ *gomega.WithT, _ *gomock.Controller) context.Context {
	h.mock.ClearEvents()
	return ctx
}

// GetMock returns the mock.
func (h *MockHelper) GetMock() *MockSender {
	return h.mock
}
