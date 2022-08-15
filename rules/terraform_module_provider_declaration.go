package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

// TerraformModuleProviderDeclarationRule checks whether local variables are sorted in alphabetic order
type TerraformModuleProviderDeclarationRule struct {
	tflint.DefaultRule
}

// NewTerraformModuleProviderDeclarationRule returns a new rule
func NewTerraformModuleProviderDeclarationRule() *TerraformModuleProviderDeclarationRule {
	return &TerraformModuleProviderDeclarationRule{}
}

// Name returns the rule name
func (r *TerraformModuleProviderDeclarationRule) Name() string {
	return "terraform_module_provider_declaration"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformModuleProviderDeclarationRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformModuleProviderDeclarationRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformModuleProviderDeclarationRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether local variables are sorted in alphabetic order
func (r *TerraformModuleProviderDeclarationRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformModuleProviderDeclarationRule) checkFile(runner tflint.Runner, file *hcl.File) error {
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
