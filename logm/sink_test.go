package logm_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/honeycombio/libhoney-go"
	"github.com/honeycombio/libhoney-go/transmission"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/outz"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
	"github.com/ibrt/golang-modules/logm"
	"github.com/ibrt/golang-modules/logm/tlogm"
)

type SinkSuite struct {
	CLK *tclkm.MockHelper
}

func TestSinkSuite(t *testing.T) {
	fixturez.RunSuite(t, &SinkSuite{})
}

func (s *SinkSuite) TestSink_Logger_Other(ctx context.Context, g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupSirupsenLogrus)
	defer outz.ResetOutputCapture()

	logger := outz.NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	sink := logm.NewSink(logger, nil)
	sink.Add(&transmission.Event{
		Timestamp: clkm.MustGet(ctx).Now(),
		Data: map[string]any{
			"name": "name",
			"k":    "v",
		},
	})

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(""))

	g.Expect(errBuf).To(Equal(fmt.Sprintf(
		"{\"k\":\"v\",\"level\":\"debug\",\"msg\":\"name\",\"name\":\"name\",\"time\":\"%v\"}\n",
		clkm.MustGet(ctx).Now().Format(time.RFC3339Nano))))
}

func (s *SinkSuite) TestSink_Logger_Debug(ctx context.Context, g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupSirupsenLogrus)
	defer outz.ResetOutputCapture()

	logger := outz.NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	sink := logm.NewSink(logger, nil)
	sink.Add(&transmission.Event{
		Timestamp: clkm.MustGet(ctx).Now(),
		Data: map[string]any{
			"debug":         true,
			"debug.message": "msg",
			"k":             "v",
			"name":          "name",
		},
	})

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(""))

	g.Expect(errBuf).To(Equal(fmt.Sprintf(
		"{\"debug\":true,\"debug.message\":\"msg\",\"k\":\"v\",\"level\":\"debug\",\"msg\":\"name: msg\",\"name\":\"name\",\"time\":\"%v\"}\n",
		clkm.MustGet(ctx).Now().Format(time.RFC3339Nano))))
}

func (s *SinkSuite) TestSink_Logger_Info(ctx context.Context, g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupSirupsenLogrus)
	defer outz.ResetOutputCapture()

	logger := outz.NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	sink := logm.NewSink(logger, nil)
	sink.Add(&transmission.Event{
		Timestamp: clkm.MustGet(ctx).Now(),
		Data: map[string]any{
			"info": true,
			"k":    "v",
		},
	})

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(""))

	g.Expect(errBuf).To(Equal(fmt.Sprintf(
		"{\"info\":true,\"k\":\"v\",\"level\":\"info\",\"msg\":\"\",\"time\":\"%v\"}\n",
		clkm.MustGet(ctx).Now().Format(time.RFC3339Nano))))
}

func (s *SinkSuite) TestSink_Logger_Warning(ctx context.Context, g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupSirupsenLogrus)
	defer outz.ResetOutputCapture()

	logger := outz.NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	sink := logm.NewSink(logger, nil)
	sink.Add(&transmission.Event{
		Timestamp: clkm.MustGet(ctx).Now(),
		Data: map[string]any{
			"warning": true,
			"k":       "v",
		},
	})

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(""))

	g.Expect(errBuf).To(Equal(fmt.Sprintf(
		"{\"k\":\"v\",\"level\":\"warning\",\"msg\":\"\",\"time\":\"%v\",\"warning\":true}\n",
		clkm.MustGet(ctx).Now().Format(time.RFC3339Nano))))
}

func (s *SinkSuite) TestSink_Logger_Error(ctx context.Context, g *WithT) {
	outz.MustBeginOutputCapture(outz.OutputSetupSirupsenLogrus)
	defer outz.ResetOutputCapture()

	logger := outz.NewLogger()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	sink := logm.NewSink(logger, nil)
	sink.Add(&transmission.Event{
		Timestamp: clkm.MustGet(ctx).Now(),
		Data: map[string]any{
			"error": true,
			"k":     "v",
		},
	})

	outBuf, errBuf := outz.MustEndOutputCapture()
	g.Expect(outBuf).To(Equal(""))

	g.Expect(errBuf).To(Equal(fmt.Sprintf(
		"{\"error\":true,\"k\":\"v\",\"level\":\"error\",\"msg\":\"\",\"time\":\"%v\"}\n",
		clkm.MustGet(ctx).Now().Format(time.RFC3339Nano))))
}

func (s *SinkSuite) TestSink_Responses(g *WithT) {
	sink := logm.NewSink(nil, nil)
	g.Expect(sink.SendResponse(transmission.Response{})).To(BeTrue())
	g.Expect(sink.Flush()).To(Succeed())

	ok := &atomic.Bool{}
	closed := &atomic.Bool{}

	go func() {
		defer func() {
			closed.Store(true)
		}()

		for range sink.TxResponses() {
			ok.Store(true)
		}
	}()

	g.Eventually(func() bool { return sink.SendResponse(transmission.Response{}) }, "1s", "10ms").Should(BeFalse())
	g.Eventually(ok.Load, "1s", "10ms").Should(BeTrue())
	g.Expect(sink.Stop()).To(Succeed())
	g.Eventually(closed.Load, "1s", "10ms").Should(BeTrue())
}

func (s *SinkSuite) TestSink_Sender(g *WithT) {
	sink := logm.NewSink(nil, tlogm.NewMockSender())
	g.Expect(sink.SendResponse(transmission.Response{})).To(BeTrue())
	g.Expect(sink.Flush()).To(Succeed())

	ok := &atomic.Bool{}
	closed := &atomic.Bool{}

	go func() {
		defer func() {
			closed.Store(true)
		}()

		for range sink.TxResponses() {
			ok.Store(true)
		}
	}()

	g.Eventually(func() bool { return sink.SendResponse(transmission.Response{}) }, "1s", "10ms").Should(BeFalse())
	g.Eventually(ok.Load, "1s", "10ms").Should(BeTrue())
	g.Expect(sink.Stop()).To(Succeed())
	g.Eventually(closed.Load, "1s", "10ms").Should(BeTrue())
}

func (s *SinkSuite) TestSink_Noop(g *WithT) {
	sink := logm.NewSink(nil, nil)
	g.Expect(sink.Start()).To(BeNil())
	g.Expect(func() { sink.Add(nil) }).ToNot(Panic())
	g.Expect(sink.Flush()).To(BeNil())
	g.Expect(sink.Stop()).To(BeNil())
}

func (s *SinkSuite) TestMustNewDefaultLogrusLogger(ctx context.Context, g *WithT) {
	{
		g.Expect(
			logm.MustNewDefaultLogrusLogger(
				cfgm.NewSingletonInjector[logm.LogConfigMixin](&logm.LogConfig{
					LogrusOutput: cfgm.DisabledValue,
					LogrusLevel:  logrus.InfoLevel,
				})(ctx))).
			To(BeNil())
	}
	{
		logger := logm.MustNewDefaultLogrusLogger(
			cfgm.NewSingletonInjector[logm.LogConfigMixin](&logm.LogConfig{
				LogrusOutput: logm.LogConfigLogrusOutputHuman,
				LogrusLevel:  logrus.WarnLevel,
			})(ctx))

		g.Expect(logger).ToNot(BeNil())
		g.Expect(logger.Level).To(Equal(logrus.WarnLevel))

		_, ok := logger.Formatter.(*outz.HumanLogFormatter)
		g.Expect(ok).To(BeTrue())
	}
	{
		logger := logm.MustNewDefaultLogrusLogger(
			cfgm.NewSingletonInjector[logm.LogConfigMixin](&logm.LogConfig{
				LogrusOutput: logm.LogConfigLogrusOutputJSON,
				LogrusLevel:  logrus.TraceLevel,
			})(ctx))

		g.Expect(logger).ToNot(BeNil())
		g.Expect(logger.Level).To(Equal(logrus.TraceLevel))

		_, ok := logger.Formatter.(*logrus.JSONFormatter)
		g.Expect(ok).To(BeTrue())
	}
}

func (s *SinkSuite) TestMustNewDefaultHoneycombSender(ctx context.Context, g *WithT) {
	{
		g.Expect(
			logm.MustNewDefaultHoneycombSender(
				cfgm.NewSingletonInjector[logm.LogConfigMixin](&logm.LogConfig{
					HoneycombAPIKey: cfgm.DisabledValue,
				})(ctx))).
			To(BeNil())
	}
	{
		g.Expect(
			logm.MustNewDefaultHoneycombSender(
				cfgm.NewSingletonInjector[logm.LogConfigMixin](&logm.LogConfig{
					HoneycombAPIKey: "test-honeycomb-api-key",
				})(ctx))).
			To(Equal(&transmission.Honeycomb{
				BatchTimeout:         libhoney.DefaultBatchTimeout,
				BlockOnSend:          true,
				MaxBatchSize:         libhoney.DefaultMaxBatchSize,
				MaxConcurrentBatches: libhoney.DefaultMaxConcurrentBatches,
				PendingWorkCapacity:  libhoney.DefaultPendingWorkCapacity,
			}))
	}
}
