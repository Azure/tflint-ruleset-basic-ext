package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// TerraformCountIndexUsageRule checks whether count.index is used as subscript of list/map
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

func (r *TerraformCountIndexUsageRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		if subErr := r.visitBlock(runner, block); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformCountIndexUsageRule) visitBlock(runner tflint.Runner, block *hclsyntax.Block) error {
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

func (r *TerraformCountIndexUsageRule) visitExp(runner tflint.Runner, exp hclsyntax.Expression) error {
	file, _ := runner.GetFile(exp.Range().Filename)
	tokens, diags := hclsyntax.LexExpression(exp.Range().SliceBytes(file.Bytes), exp.Range().Filename, exp.StartRange().Start)
	if diags.HasErrors() {
		return diags
	}
	var err error
	subscriptLevel := 0
	for i, token := range tokens {
		switch token.Type {
		case hclsyntax.TokenOBrack:
			subscriptLevel++
		case hclsyntax.TokenCBrack:
			subscriptLevel--
		case hclsyntax.TokenIdent:
			first := string(token.Bytes)
			if subscriptLevel > 0 && first == "count" && i+2 < len(tokens) {
				second := string(tokens[i+1].Bytes)
				thrid := string(tokens[i+2].Bytes)
				if second == "." && thrid == "index" {
					subErr := runner.EmitIssue(
						r,
						"`count.index` is not recommended to be used as the subscript of list/map, use for_each instead",
						hcl.Range{
							Filename: token.Range.Filename,
							Start:    token.Range.Start,
							End:      tokens[i+2].Range.End,
						},
					)
					if subErr != nil {
						err = multierror.Append(err, subErr)
					}
				}
			}
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
//		if x.Type().String() == "*hclsyntax.IndexExpr" {
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
