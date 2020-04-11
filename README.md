# terraform-provider-akc
This terraform provider will allow the creation the Azure App Configuration resources, be they simple values or Key Vault secret references.

## Usage example
```
# Configure the Akc App Configuration provider

provider "akc" {
}

# Create an App Configuration key-value

resource "akc_key_value" "test" {
  endpoint = azurerm_app_configuration.test.endpoint
  key      = "Key"
  value    = "my config value"
}

# Create an App Configuration key-value with label

resource "akc_key_value" "config_value" {
  endpoint = azurerm_app_configuration.test.endpoint
  label    = "Dev"
  key      = "Key"
  value    = "my config value"
}

# Create an App Configuration key-value with Key Vault secret reference

resource "akc_key_secret" "config_secret" {
  endpoint  = azurerm_app_configuration.test.endpoint
  label     = "Dev"
  key       = "storage-connection-string"
  secret_id = azurerm_key_vault_secret.secret.id
}

```

## Authorization
The provider uses the current Azure CLI credentials if available, and fall back to environment variables.

The identity must have be assigned the RBAC role `App Configuration Data Owner`.

## Installation
### Linux
Copy the `terraform-provider-akc_vX.X.X` to your main terraform configuration folder.
### Windows
Copy the `terraform-provider-akc_vX.X.X.exe` to your main terraform configuration folder.