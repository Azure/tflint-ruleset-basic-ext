package rules

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"sort"
)

// Arg contains attrs and nested blocks defined in a block
type Arg struct {
	Name  string
	Range hcl.Range
	Block *hclsyntax.Block
}

var headMetaArgPriority, tailMetaArgPriority = map[string]int{"for_each": 1, "count": 1, "provider": 0}, map[string]int{"lifecycle": 1, "depends_on": 0}

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
