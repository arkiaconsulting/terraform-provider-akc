package client

import (
	"fmt"
	"os"
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
	keyNoLabel  string
}

func (s *featuresTestSuite) SetupSuite() {
	s.uri = "https://testlg.azconfig.io"
	client, err := NewClientCreds(s.uri, os.Getenv("ARM_CLIENT_ID"), os.Getenv("ARM_CLIENT_SECRET"), os.Getenv("ARM_TENANT_ID"))

	if err != nil {
		panic(err)
	}

	s.client = client
}

func (s *featuresTestSuite) SetupTest() {
	s.key = uuid.New().String()
	s.label = uuid.New().String()
	s.enabled = false
	s.description = uuid.New().String()
	s.keyNoLabel = uuid.New().String()

	_, err := s.client.SetFeature(s.key, s.label, s.enabled, s.description)

	if err != nil {
		panic(fmt.Sprintf("Cannot create feature %s", s.key))
	}

	_, err = s.client.SetFeature(s.keyNoLabel, LabelNone, s.enabled, s.description)

	if err != nil {
		panic(fmt.Sprintf("Cannot create feature %s", s.key))
	}
}

func (s *featuresTestSuite) TearDownTest() {
	_, err := s.client.DeleteFeature(s.label, s.key)

	if err != nil {
		panic(fmt.Sprintf("Cannot delete feature %s", s.key))
	}

	_, err = s.client.DeleteFeature(LabelNone, s.keyNoLabel)

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

func (s *featuresTestSuite) TestFeaturesGetFeatureNoLabelShouldPass() {
	result, err := s.client.GetFeature(LabelNone, s.keyNoLabel)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.keyNoLabel, result.Key)
	assert.Equal(s.T(), "", result.Label)
	assert.Equal(s.T(), s.enabled, result.Enabled)
	assert.Equal(s.T(), s.description, result.Description)
}

func (s *featuresTestSuite) TestFeaturesDeleteFeatureShouldPass() {
	ret, _ := s.client.DeleteFeature(s.label, s.key)
	assert.Equal(s.T(), true, ret)

	_, err := s.client.GetFeature(s.label, s.key)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: toPrefixedFeature(s.key)}.Error())
}

func (s *featuresTestSuite) TestFeaturesDeleteFeatureNoLabelShouldPass() {
	name := uuid.New().String()

	_, err := s.client.SetFeature(name, LabelNone, true, "yop")
	require.Nil(s.T(), err)

	ret, _ := s.client.DeleteFeature(LabelNone, name)
	assert.Equal(s.T(), true, ret)

	_, err = s.client.GetFeature(LabelNone, name)
	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: toPrefixedFeature(name)}.Error())
}
