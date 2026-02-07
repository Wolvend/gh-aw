//go:build !integration

package workflow

import (
	"testing"

	"github.com/github/gh-aw/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildMainJob_Basic tests building a basic main job
func TestBuildMainJob_Basic(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
		AI:              "copilot",
	}

	job, err := compiler.buildMainJob(workflowData, false)
	require.NoError(t, err, "buildMainJob should succeed")
	require.NotNil(t, job)

	assert.Equal(t, string(constants.AgentJobName), job.Name)
	assert.NotEmpty(t, job.Steps, "Main job should have steps")
}

// TestBuildMainJob_WithActivation tests main job when activation job exists
func TestBuildMainJob_WithActivation(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
		AI:              "copilot",
	}

	job, err := compiler.buildMainJob(workflowData, true)
	require.NoError(t, err, "buildMainJob should succeed with activation")
	require.NotNil(t, job)

	// When activation exists, main job should depend on it
	assert.Contains(t, job.Needs, string(constants.ActivationJobName),
		"Main job should depend on activation job")
}

// TestBuildMainJob_WithPermissions tests main job permission handling
func TestBuildMainJob_WithPermissions(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
		AI:              "copilot",
		Permissions:     "contents: read\nissues: write",
	}

	job, err := compiler.buildMainJob(workflowData, false)
	require.NoError(t, err)
	require.NotNil(t, job)

	// Check permissions are set
	assert.NotEmpty(t, job.Permissions, "Main job should have permissions")
	assert.Contains(t, job.Permissions, "contents:",
		"Permissions should include contents")
}

// TestBuildMainJob_EngineSpecific tests main job with different engines
func TestBuildMainJob_EngineSpecific(t *testing.T) {
	tests := []struct {
		name   string
		engine string
	}{
		{
			name:   "copilot engine",
			engine: "copilot",
		},
		{
			name:   "claude engine",
			engine: "claude",
		},
		{
			name:   "codex engine",
			engine: "codex",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler := NewCompiler()

			workflowData := &WorkflowData{
				Name:            "Test Workflow",
				Command:         []string{"echo", "test"},
				MarkdownContent: "# Test\n\nContent",
				AI:              tt.engine,
			}

			job, err := compiler.buildMainJob(workflowData, false)
			require.NoError(t, err, "buildMainJob should succeed for engine %s", tt.engine)
			require.NotNil(t, job)
			assert.NotEmpty(t, job.Steps, "Should have steps for engine %s", tt.engine)
		})
	}
}
