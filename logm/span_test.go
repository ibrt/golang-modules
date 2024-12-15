package logm_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
	"github.com/ibrt/golang-modules/logm"
	"github.com/ibrt/golang-modules/logm/tlogm"
)

type SpanSuite struct {
	CLK *tclkm.MockHelper
	LOG *tlogm.MockHelper
}

func TestSpanSuite(t *testing.T) {
	fixturez.RunSuite(t, &SpanSuite{})
}

func (s *SpanSuite) TestSpan(ctx context.Context, g *WithT) {
	// <OUTER SPAN>
	s.CLK.GetMock().Add(time.Second)
	t1 := clkm.MustGet(ctx).Now()

	ctx, e1 := logm.MustGet(ctx).Begin("S1")
	g.Expect(logm.MustGet(ctx).GetCurrentTraceLink().Serialize()).ToNot(BeEmpty())
	logm.MustGet(ctx).SetUser(&logm.User{ID: "user-id", Email: "user-email"})
	logm.MustGet(ctx).SetMetadataKey("k01", "v01")

	s.CLK.GetMock().Add(time.Second)
	t2 := clkm.MustGet(ctx).Now()

	logm.MustGet(ctx).EmitDebug("emit: span: debug: %v",
		logm.EmitA("arg"),
		logm.EmitMetadata{"k02": "v02"})

	logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
		logm.EmitA("arg"),
		logm.EmitMetadata{"k02": "v02"})

	// <MIDDLE SPAN>
	s.CLK.GetMock().Add(time.Second)
	t3 := clkm.MustGet(ctx).Now()

	ctx, e2 := logm.MustGet(ctx).Begin("S2")
	logm.MustGet(ctx).SetPropagatingField("kp", "vp")
	logm.MustGet(ctx).SetMetadataKey("k03", "v03")
	logm.MustGet(ctx).SetErrorFlag()

	s.CLK.GetMock().Add(time.Second)
	t4 := clkm.MustGet(ctx).Now()

	logm.MustGet(ctx).EmitWarning(
		newTestCompleteError("emit: span: warning", "emit-span-warning", http.StatusBadRequest))

	// <INNER SPAN>
	s.CLK.GetMock().Add(time.Second)
	t5 := clkm.MustGet(ctx).Now()

	ctx, e3 := logm.MustGet(ctx).Begin("S3")
	logm.MustGet(ctx).SetMetadataKey("k05", "v05")

	s.CLK.GetMock().Add(time.Second)
	t6 := clkm.MustGet(ctx).Now()

	logm.MustGet(ctx).EmitTraceLink(&logm.TraceLink{
		TraceID: "link-trace-id",
		SpanID:  "link-span-id"})

	logm.MustGet(ctx).EmitError(
		newTestCompleteError("emit: span: error", "emit-span-error", http.StatusBadRequest))

	s.CLK.GetMock().Add(time.Second)

	e3()
	// </INNER SPAN>

	s.CLK.GetMock().Add(time.Second)

	e2()
	// </MIDDLE SPAN>

	s.CLK.GetMock().Add(time.Second)

	e1()
	// </OUTER SPAN>

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t2),
				"Data": And(
					HaveKeyWithValue("debug", true),
					HaveKeyWithValue("debug.message", "emit: span: debug: arg"),
					HaveKeyWithValue("debug.metadata.k02", "v02"),
					HaveKeyWithValue("name", "debug"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t2),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("name", "info"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t4),
				"Data": And(
					HaveKeyWithValue("warning", "emit-span-warning"),
					HaveKeyWithValue("warning.message", "emit: span: warning"),
					HaveKeyWithValue("warning.dump", HavePrefix("(errorz.dump)")),
					HaveKeyWithValue("warning.status", http.StatusBadRequest),
					HaveKeyWithValue("name", "warning"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.kp", "vp"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t6),
				"Data": And(
					HaveKeyWithValue("name", "link-annotation"),
					HaveKeyWithValue("meta.annotation_type", "link"),
					HaveKeyWithValue("trace.link.trace_id", "link-trace-id"),
					HaveKeyWithValue("trace.link.span_id", "link-span-id"),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.kp", "vp"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t6),
				"Data": And(
					HaveKeyWithValue("error", "emit-span-error"),
					HaveKeyWithValue("error.message", "emit: span: error"),
					HaveKeyWithValue("error.dump", HavePrefix("(errorz.dump)")),
					HaveKeyWithValue("error.status", http.StatusBadRequest),
					HaveKeyWithValue("name", "error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.kp", "vp"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t5),
				"Data": And(
					HaveKeyWithValue("name", "S3"),
					HaveKeyWithValue("duration_ms", float64(2000)),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.metadata.k05", "v05"),
					HaveKeyWithValue("scope.kp", "vp"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					HaveKey("trace.span_id"),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t3),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "S2"),
					HaveKeyWithValue("duration_ms", float64(5000)),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.metadata.k03", "v03"),
					HaveKeyWithValue("scope.kp", "vp"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					HaveKey("trace.span_id"),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "S1"),
					HaveKeyWithValue("duration_ms", float64(8000)),
					HaveKeyWithValue("scope.user", "user-id"),
					HaveKeyWithValue("scope.user.email", "user-email"),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKey("trace.trace_id"),
					Not(HaveKey("trace.parent_id")),
					HaveKey("trace.span_id"),
				),
			})),
		))
}
