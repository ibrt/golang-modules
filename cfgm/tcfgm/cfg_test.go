package tcfgm_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/cfgm/tcfgm"
)

var (
	_ cfgm.Config = (*TestConfig)(nil)
)

type TestConfig struct {
	Key string
}

func (*TestConfig) Config() {
	// intentionally empty
}

type Suite struct {
	CFG *tcfgm.Helper[*TestConfig]
}

func TestSuite(t *testing.T) {
	fixturez.RunSuite(t, &Suite{
		CFG: &tcfgm.Helper[*TestConfig]{
			ConfigLoader: func(ctx context.Context) (*TestConfig, error) {
				return &TestConfig{Key: "Value"}, nil
			},
		},
	})
}

func (s *Suite) TestMustGet(ctx context.Context, g *WithT) {
	g.Expect(cfgm.MustGet[*TestConfig](ctx)).To(Equal(&TestConfig{
		Key: "Value",
	}))
}
