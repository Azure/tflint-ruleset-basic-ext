package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformVersionsFileRule checks whether `versions.tf` only has 1 `terraform` block
type TerraformVersionsFileRule struct {
	DefaultRule
}

// NewTerraformVersionsFileRule returns a new rule
func NewTerraformVersionsFileRule() *TerraformVersionsFileRule {
	return &TerraformVersionsFileRule{}
}

// Name returns the rule name
func (r *TerraformVersionsFileRule) Name() string {
	return "terraform_versions_file"
}

func (r *TerraformVersionsFileRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	body := file.Body.(*hclsyntax.Body)
	filename := body.Range().Filename
	if filename != "versions.tf" {
		return nil
	}
	blocks := body.Blocks
	if len(blocks) != 1 || blocks[0].Type != "terraform" {
		return runner.EmitIssue(
			r,
			"`versions.tf` should have and only have 1 `terraform` block",
			hcl.Range{},
		)
	}
	return nil
}
