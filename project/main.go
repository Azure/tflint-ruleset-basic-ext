package project

import "fmt"

// Version is ruleset version
const Version string = "0.1.1"

// ReferenceLink returns the rule reference link
func ReferenceLink(name string) string {
	return fmt.Sprintf("https://github.com/Azure/tflint-ruleset-basic-ext/blob/%s/docs/rules/%s.md", Version, name)
}
