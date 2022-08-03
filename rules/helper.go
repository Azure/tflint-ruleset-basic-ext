package rules

import "github.com/hashicorp/hcl/v2"

func IsRangeEmpty(hclRange hcl.Range) bool {
	return hclRange == hcl.Range{}
}
