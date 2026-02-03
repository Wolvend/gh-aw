---
description: Test workflow to demonstrate the refactored add safe-input tool
on:
  workflow_dispatch:
permissions:
  contents: read
  pull-requests: read
  issues: read
timeout-minutes: 10
imports:
  - shared/add-workflow-safe-input.md
---

# Test Add Workflow Tool

This workflow demonstrates the refactored `safeinputs-add` tool that calls the Go function directly.

## Example Usage

Use the safeinputs-add tool with workflow: "githubnext/agentics/ci-doctor@v1.0.0"
