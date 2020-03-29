package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestNonExistingKeyValueTestSuite(t *testing.T) {
	suiteTester := new(nonExistingKeyValueWithLabelTestSuite)
	suite.Run(t, suiteTester)
}

type nonExistingKeyValueWithLabelTestSuite struct {
	suite.Suite
	uri    string
	label  string
	key    string
	value  string
	client *AppConfigClient
}

func (s *nonExistingKeyValueWithLabelTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	client, err := NewAppConfigurationClient(s.uri)

	if err != nil {
		panic(err)
	}

	s.client = client
}

func (s *nonExistingKeyValueWithLabelTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"
	s.label = "myLabel"
}

func (s *nonExistingKeyValueWithLabelTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValue(s.key)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}

	_, err = s.client.DeleteKeyValueWithLabel(s.key, s.label)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestGetKeyValueDoesntExistShouldFail() {
	_, err := s.client.GetKeyValue(s.key)

	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelDoesntExistShouldFail() {
	_, err := s.client.GetKeyValueWithLabel(s.key, s.label)

	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestCreateKeyValueWithoutLabelShouldPass() {
	result, err := s.client.SetKeyValue(s.key, s.value)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), "", result.Label)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestCreateKeyValueWithLabelShouldPass() {
	result, err := s.client.SetKeyValueWithLabel(s.key, s.value, s.label)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), s.label, result.Label)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestDeleteKeyValueDoesNotExistShouldPass() {
	isDeleted, err := s.client.DeleteKeyValue(s.key)

	require.Nil(s.T(), err)
	assert.False(s.T(), isDeleted)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelDoesNotExistShouldPass() {
	isDeleted, err := s.client.DeleteKeyValueWithLabel(s.key, s.label)

	require.Nil(s.T(), err)
	assert.False(s.T(), isDeleted)
}
