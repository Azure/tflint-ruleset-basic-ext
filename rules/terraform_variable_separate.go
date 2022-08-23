package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformVariableSeparateRule checks whether the variables are separated from other types of blocks
type TerraformVariableSeparateRule struct {
	DefaultRule
}

// NewTerraformVariableSeparateRule returns a new rule
func NewTerraformVariableSeparateRule() *TerraformVariableSeparateRule {
	return &TerraformVariableSeparateRule{}
}

// Name returns the rule name
func (r *TerraformVariableSeparateRule) Name() string {
	return "terraform_variable_separate"
}

func (r *TerraformVariableSeparateRule) CheckFile(runner tflint.Runner, file *hcl.File) error {

	blocks := file.Body.(*hclsyntax.Body).Blocks

	var firstNonVarBlockRange hcl.Range
	variableDefined := false
	for _, block := range blocks {
		switch block.Type {
		case "variable":
			if !variableDefined {
				variableDefined = true
			}
		default:
			if IsRangeEmpty(firstNonVarBlockRange) {
				firstNonVarBlockRange = block.DefRange()
			}
		}
	}

	if variableDefined && !IsRangeEmpty(firstNonVarBlockRange) {
		return runner.EmitIssue(
			r,
			"Putting variables and other types of blocks in the same file is not recommended",
			firstNonVarBlockRange,
		)
	}
	return nil
}
