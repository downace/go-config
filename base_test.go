package config

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"os"
)

type BaseTestSuite struct {
	suite.Suite
	Assert *require.Assertions
}

func (suite *BaseTestSuite) AssertFileContentsEqual(filepath string, expected string) {
	content, err := os.ReadFile(filepath)
	suite.Assert.Nil(err)
	suite.Assert.Equal(expected, string(content))
}

type MockConfigStorage struct {
	Data    []byte
	Missing bool
	Error   error
}

func (s MockConfigStorage) Load() ([]byte, bool, error) {
	if s.Error != nil {
		return nil, s.Missing, s.Error
	}

	return s.Data, s.Missing, nil
}

func (s MockConfigStorage) Save(_ []byte) error {
	if s.Error != nil {
		return s.Error
	}

	return nil
}
