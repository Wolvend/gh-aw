package workflow

import (
	"encoding/json"
	"sort"

	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/stringutil"
)

var safeOutputsConfigBuilderLog = logger.New("workflow:safe_outputs_config_builder")

// ========================================
// Safe Outputs - Configuration Builder
// ========================================
//
// This file contains the main orchestration logic for generating safe outputs
// configuration. It uses the registry pattern to eliminate repetitive
// conditional checks and delegates to specialized handlers.

// generateSafeOutputsConfig generates the safe outputs configuration JSON string
// It uses the tool registry pattern to build configuration for enabled tools
func generateSafeOutputsConfig(data *WorkflowData) string {
	// Pass the safe-outputs configuration for validation
	if data.SafeOutputs == nil {
		return ""
	}
	safeOutputsConfigBuilderLog.Print("Generating safe outputs configuration for workflow")

	// Create a simplified config object for validation
	safeOutputsConfig := make(map[string]any)

	// Use registry pattern to generate config for each tool
	if data.SafeOutputs != nil {
		for toolName, handler := range toolRegistry {
			config, shouldInclude := handler(data.SafeOutputs)
			if shouldInclude {
				safeOutputsConfigBuilderLog.Printf("Adding config for tool: %s", toolName)
				safeOutputsConfig[toolName] = config
			}
		}
	}

	// Add safe-jobs configuration from SafeOutputs.Jobs
	if len(data.SafeOutputs.Jobs) > 0 {
		for jobName, jobConfig := range data.SafeOutputs.Jobs {
			safeJobConfig := map[string]any{}

			// Add description if present
			if jobConfig.Description != "" {
				safeJobConfig["description"] = jobConfig.Description
			}

			// Add output if present
			if jobConfig.Output != "" {
				safeJobConfig["output"] = jobConfig.Output
			}

			// Add inputs information
			if len(jobConfig.Inputs) > 0 {
				inputsConfig := make(map[string]any)
				for inputName, inputDef := range jobConfig.Inputs {
					inputConfig := map[string]any{
						"type":        inputDef.Type,
						"description": inputDef.Description,
						"required":    inputDef.Required,
					}
					if inputDef.Default != "" {
						inputConfig["default"] = inputDef.Default
					}
					if len(inputDef.Options) > 0 {
						inputConfig["options"] = inputDef.Options
					}
					inputsConfig[inputName] = inputConfig
				}
				safeJobConfig["inputs"] = inputsConfig
			}

			safeOutputsConfig[jobName] = safeJobConfig
		}
	}

	// Add mentions configuration
	if data.SafeOutputs.Mentions != nil {
		mentionsConfig := make(map[string]any)

		// Handle enabled flag (simple boolean mode)
		if data.SafeOutputs.Mentions.Enabled != nil {
			mentionsConfig["enabled"] = *data.SafeOutputs.Mentions.Enabled
		}

		// Handle allow-team-members
		if data.SafeOutputs.Mentions.AllowTeamMembers != nil {
			mentionsConfig["allowTeamMembers"] = *data.SafeOutputs.Mentions.AllowTeamMembers
		}

		// Handle allow-context
		if data.SafeOutputs.Mentions.AllowContext != nil {
			mentionsConfig["allowContext"] = *data.SafeOutputs.Mentions.AllowContext
		}

		// Handle allowed list
		if len(data.SafeOutputs.Mentions.Allowed) > 0 {
			mentionsConfig["allowed"] = data.SafeOutputs.Mentions.Allowed
		}

		// Handle max
		if data.SafeOutputs.Mentions.Max != nil {
			mentionsConfig["max"] = *data.SafeOutputs.Mentions.Max
		}

		// Only add mentions config if it has any fields
		if len(mentionsConfig) > 0 {
			safeOutputsConfig["mentions"] = mentionsConfig
		}
	}

	// Add dispatch-workflow configuration
	if data.SafeOutputs.DispatchWorkflow != nil {
		dispatchWorkflowConfig := map[string]any{}

		// Include workflows list
		if len(data.SafeOutputs.DispatchWorkflow.Workflows) > 0 {
			dispatchWorkflowConfig["workflows"] = data.SafeOutputs.DispatchWorkflow.Workflows
		}

		// Include workflow files mapping (file extension for each workflow)
		if len(data.SafeOutputs.DispatchWorkflow.WorkflowFiles) > 0 {
			dispatchWorkflowConfig["workflow_files"] = data.SafeOutputs.DispatchWorkflow.WorkflowFiles
		}

		// Include max count
		maxValue := 1 // default
		if data.SafeOutputs.DispatchWorkflow.Max > 0 {
			maxValue = data.SafeOutputs.DispatchWorkflow.Max
		}
		dispatchWorkflowConfig["max"] = maxValue

		// Only add if it has fields
		if len(dispatchWorkflowConfig) > 0 {
			safeOutputsConfig["dispatch_workflow"] = dispatchWorkflowConfig
		}
	}

	configJSON, _ := json.Marshal(safeOutputsConfig)
	return string(configJSON)
}

// buildEnabledToolsSet creates a set of enabled tool names based on safe outputs configuration
// This is used by generateFilteredToolsJSON to determine which tools to include
func buildEnabledToolsSet(safeOutputs *SafeOutputsConfig) map[string]bool {
	enabledTools := make(map[string]bool)

	if safeOutputs == nil {
		return enabledTools
	}

	// Use the registry to determine which tools are enabled
	for toolName, handler := range toolRegistry {
		_, shouldInclude := handler(safeOutputs)
		if shouldInclude {
			// Skip the special issue creation tools as they're handled separately
			if toolName != "create_missing_tool_issue" && toolName != "create_missing_data_issue" {
				enabledTools[toolName] = true
			}
		}
	}

	return enabledTools
}

// addCustomJobTools adds custom job tool definitions to the filtered tools list
func addCustomJobTools(filteredTools *[]map[string]any, safeOutputs *SafeOutputsConfig) {
	if len(safeOutputs.Jobs) == 0 {
		return
	}

	safeOutputsConfigBuilderLog.Printf("Adding %d custom job tools", len(safeOutputs.Jobs))

	// Sort job names for deterministic output
	// This ensures compiled workflows have consistent tool ordering
	jobNames := make([]string, 0, len(safeOutputs.Jobs))
	for jobName := range safeOutputs.Jobs {
		jobNames = append(jobNames, jobName)
	}
	sort.Strings(jobNames)

	// Iterate over jobs in sorted order
	for _, jobName := range jobNames {
		jobConfig := safeOutputs.Jobs[jobName]
		// Normalize job name to use underscores for consistency
		normalizedJobName := stringutil.NormalizeSafeOutputIdentifier(jobName)

		// Create the tool definition for this custom job
		customTool := generateCustomJobToolDefinition(normalizedJobName, jobConfig)
		*filteredTools = append(*filteredTools, customTool)
	}
}

// addDispatchWorkflowTools adds dispatch workflow tool definitions to the filtered tools list
func addDispatchWorkflowTools(filteredTools *[]map[string]any, data *WorkflowData, markdownPath string) {
	if data.SafeOutputs.DispatchWorkflow == nil || len(data.SafeOutputs.DispatchWorkflow.Workflows) == 0 {
		return
	}

	safeOutputsConfigBuilderLog.Printf("Adding %d dispatch_workflow tools", len(data.SafeOutputs.DispatchWorkflow.Workflows))

	// Initialize WorkflowFiles map if not already initialized
	if data.SafeOutputs.DispatchWorkflow.WorkflowFiles == nil {
		data.SafeOutputs.DispatchWorkflow.WorkflowFiles = make(map[string]string)
	}

	for _, workflowName := range data.SafeOutputs.DispatchWorkflow.Workflows {
		// Find the workflow file in multiple locations
		fileResult, err := findWorkflowFile(workflowName, markdownPath)
		if err != nil {
			safeOutputsConfigBuilderLog.Printf("Warning: error finding workflow %s: %v", workflowName, err)
			// Continue with empty inputs
			tool := generateDispatchWorkflowTool(workflowName, make(map[string]any))
			*filteredTools = append(*filteredTools, tool)
			continue
		}

		// Determine which file to use - priority: .lock.yml > .yml
		var workflowPath string
		var extension string
		if fileResult.lockExists {
			workflowPath = fileResult.lockPath
			extension = ".lock.yml"
		} else if fileResult.ymlExists {
			workflowPath = fileResult.ymlPath
			extension = ".yml"
		} else {
			safeOutputsConfigBuilderLog.Printf("Warning: workflow file not found for %s (only .md exists, needs compilation)", workflowName)
			// Continue with empty inputs
			tool := generateDispatchWorkflowTool(workflowName, make(map[string]any))
			*filteredTools = append(*filteredTools, tool)
			continue
		}

		// Store the file extension for runtime use
		data.SafeOutputs.DispatchWorkflow.WorkflowFiles[workflowName] = extension

		// Extract workflow_dispatch inputs
		workflowInputs, err := extractWorkflowDispatchInputs(workflowPath)
		if err != nil {
			safeOutputsConfigBuilderLog.Printf("Warning: failed to extract inputs for workflow %s from %s: %v", workflowName, workflowPath, err)
			// Continue with empty inputs
			workflowInputs = make(map[string]any)
		}

		// Generate tool schema
		tool := generateDispatchWorkflowTool(workflowName, workflowInputs)
		*filteredTools = append(*filteredTools, tool)
	}
}
