package cfgm_test

import (
	"context"
	"os"
	"testing"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/cfgm/tcfgm"
)

var (
	_ cfgm.Config     = (*TestConfig)(nil)
	_ cfgm.Config     = (*TestConfigMixinImpl)(nil)
	_ TestConfigMixin = (*TestConfigMixinImpl)(nil)
)

type TestConfigMixin interface {
	cfgm.Config
	GetMixin() *TestConfigMixinImpl
}

type TestConfigMixinImpl struct {
	MixinKey string `env:"MIXIN_KEY_EC2B754B,required"`
}

func (*TestConfigMixinImpl) Config() {
	// intentionally empty
}

func (v *TestConfigMixinImpl) GetMixin() *TestConfigMixinImpl {
	return v
}

type TestConfig struct {
	Key string `env:"KEY_EC2B754B,required"`
	TestConfigMixinImpl
}

func (*TestConfig) Config() {
	// intentionally empty
}

type Suite struct {
	CFG *tcfgm.Helper[*TestConfig]
}

func TestSuite(t *testing.T) {
	errorz.MaybeMustWrap(os.Setenv("TEST_KEY_EC2B754B", "Value"))
	errorz.MaybeMustWrap(os.Setenv("TEST_MIXIN_KEY_EC2B754B", "MixinValue"))

	fixturez.RunSuite(t, &Suite{
		CFG: &tcfgm.Helper[*TestConfig]{
			ConfigLoader: cfgm.MustNewEnvConfigLoader[*TestConfig](&cfgm.EnvConfigLoaderOptions{
				Prefix: "TEST_",
			}),
		},
	})
}

func (*Suite) TestMustGet(ctx context.Context, g *WithT) {
	g.Expect(cfgm.MustGet[*TestConfig](ctx)).To(Equal(&TestConfig{
		Key: "Value",
		TestConfigMixinImpl: TestConfigMixinImpl{
			MixinKey: "MixinValue",
		},
	}))

	g.Expect(cfgm.MustGet[TestConfigMixin](ctx).GetMixin()).To(Equal(&TestConfigMixinImpl{
		MixinKey: "MixinValue",
	}))

	g.Expect(func() {
		cfgm.MustGet[cfgm.Config](ctx)
	}).ToNot(Panic())
}
