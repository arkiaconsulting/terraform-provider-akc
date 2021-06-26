package client

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/go-azure-helpers/authentication"
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
	client *AppConfigClient
}

func (s *existingKeyValueWithLabelTestSuite) SetupSuite() {
	s.uri = "https://arkia.azconfig.io"

	builder := authentication.Builder{
		SubscriptionID: os.Getenv("ARM_SUBSCRIPTION_ID"),
		ClientID:       os.Getenv("ARM_CLIENT_ID"),
		TenantID:       os.Getenv("ARM_TENANT_ID"),
		ClientSecret:   os.Getenv("ARM_CLIENT_SECRET"),
		Environment:    "public",
		MetadataHost:   os.Getenv("ARM_METADATA_HOST"),

		// we intentionally only support Client Secret auth for tests (since those variables are used all over)
		SupportsClientSecretAuth: true,
	}

	config, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("Error building ARM Client: %+v", err))
	}

	clientBuilder := ClientBuilder{
		AuthConfig:   config,
		AppConfigUri: s.uri,
	}

	client, err := NewAppConfigurationClient(context.TODO(), clientBuilder)

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
