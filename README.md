# terraform-provider-akc
This terraform provider will allow the creation the Azure App Configuration resources, be they simple values or Key Vault secret references.

## Usage example
### Resource
#### Configure the Akc App Configuration provider
```terraform
terraform {
  required_providers {
    akc = {
      source = "arkiaconsulting/akc"
    }
  }
}

provider "akc" {
}
```
#### Create an App Configuration key-value
```terraform
resource "akc_key_value" "test" {
  endpoint = azurerm_app_configuration.test.endpoint
  key      = "Key"
  value    = "my config value"
}
```
#### Create an App Configuration key-value with label
```terraform
resource "akc_key_value" "config_value" {
  endpoint = azurerm_app_configuration.test.endpoint
  label    = "Dev"
  key      = "Key"
  value    = "my config value"
}
```
#### Create an App Configuration key-value with Key Vault secret reference
```terraform
resource "akc_key_secret" "config_secret" {
  endpoint  = azurerm_app_configuration.test.endpoint
  label     = "Dev"
  key       = "storage-connection-string"
  secret_id = azurerm_key_vault_secret.secret.id
  latest_version = true # Trim or not the version information (default to false)
}
```

### DataSource
#### Source an existing key-value
```terraform
data "akc_key_value" "my_value" {
  endpoint  = azurerm_app_configuration.test.endpoint
  label     = "Dev"
  key       = "Key"
}
```
*Reference the resulting value using `data.akc_key_value.my_value.value`*

#### Source an existing key-secret
```terraform
data "akc_key_secret" "my_secret_id" {
  endpoint  = azurerm_app_configuration.test.endpoint
  label     = "Dev"
  key       = "Key"
}
```
*Reference the resulting secret Id using `data.akc_key_value.my_secret_id.secret_id`*

## Authorization
The provider uses the current Azure CLI credentials if available, and fall back to environment variables.

The identity must have be assigned the RBAC role `App Configuration Data Owner`, or at least `App Configuration Data Reader` in order to use the data source.

If you don't want to connect using Azure CLI credentials, you must configure the following environment variables (terraform-azurerm standard):
```sh
export ARM_CLIENT_ID=XXXXXXXX-XXX
export ARM_SUBSCRIPTION_ID=XXXXXXXX-XXX
export ARM_TENANT_ID=XXXXXXXX-XXX
export ARM_CLIENT_SECRET=XXXXXXX
```

## Installation
The provider is available on the terraform registry