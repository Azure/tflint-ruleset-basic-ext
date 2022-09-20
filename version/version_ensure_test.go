package version_test

import (
	"context"
	"fmt"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"golang.org/x/mod/semver"
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/stretchr/testify/assert"
)

func Test_EnsureVersionHasBeenBumpedUp(t *testing.T) {
	currentVersion := fmt.Sprintf("v%s", project.Version)
	assert.True(t, semver.IsValid(currentVersion))
	client := github.NewClient(nil)
	tags, _, err := client.Repositories.ListTags(context.TODO(), "Azure", "tflint-ruleset-basic-ext", &github.ListOptions{
		Page:    0,
		PerPage: 10,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, tags)
	for _, tag := range tags {
		v := tag.GetName()
		if semver.IsValid(v) {
			t.Fatalf("latest version: %s, current version %s, please update current version in project/main.go", v, currentVersion)
		}
	}
}
