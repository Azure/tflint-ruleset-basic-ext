# azurerm_arg_order

Recommend proper order for output blocks
The outputs are sorted based on their names (alphabetic order)

## Example

```hcl
terraform{}

output "instance_ip_addr" {
  value       = aws_instance.server.private_ip
  description = "The private IP address of the main server instance."
}

output "db_password" {
  value       = aws_db_instance.db.password
  description = "The password for logging in to the database."
  sensitive   = true
}

output "api_base_url" {
  value = "https://${aws_instance.example.private_dns}:8433/"

  # The EC2 instance must have an encrypted root volume.
  precondition {
    condition     = data.aws_ebs_volume.example.encrypted
    error_message = "The server's root volume is not encrypted."
  }
}
```

```
$ tflint
2 issue(s) found:

Notice: Putting outputs and other types of block in the same file is not recommended (terraform_output_order)

  on main.tf line 1:
   1: terraform{}

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_output_order.md

Notice: Recommended output order:
output "api_base_url" {
  value = "https://${aws_instance.example.private_dns}:8433/"

  # The EC2 instance must have an encrypted root volume.
  precondition {
    condition     = data.aws_ebs_volume.example.encrypted
    error_message = "The server's root volume is not encrypted."
  }
}

output "db_password" {
  value       = aws_db_instance.db.password
  description = "The password for logging in to the database."
  sensitive   = true
}

output "instance_ip_addr" {
  value       = aws_instance.server.private_ip
  description = "The private IP address of the main server instance."
} (terraform_output_order)

  on main.tf line 3:
   3: output "instance_ip_addr" {

Reference: https://github.com/terraform-linters/tflint-ruleset-basic-ext/blob/v0.0.1/docs/rules/terraform_output_order.md
```

## Why
It helps to improve the readability of terraform code by sorting output blocks in the order above.

## How To Fix
If the rule notifies that it's not recommended to put the output blocks and other types of blocks in the same file, then consider putting the output blocks in a seperate file.
If the rule recommends output order, then just copy the text with recommended output order and paste it in the tf config file to overwrite the original style of it.