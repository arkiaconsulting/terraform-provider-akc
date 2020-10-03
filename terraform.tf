terraform {
  required_providers {
    akc = {
      source = "arkiaconsulting/akc"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {
}

provider "akc" {
}

resource "azurerm_resource_group" "test" {
  name     = "tf-tests"
  location = "francecentral"
}

resource "azurerm_app_configuration" "test" {
  name                = "testlg"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "standard"
}

resource "azurerm_key_vault" "test" {
  resource_group_name = azurerm_resource_group.test.name
  name                = "testlg"
  location            = azurerm_resource_group.test.location
  sku_name            = "standard"
  tenant_id           = data.azurerm_client_config.current.tenant_id
}

resource "azurerm_key_vault_access_policy" "sp" {
  key_vault_id       = azurerm_key_vault.test.id
  tenant_id          = data.azurerm_client_config.current.tenant_id
  object_id          = data.azurerm_client_config.current.object_id
  secret_permissions = ["get", "set", "list", "delete"]
}

resource "azurerm_key_vault_secret" "secret" {
  key_vault_id = azurerm_key_vault.test.id
  name         = "my-secret"
  value        = "mySecretValue"
}

resource "akc_key_value" "test" {
  endpoint = azurerm_app_configuration.test.endpoint
  key      = "myKey"
  value    = "myValuesdfsdfds"
}

resource "akc_key_value" "test2" {
  endpoint = azurerm_app_configuration.test.endpoint
  key      = "myKey2"
  value    = "myValue2"
}

resource "akc_key_value" "test2_label" {
  endpoint = azurerm_app_configuration.test.endpoint
  label    = "myLabel"
  key      = "myKey2"
  value    = "myValue2"
}

resource "akc_key_secret" "secret1" {
  endpoint  = azurerm_app_configuration.test.endpoint
  key       = "my-secret1"
  secret_id = azurerm_key_vault_secret.secret.id
}

resource "akc_key_secret" "secret1_label" {
  endpoint  = azurerm_app_configuration.test.endpoint
  label     = "myLabel"
  key       = "my-secret12"
  secret_id = azurerm_key_vault_secret.secret.id
}

resource "akc_key_secret" "secret2" {
  endpoint       = azurerm_app_configuration.test.endpoint
  label          = "myLabel"
  key            = "my-secret12-latest"
  secret_id      = azurerm_key_vault_secret.secret.id
  latest_version = true
}

data "akc_key_value" "test2_label" {
  endpoint = azurerm_app_configuration.test.endpoint
  label    = "myLabel"
  key      = "myKey2"
}

data "akc_key_secret" "secret2" {
  endpoint = azurerm_app_configuration.test.endpoint
  label    = "myLabel"
  key      = "my-secret12-latest"
}

output "one" {
  value = data.akc_key_value.test2_label.value
}

output "two" {
  value = data.akc_key_secret.secret2.secret_id
}
