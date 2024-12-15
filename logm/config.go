package logm

import (
	"encoding"

	"github.com/ibrt/golang-utils/errorz"
	"github.com/sirupsen/logrus"

	"github.com/ibrt/golang-modules/cfgm"
)

var (
	_ encoding.TextUnmarshaler = (*LogConfigLogrusOutput)(nil)
)

// LogConfigLogrusOutput describes the acceptable values for [LogConfig.LogrusOutput].
type LogConfigLogrusOutput string

// Known LogConfigLogrusOutput values.
const (
	LogConfigLogrusOutputDisabled LogConfigLogrusOutput = cfgm.DisabledValue
	LogConfigLogrusOutputHuman    LogConfigLogrusOutput = "human"
	LogConfigLogrusOutputJSON     LogConfigLogrusOutput = "json"
)

// UnmarshalText implements the [encoding.TextUnmarshaler] interface.
func (l *LogConfigLogrusOutput) UnmarshalText(text []byte) error {
	switch v := LogConfigLogrusOutput(text); v {
	case LogConfigLogrusOutputDisabled, LogConfigLogrusOutputHuman, LogConfigLogrusOutputJSON:
		*l = v
		return nil
	default:
		return errorz.Errorf("invalid value for LogConfigLogrusOutput: '%s'", v)
	}
}

// LogConfig describes the log module configuration.
type LogConfig struct {
	HoneycombAPIKey     string                `env:"LOG_HONEYCOMB_API_KEY,required"`
	HoneycombDataset    string                `env:"LOG_HONEYCOMB_DATASET,required"`
	HoneycombSampleRate int                   `env:"LOG_HONEYCOMB_SAMPLE_RATE,required"`
	LogrusOutput        LogConfigLogrusOutput `env:"LOG_LOGRUS_OUTPUT,required"`
	LogrusLevel         logrus.Level          `env:"LOG_LOGRUS_LEVEL,required"`
}

// LogConfigMixin describes the log module configuration.
type LogConfigMixin interface {
	cfgm.Config
	GetLogConfig() *LogConfig
}