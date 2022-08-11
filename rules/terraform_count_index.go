package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"go/ast"
	"reflect"
)

// TerraformCountIndexRule checks whether count.index is used as subscript of list/map
type TerraformCountIndexRule struct {
	tflint.DefaultRule
}

// NewTerraformCountIndexRule returns a new rule
func NewTerraformCountIndexRule() *TerraformCountIndexRule {
	return &TerraformCountIndexRule{}
}

// Name returns the rule name
func (r *TerraformCountIndexRule) Name() string {
	return "terraform_count_index"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformCountIndexRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformCountIndexRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// Link returns the rule reference link
func (r *TerraformCountIndexRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether count.index is used as subscript of list/map
func (r *TerraformCountIndexRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.visitFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexRule) visitFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		if subErr := r.visitBlock(runner, block); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexRule) visitBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	var err error
	for _, attr := range block.Body.Attributes {
		if subErr := r.visitExp(runner, attr.Expr); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	for _, nestedBlock := range block.Body.Blocks {
		if subErr := r.visitBlock(runner, nestedBlock); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexRule) visitExp(runner tflint.Runner, exp hclsyntax.Expression) error {
	return r.diagnose(runner, reflect.ValueOf(exp), false)
}

func (r *TerraformCountIndexRule) diagnose(runner tflint.Runner, x reflect.Value, asIndex bool) error {
	if !ast.NotNilFilter("", x) {
		return nil
	}

	var err error

	switch x.Kind() {

	case reflect.Interface:
		err = r.diagnose(runner, x.Elem(), asIndex)

	case reflect.Map:
		if x.Len() > 0 {
			for _, key := range x.MapKeys() {
				if subErr := r.diagnose(runner, x.MapIndex(key), asIndex); subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}

	case reflect.Pointer:
		if x.Type().String() == "*hclsyntax.IndexExpr" {
			asIndex = true
		} else if asIndex && x.Type().String() == "*hclsyntax.ScopeTraversalExpr" {
			traversal, isTraversal := x.Elem().FieldByName("Traversal").Interface().(hcl.Traversal)
			if !isTraversal || len(traversal) < 2 {
				return nil
			}
			root, isRoot := traversal[0].(hcl.TraverseRoot)
			attr, isAttr := traversal[1].(hcl.TraverseAttr)
			if isRoot && isAttr && root.Name == "count" && attr.Name == "index" {
				runner.EmitIssue(
					r,
					"`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
					traversal.SourceRange(),
				)
			}
			return nil
		}
		err = r.diagnose(runner, x.Elem(), asIndex)

	case reflect.Array:
		if x.Len() > 0 {
			for i, n := 0, x.Len(); i < n; i++ {
				if subErr := r.diagnose(runner, x.Index(i), asIndex); subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}

	case reflect.Slice:
		if _, ok := x.Interface().([]byte); ok {
			return nil
		}
		if x.Len() > 0 {
			for i, n := 0, x.Len(); i < n; i++ {
				if subErr := r.diagnose(runner, x.Index(i), asIndex); subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}

	case reflect.Struct:
		t := x.Type()
		for i, n := 0, t.NumField(); i < n; i++ {
			if name := t.Field(i).Name; ast.IsExported(name) {
				if subErr := r.diagnose(runner, x.Field(i), asIndex); subErr != nil {
					err = multierror.Append(err, subErr)
				}
			}
		}
	}
	return err
}
