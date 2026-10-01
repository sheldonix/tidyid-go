package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/sheldonix/tidyid-go/v2"
)

func TestGeneratesDefaultCustomAndUppercaseIDs(t *testing.T) {
	tests := []struct {
		args      []string
		length    int
		uppercase bool
	}{
		{nil, 32, false},
		{[]string{"--size", "16"}, 16, false},
		{[]string{"-s", "3"}, 3, false},
		{[]string{"--allow-uppercase", "--size", "64"}, 64, true},
		{[]string{"-u", "-s", "3"}, 3, true},
	}
	for _, test := range tests {
		var output bytes.Buffer
		if err := Run(test.args, &output); err != nil {
			t.Fatalf("Run(%v): %v", test.args, err)
		}
		value := strings.TrimSpace(output.String())
		if !tidyid.IsValidIDOfLength(value, test.length, test.uppercase) {
			t.Fatalf("Run(%v) = %q", test.args, value)
		}
	}
}

func TestMetadataAndInvalidArguments(t *testing.T) {
	var output bytes.Buffer
	if err := Run([]string{"--version"}, &output); err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(output.String()) != tidyid.Version {
		t.Fatalf("version = %q", output.String())
	}

	output.Reset()
	if err := Run([]string{"--help"}, &output); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "--allow-uppercase") ||
		!strings.Contains(output.String(), "--size") {
		t.Fatalf("help = %q", output.String())
	}

	for _, args := range [][]string{
		{"--size"},
		{"--size", "2"},
		{"--size", "3.5"},
		{"--unknown"},
	} {
		output.Reset()
		if err := Run(args, &output); err == nil {
			t.Fatalf("Run(%v) succeeded", args)
		}
	}

	if err := Run([]string{"--size", "2"}, &output); err == nil ||
		err.Error() != "Size must be an integer between 3 and 256" {
		t.Fatalf("invalid size error = %v", err)
	}
	if err := Run([]string{"--unknown"}, &output); err == nil ||
		err.Error() != "Unknown argument --unknown" {
		t.Fatalf("unknown argument error = %v", err)
	}
}

func TestVersionAndHelpTakePrecedence(t *testing.T) {
	for _, args := range [][]string{
		{"--unknown", "--version"},
		{"--unknown", "--help"},
	} {
		var output bytes.Buffer
		if err := Run(args, &output); err != nil {
			t.Fatalf("Run(%v): %v", args, err)
		}
	}
}
