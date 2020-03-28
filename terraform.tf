provider "azurerm" {
  features {}
}

provider "akc" {}

resource "azurerm_resource_group" "test" {
  name     = "testlg"
  location = "francecentral"
}

resource "azurerm_app_configuration" "test" {
  name                = "testlg"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "free"
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
