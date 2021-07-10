package client

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestExistingKeyValueWithoutLabelTestSuite(t *testing.T) {
	suiteTester := new(existingKeyValueTestSuite)
	suite.Run(t, suiteTester)
}

type existingKeyValueTestSuite struct {
	suite.Suite
	uri    string
	key    string
	value  string
	client *Client
}

func (s *existingKeyValueTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	client, err := NewClientCreds(s.uri, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"), os.Getenv("ARM_TENANT_ID"))

	if err != nil {
		panic(err)
	}

	s.client = client
}

func (s *existingKeyValueTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"

	_, err := s.client.SetKeyValue(LabelNone, s.key, s.value)

	if err != nil {
		panic(fmt.Errorf("Cannot delete test key-value %s", err.Error()))
	}
}

func (s *existingKeyValueTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValue(LabelNone, s.key)

	if err != nil {
		panic(fmt.Errorf("Cannot delete test key-value %s", err.Error()))
	}
}

func (s *existingKeyValueTestSuite) TestGetExistingKeyValueWithoutLabelShouldPass() {
	result, err := s.client.GetKeyValue(LabelNone, s.key)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), "", result.Label)
}

func (s *existingKeyValueTestSuite) TestDeleteExistingKeyValueWithoutLabelShouldPass() {
	isDeleted, err := s.client.DeleteKeyValue(LabelNone, s.key)

	require.Nil(s.T(), err)
	assert.True(s.T(), isDeleted)
}
