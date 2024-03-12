package project

import "fmt"

// Version is ruleset version
var Version string = "0.6.0"

// ReferenceLink returns the rule reference link
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/Azure/tflint-ruleset-basic-ext/blob/v%s/docs/rules/%s.md", Version, name)
}
