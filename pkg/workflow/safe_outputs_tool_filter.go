package workflow

import (
	"encoding/json"
	"fmt"

	"github.com/github/gh-aw/pkg/logger"
)

var safeOutputsToolFilterLog = logger.New("workflow:safe_outputs_tool_filter")

// ========================================
// Safe Outputs - Tool Filtering and Enhancement
// ========================================
//
// This file handles filtering the complete tools list to only include enabled
// tools, enhancing their descriptions with configuration details, and adding
// custom job tools and dispatch workflow tools.

// generateFilteredToolsJSON filters the ALL_TOOLS array based on enabled safe outputs
// Returns a JSON string containing only the tools that are enabled in the workflow
func generateFilteredToolsJSON(data *WorkflowData, markdownPath string) (string, error) {
	if data.SafeOutputs == nil {
		return "[]", nil
	}

	safeOutputsToolFilterLog.Print("Generating filtered tools JSON for workflow")

	// Load the full tools JSON
	allToolsJSON := GetSafeOutputsToolsJSON()

	// Parse the JSON to get all tools
	var allTools []map[string]any
	if err := json.Unmarshal([]byte(allToolsJSON), &allTools); err != nil {
		safeOutputsToolFilterLog.Printf("Failed to parse safe outputs tools JSON: %v", err)
		return "", fmt.Errorf("failed to parse safe outputs tools JSON: %w", err)
	}

	// Create a set of enabled tool names using the registry pattern
	enabledTools := buildEnabledToolsSet(data.SafeOutputs)

	// Filter tools to only include enabled ones and enhance descriptions
	var filteredTools []map[string]any
	for _, tool := range allTools {
		toolName, ok := tool["name"].(string)
		if !ok {
			continue
		}
		if enabledTools[toolName] {
			// Create a copy of the tool to avoid modifying the original
			enhancedTool := make(map[string]any)
			for k, v := range tool {
				enhancedTool[k] = v
			}

			// Enhance the description with configuration details
			if description, ok := enhancedTool["description"].(string); ok {
				enhancedDescription := enhanceToolDescription(toolName, description, data.SafeOutputs)
				enhancedTool["description"] = enhancedDescription
			}

			// Add repo parameter to inputSchema if allowed-repos has entries
			addRepoParameterIfNeeded(enhancedTool, toolName, data.SafeOutputs)

			filteredTools = append(filteredTools, enhancedTool)
		}
	}

	// Add custom job tools from SafeOutputs.Jobs
	addCustomJobTools(&filteredTools, data.SafeOutputs)

	if safeOutputsToolFilterLog.Enabled() {
		safeOutputsToolFilterLog.Printf("Filtered %d tools from %d total tools (including %d custom jobs)", len(filteredTools), len(allTools), len(data.SafeOutputs.Jobs))
	}

	// Add dynamic dispatch_workflow tools
	addDispatchWorkflowTools(&filteredTools, data, markdownPath)

	// Marshal the filtered tools back to JSON with indentation for better readability
	// and to reduce merge conflicts in generated lockfiles
	filteredJSON, err := json.MarshalIndent(filteredTools, "", "  ")
	if err != nil {
		safeOutputsToolFilterLog.Printf("Failed to marshal filtered tools: %v", err)
		return "", fmt.Errorf("failed to marshal filtered tools: %w", err)
	}

	safeOutputsToolFilterLog.Printf("Successfully generated filtered tools JSON with %d tools", len(filteredTools))
	return string(filteredJSON), nil
}
