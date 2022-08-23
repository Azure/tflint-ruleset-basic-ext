package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"sort"
	"strings"
)

// TerraformRequiredProvidersDeclarationRule checks whether the required_providers block is declared in terraform block and whether the args of it are sorted in alphabetic order
type TerraformRequiredProvidersDeclarationRule struct {
	DefaultRule
}

// NewTerraformRequiredProvidersDeclarationRule returns a new rule
func NewTerraformRequiredProvidersDeclarationRule() *TerraformRequiredProvidersDeclarationRule {
	return &TerraformRequiredProvidersDeclarationRule{}
}

// Name returns the rule name
func (r *TerraformRequiredProvidersDeclarationRule) Name() string {
	return "terraform_required_providers_declaration"
}

func (r *TerraformRequiredProvidersDeclarationRule) CheckFile(runner tflint.Runner, file *hcl.File) error {
	var err error
	blocks := file.Body.(*hclsyntax.Body).Blocks
	for _, block := range blocks {
		switch block.Type {
		case "terraform":
			if subErr := r.checkTerraformRequiredProvidersDeclaration(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
	}
	return err
}

func (r *TerraformRequiredProvidersDeclarationRule) checkTerraformRequiredProvidersDeclaration(runner tflint.Runner, block *hclsyntax.Block) error {
	isRequiredProvidersDeclared := false
	var err error
	for _, nestedBlock := range block.Body.Blocks {
		switch nestedBlock.Type {
		case "required_providers":
			isRequiredProvidersDeclared = true
			err = multierror.Append(err, r.checkRequiredProvidersArgOrder(runner, nestedBlock))
		}
	}
	if isRequiredProvidersDeclared {
		return nil
	}
	return runner.EmitIssue(
		r,
		"The `required_providers` field should be declared in `terraform` block",
		block.DefRange(),
	)
}

func (r *TerraformRequiredProvidersDeclarationRule) checkRequiredProvidersArgOrder(runner tflint.Runner, block *hclsyntax.Block) error {
	file, _ := runner.GetFile(block.Range().Filename)
	var providerNames []string
	providerParamTxts := make(map[string]string)
	providerParamIssues := helper.Issues{}
	providers := block.Body.Attributes
	for providerName, providerParam := range providers {
		sortedProviderParamTxt, isProviderParamSorted := PrintSortedAttrTxt(file.Bytes, providerParam)
		providerParamTxts[providerName] = sortedProviderParamTxt
		providerNames = append(providerNames, providerName)
		if !isProviderParamSorted {
			providerParamIssues = append(providerParamIssues, &helper.Issue{
				Rule:    r,
				Message: fmt.Sprintf("Parameters of provider `%s` are expected to be sorted as follows:\n%s", providerName, sortedProviderParamTxt),
				Range:   providerParam.NameRange,
			})
		}
	}
	sort.Slice(providerNames, func(i, j int) bool {
		if providers[providerNames[i]].Range().Start.Line == providers[providerNames[j]].Range().Start.Line {
			return providers[providerNames[i]].Range().Start.Column < providers[providerNames[j]].Range().Start.Column
		}
		return providers[providerNames[i]].Range().Start.Line < providers[providerNames[j]].Range().Start.Line
	})
	if !sort.StringsAreSorted(providerNames) {
		sort.Strings(providerNames)
		var sortedProviderParamTxts []string
		for _, providerName := range providerNames {
			sortedProviderParamTxts = append(sortedProviderParamTxts, providerParamTxts[providerName])
		}
		sortedProviderParamTxt := strings.Join(sortedProviderParamTxts, "\n")
		var sortedRequiredProviderTxt string
		if RemoveSpaceAndLine(sortedProviderParamTxt) == "" {
			sortedRequiredProviderTxt = fmt.Sprintf("%s {}", block.Type)
		} else {
			sortedRequiredProviderTxt = fmt.Sprintf("%s {\n%s\n}", block.Type, sortedProviderParamTxt)
		}
		sortedRequiredProviderTxt = string(hclwrite.Format([]byte(sortedRequiredProviderTxt)))
		return runner.EmitIssue(
			r,
			fmt.Sprintf("The arguments of `required_providers` are expected to be sorted as follows:\n%s", sortedRequiredProviderTxt),
			block.DefRange(),
		)
	}
	var err error
	for _, issue := range providerParamIssues {
		if subErr := runner.EmitIssue(issue.Rule, issue.Message, issue.Range); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}
