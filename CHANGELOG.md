## v0.0.9

* Add `terraform_count_index_usage` rule
* Add `terraform_heredoc_useage` rule
* Add `terraform_module_provider_declaration` rule
* Add `terraform_required_providers_declaration` rule
* Add `terraform_sensitive_variable_no_default` rule
* Add `terraform_versions_file` rule

## v0.1.0

Fix a nil panic and incorrect check on nested block

## v0.1.1

Fix incorrect reference link

## V0.2.0

Fix install code, add CodeQL and Gosec to ci.

## v0.2.1

`terraform_sensitive_variable_no_default` won't raise error when the default value is `null`.