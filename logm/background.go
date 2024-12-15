package logm

import (
	"context"
	"sync"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/idz"

	"github.com/ibrt/golang-modules/clkm"
)

var (
	_ RawLog = (*backgroundLogImpl)(nil)
)

type backgroundLogImpl struct {
	client *libhoney.Client
}

// EmitDebug implements the [RawLog] interface.
func (bL *backgroundLogImpl) EmitDebug(ctx context.Context, format string, options ...EmitOption) {
	o := newEmitOptions(options...)
	e := newAttachableEvent(ctx, bL.client, "", "debug")
	addDebugFields(e, format, o)
	errorz.MaybeMustWrap(e.Send())
}

// EmitInfo implements the [RawLog] interface.
func (bL *backgroundLogImpl) EmitInfo(ctx context.Context, format string, options ...EmitOption) {
	o := newEmitOptions(options...)
	e := newAttachableEvent(ctx, bL.client, "", "info")
	addInfoFields(e, format, o)
	errorz.MaybeMustWrap(e.Send())
}

// EmitWarning implements the [RawLog] interface.
func (bL *backgroundLogImpl) EmitWarning(ctx context.Context, err error) {
	maybeSetIsEmitted(err)
	e := newAttachableEvent(ctx, bL.client, "", getWarningName(err))
	addWarningFields(e, err)
	errorz.MaybeMustWrap(e.Send())
}

// EmitError implements the [RawLog] interface.
func (bL *backgroundLogImpl) EmitError(ctx context.Context, err error) {
	maybeSetIsEmitted(err)
	e := newAttachableEvent(ctx, bL.client, "", getErrorName(err))
	addErrorFields(e, err)
	errorz.MaybeMustWrap(e.Send())
}

// EmitTraceLink implements the [RawLog] interface.
func (bL *backgroundLogImpl) EmitTraceLink(ctx context.Context, _ *TraceLink) {
	bL.EmitWarning(ctx, errorz.Errorf("called EmitTraceLink in background Log"))
}

// Begin implements the [RawLog] interface.
func (bL *backgroundLogImpl) Begin(ctx context.Context, name string, options ...BeginOption) (context.Context, func()) {
	o := newBeginOptions(options...)

	sL := &spanLogImpl{
		m:            &sync.Mutex{},
		b:            bL.client.NewBuilder(),
		startTime:    clkm.MustGet(ctx).Now(),
		name:         name,
		traceID:      idz.MustNewRandomUUID(),
		parentID:     "",
		spanID:       idz.MustNewRandomUUID(),
		metadata:     o.metadata,
		errMetadata:  o.errMetadata,
		hasErrorFlag: false,
	}

	sL.b.AddField("trace.trace_id", sL.traceID)
	ctx = NewSingletonInjector(sL)(ctx)

	return ctx, func() {
		sL.end(ctx)
	}
}

// SetUser implements the [RawLog] interface.
func (bL *backgroundLogImpl) SetUser(ctx context.Context, _ *User) {
	bL.EmitWarning(ctx, errorz.Errorf("called SetUser in background Log"))
}

// SetPropagatingField implements the [RawLog] interface.
func (bL *backgroundLogImpl) SetPropagatingField(ctx context.Context, _ string, _ any) {
	bL.EmitWarning(ctx, errorz.Errorf("called SetPropagatingField in background Log"))
}

// SetMetadataKey implements the [RawLog] interface.
func (bL *backgroundLogImpl) SetMetadataKey(ctx context.Context, _ string, _ any) {
	bL.EmitWarning(ctx, errorz.Errorf("called SetMetadataKey in background Log"))
}

// SetErrorMetadataKey implements the [RawLog] interface.
func (bL *backgroundLogImpl) SetErrorMetadataKey(ctx context.Context, _ string, _ any) {
	bL.EmitWarning(ctx, errorz.Errorf("called SetErrorMetadataKey in background Log"))
}

// SetErrorFlag implements the [RawLog] interface.
func (bL *backgroundLogImpl) SetErrorFlag(ctx context.Context) {
	bL.EmitWarning(ctx, errorz.Errorf("called SetErrorFlag in background Log"))
}

// GetCurrentTraceLink implements the [RawLog] interface.
func (bL *backgroundLogImpl) GetCurrentTraceLink(_ context.Context) *TraceLink {
	return nil
}

// Flush implements the [RawLog] interface.
func (bL *backgroundLogImpl) Flush(_ context.Context) {
	bL.client.Flush()
}
