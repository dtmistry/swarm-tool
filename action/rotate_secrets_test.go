package action

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RotateSecretsSuite struct {
	suite.Suite
}

func TestRotateSecretsSuite(t *testing.T) {
	s := new(RotateSecretsSuite)
	suite.Run(t, s)
}

func (s *RotateSecretsSuite) TestGetSecretNameDefault() {
	//Default Case
	name := GetSecretName("test", "")
	assert.False(s.T(), len(name) == 0)
	assert.Equal(s.T(), "temp_test", name)
}

func (s *RotateSecretsSuite) TestGetSecretName() {
	name := GetSecretName("test", "prefix")
	assert.False(s.T(), len(name) == 0)
	assert.Equal(s.T(), "prefix_test", name)
}
