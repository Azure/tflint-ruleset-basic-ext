package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"sort"
	"strings"
)

// TerraformLocalsOrderRule checks whether comments use the preferred syntax
type TerraformLocalsOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformLocalsOrderRule returns a new rule
func NewTerraformLocalsOrderRule() *TerraformLocalsOrderRule {
	return &TerraformLocalsOrderRule{}
}

// Name returns the rule name
func (r *TerraformLocalsOrderRule) Name() string {
	return "terraform_locals_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformLocalsOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformLocalsOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformLocalsOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether single line comments is used
func (r *TerraformLocalsOrderRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformLocalsOrderRule) checkFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	var err error
	for _, block := range blocks {
		switch block.Type {
		case "locals":
			if subErr := r.checkLocalsOrder(runner, block); subErr != nil {
				err = multierror.Append(subErr)
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
	runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended locals variable order:\n%s", localsHclTxt),
		block.DefRange(),
	)
	return nil
}
