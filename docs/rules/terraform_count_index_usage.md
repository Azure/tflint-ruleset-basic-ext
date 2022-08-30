# terraform_count_index_usage

Check whether count.index is used as subscript of list/map or the argument of function call

## Example

```hcl
resource "null_resource" "default" {
  count = length(var.my_list)

  triggers = {
    list_index = count.index
    list_value = var.my_list[count.index]
  }
}
```

```
$ tflint
1 issue(s) found:

Warning: `count.index` is not recommended to be used as the subscript of list/map or the argument of function call, use for_each instead (terraform_count_index_usage)

  on main.tf line 6:
  6:     list_value = var.my_list[count.index]

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_count_index_usage.md
```

## Why
see https://medium.com/@business_99069/terraform-count-vs-for-each-b7ada2c0b186

## How To Fix
Consider use for_each to traverse list/map and call the function