---
safe-inputs:
  add:
    description: "Add a workflow from a repository to .github/workflows. This tool is accessible as 'safeinputs-add'. Provide the workflow specification (e.g., 'owner/repo/workflow-name')."
    inputs:
      workflow:
        type: string
        description: "Workflow specification (e.g., 'owner/repo/workflow-name' or 'owner/repo/workflow-name@v1.0.0')"
        required: true
      engine:
        type: string
        description: "AI engine to use (copilot, claude, codex)"
        required: false
      force:
        type: boolean
        description: "Overwrite existing workflow files"
        required: false
      verbose:
        type: boolean
        description: "Enable verbose output"
        required: false
    go: |
      import (
        "github.com/github/gh-aw/pkg/cli"
      )
      
      // Extract inputs
      workflow := inputs["workflow"].(string)
      engine, _ := inputs["engine"].(string)
      force, _ := inputs["force"].(bool)
      verbose, _ := inputs["verbose"].(bool)
      
      // Call AddWorkflows function directly instead of executing gh aw add
      workflows := []string{workflow}
      addResult, err := cli.AddWorkflows(
        workflows,
        1,              // number of copies
        verbose,        // verbose flag
        engine,         // engine override
        "",             // name (empty for default)
        force,          // force flag
        "",             // append text
        false,          // create PR
        false,          // push
        false,          // no gitattributes
        "",             // workflow dir
        false,          // no stop-after
        "",             // stop-after value
      )
      
      if err != nil {
        result := map[string]any{
          "success": false,
          "error": err.Error(),
        }
        json.NewEncoder(os.Stdout).Encode(result)
        return
      }
      
      result := map[string]any{
        "success": true,
        "workflow": workflow,
        "pr_number": addResult.PRNumber,
        "pr_url": addResult.PRURL,
        "has_workflow_dispatch": addResult.HasWorkflowDispatch,
      }
      json.NewEncoder(os.Stdout).Encode(result)
    timeout: 300
---

## safeinputs-add Tool

A safe-input tool that adds workflows from repositories by calling the Go add command function directly.

**Note**: This tool has been refactored to call `cli.AddWorkflows()` directly instead of executing `gh aw add` through os/exec, which improves performance and eliminates subprocess overhead.

### Usage

```yaml
imports:
  - shared/add-workflow-safe-input.md
```

### Invocation

```
safeinputs-add with workflow: "owner/repo/workflow-name"
safeinputs-add with workflow: "owner/repo/workflow-name@v1.0.0", engine: "copilot"
safeinputs-add with workflow: "owner/repo/workflow-name", force: true
```
