package logm

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/honeycombio/libhoney-go"
	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/idz"
	"github.com/ibrt/golang-utils/memz"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
)

var (
	_ AddField               = (*testAddField)(nil)
	_ newEvent               = (*testNewEvent)(nil)
	_ error                  = (*testCompleteError)(nil)
	_ errorz.ErrorName       = (*testCompleteError)(nil)
	_ errorz.ErrorHTTPStatus = (*testCompleteError)(nil)
)

type testAddField struct {
	fields map[string]any
}

func newTestAddField() *testAddField {
	return &testAddField{
		fields: make(map[string]any),
	}
}

func (t *testAddField) AddField(k string, v any) {
	t.fields[k] = v
}

type testNewEvent struct {
	events []*libhoney.Event
}

func newTestNewEvent() *testNewEvent {
	return &testNewEvent{
		events: make([]*libhoney.Event, 0),
	}
}

func (t *testNewEvent) NewEvent() *libhoney.Event {
	e := libhoney.NewEvent()
	t.events = append(t.events, e)
	return e
}

type testCompleteError struct {
	message string
	name    string
	status  int
}

func newTestCompleteError(message, name string, status int) *testCompleteError {
	return &testCompleteError{
		message: message,
		name:    name,
		status:  status,
	}
}

func (e *testCompleteError) Error() string {
	return e.message
}

func (e *testCompleteError) GetErrorName() string {
	return e.name
}

func (e *testCompleteError) GetErrorHTTPStatus() int {
	return e.status
}

type FieldsSuite struct {
	CLK *tclkm.MockHelper
}

func TestFieldsSuite(t *testing.T) {
	fixturez.RunSuite(t, &FieldsSuite{})
}

func (*FieldsSuite) TestTraceLink(g *WithT) {
	tID := idz.MustNewRandomUUID()
	sID := idz.MustNewRandomUUID()

	g.Expect(MaybeParseTraceLink(memz.Ptr(TraceLink{
		// intentionally empty
	}).Serialize())).To(BeNil())

	g.Expect(MaybeParseTraceLink(memz.Ptr(TraceLink{
		TraceID: tID,
	}).Serialize())).To(BeNil())

	g.Expect(MaybeParseTraceLink(memz.Ptr(TraceLink{
		SpanID: sID,
	}).Serialize())).To(BeNil())

	g.Expect(MaybeParseTraceLink(memz.Ptr(TraceLink{
		TraceID: tID,
		SpanID:  sID,
	}).Serialize())).To(Equal(&TraceLink{
		TraceID: tID,
		SpanID:  sID,
	}))

	g.Expect(MaybeParseTraceLink(tID)).To(BeNil())
	g.Expect(MaybeParseTraceLink(fmt.Sprintf("%v-", tID))).To(BeNil())
	g.Expect((*TraceLink)(nil).Serialize()).To(Equal(""))
}

func (*FieldsSuite) TestMaybeAddLenField(g *WithT) {
	{
		af := newTestAddField()
		maybeAddLenField(af, "prefix", "key", "")
		g.Expect(af.fields).To(BeEmpty())
	}
	{
		af := newTestAddField()
		maybeAddLenField(af, "prefix", "key", "value")
		g.Expect(af.fields).To(Equal(map[string]any{"prefix.key": "value"}))
	}
	{
		af := newTestAddField()
		maybeAddLenField(af, "", "key", "value")
		g.Expect(af.fields).To(Equal(map[string]any{"key": "value"}))
	}
}

func (*FieldsSuite) TestMaybeAddNumericField(g *WithT) {
	{
		af := newTestAddField()
		maybeAddNumericField(af, "key", 0)
		g.Expect(af.fields).To(BeEmpty())
	}
	{
		af := newTestAddField()
		maybeAddNumericField(af, "key", 1)
		g.Expect(af.fields).To(Equal(map[string]any{"key": 1}))
	}
}

func (*FieldsSuite) TestAddLocationFields(g *WithT) {
	err := errorz.Errorf("test error")

	frames := memz.FilterSlice(errorz.GetFrames(err), func(f *errorz.Frame) bool {
		return f.ShortPackage != "logm"
	})

	af := newTestAddField()
	addLocationFields(af, err)

	g.Expect(af.fields).To(And(
		HaveKeyWithValue("location", frames[0].Summary),
		HaveKeyWithValue("location.short", frames[0].ShortLocation),
		HaveKeyWithValue("location.frames", Not(BeEmpty())),
	))
}

func (*FieldsSuite) TestMaybeAddSpanEventAnnotationFields(g *WithT) {
	{
		af := newTestAddField()
		maybeAddSpanEventAnnotationFields(af, "")
		g.Expect(af.fields).To(BeEmpty())
	}
	{
		af := newTestAddField()
		maybeAddSpanEventAnnotationFields(af, "span-id")
		g.Expect(af.fields).To(Equal(map[string]any{
			"meta.annotation_type": "span_event",
			"trace.parent_id":      "span-id",
		}))
	}
}

func (*FieldsSuite) TestMaybeAddUserFields(g *WithT) {
	{
		af := newTestAddField()
		maybeAddUserFields(af, nil)
		g.Expect(af.fields).To(BeEmpty())
	}
	{
		af := newTestAddField()
		maybeAddUserFields(af, &User{
			ID:    "user-id",
			Email: "user-email",
		})
		g.Expect(af.fields).To(Equal(map[string]any{
			"scope.user":       "user-id",
			"scope.user.email": "user-email",
		}))
	}
}

func (*FieldsSuite) TestAddMetadataFields(g *WithT) {
	type testStruct struct {
		Key string `json:"key"`
	}

	{
		af := newTestAddField()
		addMetadataFields(af, "", map[string]any{"k1": "v1"})
		g.Expect(af.fields).To(Equal(map[string]any{
			"k1": "v1",
		}))
	}
	{
		af := newTestAddField()
		addMetadataFields(af, "prefix", map[string]any{"k1": "v1", "k2": &testStruct{Key: "v2"}})
		g.Expect(af.fields).To(Equal(map[string]any{
			"prefix.k1": "v1",
			"prefix.k2": &testStruct{Key: "v2"},
		}))
	}
	{
		af := newTestAddField()
		addMetadataFields(af, "prefix", map[string]any{"k1": "v1", StandardKeyParams: &testStruct{Key: "v2"}})
		g.Expect(af.fields).To(Equal(map[string]any{
			"prefix.k1":         "v1",
			"prefix.params.key": "v2",
		}))
	}
}

func (*FieldsSuite) TestMaybeFlattenMetadataValue(g *WithT) {
	type testStruct struct {
		Key string `json:"key"`
	}

	{
		m, ok := maybeFlattenMetadataValue(nil)
		g.Expect(m).To(BeNil())
		g.Expect(ok).To(BeFalse())
	}
	{
		m, ok := maybeFlattenMetadataValue(1)
		g.Expect(m).To(BeNil())
		g.Expect(ok).To(BeFalse())
	}
	{
		m, ok := maybeFlattenMetadataValue(map[string]any{})
		g.Expect(m).To(BeNil())
		g.Expect(ok).To(BeFalse())
	}
	{
		m, ok := maybeFlattenMetadataValue(testStruct{Key: "value"})
		g.Expect(m).To(Equal(map[string]any{"key": "value"}))
		g.Expect(ok).To(BeTrue())
	}
	{
		m, ok := maybeFlattenMetadataValue(&testStruct{Key: "value"})
		g.Expect(m).To(Equal(map[string]any{"key": "value"}))
		g.Expect(ok).To(BeTrue())
	}
}

func (*FieldsSuite) TestAddDebugFields(g *WithT) {
	af := newTestAddField()
	addDebugFields(af, "fmt: %v", newEmitOptions(EmitA(1), EmitM("k", "v")))

	g.Expect(af.fields).To(And(
		HaveKeyWithValue("location", Not(BeEmpty())),
		HaveKeyWithValue("debug", true),
		HaveKeyWithValue("debug.message", "fmt: 1"),
		HaveKeyWithValue("debug.metadata.k", "v"),
	))
}

func (*FieldsSuite) TestAddInfoFields(g *WithT) {
	af := newTestAddField()
	addInfoFields(af, "fmt: %v", newEmitOptions(EmitA(1), EmitM("k", "v")))

	g.Expect(af.fields).To(And(
		HaveKeyWithValue("location", Not(BeEmpty())),
		HaveKeyWithValue("info", true),
		HaveKeyWithValue("info.message", "fmt: 1"),
		HaveKeyWithValue("info.metadata.k", "v"),
	))
}

func (*FieldsSuite) TestAddWarningFields(g *WithT) {
	{
		af := newTestAddField()
		addWarningFields(af, errorz.Errorf("test error"))

		g.Expect(af.fields).To(And(
			HaveKeyWithValue("location", Not(BeEmpty())),
			HaveKeyWithValue("warning", "generic"),
			HaveKeyWithValue("warning.message", "test error"),
			HaveKeyWithValue("warning.dump", HavePrefix("(errorz.dump)")),
			Not(HaveKey("warning.status")),
		))
	}
	{
		af := newTestAddField()
		addWarningFields(af, newTestCompleteError("test error", "name", http.StatusBadRequest))

		g.Expect(af.fields).To(And(
			HaveKeyWithValue("location", Not(BeEmpty())),
			HaveKeyWithValue("warning", "name"),
			HaveKeyWithValue("warning.message", "test error"),
			HaveKeyWithValue("warning.dump", HavePrefix("(errorz.dump)")),
			HaveKeyWithValue("warning.status", http.StatusBadRequest),
		))
	}
}

func (*FieldsSuite) TestAddErrorFields(g *WithT) {
	{
		af := newTestAddField()
		addErrorFields(af, errorz.Errorf("test error"))

		g.Expect(af.fields).To(And(
			HaveKeyWithValue("location", Not(BeEmpty())),
			HaveKeyWithValue("error", "generic"),
			HaveKeyWithValue("error.message", "test error"),
			HaveKeyWithValue("error.dump", HavePrefix("(errorz.dump)")),
			Not(HaveKey("error.status")),
		))
	}
	{
		af := newTestAddField()
		addErrorFields(af, newTestCompleteError("test error", "name", http.StatusBadRequest))

		g.Expect(af.fields).To(And(
			HaveKeyWithValue("location", Not(BeEmpty())),
			HaveKeyWithValue("error", "name"),
			HaveKeyWithValue("error.message", "test error"),
			HaveKeyWithValue("error.dump", HavePrefix("(errorz.dump)")),
			HaveKeyWithValue("error.status", http.StatusBadRequest),
		))
	}
}

func (*FieldsSuite) TestNewAttachableEvent(ctx context.Context, g *WithT) {
	ne := newTestNewEvent()
	e := newAttachableEvent(ctx, ne, "attached-span-id", "name")

	g.Expect(e.Timestamp).To(Equal(clkm.MustGet(ctx).Now()))
	g.Expect(e.Fields()).To(Equal(map[string]any{
		"meta.annotation_type": "span_event",
		"trace.parent_id":      "attached-span-id",
		"name":                 "name",
	}))
}

func (*FieldsSuite) TestNewTraceLinkEvent(ctx context.Context, g *WithT) {
	ne := newTestNewEvent()

	e := newTraceLinkEvent(
		ctx, ne, "attached-span-id",
		&TraceLink{
			TraceID: "trace-id",
			SpanID:  "span-id",
		})

	g.Expect(e.Timestamp).To(Equal(clkm.MustGet(ctx).Now()))

	g.Expect(e.Fields()).To(Equal(map[string]any{
		"meta.annotation_type": "link",
		"trace.parent_id":      "attached-span-id",
		"trace.link.trace_id":  "trace-id",
		"trace.link.span_id":   "span-id",
		"name":                 "link-annotation",
	}))
}

func (*FieldsSuite) TestNewTraceableEvent(ctx context.Context, g *WithT) {
	ne := newTestNewEvent()
	startTime := clkm.MustGet(ctx).Now().Add(-time.Second)
	e := newTraceableEvent(ctx, ne, "name", "span-id", "parent-span-id", startTime)

	g.Expect(e.Timestamp).To(Equal(startTime))

	g.Expect(e.Fields()).To(And(
		HaveKeyWithValue("location", Not(BeEmpty())),
		HaveKeyWithValue("duration_ms", BeNumerically("==", 1000)),
		HaveKeyWithValue("name", "name"),
		HaveKeyWithValue("trace.span_id", "span-id"),
		HaveKeyWithValue("trace.parent_id", "parent-span-id"),
	))
}

func (*FieldsSuite) TestGetWarningName(g *WithT) {
	g.Expect(getWarningName(nil)).To(Equal("<nil>"))
	g.Expect(getWarningName(fmt.Errorf("test error"))).To(Equal("generic"))
	g.Expect(getWarningName(newTestCompleteError("", "name", 0))).To(Equal("name"))
}

func (*FieldsSuite) TestGetErrorName(g *WithT) {
	g.Expect(getErrorName(nil)).To(Equal("<nil>"))
	g.Expect(getErrorName(fmt.Errorf("test error"))).To(Equal("generic"))
	g.Expect(getErrorName(newTestCompleteError("", "name", 0))).To(Equal("name"))
}

func (*FieldsSuite) TestIsEmitted(g *WithT) {
	{
		err := fmt.Errorf("test error")
		g.Expect(getIsEmitted(err)).To(BeFalse())
		maybeSetIsEmitted(err)
		g.Expect(getIsEmitted(err)).To(BeFalse())
	}
	{
		err := errorz.Errorf("test error")
		g.Expect(getIsEmitted(err)).To(BeFalse())
		maybeSetIsEmitted(err)
		g.Expect(getIsEmitted(err)).To(BeTrue())
	}
}
