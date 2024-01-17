package main

import (
	"github.com/Azure/tflint-ruleset-basic-ext/project"
	"github.com/Azure/tflint-ruleset-basic-ext/rules"
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "basic-ext",
			Version: project.Version,
			Rules:   rules.Rules,
		},
	})
}
