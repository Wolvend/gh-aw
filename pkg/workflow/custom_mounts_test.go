//go:build !integration

package workflow

import (
	"strings"
	"testing"
)

// TestClaudeEngineWithCustomMounts tests that custom mounts are included in AWF command for Claude engine
func TestClaudeEngineWithCustomMounts(t *testing.T) {
	t.Run("custom mounts are included in AWF command", func(t *testing.T) {
		workflowData := &WorkflowData{
			Name: "test-claude-workflow",
			EngineConfig: &EngineConfig{
				ID: "claude",
			},
			SandboxConfig: &SandboxConfig{
				Agent: &AgentSandboxConfig{
					ID: "awf",
					Mounts: []string{
						"/usr/bin/psql:/usr/bin/psql:ro",
						"/usr/lib/postgresql:/usr/lib/postgresql:ro",
					},
				},
			},
			NetworkPermissions: &NetworkPermissions{
				Firewall: &FirewallConfig{
					Enabled: true,
				},
			},
		}

		engine := NewClaudeEngine()
		steps := engine.GetExecutionSteps(workflowData, "test.log")

		if len(steps) == 0 {
			t.Fatal("Expected at least one execution step")
		}

		stepContent := strings.Join(steps[0], "\n")

		// Check that custom mounts are included
		if !strings.Contains(stepContent, "--mount /usr/bin/psql:/usr/bin/psql:ro") {
			t.Error("Expected command to contain custom mount '--mount /usr/bin/psql:/usr/bin/psql:ro'")
		}

		if !strings.Contains(stepContent, "--mount /usr/lib/postgresql:/usr/lib/postgresql:ro") {
			t.Error("Expected command to contain custom mount '--mount /usr/lib/postgresql:/usr/lib/postgresql:ro'")
		}

		// Verify standard mounts are still present
		if !strings.Contains(stepContent, "--mount /tmp:/tmp:rw") {
			t.Error("Expected command to still contain standard mount '--mount /tmp:/tmp:rw'")
		}
	})

	t.Run("custom mounts are sorted alphabetically", func(t *testing.T) {
		workflowData := &WorkflowData{
			Name: "test-claude-workflow",
			EngineConfig: &EngineConfig{
				ID: "claude",
			},
			SandboxConfig: &SandboxConfig{
				Agent: &AgentSandboxConfig{
					ID: "awf",
					Mounts: []string{
						"/var/cache:/cache:rw",
						"/etc/ssl:/etc/ssl:ro",
						"/usr/local/bin/aws:/usr/local/bin/aws:ro",
					},
				},
			},
			NetworkPermissions: &NetworkPermissions{
				Firewall: &FirewallConfig{
					Enabled: true,
				},
			},
		}

		engine := NewClaudeEngine()
		steps := engine.GetExecutionSteps(workflowData, "test.log")

		if len(steps) == 0 {
			t.Fatal("Expected at least one execution step")
		}

		stepContent := strings.Join(steps[0], "\n")

		// Find positions of mounts to verify sorting
		etcPos := strings.Index(stepContent, "--mount /etc/ssl:/etc/ssl:ro")
		usrPos := strings.Index(stepContent, "--mount /usr/local/bin/aws:/usr/local/bin/aws:ro")
		varPos := strings.Index(stepContent, "--mount /var/cache:/cache:rw")

		// Verify all mounts are present
		if etcPos == -1 {
			t.Error("Expected to find mount '/etc/ssl:/etc/ssl:ro'")
		}
		if usrPos == -1 {
			t.Error("Expected to find mount '/usr/local/bin/aws:/usr/local/bin/aws:ro'")
		}
		if varPos == -1 {
			t.Error("Expected to find mount '/var/cache:/cache:rw'")
		}

		// Verify mounts are in alphabetical order: /etc, /usr, /var
		if etcPos != -1 && usrPos != -1 && etcPos >= usrPos {
			t.Error("Expected '/etc/ssl:/etc/ssl:ro' to appear before '/usr/local/bin/aws:/usr/local/bin/aws:ro'")
		}
		if usrPos != -1 && varPos != -1 && usrPos >= varPos {
			t.Error("Expected '/usr/local/bin/aws:/usr/local/bin/aws:ro' to appear before '/var/cache:/cache:rw'")
		}
	})
}

// TestCodexEngineWithCustomMounts tests that custom mounts are included in AWF command for Codex engine
func TestCodexEngineWithCustomMounts(t *testing.T) {
	t.Run("custom mounts are included in AWF command", func(t *testing.T) {
		workflowData := &WorkflowData{
			Name: "test-codex-workflow",
			EngineConfig: &EngineConfig{
				ID: "codex",
			},
			SandboxConfig: &SandboxConfig{
				Agent: &AgentSandboxConfig{
					ID: "awf",
					Mounts: []string{
						"/usr/bin/docker:/usr/bin/docker:ro",
						"/usr/bin/kubectl:/usr/bin/kubectl:ro",
					},
				},
			},
			NetworkPermissions: &NetworkPermissions{
				Firewall: &FirewallConfig{
					Enabled: true,
				},
			},
		}

		engine := NewCodexEngine()
		steps := engine.GetExecutionSteps(workflowData, "test.log")

		if len(steps) == 0 {
			t.Fatal("Expected at least one execution step")
		}

		stepContent := strings.Join(steps[0], "\n")

		// Check that custom mounts are included
		if !strings.Contains(stepContent, "--mount /usr/bin/docker:/usr/bin/docker:ro") {
			t.Error("Expected command to contain custom mount '--mount /usr/bin/docker:/usr/bin/docker:ro'")
		}

		if !strings.Contains(stepContent, "--mount /usr/bin/kubectl:/usr/bin/kubectl:ro") {
			t.Error("Expected command to contain custom mount '--mount /usr/bin/kubectl:/usr/bin/kubectl:ro'")
		}

		// Verify standard mounts are still present
		if !strings.Contains(stepContent, "--mount /tmp:/tmp:rw") {
			t.Error("Expected command to still contain standard mount '--mount /tmp:/tmp:rw'")
		}
	})

	t.Run("no custom mounts when not specified", func(t *testing.T) {
		workflowData := &WorkflowData{
			Name: "test-codex-workflow",
			EngineConfig: &EngineConfig{
				ID: "codex",
			},
			NetworkPermissions: &NetworkPermissions{
				Firewall: &FirewallConfig{
					Enabled: true,
				},
			},
		}

		engine := NewCodexEngine()
		steps := engine.GetExecutionSteps(workflowData, "test.log")

		if len(steps) == 0 {
			t.Fatal("Expected at least one execution step")
		}

		stepContent := strings.Join(steps[0], "\n")

		// Verify standard mounts are present
		if !strings.Contains(stepContent, "--mount /tmp:/tmp:rw") {
			t.Error("Expected command to contain standard mount '--mount /tmp:/tmp:rw'")
		}

		// Custom mount should not be present
		if strings.Contains(stepContent, "--mount /usr/bin/docker:/usr/bin/docker:ro") {
			t.Error("Did not expect custom mount in output when not configured")
		}
	})
}

// TestCustomMountsAcrossEngines tests that custom mounts work consistently across all engines
func TestCustomMountsAcrossEngines(t *testing.T) {
	testCases := []struct {
		name       string
		engineName string
		engine     CodingAgentEngine
	}{
		{
			name:       "copilot",
			engineName: "copilot",
			engine:     NewCopilotEngine(),
		},
		{
			name:       "claude",
			engineName: "claude",
			engine:     NewClaudeEngine(),
		},
		{
			name:       "codex",
			engineName: "codex",
			engine:     NewCodexEngine(),
		},
	}

	customMounts := []string{
		"/usr/bin/custom-tool:/usr/bin/custom-tool:ro",
		"/data/shared:/data:ro",
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			workflowData := &WorkflowData{
				Name: "test-workflow",
				EngineConfig: &EngineConfig{
					ID: tc.engineName,
				},
				SandboxConfig: &SandboxConfig{
					Agent: &AgentSandboxConfig{
						ID:     "awf",
						Mounts: customMounts,
					},
				},
				NetworkPermissions: &NetworkPermissions{
					Firewall: &FirewallConfig{
						Enabled: true,
					},
				},
			}

			steps := tc.engine.GetExecutionSteps(workflowData, "test.log")

			if len(steps) == 0 {
				t.Fatalf("Expected at least one execution step for %s engine", tc.engineName)
			}

			stepContent := strings.Join(steps[0], "\n")

			// Verify both custom mounts are present
			for _, mount := range customMounts {
				expectedArg := "--mount " + mount
				if !strings.Contains(stepContent, expectedArg) {
					t.Errorf("%s engine: Expected command to contain '%s'", tc.engineName, expectedArg)
				}
			}

			// Verify standard mounts are still present
			if !strings.Contains(stepContent, "--mount /tmp:/tmp:rw") {
				t.Errorf("%s engine: Expected command to contain standard mount '--mount /tmp:/tmp:rw'", tc.engineName)
			}
		})
	}
}

// TestMountsSyntaxForCommonTools tests validation for common tool mount scenarios
func TestMountsSyntaxForCommonTools(t *testing.T) {
	tests := []struct {
		name    string
		mounts  []string
		wantErr bool
	}{
		{
			name: "database client binaries",
			mounts: []string{
				"/usr/bin/psql:/usr/bin/psql:ro",
				"/usr/bin/mysql:/usr/bin/mysql:ro",
				"/usr/bin/redis-cli:/usr/bin/redis-cli:ro",
			},
			wantErr: false,
		},
		{
			name: "cloud CLI tools",
			mounts: []string{
				"/usr/local/bin/aws:/usr/local/bin/aws:ro",
				"/usr/local/bin/gcloud:/usr/local/bin/gcloud:ro",
				"/usr/local/bin/az:/usr/local/bin/az:ro",
			},
			wantErr: false,
		},
		{
			name: "build tools",
			mounts: []string{
				"/usr/bin/make:/usr/bin/make:ro",
				"/usr/bin/cmake:/usr/bin/cmake:ro",
				"/usr/bin/gcc:/usr/bin/gcc:ro",
			},
			wantErr: false,
		},
		{
			name: "container tools",
			mounts: []string{
				"/usr/bin/docker:/usr/bin/docker:ro",
				"/usr/bin/kubectl:/usr/bin/kubectl:ro",
				"/usr/bin/helm:/usr/bin/helm:ro",
			},
			wantErr: false,
		},
		{
			name: "shared libraries",
			mounts: []string{
				"/usr/lib/x86_64-linux-gnu/libssl.so.3:/usr/lib/x86_64-linux-gnu/libssl.so.3:ro",
				"/usr/lib/x86_64-linux-gnu/libcrypto.so.3:/usr/lib/x86_64-linux-gnu/libcrypto.so.3:ro",
			},
			wantErr: false,
		},
		{
			name: "directories",
			mounts: []string{
				"/usr/share/ca-certificates:/usr/share/ca-certificates:ro",
				"/etc/ssl/certs:/etc/ssl/certs:ro",
			},
			wantErr: false,
		},
		{
			name: "writable cache directories",
			mounts: []string{
				"/var/cache/apt:/var/cache/apt:rw",
				"/tmp/build-cache:/cache:rw",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMountsSyntax(tt.mounts)

			if tt.wantErr && err == nil {
				t.Error("validateMountsSyntax() expected error but got none")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("validateMountsSyntax() unexpected error: %v", err)
			}
		})
	}
}
