package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const uri = "https://testlg.azconfig.io"

func TestGetKeyValueDoesntExistShouldFail(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "fdgdfgdsfgdsfgsdgfds"

	_, err := client.GetKeyValue(key)

	assert.Equal(t, kVError("%00", key, "Not found"), err)
}

func TestCreateKeyValueWithoutLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const value = "myValue"

	result, err := client.SetKeyValue(key, value)
	assert.Nil(nil, err)

	assert.Equal(t, key, result.Key)
	assert.Equal(t, value, result.Value)
	assert.Equal(t, "", result.Label)
}

func TestCreateKeyValueWithLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const value = "myValue"
	const label = "myLabel"

	result, err := client.SetKeyValueWithLabel(key, value, label)
	assert.Nil(nil, err)

	assert.Equal(t, key, result.Key)
	assert.Equal(t, value, result.Value)
	assert.Equal(t, label, result.Label)
}

func TestCreateAndGetKeyValueWithoutLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const value = "myValue"

	result, err := client.SetKeyValue(key, value)
	assert.Nil(nil, err)

	result, err = client.GetKeyValue(key)
	assert.Nil(nil, err)

	assert.Equal(t, key, result.Key)
	assert.Equal(t, value, result.Value)
	assert.Equal(t, "", result.Label)
}

func TestCreateAndDeleteKeyValueShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const value = "myValue"

	_, err := client.SetKeyValue(key, value)
	assert.Nil(nil, err)

	isDeleted, err := client.DeleteKeyValue(key)
	assert.Nil(nil, err)

	assert.True(t, isDeleted)
}

func TestDeleteKeyValueDoesNotExistShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)

	isDeleted, err := client.DeleteKeyValue("dfqdsgmlmlmdfsd")
	assert.Nil(nil, err)

	assert.False(t, isDeleted)
}

func TestGetKeyValueWithLabelDoesntExistShouldFail(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "fdgdfgdsfgdsfgsdgfds"
	const label = "mmlmlmlmlmlmsdsqdsqd"

	_, err := client.GetKeyValueWithLabel(key, label)

	assert.Equal(t, kVError(label, key, "Not found"), err)
}

func TestGetKeyValueWithLabelShouldRespectLabel(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const label = "myLabel"
	const value = "myValue"

	_, err := client.SetKeyValueWithLabel(key, value, label)
	assert.Nil(nil, err)

	_, err = client.GetKeyValue(key)

	assert.Equal(t, kVError("%00", key, "Not found"), err)
}

func TestGetKeyValueWithLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const label = "myLabel"
	const value = "myValue"

	_, err := client.SetKeyValueWithLabel(key, value, label)
	assert.Nil(nil, err)

	result, err := client.GetKeyValueWithLabel(key, label)
	assert.Nil(nil, err)

	assert.Equal(t, key, result.Key)
	assert.Equal(t, value, result.Value)
	assert.Equal(t, label, result.Label)
}

func TestDeleteKeyValueWithLabelShouldDeleteConcernedOnly(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const label = "myLabel"
	const value = "myValue"

	_, err := client.SetKeyValueWithLabel(key, value, label)
	assert.Nil(nil, err)

	_, err = client.SetKeyValueWithLabel(key, value, "otherLabel")
	assert.Nil(nil, err)

	isDeleted, err := client.DeleteKeyValueWithLabel(key, label)
	assert.Nil(nil, err)
	assert.True(t, isDeleted)

	_, err = client.GetKeyValueWithLabel(key, "otherLabel")
	assert.Nil(nil, err)
}

func TestDeleteKeyValueWithLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const label = "myLabel"
	const value = "myValue"

	_, err := client.SetKeyValueWithLabel(key, value, label)
	assert.Nil(nil, err)

	isDeleted, err := client.DeleteKeyValueWithLabel(key, label)
	assert.Nil(nil, err)
	assert.True(t, isDeleted)

	_, err = client.GetKeyValueWithLabel(key, label)
	assert.NotNil(nil, err)
}

func TestGetKeyValueWithoutLabelShouldPass(t *testing.T) {
	client := NewAppConfigurationClient(uri)
	const key = "myKey"
	const value = "myValue"

	_, err := client.SetKeyValue(key, value)
	assert.Nil(nil, err)

	_, err = client.GetKeyValue(key)
	assert.Nil(nil, err)
}
