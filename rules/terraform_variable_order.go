package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"reflect"
	"sort"
	"strings"
)

// TerraformVariableOrderRule checks whether the variables are sorted in expected order
type TerraformVariableOrderRule struct {
	tflint.DefaultRule
}

// NewTerraformVariableOrderRule returns a new rule
func NewTerraformVariableOrderRule() *TerraformVariableOrderRule {
	return &TerraformVariableOrderRule{}
}

// Name returns the rule name
func (r *TerraformVariableOrderRule) Name() string {
	return "terraform_variable_order"
}

// Enabled returns whether the rule is enabled by default
func (r *TerraformVariableOrderRule) Enabled() bool {
	return false
}

// Severity returns the rule severity
func (r *TerraformVariableOrderRule) Severity() tflint.Severity {
	return tflint.NOTICE
}

// Link returns the rule reference link
func (r *TerraformVariableOrderRule) Link() string {
	return project.ReferenceLink(r.Name())
}

// Check checks whether the variables are sorted in expected order
func (r *TerraformVariableOrderRule) Check(runner tflint.Runner) error {

	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for _, file := range files {
		if subErr := r.checkVariableOrder(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *TerraformVariableOrderRule) checkVariableOrder(runner tflint.Runner, file *hcl.File) error {
	getSortedVariableNames := func(blocks hclsyntax.Blocks, hasDefaultVal bool) []string {
		var sortedVariableNames []string
		for _, block := range blocks {
			switch block.Type {
			case "variable":
				if _, variableHasDefaultVal := block.Body.Attributes["default"]; variableHasDefaultVal == hasDefaultVal {
					sortedVariableNames = append(sortedVariableNames, block.Labels[0])
				}
			}
		}
		sort.Strings(sortedVariableNames)
		return sortedVariableNames
	}

	blocks := file.Body.(*hclsyntax.Body).Blocks
	sortedVariableNames := getSortedVariableNames(blocks, false)
	sortedVariableNames = append(sortedVariableNames, getSortedVariableNames(blocks, true)...)

	var variableNames, sortedVariableHclTxts []string
	var firstVarBlockRange hcl.Range
	variableHclTxts := make(map[string]string)
	for _, block := range blocks {
		switch block.Type {
		case "variable":
			if firstVarBlockRange.Filename == "" {
				firstVarBlockRange = block.DefRange()
			}
			variableName := block.Labels[0]
			variableNames = append(variableNames, variableName)
			variableHclTxts[variableName] = string(block.Range().SliceBytes(file.Bytes))
		}
	}

	if !reflect.DeepEqual(variableNames, sortedVariableNames) {
		for _, variableName := range sortedVariableNames {
			sortedVariableHclTxts = append(sortedVariableHclTxts, variableHclTxts[variableName])
		}
		sortedVariableHclBytes := hclwrite.Format([]byte(strings.Join(sortedVariableHclTxts, "\n\n")))
		err := runner.EmitIssue(
			r,
			fmt.Sprintf("Recommended variable order:\n%s", sortedVariableHclBytes),
			firstVarBlockRange,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
