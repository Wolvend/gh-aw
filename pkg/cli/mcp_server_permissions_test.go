//go:build integration

package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/github/gh-aw/pkg/testutil"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServer_ConditionalToolMounting(t *testing.T) {
	// Skip if the binary doesn't exist
	binaryPath := "../../gh-aw"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("Skipping test: gh-aw binary not found. Run 'make build' first.")
	}

	// Create a temporary directory with a workflow file
	tmpDir := testutil.TempDir(t, "test-*")
	workflowsDir := filepath.Join(tmpDir, ".github", "workflows")
	require.NoError(t, os.MkdirAll(workflowsDir, 0755), "Failed to create workflows directory")

	// Create a test workflow file
	workflowContent := `---
on: push
engine: copilot
---
# Test Workflow

This is a test workflow.
`
	workflowPath := filepath.Join(workflowsDir, "test.md")
	require.NoError(t, os.WriteFile(workflowPath, []byte(workflowContent), 0644), "Failed to write workflow file")

	// Change to the temporary directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Initialize git repository
	require.NoError(t, initTestGitRepo(tmpDir), "Failed to initialize git repository")

	// Get absolute path to binary
	absBinaryPath, err := filepath.Abs(filepath.Join(originalDir, binaryPath))
	require.NoError(t, err, "Failed to get absolute path")

	// Create MCP client
	client := mcp.NewClient(&mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}, nil)

	// Start the MCP server
	serverCmd := exec.Command(absBinaryPath, "mcp-server", "--cmd", absBinaryPath)
	serverCmd.Dir = tmpDir
	transport := &mcp.CommandTransport{Command: serverCmd}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	session, err := client.Connect(ctx, transport, nil)
	require.NoError(t, err, "Failed to connect to MCP server")
	defer session.Close()

	// List available tools
	listToolsResult, err := session.ListTools(ctx, &mcp.ListToolsRequest{})
	require.NoError(t, err, "Failed to list tools")
	require.NotNil(t, listToolsResult, "ListTools result should not be nil")

	// Check which tools are available
	toolNames := make(map[string]bool)
	for _, tool := range listToolsResult.Tools {
		toolNames[tool.Name] = true
	}

	// Status tool should always be available
	assert.True(t, toolNames["status"], "Status tool should be available")

	// Compile tool should always be available
	assert.True(t, toolNames["compile"], "Compile tool should be available")

	// mcp-inspect tool should always be available
	assert.True(t, toolNames["mcp-inspect"], "MCP-inspect tool should be available")

	// add, update, fix tools should always be available
	assert.True(t, toolNames["add"], "Add tool should be available")
	assert.True(t, toolNames["update"], "Update tool should be available")
	assert.True(t, toolNames["fix"], "Fix tool should be available")

	// Check if logs and audit tools are available
	// These are conditionally mounted based on GitHub Actions permissions
	hasLogs := toolNames["logs"]
	hasAudit := toolNames["audit"]

	// Log the result for debugging
	t.Logf("Logs tool available: %v", hasLogs)
	t.Logf("Audit tool available: %v", hasAudit)

	// Both should have the same availability (both require actions:read permission)
	assert.Equal(t, hasLogs, hasAudit, "Logs and audit tools should have the same availability")

	// If we have a valid GitHub token with Actions permissions, both should be available
	// If not, both should be unavailable
	// This test documents the behavior without asserting a specific outcome
	// since it depends on the test environment's GitHub token permissions
}
