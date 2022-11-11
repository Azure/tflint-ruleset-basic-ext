package rules

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformSensitiveVariableNoDefaultRule{}

// TerraformSensitiveVariableNoDefaultRule checks whether default value is set for sensitive variables
type TerraformSensitiveVariableNoDefaultRule struct {
	tflint.DefaultRule
}

// NewTerraformSensitiveVariableNoDefaultRule returns a new rule
func NewTerraformSensitiveVariableNoDefaultRule() *TerraformSensitiveVariableNoDefaultRule {
	return &TerraformSensitiveVariableNoDefaultRule{}
}

func (r *TerraformSensitiveVariableNoDefaultRule) Enabled() bool {
	return false
}

func (r *TerraformSensitiveVariableNoDefaultRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
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
		if block.Type != "variable" {
			continue
		}
		sensitive := false
		if attr, sensitiveSet := block.Body.Attributes["sensitive"]; sensitiveSet {
			val, diags := attr.Expr.Value(nil)
			if diags.HasErrors() {
				err = multierror.Append(err, diags)
			}
			sensitive = val.True()
		}
		if sensitive && withNonNullDefaultValue(block, file) {
			subErr := runner.EmitIssue(
				r,
				fmt.Sprintf("Default value is not expected to be set for sensitive variable `%s`", block.Labels[0]),
				block.Body.Attributes["default"].NameRange,
			)
			if subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func withNonNullDefaultValue(b *hclsyntax.Block, f *hcl.File) bool {
	attr, set := b.Body.Attributes["default"]
	if !set {
		return false
	}
	value := hcl.ExprAsKeyword(attr.Expr)
	return value != "null"
}
