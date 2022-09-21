package version_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/google/go-github/v47/github"
	"github.com/stretchr/testify/assert"
	"github.com/terraform-linters/tflint-ruleset-basic-ext/project"
	"golang.org/x/mod/semver"
	"golang.org/x/oauth2"
)

func Test_EnsureVersionHasBeenBumpedUp(t *testing.T) {
	currentVersion := fmt.Sprintf("v%s", project.Version)
	assert.True(t, semver.IsValid(currentVersion))
	client := gitClient()
	tags, _, err := client.Repositories.ListTags(context.TODO(), "Azure", "tflint-ruleset-basic-ext", &github.ListOptions{
		Page:    0,
		PerPage: 10,
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, tags)
	for _, tag := range tags {
		v := tag.GetName()
		if semver.IsValid(v) && semver.Compare(v, currentVersion) >= 0 {
			t.Fatalf("latest version: %s, current version %s, please update current version in project/main.go", v, currentVersion)
		}
	}
}

func gitClient() *github.Client {
	var client *github.Client
	token := os.Getenv("TOKEN")
	if token != "" {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		)
		tc := oauth2.NewClient(context.TODO(), ts)
		client = github.NewClient(tc)
	} else {
		client = github.NewClient(nil)
	}
	return client
}
