package logm

import (
	"context"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/injectz"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/clkm"
)

type contextKey int

const (
	logContextKey contextKey = iota
)

var (
	_ injectz.Initializer = Initializer
)

// Log describes the module (with cached context).
type Log interface {
	EmitDebug(format string, options ...EmitOption)
	EmitInfo(format string, options ...EmitOption)
	EmitWarning(err error)
	EmitError(err error)
	EmitTraceLink(traceLink *TraceLink)
	Begin(name string, options ...BeginOption) (context.Context, func())
	SetUser(user *User)
	SetPropagatingField(k string, v any)
	SetMetadataKey(k string, v any)
	SetErrorMetadataKey(k string, v any)
	SetErrorFlag()
	GetCurrentTraceLink() *TraceLink
	Flush()
}

// RawLog describes the module.
type RawLog interface {
	EmitDebug(ctx context.Context, format string, options ...EmitOption)
	EmitInfo(ctx context.Context, format string, options ...EmitOption)
	EmitWarning(ctx context.Context, err error)
	EmitError(ctx context.Context, err error)
	EmitTraceLink(ctx context.Context, linkAnnotation *TraceLink)
	Begin(ctx context.Context, name string, options ...BeginOption) (context.Context, func())
	SetErrorFlag(ctx context.Context)
	SetUser(ctx context.Context, user *User)
	SetPropagatingField(ctx context.Context, k string, v any)
	SetMetadataKey(ctx context.Context, k string, v any)
	SetErrorMetadataKey(ctx context.Context, k string, v any)
	GetCurrentTraceLink(ctx context.Context) *TraceLink
	Flush(ctx context.Context)
}

// Initializer initializes.
func Initializer(ctx context.Context) (injectz.Injector, injectz.Releaser) {
	clkm.MustGet(ctx)
	logCfg := cfgm.MustGet[LogConfigMixin](ctx).GetLogConfig()

	client, err := libhoney.NewClient(libhoney.ClientConfig{
		APIKey:     logCfg.HoneycombAPIKey,
		Dataset:    logCfg.HoneycombDataset,
		SampleRate: logCfg.HoneycombSampleRate,
		Transmission: NewSink(
			MustNewDefaultLogrusLogger(ctx),
			MustNewDefaultHoneycombSender(ctx)),
	})
	errorz.MaybeMustWrap(err)

	// client.AddField("service.runId", cfgz.GetRunID())
	// client.AddField("service.version", cfgz.GetVersion())
	// client.AddField("service.name", cfgz.GetLambdaFunctionName())

	return NewSingletonInjector(NewLogFromClient(client)), func() { client.Close() }
}

// NewLogFromClient initializes a new RawLog.
func NewLogFromClient(client *libhoney.Client) RawLog {
	return &backgroundLogImpl{
		client: client,
	}
}

// NewSingletonInjector injects.
func NewSingletonInjector(log RawLog) injectz.Injector {
	return injectz.NewSingletonInjector(logContextKey, log)
}

// Get extracts, panics if not found.
func Get(ctx context.Context) Log {
	return &adapterLogImpl{
		ctx: ctx,
		log: ctx.Value(logContextKey).(RawLog),
	}
}
