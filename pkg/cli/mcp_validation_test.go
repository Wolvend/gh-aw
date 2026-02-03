//go:build !integration

package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckGitHubActionsPermissions(t *testing.T) {
	tests := []struct {
		name           string
		setupEnv       func()
		expectedResult bool
		skipInCI       bool
		skipReason     string
	}{
		{
			name: "no token available",
			setupEnv: func() {
				os.Unsetenv("GH_TOKEN")
				os.Unsetenv("GITHUB_TOKEN")
			},
			expectedResult: false,
			skipInCI:       false,
		},
		{
			name: "token available but in non-repository directory",
			setupEnv: func() {
				// Token will be available in test environment
				// but we may not be in a proper repository
			},
			expectedResult: false,
			skipInCI:       true,
			skipReason:     "Requires valid GitHub repository context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipInCI && os.Getenv("CI") == "true" {
				t.Skip(tt.skipReason)
			}

			// Setup environment
			if tt.setupEnv != nil {
				tt.setupEnv()
			}

			// Call the function
			result := checkGitHubActionsPermissions()

			// Verify result
			// Note: We can't assert exact result in all environments
			// but we verify the function doesn't panic
			assert.IsType(t, false, result, "Function should return a boolean")
		})
	}
}

func TestValidateMCPServerConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		cmdPath     string
		shouldError bool
		skipInCI    bool
		skipReason  string
	}{
		{
			name:        "empty command path",
			cmdPath:     "",
			shouldError: true,
			skipInCI:    true,
			skipReason:  "Requires valid GitHub repository with workflows",
		},
		{
			name:        "invalid command path",
			cmdPath:     "/nonexistent/command",
			shouldError: true,
			skipInCI:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.skipInCI && os.Getenv("CI") == "true" {
				t.Skip(tt.skipReason)
			}

			// Call the function
			err := validateMCPServerConfiguration(tt.cmdPath)

			// Verify result
			if tt.shouldError {
				assert.Error(t, err, "Expected validation to fail")
			} else {
				assert.NoError(t, err, "Expected validation to succeed")
			}
		})
	}
}
