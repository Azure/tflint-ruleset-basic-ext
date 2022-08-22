package rules

import (
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
)

// TODO: 为何我们不直接使用 tflint.Rule？
type Rule interface {
	tflint.Rule
}

type DefaultRule struct {
	tflint.DefaultRule
	Rulename  string
	CheckFile func(runner tflint.Runner, file *hcl.File) error
}

func (r *DefaultRule) Check(runner tflint.Runner) error {
	files, err := runner.GetFiles()
	if err != nil {
		return err
	}
	for filename, file := range files {
		if ignoreFile(runner, filename, r.Rulename) {
			continue
		}
		if subErr := r.CheckFile(runner, file); subErr != nil {
			err = multierror.Append(err, subErr)
		}
	}
	return err
}

func (r *DefaultRule) Link() string              { return project.ReferenceLink(r.Rulename) }
func (r *DefaultRule) Enabled() bool             { return false }
func (r *DefaultRule) Severity() tflint.Severity { return tflint.NOTICE }
func (r *DefaultRule) Name() string              { return "" }
