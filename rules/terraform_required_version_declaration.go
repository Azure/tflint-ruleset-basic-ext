package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

var _ tflint.Rule = &TerraformRequiredVersionDeclarationRule{}

// TerraformRequiredVersionDeclarationRule checks whether required_version field is declared at the beginning of terraform block
type TerraformRequiredVersionDeclarationRule struct {
	tflint.DefaultRule
}

// NewTerraformRequiredVersionDeclarationRule returns a new rule
func NewTerraformRequiredVersionDeclarationRule() *TerraformRequiredVersionDeclarationRule {
	return &TerraformRequiredVersionDeclarationRule{}
}

func (r *TerraformRequiredVersionDeclarationRule) Enabled() bool {
	return false
}

func (r *TerraformRequiredVersionDeclarationRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

func (r *TerraformRequiredVersionDeclarationRule) Check(runner tflint.Runner) error {
	return ForFiles(runner, r.CheckFile)
}

// Name returns the rule name
func (r *TerraformRequiredVersionDeclarationRule) Name() string {
	return "terraform_required_version_declaration"
}

func (r *TerraformRequiredVersionDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	var err error
	blocks := file.Body.(*hclsyntax.Body).Blocks
	for _, block := range blocks {
		switch block.Type {
		case "terraform":
			if subErr := r.checkTerraformRequiredVersionDeclaration(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func (r *TerraformRequiredVersionDeclarationRule) checkTerraformRequiredVersionDeclaration(runner tflint.Runner, block *hclsyntax.Block) error {
	comparePos := func(pos1 hcl.Pos, pos2 hcl.Pos) bool {
		if pos1.Line != pos2.Line {
			return pos1.Line < pos2.Line
		}
		return pos1.Line < pos2.Line
	}
	msg := "The `required_version` field should be declared at the beginning of `terraform` block"
	requiredVersionAttr, requiredVersionDefined := block.Body.Attributes["required_version"]
	if !requiredVersionDefined {
		return runner.EmitIssue(
			r,
			msg,
			block.DefRange(),
		)
	}
	for _, attr := range block.Body.Attributes {
		if attr.Name != "required_version" && comparePos(attr.Range().Start, requiredVersionAttr.Range().Start) {
			return runner.EmitIssue(
				r,
				msg,
				requiredVersionAttr.NameRange,
			)
		}
	}
	for _, nestedBlock := range block.Body.Blocks {
		if comparePos(nestedBlock.Range().Start, requiredVersionAttr.Range().Start) {
			return runner.EmitIssue(
				r,
				msg,
				requiredVersionAttr.NameRange,
			)
		}
	}
	return nil
}
