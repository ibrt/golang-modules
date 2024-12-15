package logm_test

import (
	"testing"

	"github.com/caarlos0/env/v11"
	"github.com/ibrt/golang-utils/envz"
	"github.com/ibrt/golang-utils/fixturez"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"github.com/ibrt/golang-modules/cfgm"
	"github.com/ibrt/golang-modules/logm"
)

type ConfigSuite struct {
	// intentionally empty
}

func TestConfigSuite(t *testing.T) {
	fixturez.RunSuite(t, &ConfigSuite{})
}

func (*ConfigSuite) TestLogConfig(g *WithT) {
	envz.WithEnv(
		map[string]string{
			"LOG_HONEYCOMB_API_KEY":     cfgm.DisabledValue,
			"LOG_HONEYCOMB_DATASET":     "test",
			"LOG_HONEYCOMB_SAMPLE_RATE": "1",
			"LOG_LOGRUS_OUTPUT":         cfgm.DisabledValue,
			"LOG_LOGRUS_LEVEL":          logrus.InfoLevel.String(),
		},
		func() {
			logCfg, err := env.ParseAs[logm.LogConfig]()
			g.Expect(err).ToNot(HaveOccurred())
			g.Expect(logCfg).To(Equal(logm.LogConfig{
				HoneycombAPIKey:     cfgm.DisabledValue,
				HoneycombDataset:    "test",
				HoneycombSampleRate: 1,
				LogrusOutput:        cfgm.DisabledValue,
				LogrusLevel:         logrus.InfoLevel,
			}))
		})

	envz.WithEnv(
		map[string]string{
			"LOG_HONEYCOMB_API_KEY":     cfgm.DisabledValue,
			"LOG_HONEYCOMB_DATASET":     "test",
			"LOG_HONEYCOMB_SAMPLE_RATE": "1",
			"LOG_LOGRUS_OUTPUT":         "invalid",
			"LOG_LOGRUS_LEVEL":          logrus.InfoLevel.String(),
		},
		func() {
			_, err := env.ParseAs[logm.LogConfig]()
			g.Expect(err).To(MatchError(`env: parse error on field "LogrusOutput" of type "logm.LogConfigLogrusOutput": invalid value for LogConfigLogrusOutput: 'invalid'`))
		})
}
