# azurerm_arg_order

Recommend proper order for variable blocks
The variables with default value set are placed prior to those without default values
Then the variables are sorted based on their names (alphabetic order)

## Example

```hcl
terraform{}

variable "image_id" {
  type = string
}

variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "availability_zone_names" {
  type    = list(string)
  default = ["us-west-1a"]
}
```

```
$ tflint
2 issue(s) found:

Notice: Putting variables and other types of blocks in the same file is not recommended (terraform_variable_order)

  on main.tf line 1:
   1: terraform{}

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_variable_order.md

Notice: Recommended variable order:
variable "docker_ports" {
  type = list(object({
    internal = number
    external = number
    protocol = string
  }))
  default = [
    {
      internal = 8300
      external = 8300
      protocol = "tcp"
    }
  ]
}

variable "image_id" {
  type = string
} (terraform_variable_order)

  on main.tf line 3:
   3: variable "image_id" {

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_variable_order.md
```

## Why
It helps to improve the readability of terraform code by sorting variable blocks in the order above.

## How To Fix
If the rule notifies that it's not recommended to put the variable blocks and other types of blocks in the same file, then consider putting the variable blocks in a seperate file.
If the rule recommends variable order, then just copy the text with recommended variable order and paste it in the tf config file to overwrite the original style of it.