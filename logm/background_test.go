package logm_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
	"github.com/ibrt/golang-modules/logm"
	"github.com/ibrt/golang-modules/logm/tlogm"
)

type BackgroundSuite struct {
	CLK *tclkm.MockHelper
	LOG *tlogm.MockHelper
}

func TestBackgroundSuite(t *testing.T) {
	fixturez.RunSuite(t, &BackgroundSuite{})
}

func (s *BackgroundSuite) TestBackgroundEmitDebug(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitDebug("emit: background: debug: %v",
		logm.EmitA("arg"),
		logm.EmitM("ebdk", "ebdv"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("debug", true),
					HaveKeyWithValue("debug.message", "emit: background: debug: arg"),
					HaveKeyWithValue("debug.metadata.ebdk", "ebdv"),
					HaveKeyWithValue("name", "debug"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundEmitInfo(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitInfo("emit: background: info: %v",
		logm.EmitA("arg"),
		logm.EmitM("ebik", "ebiv"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("info", true),
					HaveKeyWithValue("info.message", "emit: background: info: arg"),
					HaveKeyWithValue("info.metadata.ebik", "ebiv"),
					HaveKeyWithValue("name", "info"),
				),
			}))))
}

/*
func (s *BackgroundSuite) TestBackgroundEmitWarning(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitWarning(
		errorz.Errorf("emit: background: warning",
			errorz.ID("emit-background-warning"),
			errorz.Status(http.StatusBadRequest),
			errorz.M("ebwk", "ebwv")))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning-emit-background-warning"),
					HaveKeyWithValue("warning.message", "emit: background: warning"),
					HaveKeyWithValue("warning.id", errorz.ID("emit-background-warning")),
					HaveKeyWithValue("warning.status", errorz.Status(http.StatusBadRequest)),
					HaveKeyWithValue("warning.metadata.ebwk", "ebwv"),
					HaveKeyWithValue("name", "warning-emit-background-warning"),
				),
			}))))

	s.LOG.GetMock().ClearEvents()
	logm.MustGet(ctx).EmitWarning(fmt.Errorf("emit: background: warning"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "emit: background: warning"),
					Not(HaveKey("warning.status")),
					Not(HaveKey("warning.metadata")),
					HaveKeyWithValue("name", "warning"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundEmitError(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitError(
		errorz.Errorf("emit: background: error",
			errorz.ID("emit-background-error"),
			errorz.Status(http.StatusBadRequest),
			errorz.M("ebek", "ebev")))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("error", "error-emit-background-error"),
					HaveKeyWithValue("error.message", "emit: background: error"),
					HaveKeyWithValue("error.id", errorz.ID("emit-background-error")),
					HaveKeyWithValue("error.status", errorz.Status(http.StatusBadRequest)),
					HaveKeyWithValue("error.metadata.ebek", "ebev"),
					HaveKeyWithValue("name", "error-emit-background-error"),
				),
			}))))

	s.LOG.GetMock().ClearEvents()
	logm.MustGet(ctx).EmitError(fmt.Errorf("emit: background: error"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("error", "error"),
					HaveKeyWithValue("error.message", "emit: background: error"),
					Not(HaveKey("error.status")),
					Not(HaveKey("error.metadata")),
					HaveKeyWithValue("name", "error"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundEmitTraceLink(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitTraceLink(&logm.TraceLink{
		TraceID: idz.MustNewRandomUUID(),
		SpanID:  idz.MustNewRandomUUID(),
	})

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called EmitTraceLink in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundSetUser(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).SetUser(&logm.User{
		ID: idz.MustNewRandomUUID(),
	})

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called SetUser in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundSetPropagatingField(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).SetPropagatingField("k", "v")

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called SetPropagatingField in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundSetMetadataKey(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).SetMetadataKey("k", "v")

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called SetMetadataKey in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundSetErrorMetadataKey(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).SetErrorMetadataKey("k", "v")

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called SetErrorMetadataKey in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundSetErrorFlag(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).SetErrorFlag()

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "warning"),
					HaveKeyWithValue("warning.message", "called SetErrorFlag in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundGetCurrentTraceLink(ctx context.Context, g *WithT) {
	g.Expect(logm.MustGet(ctx).GetCurrentTraceLink()).To(BeNil())
}
*/
