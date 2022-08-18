package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformSensitiveVariableNoDefaultRule checks whether default value is set for sensitive variables
type TerraformSensitiveVariableNoDefaultRule struct {
	DefaultRule
}

// NewTerraformSensitiveVariableNoDefaultRule returns a new rule
func NewTerraformSensitiveVariableNoDefaultRule() *TerraformSensitiveVariableNoDefaultRule {
	r := &TerraformSensitiveVariableNoDefaultRule{}
	r.DefaultRule.Rulename = r.Name()
	r.DefaultRule.CheckFile = r.CheckFile
	return r
}

// Name returns the rule name
func (r *TerraformSensitiveVariableNoDefaultRule) Name() string {
	return "terraform_sensitive_variable_no_default"
}

// Severity returns the rule severity
func (r *TerraformSensitiveVariableNoDefaultRule) Severity() tflint.Severity {
	return tflint.WARNING
}

func (r *TerraformSensitiveVariableNoDefaultRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
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
