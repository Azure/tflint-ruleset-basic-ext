package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformCountIndexUsageRule checks whether count.index is used as subscript of list/map or the argument of function call
type TerraformCountIndexUsageRule struct {
	DefaultRule
}

// NewTerraformCountIndexUsageRule returns a new rule
func NewTerraformCountIndexUsageRule() *TerraformCountIndexUsageRule {
	return &TerraformCountIndexUsageRule{}
}

// Name returns the rule name
func (r *TerraformCountIndexUsageRule) Name() string {
	return "terraform_count_index_usage"
}

// Severity returns the rule severity
func (r *TerraformCountIndexUsageRule) Severity() tflint.Severity {
	return tflint.WARNING
}

// CheckFile checks whether count.index is used as subscript of list/map or the argument of function call
func (r *TerraformCountIndexUsageRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	w := &countIndexCheckWalker{
		f:      file,
		runner: runner,
		rule:   r,
	}
	for _, block := range blocks {
		diag := hclsyntax.Walk(block, w)
		if diag.HasErrors() {
			err = multierror.Append(err, diag)
		}
	}
	return err
}

type countIndexCheckWalker struct {
	f      *hcl.File
	issues helper.Issues
	runner tflint.Runner
	rule   *TerraformCountIndexUsageRule
}

func (w *countIndexCheckWalker) Enter(node hclsyntax.Node) hcl.Diagnostics {
	switch node.(type) {
	case *hclsyntax.FunctionCallExpr, *hclsyntax.IndexExpr:
		err := w.checkCountIndex(node)
		if err != nil {
			return hcl.Diagnostics{
				&hcl.Diagnostic{
					Detail:     err.Error(),
					Expression: node.(hcl.Expression),
				},
			}
		}
	}
	return nil
}

func (w *countIndexCheckWalker) Exit(hclsyntax.Node) hcl.Diagnostics {
	return nil
}

func (w *countIndexCheckWalker) checkCountIndex(node hclsyntax.Node) error {
	tokens, diags := hclsyntax.LexExpression(node.Range().SliceBytes(w.f.Bytes), node.Range().Filename, node.Range().Start)
	if diags.HasErrors() {
		return diags
	}
	var err error
	for i := range tokens {
		if !(i+2 < len(tokens) && string(tokens[i].Bytes) == "count" && string(tokens[i+1].Bytes) == "." && string(tokens[i+2].Bytes) == "index") {
			continue
		}
		issue := &helper.Issue{
			Rule:    nil,
			Message: "`count.index` is not recommended to be used as the subscript of list/map or the argument of function call, use for_each instead",
			Range: hcl.Range{
				Filename: tokens[i].Range.Filename,
				Start:    tokens[i].Range.Start,
				End:      tokens[i+2].Range.End,
			},
		}
		isRuleExists := false
		for _, existedIssue := range w.issues {
			if *issue == *existedIssue {
				isRuleExists = true
				break
			}
		}
		if !isRuleExists {
			w.runner.EmitIssue(w.rule, issue.Message, issue.Range)
			w.issues = append(w.issues, issue)
		}
	}
	return err
}

//func (r *TerraformCountIndexUsageRule) diagnose(runner tflint.Runner, x reflect.Value, asIndex bool) error {
//	if !ast.NotNilFilter("", x) {
//		return nil
//	}
//
//	var err error
//
//	switch x.Kind() {
//
//	case reflect.Interface:
//		err = r.diagnose(runner, x.Elem(), asIndex)
//
//	case reflect.Map:
//		if x.Len() > 0 {
//			for _, key := range x.MapKeys() {
//				if subErr := r.diagnose(runner, x.MapIndex(key), asIndex); subErr != nil {
//					err = multierror.Append(err, subErr)
//				}
//			}
//		}
//
//	case reflect.Pointer:
//		if x. == "*hclsyntax.IndexExpr" {
//			asIndex = true
//		} else if asIndex && x.Type().String() == "*hclsyntax.ScopeTraversalExpr" {
//			traversal, isTraversal := x.Elem().FieldByName("Traversal").Interface().(hcl.Traversal)
//			if !isTraversal || len(traversal) < 2 {
//				return nil
//			}
//			root, isRoot := traversal[0].(hcl.TraverseRoot)
//			attr, isAttr := traversal[1].(hcl.TraverseAttr)
//			if isRoot && isAttr && root.Name == "count" && attr.Name == "index" {
//				runner.EmitIssue(
//					r,
//					"`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
//					traversal.SourceRange(),
//				)
//			}
//			return nil
//		}
//		err = r.diagnose(runner, x.Elem(), asIndex)
//
//	case reflect.Array:
//		if x.Len() > 0 {
//			for i, n := 0, x.Len(); i < n; i++ {
//				if subErr := r.diagnose(runner, x.Index(i), asIndex); subErr != nil {
//					err = multierror.Append(err, subErr)
//				}
//			}
//		}
//
//	case reflect.Slice:
//		if _, ok := x.Interface().([]byte); ok {
//			return nil
//		}
//		if x.Len() > 0 {
//			for i, n := 0, x.Len(); i < n; i++ {
//				if subErr := r.diagnose(runner, x.Index(i), asIndex); subErr != nil {
//					err = multierror.Append(err, subErr)
//				}
//			}
//		}
//
//	case reflect.Struct:
//		t := x.Type()
//		for i, n := 0, t.NumField(); i < n; i++ {
//			if name := t.Field(i).Name; ast.IsExported(name) {
//				if subErr := r.diagnose(runner, x.Field(i), asIndex); subErr != nil {
//					err = multierror.Append(err, subErr)
//				}
//			}
//		}
//	}
//	return err
//}
