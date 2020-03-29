package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	client *AppConfigClient
}

func (s *existingKeyValueWithLabelTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	s.client = NewAppConfigurationClient(s.uri)
}

func (s *existingKeyValueWithLabelTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"
	s.label = "myLabel"

	_, err := s.client.SetKeyValueWithLabel(s.key, s.value, s.label)

	if err != nil {
		panic("Cannot create test key-value with label")
	}
}

func (s *existingKeyValueWithLabelTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValueWithLabel(s.key, s.label)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}
}

func (s *existingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelShouldRespectLabel() {
	_, err := s.client.GetKeyValue(s.key)

	assert.NotNil(s.T(), err)
	assert.EqualError(s.T(), KVNotFoundError, err.Error())
	assert.Equal(s.T(), s.key, err.(AppConfigClientError).Info)

	_, err = s.client.GetKeyValueWithLabel(s.key, "anyOtherLabel")

	assert.NotNil(s.T(), err)
	assert.EqualError(s.T(), KVNotFoundError, err.Error())
	assert.Equal(s.T(), s.key, err.(AppConfigClientError).Info)
}

func (s *existingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelShouldPass() {
	result, err := s.client.GetKeyValueWithLabel(s.key, s.label)
	assert.Nil(s.T(), err)

	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), s.label, result.Label)
}

func (s *existingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelShouldDeleteConcernedOnly() {
	_, err := s.client.SetKeyValueWithLabel(s.key, s.value, "otherLabel")
	assert.Nil(s.T(), err)

	isDeleted, err := s.client.DeleteKeyValueWithLabel(s.key, s.label)
	assert.Nil(s.T(), err)
	assert.True(s.T(), isDeleted)

	_, err = s.client.GetKeyValueWithLabel(s.key, "otherLabel")
	assert.Nil(s.T(), err)
}

func (s *existingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelShouldPass() {
	isDeleted, err := s.client.DeleteKeyValueWithLabel(s.key, s.label)
	assert.Nil(s.T(), err)
	assert.True(s.T(), isDeleted)

	_, err = s.client.GetKeyValueWithLabel(s.key, s.label)
	assert.NotNil(s.T(), err)
}
