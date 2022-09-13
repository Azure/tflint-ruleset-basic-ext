package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"reflect"
	"sort"
	"strings"
)

// TerraformVariableOrderRule checks whether the variables are sorted in expected order
type TerraformVariableOrderRule struct {
	DefaultRule
}

// NewTerraformVariableOrderRule returns a new rule
func NewTerraformVariableOrderRule() *TerraformVariableOrderRule {
	return &TerraformVariableOrderRule{}
}

// Name returns the rule name
func (r *TerraformVariableOrderRule) Name() string {
	return "terraform_variable_order"
}

// CheckFile checks whether the variables are sorted in expected order
func (r *TerraformVariableOrderRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	blocks := file.Body.(*hclsyntax.Body).Blocks
	sortedVariableNames := r.getSortedVariableNames(blocks, false)
	sortedVariableNames = append(sortedVariableNames, r.getSortedVariableNames(blocks, true)...)

	var variableNames []string
	var firstVarBlockRange hcl.Range
	variableHclTxts := make(map[string]string)
	for _, block := range blocks {
		switch block.Type {
		case "variable":
			if IsRangeEmpty(firstVarBlockRange) {
				firstVarBlockRange = block.DefRange()
			}
			variableName := block.Labels[0]
			variableNames = append(variableNames, variableName)
			variableHclTxts[variableName] = string(block.Range().SliceBytes(file.Bytes))
		}
	}

	if reflect.DeepEqual(variableNames, sortedVariableNames) {
		return nil
	}
	var sortedVariableHclTxts []string
	for _, variableName := range sortedVariableNames {
		sortedVariableHclTxts = append(sortedVariableHclTxts, variableHclTxts[variableName])
	}
	sortedVariableHclBytes := hclwrite.Format([]byte(strings.Join(sortedVariableHclTxts, "\n\n")))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended variable order:\n%s", sortedVariableHclBytes),
		firstVarBlockRange,
	)
}

func (r *TerraformVariableOrderRule) getSortedVariableNames(blocks hclsyntax.Blocks, hasDefaultVal bool) []string {
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
