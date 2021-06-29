package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestExistingKeyValueWithLabelTestSuite(t *testing.T) {
	suiteTester := new(existingKeyValueWithLabelTestSuite)
	suite.Run(t, suiteTester)
}

type existingKeyValueWithLabelTestSuite struct {
	suite.Suite
	uri    string
	label  string
	key    string
	value  string
	client *Client
}

func (s *existingKeyValueWithLabelTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	client, err := NewClientCli(s.uri)

	if err != nil {
		panic(err)
	}

	s.client = client
}

func (s *existingKeyValueWithLabelTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"
	s.label = "myLabel"

	_, err := s.client.SetKeyValue(s.label, s.key, s.value)

	if err != nil {
		panic("Cannot create test key-value with label")
	}
}

func (s *existingKeyValueWithLabelTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValue(s.label, s.key)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}
}

func (s *existingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelShouldRespectLabel() {
	_, err := s.client.GetKeyValue(LabelNone, s.key)

	require.NotNil(s.T(), err)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())

	_, err = s.client.GetKeyValue("anyOtherLabel", s.key)

	require.NotNil(s.T(), err)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())
}

func (s *existingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelShouldPass() {
	result, err := s.client.GetKeyValue(s.label, s.key)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), s.label, result.Label)
}

func (s *existingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelShouldDeleteConcernedOnly() {
	_, err := s.client.SetKeyValue("otherLabel", s.key, s.value)

	require.Nil(s.T(), err)

	isDeleted, err := s.client.DeleteKeyValue(s.label, s.key)
	require.Nil(s.T(), err)
	assert.True(s.T(), isDeleted)

	_, err = s.client.GetKeyValue("otherLabel", s.key)
	require.Nil(s.T(), err)
}

func (s *existingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelShouldPass() {
	isDeleted, err := s.client.DeleteKeyValue(s.label, s.key)

	require.Nil(s.T(), err)
	assert.True(s.T(), isDeleted)

	_, err = s.client.GetKeyValue(s.label, s.key)
	require.NotNil(s.T(), err)
}
