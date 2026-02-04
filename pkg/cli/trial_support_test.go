//go:build !integration

package cli

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCleanupTrialSecrets_RequiresEnvironmentVariable(t *testing.T) {
	tests := []struct {
		name           string
		tracker        *TrialSecretTracker
		envVarSet      bool
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:          "nil tracker skips cleanup",
			tracker:       nil,
			envVarSet:     false,
			expectedError: false,
		},
		{
			name: "without env var returns error",
			tracker: &TrialSecretTracker{
				RepoSlug:     "test/repo",
				AddedSecrets: map[string]bool{"TEST_SECRET": true},
			},
			envVarSet:      false,
			expectedError:  true,
			expectedErrMsg: "secret deletion is disabled",
		},
		{
			name: "with env var allows cleanup",
			tracker: &TrialSecretTracker{
				RepoSlug:     "test/repo",
				AddedSecrets: map[string]bool{},
			},
			envVarSet:     true,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			if tt.envVarSet {
				os.Setenv("GH_AW_ALLOW_SECRET_DELETION", "1")
				defer os.Unsetenv("GH_AW_ALLOW_SECRET_DELETION")
			} else {
				os.Unsetenv("GH_AW_ALLOW_SECRET_DELETION")
			}

			// Execute
			err := cleanupTrialSecrets("test/repo", tt.tracker, false)

			// Verify
			if tt.expectedError {
				require.Error(t, err, "Expected an error but got none")
				assert.Contains(t, err.Error(), tt.expectedErrMsg, "Error message should contain expected text")
			} else {
				assert.NoError(t, err, "Expected no error but got: %v", err)
			}
		})
	}
}

func TestCleanupTrialSecrets_EnvironmentVariableFormat(t *testing.T) {
	tests := []struct {
		name          string
		envValue      string
		expectedError bool
	}{
		{
			name:          "value '1' enables deletion",
			envValue:      "1",
			expectedError: false,
		},
		{
			name:          "empty value disables deletion",
			envValue:      "",
			expectedError: true,
		},
		{
			name:          "value 'true' disables deletion (not '1')",
			envValue:      "true",
			expectedError: true,
		},
		{
			name:          "value 'yes' disables deletion (not '1')",
			envValue:      "yes",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			tracker := &TrialSecretTracker{
				RepoSlug:     "test/repo",
				AddedSecrets: map[string]bool{},
			}

			if tt.envValue != "" {
				os.Setenv("GH_AW_ALLOW_SECRET_DELETION", tt.envValue)
				defer os.Unsetenv("GH_AW_ALLOW_SECRET_DELETION")
			} else {
				os.Unsetenv("GH_AW_ALLOW_SECRET_DELETION")
			}

			// Execute
			err := cleanupTrialSecrets("test/repo", tracker, false)

			// Verify
			if tt.expectedError {
				require.Error(t, err, "Expected an error but got none")
				assert.Contains(t, err.Error(), "secret deletion is disabled", "Error should mention secret deletion is disabled")
			} else {
				// Note: This may fail if gh CLI is not available or the repo doesn't exist,
				// but that's okay for this test - we're just checking the env var guard
				// The error (if any) should not be about secret deletion being disabled
				if err != nil {
					assert.NotContains(t, err.Error(), "secret deletion is disabled", "Error should not be about secret deletion being disabled")
				}
			}
		})
	}
}

func TestNewTrialSecretTracker(t *testing.T) {
	repoSlug := "test/repo"
	tracker := NewTrialSecretTracker(repoSlug)

	require.NotNil(t, tracker, "Tracker should not be nil")
	assert.Equal(t, repoSlug, tracker.RepoSlug, "RepoSlug should match")
	assert.NotNil(t, tracker.AddedSecrets, "AddedSecrets map should be initialized")
	assert.Empty(t, tracker.AddedSecrets, "AddedSecrets map should be empty initially")
}
