package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformVarNameConvention(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "1. a variable name with property name plus suffix is invalid",
			Content: `
resource "azurerm_container_group" "example" {
  name      = var.name_cg
  location  = var.locations.location_cg

  container {
    cpu    = var.cpu_cg
    image  = "mcr.microsoft.com/azuredocs/aci-tutorial-sidecar"
    memory = "1.5"
    name   = "sidecar"
  }
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`name` Variable:`name_cg` is invalid, expected var name:${prefix}name",
				},
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`location` Variable:`location_cg` is invalid, expected var name:${prefix}location",
				},
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`cpu` Variable:`cpu_cg` is invalid, expected var name:${prefix}cpu",
				},
			},
		},
		{
			Name: "2. only validate true/false var name in ConditionalExpr",
			Content: `
resource "azurerm_resource_group" "test" {
  name = var.name_flag != "" ? var.name_true : var.name_false
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`name` Variable:`name_true` is invalid, expected var name:${prefix}name",
				},
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`name` Variable:`name_false` is invalid, expected var name:${prefix}name",
				},
			},
		},
		{
			Name: "3. a variable name with different value or partial value from property is invalid",
			Content: `
resource "azurerm_resource_group" "test" {
  location = var.address
  tags     = var.tag
}`,
			Expected: helper.Issues{
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`location` Variable:`address` is invalid, expected var name:${prefix}location",
				},
				{
					Rule:    NewTerraformVarNameConventionRule(),
					Message: "Property:`tags` Variable:`tag` is invalid, expected var name:${prefix}tags",
				},
			},
		},
		{
			Name: "4. correct cases",
			Content: `
resource "azurerm_resource_group" "test" {
  # A valid variable name is same as property name or property name with prefix
  name = var.resource_group_name
  location = var.location
  tags = var.test_tags
}

resource "azurerm_var_validation" "test" {
  # Direct usage
  name = var.test_name

  # Nested
  location = var.nested.location

  # List
  list = var.test_list[2]

  # Nested list
  region = var.nested_list[1].region

  # Nested property
  address {
    street = var.test_street
  }

  # Vars inside List are excluded
  ip_addresses = ["0.0.0.0", var.ip_addr_1]

  # Vars inside TypeMap are excluded
  tags = {
    a = var.inner_type_map_test
  }

  # Vars in string template are excluded
  string_template = "Test_${var.str_template_test}"

  # Vars in function call are excluded
  sub_string = substr(var.sub_str_test, 0, 1)

  # For conditional expression, vars in condition part are excluded. Only vars in true/false value part are included, and excluded pattern like string template are still excluded.
  condition_prop = var.condition_flag != "" ? var.true_condition_prop : "Test_${false_prop}"

  # Built-in properties are excluded
  # https://www.terraform.io/language/resources/syntax#meta-arguments and https://www.terraform.io/language/resources/syntax#operation-timeouts
  depends_on = var.depends
  timeouts {
    read = var.timeout_read_1
  }
  lifecycle {
    ignore_changes = var.ignored_props
  }
}`,
			Expected: helper.Issues{},
		},
	}
	rule := NewRule(NewTerraformVarNameConventionRule())
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}
			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
