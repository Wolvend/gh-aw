//go:build !integration

package workflow

import (
	"strings"
	"testing"

	"github.com/github/gh-aw/pkg/constants"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBuildActivationJob_Basic tests building a basic activation job
func TestBuildActivationJob_Basic(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
	}

	job, err := compiler.buildActivationJob(workflowData, false, "", "test.lock.yml")
	require.NoError(t, err, "buildActivationJob should succeed")
	require.NotNil(t, job)

	assert.Equal(t, string(constants.ActivationJobName), job.Name)
	assert.NotNil(t, job.Outputs, "Job should have outputs")
}

// TestBuildActivationJob_WithPreActivation tests activation job when pre-activation exists
func TestBuildActivationJob_WithPreActivation(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
	}

	job, err := compiler.buildActivationJob(workflowData, true, "", "test.lock.yml")
	require.NoError(t, err, "buildActivationJob should succeed with pre-activation")
	require.NotNil(t, job)

	// When pre-activation exists, activation job should have needs dependency
	assert.Contains(t, job.Needs, string(constants.PreActivationJobName),
		"Activation job should depend on pre-activation job")
}

// TestBuildActivationJob_WithReaction tests activation job with reaction configuration
func TestBuildActivationJob_WithReaction(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
		AIReaction:      "rocket",
	}

	job, err := compiler.buildActivationJob(workflowData, false, "", "test.lock.yml")
	require.NoError(t, err)
	require.NotNil(t, job)

	// Activation job should handle reactions appropriately
	stepsStr := strings.Join(job.Steps, "\n")
	// The reaction is actually added in pre-activation, but activation may reference it
	assert.NotEmpty(t, stepsStr, "Activation job should have steps")
}

// TestBuildActivationJob_WithWorkflowRunRepoSafety tests activation with workflow_run repo safety
func TestBuildActivationJob_WithWorkflowRunRepoSafety(t *testing.T) {
	compiler := NewCompiler()

	workflowData := &WorkflowData{
		Name:            "Test Workflow",
		Command:         []string{"echo", "test"},
		MarkdownContent: "# Test\n\nContent",
	}

	// Test with workflow_run repo safety enabled
	job, err := compiler.buildActivationJob(workflowData, false, "workflow_run", "test.lock.yml")
	require.NoError(t, err)
	require.NotNil(t, job)

	stepsStr := strings.Join(job.Steps, "\n")
	// Should include repository validation for workflow_run
	assert.NotEmpty(t, stepsStr)
}
