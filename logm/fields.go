package logm

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/jsonz"
	"github.com/ibrt/golang-utils/memz"

	"github.com/ibrt/golang-modules/clkm"
)

// Some commonly used field keys.
const (
	StandardKeyParams          = "params"
	StandardKeySecondaryParams = "secondaryParams"
)

type errorMetadataKey int

const (
	isEmittedErrorMetadataKey errorMetadataKey = iota
)

var (
	traceLinkRegexp = regexp.MustCompile(`^([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})-([\da-f]{8}-[\da-f]{4}-[\da-f]{4}-[\da-f]{4}-[\da-f]{12})$`)
)

// TraceLink describes a link annotation.
type TraceLink struct {
	TraceID string
	SpanID  string
}

// Serialize the [*TraceLink].
func (l *TraceLink) Serialize() string {
	if l == nil {
		return ""
	}
	if l.TraceID != "" && l.SpanID != "" {
		return fmt.Sprintf("%v-%v", l.TraceID, l.SpanID)
	}
	return ""
}

// MaybeParseTraceLink parses a serialized [*TraceLink].
func MaybeParseTraceLink(traceLink string) *TraceLink {
	matches := traceLinkRegexp.FindStringSubmatch(traceLink)

	if len(matches) != 3 {
		return nil
	}

	return &TraceLink{
		TraceID: matches[1],
		SpanID:  matches[2],
	}
}

// User describes a user.
type User struct {
	ID    string
	Email string
}

type fieldValueLen interface {
	string | map[string]any | []string | EmitMetadata
}

type fieldValueNumeric interface {
	int | float64
}

type addField interface {
	AddField(string, any)
}

type newEvent interface {
	NewEvent() *libhoney.Event
}

func maybeAddLenField[T fieldValueLen](af addField, prefix, key string, value T) {
	if len(value) == 0 {
		return
	}

	if prefix != "" {
		if !strings.HasSuffix(prefix, ".") {
			prefix += "."
		}
	}

	af.AddField(prefix+key, value)
}

func maybeAddNumericField[T fieldValueNumeric](af addField, key string, value T) {
	if value != 0 {
		af.AddField(key, value)
	}
}

func addLocationFields(af addField, framesSource error) {
	frames := memz.FilterSlice(errorz.GetFrames(framesSource), func(f *errorz.Frame) bool {
		return f.ShortPackage != "logm"
	})

	if len(frames) > 0 {
		maybeAddLenField(af, "", "location", frames[0].Summary)
		maybeAddLenField(af, "", "location.short", frames[0].ShortLocation)
		maybeAddLenField(af, "", "location.frames", memz.TransformSlice(frames, func(_ int, f *errorz.Frame) string { return f.Summary }))
	}
}

func maybeAddSpanEventAnnotationFields(af addField, spanID string) {
	if spanID != "" {
		af.AddField("meta.annotation_type", "span_event")
		af.AddField("trace.parent_id", spanID)
	}
}

func maybeAddUserFields(af addField, user *User) {
	if user != nil {
		maybeAddLenField(af, "", "scope.user", user.ID)
		maybeAddLenField(af, "", "scope.user.email", user.Email)
	}
}

func addMetadataFields(af addField, prefix string, metadata map[string]any) {
	if prefix != "" {
		if !strings.HasSuffix(prefix, ".") {
			prefix += "."
		}
	}

	for k, v := range metadata {
		if k == StandardKeyParams || k == StandardKeySecondaryParams {
			k = StandardKeyParams

			if m, ok := maybeFlattenMetadataValue(v); ok {
				for kk, vv := range m {
					af.AddField(prefix+k+"."+kk, vv)
				}
				continue
			}
		}

		af.AddField(prefix+k, v)
	}
}

func maybeFlattenMetadataValue(v any) (map[string]any, bool) {
	m, err := errorz.Catch1(func() (map[string]any, error) {
		if reflect.Indirect(reflect.ValueOf(v)).Type().Kind() != reflect.Struct {
			return nil, errorz.Errorf("not a struct")
		}

		return jsonz.Unmarshal[map[string]any](jsonz.MustMarshal(v))
	})
	if err != nil || m == nil {
		return nil, false
	}

	return m, true
}

func addDebugFields(af addField, format string, o *emitOptions) {
	addLocationFields(af, nil)
	af.AddField("debug", true)
	maybeAddLenField(af, "", "debug.message", fmt.Sprintf(format, o.args...))
	addMetadataFields(af, "debug.metadata", o.metadata)
}

func addInfoFields(af addField, format string, o *emitOptions) {
	addLocationFields(af, nil)
	af.AddField("info", true)
	maybeAddLenField(af, "", "info.message", fmt.Sprintf(format, o.args...))
	addMetadataFields(af, "info.metadata", o.metadata)
}

func addWarningFields(af addField, err error) {
	addLocationFields(af, err)
	af.AddField("warning", getWarningName(err))
	af.AddField("warning.message", err.Error())
	af.AddField("warning.dump", errorz.SDump(err))
	maybeAddNumericField(af, "warning.status", errorz.GetHTTPStatus(err, 0))
}

func addErrorFields(af addField, err error) {
	addLocationFields(af, err)
	af.AddField("error", getErrorName(err))
	af.AddField("error.message", err.Error())
	af.AddField("error.dump", errorz.SDump(err))
	maybeAddNumericField(af, "error.status", errorz.GetHTTPStatus(err, 0))
}

func newAttachableEvent(ctx context.Context, ne newEvent, attachedSpanID, name string) *libhoney.Event {
	e := ne.NewEvent()
	e.Timestamp = clkm.MustGet(ctx).Now()
	maybeAddSpanEventAnnotationFields(e, attachedSpanID)
	e.AddField("name", name)
	return e
}

func newTraceLinkEvent(ctx context.Context, ne newEvent, attachedSpanID string, traceLink *TraceLink) *libhoney.Event {
	e := ne.NewEvent()
	e.Timestamp = clkm.MustGet(ctx).Now()
	e.AddField("meta.annotation_type", "link")
	e.AddField("trace.parent_id", attachedSpanID)
	e.AddField("trace.link.trace_id", traceLink.TraceID)
	e.AddField("trace.link.span_id", traceLink.SpanID)
	e.AddField("name", "link-annotation")
	return e
}

func newTraceableEvent(ctx context.Context, ne newEvent, name, spanID, parentSpanID string, startTime time.Time) *libhoney.Event {
	e := ne.NewEvent()
	e.Timestamp = startTime
	addLocationFields(e, nil)
	e.AddField("duration_ms", float64(clkm.MustGet(ctx).Now().Sub(startTime))/float64(time.Millisecond))
	e.AddField("name", name)
	e.AddField("trace.span_id", spanID)
	maybeAddLenField(e, "", "trace.parent_id", parentSpanID)
	return e
}

func getWarningName(err error) string {
	return errorz.GetName(err, "generic")
}

func getErrorName(err error) string {
	return errorz.GetName(err, "generic")
}

func maybeSetIsEmitted(err error) {
	errorz.MaybeSetMetadata(err, isEmittedErrorMetadataKey, true)
}

func getIsEmitted(err error) bool {
	isEmitted, ok := errorz.MaybeGetMetadata[bool](err, isEmittedErrorMetadataKey)
	return ok && isEmitted
}
