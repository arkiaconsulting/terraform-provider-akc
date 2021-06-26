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

func TestExistingKeyValueWithoutLabelTestSuite(t *testing.T) {
	suiteTester := new(existingKeyValueTestSuite)
	suite.Run(t, suiteTester)
}

type existingKeyValueTestSuite struct {
	suite.Suite
	uri    string
	key    string
	value  string
	client *AppConfigClient
}

func (s *existingKeyValueTestSuite) SetupSuite() {
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

func (s *existingKeyValueTestSuite) SetupTest() {
	s.key = "myKey"
	s.value = "myValue"

	_, err := s.client.SetKeyValue(LabelNone, s.key, s.value)

	if err != nil {
		panic("Cannot create test key-value")
	}
}

func (s *existingKeyValueTestSuite) TearDownTest() {
	_, err := s.client.DeleteKeyValue(LabelNone, s.key)

	if err != nil {
		panic("Cannot delete test key-value")
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
