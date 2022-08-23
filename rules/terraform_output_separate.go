package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformOutputSeparateRule checks whether the outputs are separated from other types of blocks
type TerraformOutputSeparateRule struct {
	DefaultRule
}

// NewTerraformOutputSeparateRule returns a new rule
func NewTerraformOutputSeparateRule() *TerraformOutputSeparateRule {
	return &TerraformOutputSeparateRule{}
}

// Name returns the rule name
func (r *TerraformOutputSeparateRule) Name() string {
	return "terraform_output_separate"
}

func (r *TerraformOutputSeparateRule) CheckFile(runner tflint.Runner, file *hcl.File) error {

	blocks := file.Body.(*hclsyntax.Body).Blocks

	var firstNonOutputBlockRange hcl.Range
	outputDefined := false
	for _, block := range blocks {
		switch block.Type {
		case "output":
			if !outputDefined {
				outputDefined = true
			}
		default:
			if IsRangeEmpty(firstNonOutputBlockRange) {
				firstNonOutputBlockRange = block.DefRange()
			}
		}
	}

	if outputDefined && !IsRangeEmpty(firstNonOutputBlockRange) {
		return runner.EmitIssue(
			r,
			"Putting outputs and other types of blocks in the same file is not recommended",
			firstNonOutputBlockRange,
		)
	}
	return nil
}
