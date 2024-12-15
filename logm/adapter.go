package logm

import (
	"context"
)

var (
	_ Log = (*adapterLogImpl)(nil)
)

type adapterLogImpl struct {
	ctx context.Context
	log RawLog
}

// EmitDebug implements the Log interface.
func (l *adapterLogImpl) EmitDebug(format string, options ...EmitOption) {
	l.log.EmitDebug(l.ctx, format, options...)
}

// EmitInfo implements the Log interface.
func (l *adapterLogImpl) EmitInfo(format string, options ...EmitOption) {
	l.log.EmitInfo(l.ctx, format, options...)
}

// EmitWarning implements the Log interface.
func (l *adapterLogImpl) EmitWarning(err error) {
	l.log.EmitWarning(l.ctx, err)
}

// EmitError implements the Log interface.
func (l *adapterLogImpl) EmitError(err error) {
	l.log.EmitError(l.ctx, err)
}

// EmitTraceLink implements the Log interface.
func (l *adapterLogImpl) EmitTraceLink(link *TraceLink) {
	l.log.EmitTraceLink(l.ctx, link)
}

// Begin implements the Log interface.
func (l *adapterLogImpl) Begin(name string, options ...BeginOption) (context.Context, func()) {
	return l.log.Begin(l.ctx, name, options...)
}

// SetErrorFlag implements the Log interface.
func (l *adapterLogImpl) SetErrorFlag() {
	l.log.SetErrorFlag(l.ctx)
}

// SetUser implements the Log interface.
func (l *adapterLogImpl) SetUser(user *User) {
	l.log.SetUser(l.ctx, user)
}

// SetPropagatingField implements the Log interface.
func (l *adapterLogImpl) SetPropagatingField(k string, v any) {
	l.log.SetPropagatingField(l.ctx, k, v)
}

// SetMetadataKey implements the Log interface.
func (l *adapterLogImpl) SetMetadataKey(k string, v any) {
	l.log.SetMetadataKey(l.ctx, k, v)
}

// SetErrorMetadataKey implements the Log interface.
func (l *adapterLogImpl) SetErrorMetadataKey(k string, v any) {
	l.log.SetErrorMetadataKey(l.ctx, k, v)
}

// GetCurrentTraceLink implements the Log interface.
func (l *adapterLogImpl) GetCurrentTraceLink() *TraceLink {
	return l.log.GetCurrentTraceLink(l.ctx)
}

// Flush implements the RawLog interface.
func (l *adapterLogImpl) Flush() {
	l.log.Flush(l.ctx)
}
