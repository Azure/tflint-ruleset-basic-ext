package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

// DefaultRule is the default template for rules in this plugin
type DefaultRule struct {
	tflint.DefaultRule
	Rulename  string
	CheckFile func(runner tflint.Runner, file *hcl.File) error
}

// Check checks whether the tf config files match given rules
func (r *DefaultRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for filename, file := range files {
		if ignoreFile(filename, r.Rulename) {
			continue
		}
		if subErr := r.CheckFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

// Link returns the rule reference link
func (r *DefaultRule) Link() string { return project.ReferenceLink(r.Rulename) }

// Enabled returns whether the rule is enabled by default
func (r *DefaultRule) Enabled() bool { return false }

// Severity returns the rule severity
func (r *DefaultRule) Severity() tflint.Severity { return tflint.NOTICE }

// Name returns the rule name
func (r *DefaultRule) Name() string { return "" }

func (r *DefaultRule) setName(n string) {
	r.Rulename = n
}

func (r *DefaultRule) setCheckFunc(f func(runner tflint.Runner, file *hcl.File) error) {
	r.CheckFile = f
}

type myRule interface {
	tflint.Rule
	Name() string
	CheckFile(runner tflint.Runner, file *hcl.File) error
	setCheckFunc(f func(runner tflint.Runner, file *hcl.File) error)
	setName(string)
}

// NewRule returns a rule to be checked
func NewRule(r myRule) tflint.Rule {
	r.setName(r.Name())
	r.setCheckFunc(r.CheckFile)
	return r
}

func buildRules() {
	myRules := []myRule{
		NewTerraformVariableOrderRule(),
		NewTerraformVariableSeparateRule(),
		NewTerraformOutputSeparateRule(),
		NewTerraformOutputOrderRule(),
		NewTerraformLocalsOrderRule(),
		NewTerraformResourceDataArgLayoutRule(),
		NewTerraformCountIndexUsageRule(),
		NewTerraformHeredocUsageRule(),
		NewTerraformSensitiveVariableNoDefaultRule(),
		NewTerraformVersionsFileRule(),
		NewTerraformRequiredVersionDeclarationRule(),
		NewTerraformRequiredProvidersDeclarationRule(),
		NewTerraformModuleProviderDeclarationRule(),
		NewTerraformVarNameConventionRule(),
	}
	for _, rule := range myRules {
		Rules = append(Rules, NewRule(rule))
	}
}
