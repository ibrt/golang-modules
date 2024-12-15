package tclkm_test

import (
	"context"
	"testing"
	"time"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
)

type RealSuite struct {
	CLK *tclkm.RealHelper
}

func TestRealSuite(t *testing.T) {
	fixturez.RunSuite(t, &RealSuite{})
}

func (s *RealSuite) TestRealHelper(ctx context.Context, g *WithT) {
	g.Expect(clkm.MustGet(ctx)).NotTo(BeNil())
	g.Expect(clkm.MustGet(ctx)).NotTo(BeZero())
}

type MockSuite struct {
	CLK *tclkm.MockHelper
}

func TestMockSuite(t *testing.T) {
	fixturez.RunSuite(t, &MockSuite{})
}

func (s *MockSuite) TestMockHelper(ctx context.Context, g *WithT) {
	now := time.Now().Add(-time.Minute)
	s.CLK.GetMock().Set(now)
	g.Expect(clkm.MustGet(ctx).Now()).To(Equal(now))
}
