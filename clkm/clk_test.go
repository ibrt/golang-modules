package clkm_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/clkm"
	"github.com/ibrt/golang-modules/clkm/tclkm"
)

type Suite struct {
	CLK *tclkm.RealHelper
}

func TestSuite(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
}

func (s *Suite) TestMustGet(ctx context.Context, g *WithT) {
	g.Expect(clkm.MustGet(ctx)).NotTo(BeNil())
	g.Expect(clkm.MustGet(ctx)).NotTo(BeZero())
}
