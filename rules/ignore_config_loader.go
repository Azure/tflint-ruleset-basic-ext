package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"log"
	"os"
	"regexp"
	"strings"
)

var ignores map[string][]*regexp.Regexp
var retains map[string][]*regexp.Regexp
var ignoreConfigFile = ".tflintignore.basic-ext.json"

func loadIgnoreConfig() {
	ignores = make(map[string][]*regexp.Regexp)
	retains = make(map[string][]*regexp.Regexp)
	if _, err := os.Stat(ignoreConfigFile); err != nil {
		if !os.IsNotExist(err) {
			log.Panicf("Ignore config file %s check failed:\n%s", ignoreConfigFile, err)
		}
		return
	}
	ruleIgnoreRegExps := make(map[string][]string)
	if err := hclsimple.DecodeFile(ignoreConfigFile, nil, &ruleIgnoreRegExps); err != nil {
		log.Panicf("Ignore config file %s read error:\n%s", ignoreConfigFile, err)
	}
	for ruleName, regExps := range ruleIgnoreRegExps {
		registerPatterns(ruleName, regExps)
	}
}

func registerPatterns(ruleName string, regExps []string) {
	existingRules := getExistedRules()
	if _, isRuleExist := existingRules[ruleName]; !isRuleExist {
		log.Panicf("Rule %s in %s not found", ruleName, ignoreConfigFile)
	}
	var err error
	for _, regExp := range regExps {
		isIgnorePattern := true
		if strings.HasPrefix(regExp, "!") {
			isIgnorePattern = false
			regExp = regExp[1:]
		}
		pattern, subErr := regexp.Compile(regExp)
		if subErr != nil {
			err = multierror.Append(err, subErr)
			continue
		}
		if isIgnorePattern {
			ignores[ruleName] = append(ignores[ruleName], pattern)
		} else {
			retains[ruleName] = append(retains[ruleName], pattern)
		}
	}
	if err != nil {
		log.Panicf("Error found when compile regexps in %s:\n%s", ignoreConfigFile, err)
	}
}

func ignoreFile(filename string, rulename string) bool {
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
