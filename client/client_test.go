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

func TestNonExistingKeyValueTestSuite(t *testing.T) {
	suiteTester := new(nonExistingKeyValueWithLabelTestSuite)
	suite.Run(t, suiteTester)
}

type nonExistingKeyValueWithLabelTestSuite struct {
	suite.Suite
	uri       string
	label     string
	key       string
	value     string
	secretURI string
	client    *AppConfigClient
}

func (s *nonExistingKeyValueWithLabelTestSuite) SetupSuite() {
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

func (s *nonExistingKeyValueWithLabelTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"
	s.label = "myLabel"
	s.secretURI = "https://testlg.vault.azure.net/secrets/my-secret"
}

func (s *nonExistingKeyValueWithLabelTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValue(LabelNone, s.key)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}

	_, err = s.client.DeleteKeyValue(s.label, s.key)

	if err != nil {
		panic("Cannot delete test key-value with label")
	}
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestGetKeyValueDoesntExistShouldFail() {
	_, err := s.client.GetKeyValue(LabelNone, s.key)

	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestGetKeyValueWithLabelDoesntExistShouldFail() {
	_, err := s.client.GetKeyValue(s.label, s.key)

	require.EqualError(s.T(), err, AppConfigClientError{Message: KVNotFoundError.Message, Info: s.key}.Error())
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestCreateKeyValueWithoutLabelShouldPass() {
	result, err := s.client.SetKeyValue(LabelNone, s.key, s.value)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), "", result.Label)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestCreateKeyValueWithLabelShouldPass() {
	result, err := s.client.SetKeyValue(s.label, s.key, s.value)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), s.value, result.Value)
	assert.Equal(s.T(), s.label, result.Label)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestDeleteKeyValueDoesNotExistShouldPass() {
	isDeleted, err := s.client.DeleteKeyValue(LabelNone, s.key)

	require.Nil(s.T(), err)
	assert.False(s.T(), isDeleted)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestDeleteKeyValueWithLabelDoesNotExistShouldPass() {
	isDeleted, err := s.client.DeleteKeyValue(s.label, s.key)

	require.Nil(s.T(), err)
	assert.False(s.T(), isDeleted)
}

func (s *nonExistingKeyValueWithLabelTestSuite) TestCreateKeyValueSecretShouldPass() {
	result, err := s.client.SetKeyValueSecret(s.key, s.secretURI, LabelNone)

	require.Nil(s.T(), err)
	assert.Equal(s.T(), s.key, result.Key)
	assert.Equal(s.T(), fmt.Sprintf("{\"uri\":\"%s\"}", s.secretURI), result.Value)
	assert.Equal(s.T(), "", result.Label)
}
