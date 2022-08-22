package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// Rules is a list of all rules
var Rules = []tflint.Rule{
	// TODO: 一个例子，没有测试过
	//NewRule[*TerraformCountIndexUsageRule](),
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
