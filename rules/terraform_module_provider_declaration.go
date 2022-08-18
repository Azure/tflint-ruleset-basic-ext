package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformModuleProviderDeclarationRule checks whether local variables are sorted in alphabetic order
type TerraformModuleProviderDeclarationRule struct {
	DefaultRule
}

// NewTerraformModuleProviderDeclarationRule returns a new rule
func NewTerraformModuleProviderDeclarationRule() *TerraformModuleProviderDeclarationRule {
	r := &TerraformModuleProviderDeclarationRule{}
	r.DefaultRule.Rulename = r.Name()
	r.DefaultRule.CheckFile = r.CheckFile
	return r
}

// Name returns the rule name
func (r *TerraformModuleProviderDeclarationRule) Name() string {
	return "terraform_module_provider_declaration"
}

// Severity returns the rule severity
func (r *TerraformModuleProviderDeclarationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformModuleProviderDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		switch block.Type {
		case "provider":
			isProviderErrorUsage := false
			if len(block.Body.Attributes) != 1 {
				isProviderErrorUsage = true
			} else if _, isAliasDeclared := block.Body.Attributes["alias"]; !isAliasDeclared {
				isProviderErrorUsage = true
			}
			if len(block.Body.Blocks) > 0 {
				isProviderErrorUsage = true
			}
			if isProviderErrorUsage {
				subErr := runner.EmitIssue(
					r,
					"Provider block in terraform module is expected to have and only have `alias` declared",
					block.DefRange(),
				)
				if subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}
	}
	return err
}
