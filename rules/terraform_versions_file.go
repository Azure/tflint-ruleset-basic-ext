package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

// TerraformVersionsFileRule checks whether `versions.tf` only has 1 `terraform` block
type TerraformVersionsFileRule struct {
	tflint.DefaultRule
}

// NewTerraformVersionsFileRule returns a new rule
func NewTerraformVersionsFileRule() *TerraformVersionsFileRule {
	return &TerraformVersionsFileRule{}
}

// Name returns the rule name
func (r *TerraformVersionsFileRule) Name() string {
	return "terraform_versions_file"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVersionsFileRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVersionsFileRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVersionsFileRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether `versions.tf` only has 1 `terraform` block
func (r *TerraformVersionsFileRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	if file, versionsFileExists := files["versions.tf"]; versionsFileExists {
		blocks := file.Body.(*hclsyntax.Body).Blocks
		if len(blocks) != 1 || blocks[0].Type != "terraform" {
			return runner.EmitIssue(
				r,
				"`versions.tf` should have and only have 1 `terraform` block",
				hcl.Range{},
			)
		}
	}
	return nil
}
