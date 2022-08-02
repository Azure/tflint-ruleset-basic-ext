package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

// TerraformVariableSeparateRule checks whether the variables are separated from other types of blocks
type TerraformVariableSeparateRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableSeparateRule returns a new rule
func NewTerraformVariableSeparateRule() *TerraformVariableSeparateRule {
	return &TerraformVariableSeparateRule{}
}

// Name returns the rule name
func (r *TerraformVariableSeparateRule) Name() string {
	return "terraform_variable_separate"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariableSeparateRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVariableSeparateRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariableSeparateRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are separated from other types of blocks
func (r *TerraformVariableSeparateRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkVariableSeparate(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformVariableSeparateRule) checkVariableSeparate(runner tflint.Runner, file *hcl.File) error {

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
			if firstNonVarBlockRange.Filename == "" {
				firstNonVarBlockRange = block.DefRange()
			}
		}
	}

	if variableDefined && firstNonVarBlockRange.Filename != "" {
		err := runner.EmitIssue(
			r,
			"Putting variables and other types of blocks in the same file is not recommended",
			firstNonVarBlockRange,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
