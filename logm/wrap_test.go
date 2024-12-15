package logm_test

import (
	"context"
	"testing"
	"time"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
	"github.com/ibrt/golang-modules/logm"
	"github.com/ibrt/golang-modules/logm/tlogm"
)

type WrapSuite struct {
	Clock *tclkm.MockHelper
	Log   *tlogm.MockHelper
}

func TestWrapSuite(t *testing.T) {
	fixturez.RunSuite(t, &WrapSuite{})
}

func (s *WrapSuite) TestWrap0_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(logm.Wrap0(ctx, "wrap0",
		func(ctx context.Context) error {
			logm.MustGet(ctx).SetMetadataKey("k01", "v01")

			logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
				logm.EmitA("arg"),
				logm.EmitMetadata{"k02": "v02"})

			s.Clock.GetMock().Add(time.Second)
			return nil
		},
		logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))).
		To(Succeed())

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap0"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap0_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(logm.Wrap0(ctx, "wrap0",
		func(ctx context.Context) error {
			s.Clock.GetMock().Add(time.Second)
			return errorz.Errorf("test error")
		},
		logm.BeginErrM("kew", "vew"),
		logm.BeginErrM(logm.StandardKeyParams,
			struct {
				K string `json:"kews"`
			}{
				K: "vews",
			}),
		logm.BeginErrM(logm.StandardKeySecondaryParams,
			struct {
				K string `json:"kewss"`
			}{
				K: "vewss",
			}))).
		To(MatchError("test error"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap0"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKeyWithValue("scope.metadata.params.kews", "vews"),
					HaveKeyWithValue("scope.metadata.params.kewss", "vewss"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap0_Error_Emitted(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		logm.Wrap0(ctx, "wrap0_outer",
			func(ctx context.Context) error {
				s.Clock.GetMock().Add(time.Second)
				return logm.Wrap0(ctx, "wrap0_inner",
					func(ctx context.Context) error {
						s.Clock.GetMock().Add(time.Second)
						return errorz.Errorf("test error")
					},
					logm.BeginErrM("kew", "vew"),
					logm.BeginErrM(logm.StandardKeyParams,
						struct {
							K string `json:"kews"`
						}{
							K: "vews",
						}))
			})).
		To(MatchError("test error"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(2 * time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap0_inner"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKeyWithValue("scope.metadata.params.kews", "vews"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					HaveKey("trace.parent_id"),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap0_outer"),
					HaveKeyWithValue("duration_ms", float64(2000)),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap0_Panic(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(logm.Wrap0(ctx, "wrap0",
		func(ctx context.Context) error {
			s.Clock.GetMock().Add(time.Second)
			panic(errorz.Errorf("test error"))
		})).
		To(MatchError("test error"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap0"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap0Panic_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			logm.Wrap0Panic(ctx, "wrap0",
				func(ctx context.Context) error {
					logm.MustGet(ctx).SetMetadataKey("k01", "v01")

					logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
						logm.EmitA("arg"),
						logm.EmitMetadata{"k02": "v02"})

					s.Clock.GetMock().Add(time.Second)
					return nil
				},
				logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))
		}).
		ToNot(Panic())

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap0"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap0Panic_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			logm.Wrap0Panic(ctx, "wrap0",
				func(ctx context.Context) error {
					s.Clock.GetMock().Add(time.Second)
					return errorz.Errorf("test error")
				},
				logm.BeginErrM("kew", "vew"),
				logm.BeginErrM(logm.StandardKeyParams,
					struct {
						K string `json:"kews"`
					}{
						K: "vews",
					}))
		}).
		To(PanicWith(MatchError("test error")))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap0"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKeyWithValue("scope.metadata.params.kews", "vews"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap1_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out, err := logm.Wrap1(ctx, "wrap1",
		func(ctx context.Context) (string, error) {
			logm.MustGet(ctx).SetMetadataKey("k01", "v01")

			logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
				logm.EmitA("arg"),
				logm.EmitMetadata{"k02": "v02"})

			s.Clock.GetMock().Add(time.Second)
			return "test", nil
		},
		logm.BeginMetadata{"kw": "vw"}, logm.BeginErrMetadata{"kew": "vew"})
	g.Expect(err).To(Succeed())
	g.Expect(out).To(Equal("test"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap1"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap1_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out, err := logm.Wrap1(ctx, "wrap1",
		func(ctx context.Context) (string, error) {
			s.Clock.GetMock().Add(time.Second)
			return "", errorz.Errorf("test error")
		}, logm.BeginErrMetadata{"kew": "vew"})
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap1"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap1_Panic(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out, err := logm.Wrap1(ctx, "wrap1",
		func(ctx context.Context) (string, error) {
			s.Clock.GetMock().Add(time.Second)
			panic(errorz.Errorf("test error"))
		})
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap1"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap1Panic_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			out := logm.Wrap1Panic(ctx, "wrap1",
				func(ctx context.Context) (string, error) {
					logm.MustGet(ctx).SetMetadataKey("k01", "v01")

					logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
						logm.EmitA("arg"),
						logm.EmitMetadata{"k02": "v02"})

					s.Clock.GetMock().Add(time.Second)
					return "test", nil
				},
				logm.BeginMetadata{"kw": "vw"}, logm.BeginErrMetadata{"kew": "vew"})
			g.Expect(out).To(Equal("test"))
		}).
		ToNot(Panic())

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap1"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap1Panic_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			logm.Wrap1Panic(ctx, "wrap1",
				func(ctx context.Context) (string, error) {
					s.Clock.GetMock().Add(time.Second)
					return "", errorz.Errorf("test error")
				}, logm.BeginErrMetadata{"kew": "vew"})
		}).
		To(PanicWith(MatchError("test error")))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap1"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap2_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, err := logm.Wrap2(ctx, "wrap2",
		func(ctx context.Context) (string, string, error) {
			logm.MustGet(ctx).SetMetadataKey("k01", "v01")

			logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
				logm.EmitA("arg"),
				logm.EmitMetadata{"k02": "v02"})

			s.Clock.GetMock().Add(time.Second)
			return "t1", "t2", nil
		},
		logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))
	g.Expect(err).To(Succeed())
	g.Expect(out1).To(Equal("t1"))
	g.Expect(out2).To(Equal("t2"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap2"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap2_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, err := logm.Wrap2(ctx, "wrap2",
		func(ctx context.Context) (string, string, error) {
			s.Clock.GetMock().Add(time.Second)
			return "", "", errorz.Errorf("test error")
		}, logm.BeginErrM("kew", "vew"), logm.BeginErrM(logm.StandardKeyParams, "p"))
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out1).To(Equal(""))
	g.Expect(out2).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap2"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKeyWithValue("scope.metadata.params", "p"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap2_Panic(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, err := logm.Wrap2(ctx, "wrap2",
		func(ctx context.Context) (string, string, error) {
			s.Clock.GetMock().Add(time.Second)
			panic(errorz.Errorf("test error"))
		})
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out1).To(Equal(""))
	g.Expect(out2).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap2"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap2Panic_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			out1, out2 := logm.Wrap2Panic(ctx, "wrap2",
				func(ctx context.Context) (string, string, error) {
					logm.MustGet(ctx).SetMetadataKey("k01", "v01")

					logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
						logm.EmitA("arg"),
						logm.EmitMetadata{"k02": "v02"})

					s.Clock.GetMock().Add(time.Second)
					return "t1", "t2", nil
				},
				logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))
			g.Expect(out1).To(Equal("t1"))
			g.Expect(out2).To(Equal("t2"))
		}).
		ToNot(Panic())

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap2"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap2Panic_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			logm.Wrap2Panic(ctx, "wrap2",
				func(ctx context.Context) (string, string, error) {
					s.Clock.GetMock().Add(time.Second)
					return "", "", errorz.Errorf("test error")
				}, logm.BeginErrM("kew", "vew"), logm.BeginErrM(logm.StandardKeyParams, "p"))
		}).
		To(PanicWith(MatchError("test error")))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap2"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKeyWithValue("scope.metadata.params", "p"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap3_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, out3, err := logm.Wrap3(ctx, "wrap3",
		func(ctx context.Context) (string, string, string, error) {
			logm.MustGet(ctx).SetMetadataKey("k01", "v01")

			logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
				logm.EmitA("arg"),
				logm.EmitMetadata{"k02": "v02"})

			s.Clock.GetMock().Add(time.Second)
			return "t1", "t2", "t3", nil
		},
		logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))
	g.Expect(err).To(Succeed())
	g.Expect(out1).To(Equal("t1"))
	g.Expect(out2).To(Equal("t2"))
	g.Expect(out3).To(Equal("t3"))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap3"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap3_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, out3, err := logm.Wrap3(ctx, "wrap3",
		func(ctx context.Context) (string, string, string, error) {
			logm.MustGet(ctx).SetErrorMetadataKey("kew", "vew")
			s.Clock.GetMock().Add(time.Second)
			return "", "", "", errorz.Errorf("test error")
		})
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out1).To(Equal(""))
	g.Expect(out2).To(Equal(""))
	g.Expect(out3).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap3"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap3_Panic(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	out1, out2, out3, err := logm.Wrap3(ctx, "wrap3",
		func(ctx context.Context) (string, string, string, error) {
			s.Clock.GetMock().Add(time.Second)
			panic(errorz.Errorf("test error"))
		})
	g.Expect(err).To(MatchError("test error"))
	g.Expect(out1).To(Equal(""))
	g.Expect(out2).To(Equal(""))
	g.Expect(out3).To(Equal(""))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap3"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap3Panic_Success(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			out1, out2, out3 := logm.Wrap3Panic(ctx, "wrap3",
				func(ctx context.Context) (string, string, string, error) {
					logm.MustGet(ctx).SetMetadataKey("k01", "v01")

					logm.MustGet(ctx).EmitInfo("emit: span: info: %v",
						logm.EmitA("arg"),
						logm.EmitMetadata{"k02": "v02"})

					s.Clock.GetMock().Add(time.Second)
					return "t1", "t2", "t3", nil
				},
				logm.BeginM("kw", "vw"), logm.BeginErrM("kew", "vew"))
			g.Expect(out1).To(Equal("t1"))
			g.Expect(out2).To(Equal("t2"))
			g.Expect(out3).To(Equal("t3"))
		}).
		ToNot(Panic())

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: span: info: arg"),
					HaveKeyWithValue("info.metadata.k02", "v02"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "info"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("name", "wrap3"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.k01", "v01"),
					HaveKeyWithValue("scope.metadata.kw", "vw"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
					Not(HaveKey("scope.metadata.kew")),
				),
			})),
		))
}

func (s *WrapSuite) TestWrap3Panic_Error(ctx context.Context, g *WithT) {
	t1 := clkm.MustGet(ctx).Now()

	g.Expect(
		func() {
			logm.Wrap3Panic(ctx, "wrap3",
				func(ctx context.Context) (string, string, string, error) {
					logm.MustGet(ctx).SetErrorMetadataKey("kew", "vew")
					s.Clock.GetMock().Add(time.Second)
					return "", "", "", errorz.Errorf("test error")
				})
		}).
		To(PanicWith(MatchError("test error")))

	g.Expect(s.Log.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1.Add(time.Second)),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "test error"),
					HaveKeyWithValue("meta.annotation_type", "span_event"),
					HaveKeyWithValue("name", "error"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.parent_id"),
					Not(HaveKey("trace.span_id")),
				),
			})),
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(t1),
				"Data": And(
					HaveKeyWithValue("error", true),
					HaveKeyWithValue("name", "wrap3"),
					HaveKeyWithValue("duration_ms", float64(1000)),
					HaveKeyWithValue("scope.metadata.kew", "vew"),
					HaveKey("trace.trace_id"),
					HaveKey("trace.span_id"),
					Not(HaveKey("trace.parent_id")),
				),
			})),
		))
}
