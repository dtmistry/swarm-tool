package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RotateSecretsTestSuite struct {
	suite.Suite
}

func TestRotateSecretsTestSuite(t *testing.T) {
	s := new(RotateSecretsTestSuite)
	suite.Run(t, s)
}

func (s *RotateSecretsTestSuite) TestGetSecretNameDefault() {
	//Default Case
	name := GetSecretName("test", "")
	assert.False(s.T(), len(name) == 0)
	assert.Equal(s.T(), "temp_test", name)
}

func (s *RotateSecretsTestSuite) TestGetSecretName() {
	name := GetSecretName("test", "prefix")
	assert.False(s.T(), len(name) == 0)
	assert.Equal(s.T(), "prefix_test", name)
}
