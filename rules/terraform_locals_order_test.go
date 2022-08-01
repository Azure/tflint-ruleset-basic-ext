package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
	"testing"
)

func Test_TerraformLocalsOrderRule(t *testing.T) {
	cases := []struct {
		Name     string
		JSON     bool
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "common",
			Content: `
locals {
  instance_ids = concat(aws_instance.blue.*.id, aws_instance.green.*.id)
  common_tags = {
    Service = local.service_name
    Owner   = local.owner
  }
}

locals {
  service_name = "forum"
  owner        = "Community Team"
}`,
			Expected: helper.Issues{
				{
					Rule: NewTerraformLocalsOrderRule(),
					Message: `Recommended locals variable order:
locals {
  common_tags = {
    Service = local.service_name
    Owner   = local.owner
  }
  instance_ids = concat(aws_instance.blue.*.id, aws_instance.green.*.id)
}`,
				},
				{
					Rule: NewTerraformLocalsOrderRule(),
					Message: `Recommended locals variable order:
locals {
  owner        = "Community Team"
  service_name = "forum"
}`,
				},
			},
		},
	}
	rule := NewTerraformLocalsOrderRule()

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			filename := "config.tf"
			if tc.JSON {
				filename = "config.tf.json"
			}
			runner := helper.TestRunner(t, map[string]string{filename: tc.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			AssertIssues(t, tc.Expected, runner.Issues)
		})
	}
}
