package rules

import (
	"fmt"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"github.com/zclconf/go-cty/cty"
	"os"
	"regexp"
)

var ignores, retains = make(map[string][]*regexp.Regexp), make(map[string][]*regexp.Regexp)
var ignoreConfigLoaded = false

// BasicExtIgnoreConfigRule checks whether count.index is used as subscript of list/map
type BasicExtIgnoreConfigRule struct {
	tflint.DefaultRule
}

// NewBasicExtIgnoreConfigRule returns a new rule
func NewBasicExtIgnoreConfigRule() *BasicExtIgnoreConfigRule {
	return &BasicExtIgnoreConfigRule{}
}

// Name returns the rule name
func (r *BasicExtIgnoreConfigRule) Name() string {
	return "basic_ext_ignore_config"
}

// Enabled returns whether the rule is enabled by default
func (r *BasicExtIgnoreConfigRule) Enabled() bool {
	return true
}

// Severity returns the rule severity
func (r *BasicExtIgnoreConfigRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link
func (r *BasicExtIgnoreConfigRule) Link() string {
	return project.ReferenceLink(r.Name())
}

func (r *BasicExtIgnoreConfigRule) Check(runner tflint.Runner) error {
	ignoreConfigFile := ".tflintignore.basic-ext"
	if _, err := os.Stat(ignoreConfigFile); err != nil {
		if !os.IsNotExist(err) {
			return runner.EmitIssue(r, fmt.Sprintf("Load basic-ext ignore config file failed:\n%s", err), hcl.Range{})
		}
		return nil
	} else {
		bytes, err := os.ReadFile(ignoreConfigFile)
		if err != nil {
			return runner.EmitIssue(r, fmt.Sprintf("Load basic-ext ignore config file failed:\n%s", err), hcl.Range{})
		}
		file, diags := hclsyntax.ParseConfig(bytes, ignoreConfigFile, hcl.InitialPos)
		if diags.HasErrors() {
			return runner.EmitIssue(r, fmt.Sprintf("Parse basic-ext ignore config file failed:\n%s", diags), hcl.Range{})
		}
		blocks := file.Body.(*hclsyntax.Body).Blocks
		for _, block := range blocks {
			if subErr := r.handleBlock(runner, block); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		}
		return err
	}
}

func (r *BasicExtIgnoreConfigRule) handleBlock(runner tflint.Runner, block *hclsyntax.Block) error {
	var patterns map[string][]*regexp.Regexp
	switch block.Type {
	case "ignore":
		patterns = ignores
	case "retain":
		patterns = retains
	default:
		return runner.EmitIssue(
			r,
			fmt.Sprintf("block type `%s` not supported in basic-ext ignore config!", block.Type),
			block.DefRange(),
		)
	}
	args := block.Body.Attributes
	includeRulesArg, includeRulesDeclared := args["include_rules"]
	excludeRulesArg, excludeRulesDeclared := args["exclude_rules"]
	patternsArg, patternsDeclared := args["patterns"]
	if !patternsDeclared {
		return runner.EmitIssue(
			r,
			fmt.Sprintf("`patterns` field must be declared in block `%s`", block.Type),
			block.DefRange(),
		)
	}
	patternStrs, patternRanges, err := r.extractListExp(runner, patternsArg.Expr)
	if err != nil {
		return err
	}
	var patternRegExps []*regexp.Regexp
	for i, patternStr := range patternStrs {
		if patternRegExp, subErr := regexp.Compile(patternStr); subErr != nil {
			if subErr := runner.EmitIssue(r, fmt.Sprintf("Pattern `%s` is not valid:\n5%s", patternStr, subErr), patternRanges[i]); subErr != nil {
				err = multierror.Append(err, subErr)
			}
		} else {
			patternRegExps = append(patternRegExps, patternRegExp)
		}
	}
	if !includeRulesDeclared && !excludeRulesDeclared {
		for _, rule := range Rules {
			patterns[rule.Name()] = append(patterns[rule.Name()], patternRegExps...)
		}
	} else if includeRulesDeclared && excludeRulesDeclared {
		return runner.EmitIssue(r, "`include_rules` and `exclude_rules` can't be declared within the same block", block.DefRange())
	} else {
		var subErr error
		var rules []string
		var ruleRanges []hcl.Range
		if includeRulesDeclared {
			rules, ruleRanges, subErr = r.extractListExp(runner, includeRulesArg.Expr)
		} else {
			rules, ruleRanges, subErr = r.extractListExp(runner, excludeRulesArg.Expr)
		}
		if subErr != nil {
			err = multierror.Append(err, subErr)
		} else {
			existedRules := getExistedRules()
			for i, ruleName := range rules {
				if _, isRuleExisted := existedRules[ruleName]; !isRuleExisted {
					if subErr := runner.EmitIssue(r, fmt.Sprintf("Rule `%s` doesn't exist", ruleName), ruleRanges[i]); subErr != nil {
						err = multierror.Append(err, subErr)
					}
				} else {
					patterns[ruleName] = append(patterns[ruleName], patternRegExps...)
				}
			}
		}
	}
	return err
}

func (r *BasicExtIgnoreConfigRule) extractListExp(runner tflint.Runner, exp hclsyntax.Expression) ([]string, []hcl.Range, error) {
	msg := "The expression is expected to be a string list"
	listExp, isList := exp.(*hclsyntax.TupleConsExpr)
	if !isList {
		return nil, nil, runner.EmitIssue(r, msg, exp.Range())
	}
	var err error
	var strs []string
	var ranges []hcl.Range
	for _, subExp := range listExp.Exprs {
		subExp, isTemplateExpr := subExp.(*hclsyntax.TemplateExpr)
		if !isTemplateExpr {
			return nil, nil, runner.EmitIssue(r, msg, exp.Range())
		}
		val, diags := subExp.Parts[0].Value(nil)
		if diags.HasErrors() {
			err = multierror.Append(err, diags)
		} else {
			if val.Type() != cty.String {
				return nil, nil, runner.EmitIssue(r, msg, exp.Range())
			}
			strs = append(strs, val.AsString())
			ranges = append(ranges, subExp.Parts[0].Range())
		}
	}
	return strs, ranges, err
}

func loadIgnoreConfig(runner tflint.Runner) {
	if ignoreConfigLoaded {
		return
	}
	r := NewBasicExtIgnoreConfigRule()
	r.Check(runner)
	ignoreConfigLoaded = true
}

func ignoreFile(runner tflint.Runner, filename string, rulename string) bool {
	loadIgnoreConfig(runner)
	isIgnore := false
	ignorePatterns, isRuleIgnorePatternDefined := ignores[rulename]
	if isRuleIgnorePatternDefined {
		for _, ignorePattern := range ignorePatterns {
			if ignorePattern.MatchString(filename) {
				isIgnore = true
				break
			}
		}
	}
	retainPatterns, isRuleRetainPatternDefined := retains[rulename]
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
