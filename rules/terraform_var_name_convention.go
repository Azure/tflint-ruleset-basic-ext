package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"strings"
)

var ignoredVariables = map[string]struct{}{}

var builtInProperties = map[string]struct{}{
	// Built-in Properties are from https://www.terraform.io/language/resources/syntax#meta-arguments and https://www.terraform.io/language/resources/syntax#operation-timeouts
	"depends_on":  {},
	"count":       {},
	"for_each":    {},
	"provider":    {},
	"lifecycle":   {},
	"provisioner": {},
	"timeouts":    {},
}

// TerraformVarNameConventionRule checks whether the var name is valid
type TerraformVarNameConventionRule struct {
	DefaultRule
}

// NewTerraformVarNameConventionRule returns a new rule
func NewTerraformVarNameConventionRule() *TerraformVarNameConventionRule {
	return &TerraformVarNameConventionRule{}
}

// Name returns the rule name
func (r *TerraformVarNameConventionRule) Name() string {
	return "terraform_var_name_convention"
}

// Severity returns the rule severity
func (r *TerraformVarNameConventionRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// CheckFile checks whether the var name is valid
func (r *TerraformVarNameConventionRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	var err error
	blocks := file.Body.(*hclsyntax.Body).Blocks
	for _, block := range blocks {
		switch block.Type {
		case "data", "resource":
			if subErr := r.visitBlock(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func (r *TerraformVarNameConventionRule) visitBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	var err error
	for _, attribute := range block.Body.Attributes {
		expression := attribute.Expr
		switch expression.(type) {
		// these 2 cases are skipped for now
		case *hclsyntax.FunctionCallExpr:
			// Function call
		case *hclsyntax.ObjectConsExpr:
			// TypeMap
		case *hclsyntax.ConditionalExpr:
			if subErr := r.validateExpression(runner, expression.(*hclsyntax.ConditionalExpr).TrueResult, attribute.Name); subErr != nil {
				err = multierror.Append(err, subErr)
			}
			if subErr := r.validateExpression(runner, expression.(*hclsyntax.ConditionalExpr).FalseResult, attribute.Name); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		case *hclsyntax.ScopeTraversalExpr:
			if subErr := r.validateExpression(runner, expression, attribute.Name); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}

	for _, nestedBlock := range block.Body.Blocks {
		if _, ok := builtInProperties[nestedBlock.Type]; !ok {
			if subErr := r.visitBlock(runner, nestedBlock); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}

	return err
}

func (r *TerraformVarNameConventionRule) validateExpression(runner tflint.Runner, expression hclsyntax.Expression, propertyName string) error {
	if _, ok := builtInProperties[propertyName]; !ok {
		variables := expression.Variables()
		if len(variables) == 1 {
			variable := variables[0]
			if variable.RootName() == "var" {
				variableName := ""
				// Find last named variable in statement like a[0].b[1]
				for i := len(variable) - 1; i >= 0 && variableName == ""; i-- {
					switch variable[i].(type) {
					case hcl.TraverseAttr:
						variableName = variable[i].(hcl.TraverseAttr).Name
					}
				}
				return r.validate(runner, variableName, propertyName, variable[0].SourceRange())
			}
		}
	}
	return nil
}

func (r *TerraformVarNameConventionRule) validate(runner tflint.Runner, variableName string, propertyName string, location hcl.Range) error {
	var err error

	if variableName == "" || propertyName == "" {
		if subErr := runner.EmitIssue(r, "Variable or Property name is empty", location); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}

	if index := strings.LastIndex(variableName, propertyName); index == -1 || index != len(variableName)-len(propertyName) {
		var subErr error
		if _, isIgnore := ignoredVariables[variableName]; !isIgnore {
			subErr = runner.EmitIssue(r, fmt.Sprintf("Property:`%s` Variable:`%s` is invalid, expected var name:${prefix}%s", propertyName, variableName, propertyName), location)
		}
		if subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
