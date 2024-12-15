package logm_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	"github.com/ibrt/golang-utils/idz"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
	"github.com/ibrt/golang-modules/logm"
	"github.com/ibrt/golang-modules/logm/tlogm"
)

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

func (s *BackgroundSuite) TestBackgroundEmitWarning(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitWarning(
		newTestCompleteError(
			"emit: background: warning", "emit-background-warning", http.StatusBadRequest))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "emit-background-warning"),
					HaveKeyWithValue("warning.message", "emit: background: warning"),
					HaveKeyWithValue("warning.dump", HavePrefix("(errorz.dump)")),
					HaveKeyWithValue("warning.status", http.StatusBadRequest),
					HaveKeyWithValue("name", "warning"),
				),
			}))))

	s.LOG.GetMock().ClearEvents()
	logm.MustGet(ctx).EmitWarning(fmt.Errorf("emit: background: warning"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("warning", "generic"),
					HaveKeyWithValue("warning.message", "emit: background: warning"),
					HaveKeyWithValue("warning.dump", HavePrefix("(errorz.dump)")),
					Not(HaveKey("warning.status")),
					HaveKeyWithValue("name", "warning"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundEmitError(ctx context.Context, g *WithT) {
	logm.MustGet(ctx).EmitError(
		newTestCompleteError(
			"emit: background: error", "emit-background-error", http.StatusBadRequest))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("error", "emit-background-error"),
					HaveKeyWithValue("error.message", "emit: background: error"),
					HaveKeyWithValue("error.dump", HavePrefix("(errorz.dump)")),
					HaveKeyWithValue("error.status", http.StatusBadRequest),
					HaveKeyWithValue("name", "error"),
				),
			}))))

	s.LOG.GetMock().ClearEvents()
	logm.MustGet(ctx).EmitError(fmt.Errorf("emit: background: error"))

	g.Expect(s.LOG.GetMock().GetEvents()).
		To(HaveExactElements(
			PointTo(MatchFields(IgnoreExtras, Fields{
				"Timestamp": Equal(clkm.MustGet(ctx).Now()),
				"Data": And(
					HaveKeyWithValue("error", "generic"),
					HaveKeyWithValue("error.message", "emit: background: error"),
					HaveKeyWithValue("error.dump", HavePrefix("(errorz.dump)")),
					Not(HaveKey("error.status")),
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
					HaveKeyWithValue("warning", "generic"),
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
					HaveKeyWithValue("warning", "generic"),
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
					HaveKeyWithValue("warning", "generic"),
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
					HaveKeyWithValue("warning", "generic"),
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
					HaveKeyWithValue("warning", "generic"),
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
					HaveKeyWithValue("warning", "generic"),
					HaveKeyWithValue("warning.message", "called SetErrorFlag in background Log"),
				),
			}))))
}

func (s *BackgroundSuite) TestBackgroundGetCurrentTraceLink(ctx context.Context, g *WithT) {
	g.Expect(logm.MustGet(ctx).GetCurrentTraceLink()).To(BeNil())
}
