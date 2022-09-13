# terraform_var_name_convention

Check whether the variable name is valid

## Example

```hcl
resource "azurerm_resource_group" "test" {
  name = var.name_rg
  location = "westeurope"
}
```

```
$ tflint
1 issue(s) found:

Warning: Property:`name` Variable:`name_rg` is invalid, expected var name:${prefix}name (terraform_var_name_convention)

  on test.tf line 2:
   2:   name = var.name_rg

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_var_name_convention.md
```

## Why
For ScopeTraversalExpr and ConditionalExpr, the variable name should take the form of ${prefix}${property name}

## How To Fix
Rename the variable to fit the naming convention above