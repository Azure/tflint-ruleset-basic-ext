package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"sort"
	"strings"
)

// TerraformLocalsOrderRule checks whether local variables are sorted in alphabetic order
type TerraformLocalsOrderRule struct {
	DefaultRule
}

// NewTerraformLocalsOrderRule returns a new rule
func NewTerraformLocalsOrderRule() *TerraformLocalsOrderRule {
	return &TerraformLocalsOrderRule{}
}

// Name returns the rule name
func (r *TerraformLocalsOrderRule) Name() string {
	return "terraform_locals_order"
}

func (r *TerraformLocalsOrderRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		switch block.Type {
		case "locals":
			if subErr := r.checkLocalsOrder(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func (r *TerraformLocalsOrderRule) checkLocalsOrder(runner tflint.Runner, block *hclsyntax.Block) error {
	file, err := runner.GetFile(block.Range().Filename)
	if err != nil {
		return err
	}
	var attrNames, localsHclTxts []string
	attrStartPos := make(map[string]hcl.Pos)
	attrHclTxts := make(map[string]string)
	for attrName, attr := range block.Body.Attributes {
		attrNames = append(attrNames, attrName)
		attrStartPos[attrName] = attr.NameRange.Start
		attrHclTxts[attrName] = string(attr.SrcRange.SliceBytes(file.Bytes))
	}
	sort.Slice(attrNames, func(i, j int) bool {
		if attrStartPos[attrNames[i]].Line == attrStartPos[attrNames[j]].Line {
			return attrStartPos[attrNames[i]].Column < attrStartPos[attrNames[j]].Column
		}
		return attrStartPos[attrNames[i]].Line < attrStartPos[attrNames[j]].Line
	})
	if sort.StringsAreSorted(attrNames) {
		return nil
	}
	sort.Strings(attrNames)
	for _, attrName := range attrNames {
		localsHclTxts = append(localsHclTxts, attrHclTxts[attrName])
	}
	localsHclTxt := strings.Join(localsHclTxts, "\n")
	localsHclTxt = fmt.Sprintf("%s {\n%s\n}", block.Type, localsHclTxt)
	localsHclTxt = string(hclwrite.Format([]byte(localsHclTxt)))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended locals variable order:\n%s", localsHclTxt),
		block.DefRange(),
	)
}
