package rules

import (
	"fmt"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"sort"
	"strings"
)

// TerraformOutputOrderRule checks whether the outputs are sorted in expected order
type TerraformOutputOrderRule struct {
	DefaultRule
}

// NewTerraformOutputOrderRule returns a new rule
func NewTerraformOutputOrderRule() *TerraformOutputOrderRule {
	return &TerraformOutputOrderRule{}
}

// Name returns the rule name
func (r *TerraformOutputOrderRule) Name() string {
	return "terraform_output_order"
}

func (r *TerraformOutputOrderRule) CheckFile(runner tflint.Runner, file *hcl.File) error {

	blocks := file.Body.(*hclsyntax.Body).Blocks

	var outputNames []string
	var firstOutputBlockRange hcl.Range
	outputHclTxts := make(map[string]string)
	for _, block := range blocks {
		switch block.Type {
		case "output":
			if IsRangeEmpty(firstOutputBlockRange) {
				firstOutputBlockRange = block.DefRange()
			}
			outputName := block.Labels[0]
			outputNames = append(outputNames, outputName)
			outputHclTxts[outputName] = string(block.Range().SliceBytes(file.Bytes))
		}
	}

	if sort.StringsAreSorted(outputNames) {
		return nil
	}
	sort.Strings(outputNames)
	var sortedOutputHclTxts []string
	for _, outputName := range outputNames {
		sortedOutputHclTxts = append(sortedOutputHclTxts, outputHclTxts[outputName])
	}
	sortedOutputHclBytes := hclwrite.Format([]byte(strings.Join(sortedOutputHclTxts, "\n\n")))
	return runner.EmitIssue(
		r,
		fmt.Sprintf("Recommended output order:\n%s", sortedOutputHclBytes),
		firstOutputBlockRange,
	)
}
