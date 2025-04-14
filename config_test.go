// Package config provides simple persistent storage for application configuration
package config

import (
	"errors"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type ConfigTestSuite struct {
	BaseTestSuite
}

func (suite *ConfigTestSuite) TearDownSubTest() {
	_ = os.Remove("config.yaml")
}

func TestExampleTestSuite(t *testing.T) {
	s := ConfigTestSuite{}
	s.Assert = require.New(t)
	suite.Run(t, &s)
}

type TestSubConfig struct {
	Key1 string `yaml:"key1"`
	Key2 int    `yaml:"key2"`
}

type TestConfig struct {
	StringProp string        `yaml:"stringProp"`
	IntProp    int           `yaml:"intProp"`
	BoolProp   bool          `yaml:"boolProp"`
	SubConfig  TestSubConfig `yaml:"subConfig"`
}

type NonSerializableConfig struct {
	FuncProp func()
}

func (suite *ConfigTestSuite) TestLoad() {
	suite.Run("NormalLoad", func() {
		storage := MockConfigStorage{
			Data: []byte(`stringProp: foo
intProp: 11
boolProp: true
subConfig:
   key1: bar
   key2: 22
`),
		}

		conf := NewConfig(TestConfig{}, storage, YamlSerializer[TestConfig]{})
		suite.Assert.Nil(conf.Load())

		suite.Assert.Equal(TestConfig{
			StringProp: "foo",
			IntProp:    11,
			BoolProp:   true,
			SubConfig: TestSubConfig{
				Key1: "bar",
				Key2: 22,
			},
		}, conf.Data)
	})

	suite.Run("StorageLoadError", func() {
		defaultConfig := TestConfig{StringProp: "default"}

		storage := MockConfigStorage{
			Error: errors.New("unknown error"),
		}

		conf := NewConfig(defaultConfig, storage, YamlSerializer[TestConfig]{})
		suite.Assert.EqualError(conf.Load(), "unknown error")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})

	suite.Run("MissingData", func() {
		defaultConfig := TestConfig{StringProp: "default"}

		storage := MockConfigStorage{
			Missing: true,
		}

		conf := NewConfig(defaultConfig, storage, YamlSerializer[TestConfig]{})
		suite.Assert.Nil(conf.Load())

		suite.Assert.Equal(defaultConfig, conf.Data)
	})

	suite.Run("YamlParseError", func() {
		defaultConfig := TestConfig{StringProp: "default"}

		storage := MockConfigStorage{
			Data: []byte(`stringProp: foo
intProp: 11
boolProp # error here
subConfig:
   key1: bar
   key2: 22
`),
		}

		conf := NewConfig(defaultConfig, storage, YamlSerializer[TestConfig]{})
		suite.Assert.EqualError(conf.Load(), "yaml: line 3: could not find expected ':'")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})

	suite.Run("DeserializeError", func() {
		defaultConfig := TestConfig{StringProp: "default"}

		storage := MockConfigStorage{
			Data: []byte(`stringProp: foo
intProp: 11
boolProp: yes! # error here: bool value expected
subConfig:
   key1: bar
   key2: 22
`),
		}

		conf := NewConfig(defaultConfig, storage, YamlSerializer[TestConfig]{})
		suite.Assert.EqualError(conf.Load(), "yaml: unmarshal errors:\n  line 3: cannot unmarshal !!str `yes!` into bool")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})
}

func (suite *ConfigTestSuite) TestSave() {
	suite.Run("NormalSave", func() {
		conf := NewConfigMinimal(TestConfig{
			StringProp: "default",
			IntProp:    11,
			BoolProp:   true,
			SubConfig: TestSubConfig{
				Key1: "value 1",
				Key2: 22,
			},
		})

		conf.Data.StringProp = "not default"
		conf.Data.SubConfig.Key2 = 33

		suite.Assert.Nil(conf.Save())

		suite.AssertFileContentsEqual("config.yaml", `stringProp: not default
intProp: 11
boolProp: true
subConfig:
    key1: value 1
    key2: 33
`)
	})

	suite.Run("SaveError", func() {
		conf := NewConfigMinimal(NonSerializableConfig{
			FuncProp: func() {},
		})

		suite.Assert.EqualError(conf.Save(), "cannot marshal type: func()")
	})
}

func (suite *ConfigTestSuite) TestTransaction() {
	suite.Run("NormalSave", func() {
		conf := NewConfigMinimal(TestConfig{StringProp: "default"})

		err := conf.Transaction(func(c *TestConfig) error {
			c.BoolProp = true
			c.SubConfig.Key2 = 11
			return nil
		})

		suite.Assert.Nil(err)

		suite.AssertFileContentsEqual("config.yaml", `stringProp: default
intProp: 0
boolProp: true
subConfig:
    key1: ""
    key2: 11
`)
	})

	suite.Run("TransactionError", func() {
		defaultConfig := TestConfig{StringProp: "default"}
		conf := NewConfigMinimal(defaultConfig)

		err := conf.Transaction(func(c *TestConfig) error {
			c.BoolProp = false
			c.SubConfig.Key2 = 44
			return errors.New("unknown error")
		})

		suite.Assert.EqualError(err, "unknown error")

		suite.Assert.NoFileExists("config.yaml")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})

	suite.Run("TransactionPanic", func() {
		defaultConfig := TestConfig{StringProp: "default"}
		conf := NewConfigMinimal(defaultConfig)

		suite.Assert.PanicsWithValue("unknown panic", func() {
			_ = conf.Transaction(func(c *TestConfig) error {
				c.BoolProp = false
				c.SubConfig.Key2 = 44
				panic("unknown panic")
			})
		})

		suite.Assert.NoFileExists("config.yaml")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})

	suite.Run("SaveError", func() {
		defaultConfig := TestConfig{StringProp: "default"}
		storage := MockConfigStorage{
			Error: errors.New("unknown error"),
		}
		conf := NewConfig(defaultConfig, storage, YamlSerializer[TestConfig]{})

		err := conf.Transaction(func(c *TestConfig) error {
			c.BoolProp = false
			c.SubConfig.Key2 = 44
			return nil
		})

		suite.Assert.EqualError(err, "unknown error")

		suite.Assert.NoFileExists("config.yaml")

		suite.Assert.Equal(defaultConfig, conf.Data)
	})
}
