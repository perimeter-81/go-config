package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
	"time"
)

type internalStruct struct {
	InternalField string   `yaml:"internal_field" env:"INTERNAL_FIELD" default:"dudu"`
	Strings       []string `yaml:"strings" env:"STRINGS"  default:"\"4\",\"5\""`
}

type testConfig struct {
	Some                string         `yaml:"somevar" env:"somevar" default:"koko"`
	Other               bool           `yaml:"othervar" env:"othervar" default:"true"`
	Port                int            `yaml:"port" env:"PORT"  default:"-1"`
	NotAnEnv            string         `yaml:"notAnEnv" env:"" default:"popo"`
	DatabaseURL         string         `yaml:"database" env:"DATABASE_URL" default:"postgres://localhost:5432/db"`
	Strings             []string       `yaml:"strongs" env:"STRINGS"`
	SepStrings          []string       `yaml:"strings" env:"SEPSTRINGS" envSeparator:":"`
	Numbers             []int          `yaml:"numbers" env:"NUMBERS"`
	Numbers64           []int64        `yaml:"numbers64" env:"NUMBERS64"`
	Bools               []bool         `yaml:"bools" env:"BOOLS"`
	Duration            time.Duration  `yaml:"duration" env:"DURATION" default:"5s"`
	Float32             float32        `yaml:"float32" env:"FLOAT32"`
	Float64             float64        `yaml:"float64" env:"FLOAT64"`
	Float32s            []float32      `yaml:"float32s" env:"FLOAT32S"`
	Float64s            []float64      `yaml:"float64s" env:"FLOAT64S"`
	InternalStructField internalStruct `yaml:"internal"`
}

func TestMixedConfig(t *testing.T) {
	readBuf, err := ioutil.ReadFile("testdata/test.yml")
	require.NoError(t, err)
	testCfg := &testConfig{}
	require.NoError(t, Parse(readBuf, testCfg))
	assert.Equal(t, "a", testCfg.Some)
	assert.Equal(t, true, testCfg.Other)
	assert.Equal(t, 666, testCfg.Port)
	assert.Equal(t, "http://cool_db", testCfg.DatabaseURL)
	assert.Equal(t, []string{"a", "b"}, testCfg.Strings)
	assert.Equal(t, []int{1, 2}, testCfg.Numbers)
	assert.Equal(t, []int64{1, 2}, testCfg.Numbers64)
	assert.Equal(t, []bool{true, true}, testCfg.Bools)
	assert.Equal(t, "cupcake", testCfg.InternalStructField.InternalField)
	assert.Equal(t, []string{"a", "b"}, testCfg.InternalStructField.Strings)
}
