package workflow

import (
	"fmt"
	"sort"

	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/stringutil"
)

var safeOutputsToolDefinitionsLog = logger.New("workflow:safe_outputs_tool_definitions")

// ========================================
// Safe Outputs - MCP Tool Definition Generation
// ========================================
//
// This file contains functions for generating MCP tool definitions from
// workflow configurations. It handles both custom job tools and dispatch
// workflow tools, creating proper JSON Schema inputSchema definitions.

// generateCustomJobToolDefinition creates an MCP tool definition for a custom safe-output job
// Returns a map representing the tool definition in MCP format with name, description, and inputSchema
func generateCustomJobToolDefinition(jobName string, jobConfig *SafeJobConfig) map[string]any {
	safeOutputsToolDefinitionsLog.Printf("Generating tool definition for custom job: %s", jobName)

	// Build the tool definition
	tool := map[string]any{
		"name": jobName,
	}

	// Add description if present
	if jobConfig.Description != "" {
		tool["description"] = jobConfig.Description
	} else {
		// Provide a default description if none is specified
		tool["description"] = fmt.Sprintf("Execute the %s custom job", jobName)
	}

	// Build the input schema
	inputSchema := map[string]any{
		"type":       "object",
		"properties": make(map[string]any),
	}

	// Track required fields
	var requiredFields []string

	// Add each input to the schema
	if len(jobConfig.Inputs) > 0 {
		properties := inputSchema["properties"].(map[string]any)

		for inputName, inputDef := range jobConfig.Inputs {
			property := map[string]any{}

			// Add description
			if inputDef.Description != "" {
				property["description"] = inputDef.Description
			}

			// Convert type to JSON Schema type
			switch inputDef.Type {
			case "choice":
				// Choice inputs are strings with enum constraints
				property["type"] = "string"
				if len(inputDef.Options) > 0 {
					property["enum"] = inputDef.Options
				}
			case "boolean":
				property["type"] = "boolean"
			case "number":
				property["type"] = "number"
			case "string", "":
				// Default to string if type is not specified
				property["type"] = "string"
			default:
				// For any unknown type, default to string
				property["type"] = "string"
			}

			// Add default value if present
			if inputDef.Default != nil {
				property["default"] = inputDef.Default
			}

			// Track required fields
			if inputDef.Required {
				requiredFields = append(requiredFields, inputName)
			}

			properties[inputName] = property
		}
	}

	// Add required fields array if any inputs are required
	if len(requiredFields) > 0 {
		sort.Strings(requiredFields)
		inputSchema["required"] = requiredFields
	}

	// Prevent additional properties to maintain schema strictness
	inputSchema["additionalProperties"] = false

	tool["inputSchema"] = inputSchema

	safeOutputsToolDefinitionsLog.Printf("Generated tool definition for %s with %d inputs, %d required",
		jobName, len(jobConfig.Inputs), len(requiredFields))

	return tool
}

// generateDispatchWorkflowTool generates an MCP tool definition for a specific workflow
// The tool will be named after the workflow and accept the workflow's defined inputs
func generateDispatchWorkflowTool(workflowName string, workflowInputs map[string]any) map[string]any {
	// Normalize workflow name to use underscores for tool name
	toolName := stringutil.NormalizeSafeOutputIdentifier(workflowName)

	// Build the description
	description := fmt.Sprintf("Dispatch the '%s' workflow with workflow_dispatch trigger. This workflow must support workflow_dispatch and be in .github/workflows/ directory in the same repository.", workflowName)

	// Build input schema properties
	properties := make(map[string]any)
	required := []string{} // No required fields by default

	// Convert GitHub Actions workflow_dispatch inputs to MCP tool schema
	for inputName, inputDef := range workflowInputs {
		inputDefMap, ok := inputDef.(map[string]any)
		if !ok {
			continue
		}

		// Extract input properties
		inputType := "string" // Default type
		inputDescription := fmt.Sprintf("Input parameter '%s' for workflow %s", inputName, workflowName)
		inputRequired := false

		if desc, ok := inputDefMap["description"].(string); ok && desc != "" {
			inputDescription = desc
		}

		if req, ok := inputDefMap["required"].(bool); ok {
			inputRequired = req
		}

		// GitHub Actions workflow_dispatch supports: string, number, boolean, choice, environment
		// Map these to JSON schema types
		if typeStr, ok := inputDefMap["type"].(string); ok {
			switch typeStr {
			case "number":
				inputType = "number"
			case "boolean":
				inputType = "boolean"
			case "choice":
				inputType = "string"
				// Add enum if options are provided
				if options, ok := inputDefMap["options"].([]any); ok && len(options) > 0 {
					properties[inputName] = map[string]any{
						"type":        inputType,
						"description": inputDescription,
						"enum":        options,
					}
					if inputRequired {
						required = append(required, inputName)
					}
					continue
				}
			case "environment":
				inputType = "string"
			}
		}

		properties[inputName] = map[string]any{
			"type":        inputType,
			"description": inputDescription,
		}

		// Add default value if provided
		if defaultVal, ok := inputDefMap["default"]; ok {
			properties[inputName].(map[string]any)["default"] = defaultVal
		}

		if inputRequired {
			required = append(required, inputName)
		}
	}

	// Add internal workflow_name parameter (hidden from description but used internally)
	// This will be injected by the safe output handler

	// Build the complete tool definition
	tool := map[string]any{
		"name":           toolName,
		"description":    description,
		"_workflow_name": workflowName, // Internal metadata for handler routing
		"inputSchema": map[string]any{
			"type":                 "object",
			"properties":           properties,
			"additionalProperties": false,
		},
	}

	if len(required) > 0 {
		tool["inputSchema"].(map[string]any)["required"] = required
	}

	return tool
}
