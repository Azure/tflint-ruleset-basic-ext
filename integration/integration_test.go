package integration

import (
	"bytes"
	"os"
	"os/exec"
	"testing"
)

func TestIntegration(t *testing.T) {
	cases := []struct {
		Name    string
		Command *exec.Cmd
		Dir     string
	}{
		{
			Name:    "basic",
			Command: exec.Command("tflint", "--format", "json", "--force"),
			Dir:     "basic",
		},
	}

	dir, _ := os.Getwd()
	defer os.Chdir(dir)

	for _, tc := range cases {
		testDir := dir + "/" + tc.Dir
		os.Chdir(testDir)

		var stdout, stderr bytes.Buffer
		tc.Command.Stdout = &stdout
		tc.Command.Stderr = &stderr
		if err := tc.Command.Run(); err != nil {
			t.Fatalf("Failed `%s`: %s, stdout=%s stderr=%s", tc.Name, err, stdout.String(), stderr.String())
		}
	}
}
