package logm

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/idz"

	"github.com/ibrt/golang-modules/clkm"
)

var (
	_ RawLog = (*spanLogImpl)(nil)
)

type spanLogImpl struct {
	m            *sync.Mutex
	b            *libhoney.Builder
	startTime    time.Time
	name         string
	traceID      string
	parentID     string
	spanID       string
	metadata     map[string]any
	errMetadata  map[string]any
	hasErrorFlag bool
}

// EmitDebug implements the RawLog interface.
func (sL *spanLogImpl) EmitDebug(ctx context.Context, format string, options ...EmitOption) {
	sL.m.Lock()
	defer sL.m.Unlock()

	o := newEmitOptions(options...)
	e := newAttachableEvent(ctx, sL.b, sL.spanID, "debug")
	addDebugFields(e, format, o)
	errorz.MaybeMustWrap(e.Send())
}

// EmitInfo implements the RawLog interface.
func (sL *spanLogImpl) EmitInfo(ctx context.Context, format string, options ...EmitOption) {
	sL.m.Lock()
	defer sL.m.Unlock()

	o := newEmitOptions(options...)
	e := newAttachableEvent(ctx, sL.b, sL.spanID, "info")
	addInfoFields(e, format, o)
	errorz.MaybeMustWrap(e.Send())
}

// EmitWarning implements the RawLog interface.
func (sL *spanLogImpl) EmitWarning(ctx context.Context, err error) {
	sL.m.Lock()
	defer sL.m.Unlock()

	maybeSetIsEmitted(err)
	e := newAttachableEvent(ctx, sL.b, sL.spanID, "warning")
	addWarningFields(e, err)
	errorz.MaybeMustWrap(e.Send())
}

// EmitError implements the RawLog interface.
func (sL *spanLogImpl) EmitError(ctx context.Context, err error) {
	sL.m.Lock()
	defer sL.m.Unlock()

	maybeSetIsEmitted(err)
	sL.hasErrorFlag = true
	e := newAttachableEvent(ctx, sL.b, sL.spanID, "error")
	addErrorFields(e, err)
	errorz.MaybeMustWrap(e.Send())
}

// EmitTraceLink implements the RawLog interface.
func (sL *spanLogImpl) EmitTraceLink(ctx context.Context, traceLink *TraceLink) {
	sL.m.Lock()
	defer sL.m.Unlock()

	if traceLink != nil && traceLink.Serialize() != "" {
		e := newTraceLinkEvent(ctx, sL.b, sL.spanID, traceLink)
		errorz.MaybeMustWrap(e.Send())
	}
}

// Begin implements the RawLog interface.
func (sL *spanLogImpl) Begin(ctx context.Context, name string, options ...BeginOption) (context.Context, func()) {
	sL.m.Lock()
	defer sL.m.Unlock()

	o := newBeginOptions(options...)

	nsL := &spanLogImpl{
		m:            &sync.Mutex{},
		b:            sL.b.Clone(),
		startTime:    clkm.MustGet(ctx).Now(),
		name:         name,
		traceID:      sL.traceID,
		parentID:     sL.spanID,
		spanID:       idz.MustNewRandomUUID(),
		metadata:     o.metadata,
		errMetadata:  o.errMetadata,
		hasErrorFlag: false,
	}

	ctx = NewSingletonInjector(nsL)(ctx)

	return ctx, func() {
		nsL.end(ctx)
	}
}

func (sL *spanLogImpl) end(ctx context.Context) {
	sL.m.Lock()
	defer sL.m.Unlock()

	e := newTraceableEvent(ctx, sL.b, sL.name, sL.spanID, sL.parentID, sL.startTime)
	addMetadataFields(e, "scope.metadata", sL.metadata)

	if sL.hasErrorFlag {
		e.AddField("error", true)
		addMetadataFields(e, "scope.metadata", sL.errMetadata)
	}

	errorz.MaybeMustWrap(e.Send())
}

// SetUser implements the RawLog interface.
func (sL *spanLogImpl) SetUser(_ context.Context, user *User) {
	sL.m.Lock()
	defer sL.m.Unlock()
	maybeAddUserFields(sL.b, user)
}

// SetPropagatingField implements the RawLog interface.
func (sL *spanLogImpl) SetPropagatingField(_ context.Context, k string, v any) {
	sL.m.Lock()
	defer sL.m.Unlock()
	sL.b.AddField(fmt.Sprintf("scope.%v", strings.TrimPrefix(k, "scope.")), v)
}

// SetMetadataKey implements the RawLog interface.
func (sL *spanLogImpl) SetMetadataKey(_ context.Context, k string, v any) {
	sL.m.Lock()
	defer sL.m.Unlock()
	sL.metadata[k] = v
}

// SetErrorMetadataKey implements the RawLog interface.
func (sL *spanLogImpl) SetErrorMetadataKey(_ context.Context, k string, v any) {
	sL.m.Lock()
	defer sL.m.Unlock()
	sL.errMetadata[k] = v
}

// SetErrorFlag implements the RawLog interface.
func (sL *spanLogImpl) SetErrorFlag(_ context.Context) {
	sL.m.Lock()
	defer sL.m.Unlock()
	sL.hasErrorFlag = true
}

// GetCurrentTraceLink implements the RawLog interface.
func (sL *spanLogImpl) GetCurrentTraceLink(_ context.Context) *TraceLink {
	sL.m.Lock()
	defer sL.m.Unlock()

	return &TraceLink{
		TraceID: sL.traceID,
		SpanID:  sL.spanID,
	}
}

// Flush implements the RawLog interface.
func (sL *spanLogImpl) Flush(_ context.Context) {
	// do nothing: flushing is only enabled on the background log
}
