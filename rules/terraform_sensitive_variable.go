package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

type TerraformSensitiveVariableRule struct {
	tflint.DefaultRule
}

// NewTerraformSensitiveVariableRule returns a new rule
func NewTerraformSensitiveVariableRule() *TerraformSensitiveVariableRule {
	return &TerraformSensitiveVariableRule{}
}

// Name returns the rule name
func (r *TerraformSensitiveVariableRule) Name() string {
	return "terraform_sensitive_variable"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformSensitiveVariableRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformSensitiveVariableRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformSensitiveVariableRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether single line comments is used
func (r *TerraformSensitiveVariableRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkSensitiveVar(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformSensitiveVariableRule) checkSensitiveVar(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		switch block.Type {
		case "variable":
			isSensitive := false
			if sensitiveAttr, isSensitiveSet := block.Body.Attributes["sensitive"]; isSensitiveSet {
				val, diags := sensitiveAttr.Expr.Value(nil)
				if diags.HasErrors() {
					err = multierror.Append(err, diags)
				}
				isSensitive = val.True()
			}
			if defaultAttr, isDefaultSet := block.Body.Attributes["default"]; isSensitive && isDefaultSet {
				subErr := runner.EmitIssue(
					r,
					fmt.Sprintf("Default value is not expected to be set for sensitive variable `%s`", block.Labels[0]),
					defaultAttr.NameRange,
				)
				if subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}
	}
	return err
}
