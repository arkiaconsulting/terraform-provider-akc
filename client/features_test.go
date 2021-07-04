package client

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestFeaturesTestSuite(t *testing.T) {
	suiteTester := new(featuresTestSuite)
	suite.Run(t, suiteTester)
}

type featuresTestSuite struct {
	suite.Suite
	uri         string
	label       string
	key         string
	enabled     bool
	description string
	client      *Client
}

func (s *featuresTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	client, err := NewClientCli(s.uri)

	if err != nil {
		panic(err)
	}

	s.client = client
}

func (s *featuresTestSuite) SetupTest() {
	s.key = uuid.New().String()
	s.label = "myLabel"
	s.enabled = false
	s.description = "my cool description"

	_, err := s.client.SetFeature(s.key, s.label, s.enabled, s.description)

	if err != nil {
		panic(fmt.Sprintf("Cannot create feature %s", s.key))
	}
}

func (s *featuresTestSuite) TearDownTest() {
	_, err := s.client.DeleteFeature(s.label, s.key)

	if err != nil {
		panic(fmt.Sprintf("Cannot delete feature %s", s.key))
	}
}

func (s *featuresTestSuite) TestFeaturesGetNonExistingFeatureShouldPass() {
	_, err := s.client.GetFeature(s.label, "idontexist")

	require.NotNil(s.T(), err)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: toPrefixedFeature("idontexist")}.Error())

	_, err = s.client.GetFeature(LabelNone, s.key)

	require.NotNil(s.T(), err)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: toPrefixedFeature(s.key)}.Error())
}

func (s *featuresTestSuite) TestFeaturesGetFeatureShouldPass() {
	result, err := s.client.GetFeature(s.label, s.key)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.label, result.Label)
	assert.Equal(s.T(), s.enabled, result.Enabled)
	assert.Equal(s.T(), s.description, result.Description)
}

func (s *featuresTestSuite) TestFeaturesDeleteFeatureShouldPass() {
	ret, _ := s.client.DeleteFeature(s.label, s.key)
	assert.Equal(s.T(), true, ret)

	_, err := s.client.GetFeature(s.label, s.key)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: toPrefixedFeature(s.key)}.Error())
}
