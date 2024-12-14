package cfgm_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/cfgm/tcfgm"
)

type testConfig struct {
	Key string `env:"KEY_EC2B754B"`
}

// Config implements the [cfgm.Config] interface.
func (cfg testConfig) Config() {
	// intentionally empty
}

type Suite struct {
	CFG *tcfgm.Helper[*testConfig]
}

func TestSuite(t *testing.T) {
	fixturez.RunSuite(t, &Suite{
		CFG: &tcfgm.Helper[*testConfig]{
			ConfigLoader: func(ctx context.Context) (*testConfig, error) {
				return &testConfig{Key: "Value"}, nil
			},
		},
	})
}

func (s *Suite) TestMustGet(ctx context.Context, g *WithT) {
	g.Expect(cfgm.MustGet[*testConfig](ctx)).To(Equal(&testConfig{
		Key: "Value",
	}))
}
