package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Rules is a list of all rules
var Rules = []tflint.Rule{
	NewBasicExtIgnoreConfigRule(),
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
}
