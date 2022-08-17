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

// Arg contains attrs and nested blocks defined in a block
type Arg struct {
	Name  string
	Range hcl.Range
	Attr  *hclsyntax.Attribute
	Block *hclsyntax.Block
}

var headMetaArgPriority, tailMetaArgPriority = map[string]int{"for_each": 1, "count": 1, "provider": 0}, map[string]int{"lifecycle": 1, "depends_on": 0}

var ignoreConfigLoad bool

// IsHeadMeta checks whether a name represents a type of head Meta arg
func IsHeadMeta(argName string) bool {
	_, isHeadMeta := headMetaArgPriority[argName]
	return isHeadMeta
}

// IsTailMeta checks whether a name represents a type of tail Meta arg
func IsTailMeta(argName string) bool {
	_, isTailMeta := tailMetaArgPriority[argName]
	return isTailMeta
}

// GetHeadMetaPriority gets the priority of a head Meta arg
func GetHeadMetaPriority(argName string) int {
	return headMetaArgPriority[argName]
}

// GetTailMetaPriority gets the priority of a tail Meta arg
func GetTailMetaPriority(argName string) int {
	return tailMetaArgPriority[argName]
}

// GetArgsWithOriginalOrder returns the args with original order
func GetArgsWithOriginalOrder(args []Arg) []Arg {
	argsWithOriginalOrder := make([]Arg, len(args))
	copy(argsWithOriginalOrder, args)
	sort.Slice(argsWithOriginalOrder, func(i, j int) bool {
		if argsWithOriginalOrder[i].Range.Start.Line == argsWithOriginalOrder[j].Range.Start.Line {
			return argsWithOriginalOrder[i].Range.Start.Column < argsWithOriginalOrder[j].Range.Start.Column
		}
		return argsWithOriginalOrder[i].Range.Start.Line < argsWithOriginalOrder[j].Range.Start.Line
	})
	return argsWithOriginalOrder
}

// IsRangeEmpty checks whether a range is empty
func IsRangeEmpty(hclRange hcl.Range) bool {
	return hclRange == hcl.Range{}
}

// PrintSortedAttrTxt print the sorted hcl text of an attribute
func PrintSortedAttrTxt(src []byte, attr *hclsyntax.Attribute) (string, bool) {
	isSorted := true
	exp, isMap := attr.Expr.(*hclsyntax.ObjectConsExpr)
	if !isMap {
		return string(attr.Range().SliceBytes(src)), isSorted
	}
	var keys []string
	itemTxts := make(map[string]string)
	for _, item := range exp.Items {
		key := string(item.KeyExpr.Range().SliceBytes(src))
		itemTxt := fmt.Sprintf("%s = %s", key, string(item.ValueExpr.Range().SliceBytes(src)))
		keys = append(keys, key)
		itemTxts[key] = itemTxt
	}
	isSorted = sort.StringsAreSorted(keys)
	if !isSorted {
		sort.Strings(keys)
	}
	var sortedItemTxts []string
	for _, key := range keys {
		sortedItemTxts = append(sortedItemTxts, itemTxts[key])
	}
	sortedExpTxt := strings.Join(sortedItemTxts, "\n")
	var sortedAttrTxt string
	if RemoveSpaceAndLine(sortedExpTxt) == "" {
		sortedAttrTxt = fmt.Sprintf("%s = {}", attr.Name)
	} else {
		sortedAttrTxt = fmt.Sprintf("%s = {\n%s\n}", attr.Name, sortedExpTxt)
	}
	sortedAttrTxt = string(hclwrite.Format([]byte(sortedAttrTxt)))
	return sortedAttrTxt, isSorted
}

// RemoveSpaceAndLine remove space, "\t" and "\n" from the given string
func RemoveSpaceAndLine(str string) string {
	newStr := strings.ReplaceAll(str, " ", "")
	newStr = strings.ReplaceAll(newStr, "\t", "")
	newStr = strings.ReplaceAll(newStr, "\n", "")
	return newStr
}

func ignoreFile(filename string, rule tflint.Rule) bool {
	isIgnore := false
	ruleName := rule.Name()
	ignorePatterns, isRuleIgnorePatternDefined := ignores[ruleName]
	if isRuleIgnorePatternDefined {
		for _, ignorePattern := range ignorePatterns {
			if ignorePattern.MatchString(filename) {
				isIgnore = true
				break
			}
		}
	}
	retainPatterns, isRuleRetainPatternDefined := retains[ruleName]
	if isRuleRetainPatternDefined {
		for _, retainPattern := range retainPatterns {
			if retainPattern.MatchString(filename) {
				isIgnore = false
				break
			}
		}
	}
	return isIgnore
}

func getExistedRules() map[string]tflint.Rule {
	rules := make(map[string]tflint.Rule)
	for _, rule := range Rules {
		if rule.Name() == "basic_ext_ignore_config" {
			continue
		}
		rules[rule.Name()] = rule
	}
	return rules
}
