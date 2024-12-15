package cfgm_test

import (
	"context"
	"testing"

	"github.com/ibrt/golang-utils/envz"
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
	MixinKey string `env:"MIXIN_KEY_EC2B754B,notEmpty"`
}

func (*TestConfigMixinImpl) Config() {
	// intentionally empty
}

func (v *TestConfigMixinImpl) GetMixin() *TestConfigMixinImpl {
	return v
}

type TestConfig struct {
	Key string `env:"KEY_EC2B754B,notEmpty"`
	TestConfigMixinImpl
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
				return &TestConfig{
					Key: "Value",
					TestConfigMixinImpl: TestConfigMixinImpl{
						MixinKey: "MixinValue",
					},
				}, nil
			},
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

func (*Suite) TestEnvConfigLoader(ctx context.Context, g *WithT) {
	envz.WithEnv(
		map[string]string{
			"TEST_KEY_EC2B754B":       "Value",
			"TEST_MIXIN_KEY_EC2B754B": "MixinValue",
		},
		func() {
			cfg, err := cfgm.MustNewEnvConfigLoader[*TestConfig](&cfgm.EnvConfigLoaderOptions{Prefix: "TEST_"})(ctx)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(cfg).To(Equal(&TestConfig{
				Key: "Value",
				TestConfigMixinImpl: TestConfigMixinImpl{
					MixinKey: "MixinValue",
				},
			}))
		})

	envz.WithEnv(
		map[string]string{
			"KEY_EC2B754B":       "Value",
			"MIXIN_KEY_EC2B754B": "MixinValue",
		},
		func() {
			cfg, err := cfgm.MustNewEnvConfigLoader[*TestConfig](nil)(ctx)
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(cfg).To(Equal(&TestConfig{
				Key: "Value",
				TestConfigMixinImpl: TestConfigMixinImpl{
					MixinKey: "MixinValue",
				},
			}))
		})

	envz.WithEnv(
		map[string]string{
			"KEY_EC2B754B":       "",
			"MIXIN_KEY_EC2B754B": "",
		},
		func() {
			_, err := cfgm.MustNewEnvConfigLoader[*TestConfig](nil)(ctx)
			g.Expect(err).To(MatchError("env: environment variable \"KEY_EC2B754B\" should not be empty; environment variable \"MIXIN_KEY_EC2B754B\" should not be empty"))
		})
}
