package logm

import (
	"context"
	"fmt"
	"time"

	"github.com/honeycombio/libhoney-go"
	"github.com/honeycombio/libhoney-go/transmission"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/outz"
	"github.com/sirupsen/logrus"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/clkm"
)

var (
	_ transmission.Sender = (*Sink)(nil)
)

var (
	messageKeys = []string{
		"debug.message",
		"info.message",
		"warning.message",
		"error.message",
	}
)

// MustNewDefaultLogrusLogger initializes a default [*logrus.Logger] using the [LogConfigMixin] from context.
func MustNewDefaultLogrusLogger(ctx context.Context) *logrus.Logger {
	if logCfg := cfgm.MustGet[LogConfigMixin](ctx).GetLogConfig(); logCfg.LogrusOutput != LogConfigLogrusOutputDisabled {
		logger := outz.NewLogger()
		logger.SetLevel(logCfg.LogrusLevel)

		switch v := logCfg.LogrusOutput; v {
		case LogConfigLogrusOutputHuman:
			logger.SetFormatter(outz.NewHumanLogFormatter().SetInitTime(clkm.MustGet(ctx).Now()))
		case LogConfigLogrusOutputJSON:
			logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
		default:
			errorz.MustErrorf("invalid value for LogConfigLogrusOutput: '%s'", v)
		}

		return logger
	}

	return nil
}

// MustNewDefaultHoneycombSender initializes a default [transmission.Sender] using the [LogConfigMixin] from context.
func MustNewDefaultHoneycombSender(ctx context.Context) transmission.Sender {
	if apiKey := cfgm.MustGet[LogConfigMixin](ctx).GetLogConfig().HoneycombAPIKey; apiKey != "" && apiKey != cfgm.DisabledValue {
		return &transmission.Honeycomb{
			BatchTimeout:         libhoney.DefaultBatchTimeout,
			BlockOnSend:          true,
			MaxBatchSize:         libhoney.DefaultMaxBatchSize,
			MaxConcurrentBatches: libhoney.DefaultMaxConcurrentBatches,
			PendingWorkCapacity:  libhoney.DefaultPendingWorkCapacity,
		}
	}

	return nil
}

// Sink describes a sink.
type Sink struct {
	l *logrus.Logger
	s transmission.Sender
	c chan transmission.Response
}

// NewSink initializes a new [*Sink].
func NewSink(logger *logrus.Logger, sender transmission.Sender) *Sink {
	s := &Sink{
		l: logger,
		s: sender,
	}

	if sender == nil {
		s.c = make(chan transmission.Response)
	}

	return s
}

// Add implements the transmission.Sender interface.
func (s *Sink) Add(e *transmission.Event) {
	if s.l != nil {
		logrus.NewEntry(s.l).
			WithTime(e.Timestamp).
			WithFields(e.Data).
			Log(s.getLevel(e.Data), s.getMessage(e.Data))
	}

	if s.s != nil {
		s.s.Add(e)
	}
}

func (s *Sink) getLevel(data map[string]any) logrus.Level {
	if _, ok := data["debug"]; ok {
		return logrus.DebugLevel
	} else if _, ok := data["info"]; ok {
		return logrus.InfoLevel
	} else if _, ok := data["warning"]; ok {
		return logrus.WarnLevel
	} else if _, ok := data["error"]; ok {
		return logrus.ErrorLevel
	} else {
		return logrus.DebugLevel
	}
}

func (s *Sink) getMessage(data map[string]any) string {
	msg := ""

	for _, k := range messageKeys {
		if v, ok := data[k].(string); ok && v != "" {
			msg = v
			break
		}
	}

	if name, ok := data["name"].(string); ok && name != "" {
		if msg != "" {
			return fmt.Sprintf("%v: %v", name, msg)
		}
		return name
	}

	return msg
}

// Start implements the transmission.Sender interface.
func (s *Sink) Start() error {
	if s.s != nil {
		return s.s.Start()
	}

	return nil
}

// Stop implements the transmission.Sender interface.
func (s *Sink) Stop() error {
	if s.s != nil {
		return s.s.Stop()
	}

	close(s.c)
	return nil
}

// Flush implements the transmission.Sender interface.
func (s *Sink) Flush() error {
	if s.s != nil {
		return s.s.Flush()
	}

	return nil
}

// TxResponses implements the transmission.Sender interface.
func (s *Sink) TxResponses() chan transmission.Response {
	if s.s != nil {
		return s.s.TxResponses()
	}

	return s.c
}

// SendResponse implements the transmission.Sender interface.
func (s *Sink) SendResponse(response transmission.Response) bool {
	if s.s != nil {
		return s.s.SendResponse(response)
	}

	select {
	case s.c <- response:
		return false
	default:
		return true
	}
}
