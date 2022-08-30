package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_TerraformResourceDataArgLayout(t *testing.T) {

	cases := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. separate attributes from blocks",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  diagnostics {
    log_analytics {
      workspace_id  = "test"
      workspace_key = "test"
    }
  }
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  dns_config {
    nameservers = []
  }
  tags = {
    Name = "test"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "test"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  diagnostics {
    log_analytics {
      workspace_id  = "test"
      workspace_key = "test"
    }
  }
  dns_config {
    nameservers = []
  }
}`,
				},
			},
		},
		{
			Name: "2. check for nested block",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "test"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
    memory = "1.5"
    ports {
      port     = 443
      protocol = "TCP"
    }
    name   = "hello-world"
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
container {
  cpu    = "0.5"
  image  = "mcr.microsoft.com/azuredocs/aci-helloworld:latest"
  memory = "1.5"
  name   = "hello-world"

  ports {
    port     = 443
    protocol = "TCP"
  }
}`,
				},
			},
		},
		{
			Name: "3. Gap between different types of args",
			Content: `
resource "azurerm_container_group" "example" {
  count               = 4
  provider            = azurerm.europe
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "container ${count.index}"
  }
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  diagnostics {
    log_analytics {
      workspace_id  = "test"
      workspace_key = "test"
    }
  }
  dns_config {
    nameservers = []
  }
  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_container_group" "example" {
  count    = 4
  provider = azurerm.europe

  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "container ${count.index}"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  diagnostics {
    log_analytics {
      workspace_id  = "test"
      workspace_key = "test"
    }
  }
  dns_config {
    nameservers = []
  }

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
				},
			},
		},
		{
			Name: "4. Meta Arg",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  count               = 4
  name                = "example-continst"
  provider            = azurerm.europe
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  depends_on = [
    azurerm_resource_group.example
  ]
  tags = {
    Name = "container ${count.index}"
  }
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  for_each            = local.containers
  name                = "example-continst"
  provider            = azurerm.europe
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  depends_on = [
    azurerm_resource_group.example
  ]
  tags = {
    Name = "container ${each.key}"
  }
  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
  lifecycle {
    create_before_destroy = true
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_container_group" "example" {
  count    = 4
  provider = azurerm.europe

  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "container ${count.index}"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
				},
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
resource "azurerm_container_group" "example" {
  for_each = local.containers
  provider = azurerm.europe

  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "container ${each.key}"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
				},
			},
		},

		{
			Name: "5. dynamic block",
			Content: `
resource "azurerm_container_group" "example" {
  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name

  dynamic "container" {
    content {
	  name   = container.value["name"]
	  image  = container.value["image"]
	  cpu 	 = container.value["cpu"]
	  ports {
		port     = 443
		protocol = "TCP"
	  }
      memory = container.value["memory"]
    }
	for_each = var.containers
  }
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
dynamic "container" {
  for_each = var.containers

  content {
    name   = container.value["name"]
    image  = container.value["image"]
    cpu 	 = container.value["cpu"]
    memory = container.value["memory"]

    ports {
	  port     = 443
	  protocol = "TCP"
    }
  }
}`,
				},
			},
		},

		{
			Name: "6. correct layout",
			Content: `
resource "azurerm_container_group" "example" {
  count    = 4
  provider = azurerm.europe

  location            = azurerm_resource_group.example.location
  name                = "example-continst"
  ip_address_type     = "Public"
  dns_name_label      = "aci-label"
  os_type             = "Linux"
  resource_group_name = azurerm_resource_group.example.name
  tags = {
    Name = "container ${count.index}"
  }

  container {
    cpu    = "0.5"
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }

  lifecycle {
    create_before_destroy = true
  }
  depends_on = [
    azurerm_resource_group.example
  ]
}`,
			Expected: helper.Issues{},
		},
		{
			Name: "7. datasource",
			Content: `
data "azurerm_resources" "example" {
  testblock {
    environment = "production"
    role        = "webserver"
  }
  resource_group_name = "example-resources"
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformResourceDataArgLayoutRule(),
					Message: `Arguments are expected to be arranged in following Layout:
data "azurerm_resources" "example" {
  resource_group_name = "example-resources"

  testblock {
    environment = "production"
    role        = "webserver"
  }
}`,
				},
			},
		},
		{
			Name: "8. empty block",
			Content: `
resource "azurerm_virtual_network" "vnet" {}`,
			Expected: helper.Issues{},
		},
	}

	rule := NewRule(NewTerraformResourceDataArgLayoutRule())

	for _, tc := range cases {
		runner := helper.TestRunner(t, map[string]string{"config.tf": tc.Content})
		t.Run(tc.Name, func(t *testing.T) {
			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
