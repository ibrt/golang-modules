package logm

import (
	"context"
)

var (
	_ Log = (*adapterLogImpl)(nil)
)

type adapterLogImpl struct {
	ctx    context.Context
	rawLog RawLog
}

// EmitDebug implements the Log interface.
func (l *adapterLogImpl) EmitDebug(format string, options ...EmitOption) {
	l.rawLog.EmitDebug(l.ctx, format, options...)
}

// EmitInfo implements the Log interface.
func (l *adapterLogImpl) EmitInfo(format string, options ...EmitOption) {
	l.rawLog.EmitInfo(l.ctx, format, options...)
}

// EmitWarning implements the Log interface.
func (l *adapterLogImpl) EmitWarning(err error) {
	l.rawLog.EmitWarning(l.ctx, err)
}

// EmitError implements the Log interface.
func (l *adapterLogImpl) EmitError(err error) {
	l.rawLog.EmitError(l.ctx, err)
}

// EmitTraceLink implements the Log interface.
func (l *adapterLogImpl) EmitTraceLink(link *TraceLink) {
	l.rawLog.EmitTraceLink(l.ctx, link)
}

// Begin implements the Log interface.
func (l *adapterLogImpl) Begin(name string, options ...BeginOption) (context.Context, func()) {
	return l.rawLog.Begin(l.ctx, name, options...)
}

// SetErrorFlag implements the Log interface.
func (l *adapterLogImpl) SetErrorFlag() {
	l.rawLog.SetErrorFlag(l.ctx)
}

// SetUser implements the Log interface.
func (l *adapterLogImpl) SetUser(user *User) {
	l.rawLog.SetUser(l.ctx, user)
}

// SetPropagatingField implements the Log interface.
func (l *adapterLogImpl) SetPropagatingField(k string, v any) {
	l.rawLog.SetPropagatingField(l.ctx, k, v)
}

// SetMetadataKey implements the Log interface.
func (l *adapterLogImpl) SetMetadataKey(k string, v any) {
	l.rawLog.SetMetadataKey(l.ctx, k, v)
}

// SetErrorMetadataKey implements the Log interface.
func (l *adapterLogImpl) SetErrorMetadataKey(k string, v any) {
	l.rawLog.SetErrorMetadataKey(l.ctx, k, v)
}

// GetCurrentTraceLink implements the Log interface.
func (l *adapterLogImpl) GetCurrentTraceLink() *TraceLink {
	return l.rawLog.GetCurrentTraceLink(l.ctx)
}

// Flush implements the RawLog interface.
func (l *adapterLogImpl) Flush() {
	l.rawLog.Flush(l.ctx)
}
