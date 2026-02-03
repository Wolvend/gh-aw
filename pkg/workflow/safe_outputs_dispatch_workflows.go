package workflow

import (
	"github.com/github/gh-aw/pkg/logger"
)

var safeOutputsDispatchWorkflowLog = logger.New("workflow:safe_outputs_dispatch_workflows")

// ========================================
// Safe Outputs - Dispatch Workflow File Handling
// ========================================
//
// This file handles workflow file discovery and mapping for dispatch-workflow
// safe output configuration. It populates the WorkflowFiles map which stores
// the file extension (.lock.yml or .yml) for each workflow to be dispatched.

// populateDispatchWorkflowFiles populates the WorkflowFiles map for dispatch-workflow configuration.
// This must be called before generateSafeOutputsConfig to ensure workflow file extensions are available.
func populateDispatchWorkflowFiles(data *WorkflowData, markdownPath string) {
	if data.SafeOutputs == nil || data.SafeOutputs.DispatchWorkflow == nil {
		return
	}

	if len(data.SafeOutputs.DispatchWorkflow.Workflows) == 0 {
		return
	}

	safeOutputsDispatchWorkflowLog.Printf("Populating workflow files for %d dispatch workflows", len(data.SafeOutputs.DispatchWorkflow.Workflows))

	// Initialize WorkflowFiles map if not already initialized
	if data.SafeOutputs.DispatchWorkflow.WorkflowFiles == nil {
		data.SafeOutputs.DispatchWorkflow.WorkflowFiles = make(map[string]string)
	}

	for _, workflowName := range data.SafeOutputs.DispatchWorkflow.Workflows {
		// Find the workflow file
		fileResult, err := findWorkflowFile(workflowName, markdownPath)
		if err != nil {
			safeOutputsDispatchWorkflowLog.Printf("Warning: error finding workflow %s: %v", workflowName, err)
			continue
		}

		// Determine which file to use - priority: .lock.yml > .yml
		var extension string
		if fileResult.lockExists {
			extension = ".lock.yml"
		} else if fileResult.ymlExists {
			extension = ".yml"
		} else {
			safeOutputsDispatchWorkflowLog.Printf("Warning: workflow file not found for %s (only .md exists, needs compilation)", workflowName)
			continue
		}

		// Store the file extension for runtime use
		data.SafeOutputs.DispatchWorkflow.WorkflowFiles[workflowName] = extension
		safeOutputsDispatchWorkflowLog.Printf("Mapped workflow %s to extension %s", workflowName, extension)
	}
}
