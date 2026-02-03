// This file contains MCP (Model Context Protocol) validation functions.
// This file consolidates validation logic for MCP server configurations.

package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/github/gh-aw/pkg/console"
	"github.com/github/gh-aw/pkg/logger"
	"github.com/github/gh-aw/pkg/parser"
	"github.com/github/gh-aw/pkg/workflow"
)

var mcpValidationLog = logger.New("cli:mcp_validation")

// validateServerSecrets checks if required environment variables/secrets are available
func validateServerSecrets(config parser.MCPServerConfig, verbose bool, useActionsSecrets bool) error {
	mcpValidationLog.Printf("Validating server secrets: server=%s, type=%s, useActionsSecrets=%v", config.Name, config.Type, useActionsSecrets)

	// Extract secrets from the config
	requiredSecrets := extractSecretsFromConfig(config)

	// Special case: Check for GH_AW_GITHUB_TOKEN when GitHub tool is in remote mode
	if config.Name == "github" && config.Type == "http" {
		mcpValidationLog.Print("GitHub remote mode detected, checking for GH_AW_GITHUB_TOKEN")
		// GitHub remote mode requires GH_AW_GITHUB_TOKEN secret
		// Check if a custom token is already specified in the env
		hasCustomToken := false
		for _, value := range config.Env {
			if strings.Contains(value, "secrets.") && !strings.Contains(value, "GH_AW_GITHUB_TOKEN") {
				// Custom token specified, no need to check GH_AW_GITHUB_TOKEN
				hasCustomToken = true
				break
			}
		}

		if !hasCustomToken {
			// Add GH_AW_GITHUB_TOKEN to required secrets if not already present
			alreadyPresent := false
			for _, secret := range requiredSecrets {
				if secret.Name == "GH_AW_GITHUB_TOKEN" {
					alreadyPresent = true
					break
				}
			}
			if !alreadyPresent {
				requiredSecrets = append(requiredSecrets, SecretInfo{
					Name:   "GH_AW_GITHUB_TOKEN",
					EnvKey: "GITHUB_TOKEN",
				})
			}
		}
	}

	if len(requiredSecrets) == 0 {
		mcpValidationLog.Printf("No required secrets found, validating %d environment variables", len(config.Env))
		// No secrets required, proceed with normal env var validation
		for key, value := range config.Env {
			// Check if value contains variable references
			if strings.Contains(value, "${") {
				// Extract variable name (simplified parsing)
				if strings.Contains(value, "secrets.") {
					// This should have been caught by extractSecretsFromConfig
					continue
				}
				if strings.Contains(value, "GH_TOKEN") || strings.Contains(value, "GITHUB_TOKEN") || strings.Contains(value, "GITHUB_PERSONAL_ACCESS_TOKEN") {
					if token, err := parser.GetGitHubToken(); err != nil {
						return fmt.Errorf("GitHub token not found in environment (set GH_TOKEN or GITHUB_TOKEN)")
					} else {
						config.Env[key] = token
					}
				}
				// Handle our placeholder for GitHub token requirement
				if strings.Contains(value, "GITHUB_TOKEN_REQUIRED") {
					if token, err := parser.GetGitHubToken(); err != nil {
						return fmt.Errorf("GitHub token required but not available: %w", err)
					} else {
						config.Env[key] = token
					}
				}
			} else {
				// For direct environment variable values (not containing ${}),
				// check if they represent actual token values
				if value == "" {
					return fmt.Errorf("environment variable '%s' has empty value", key)
				}
				// If value contains "GITHUB_TOKEN_REQUIRED", treat it as needing validation
				if strings.Contains(value, "GITHUB_TOKEN_REQUIRED") {
					if token, err := parser.GetGitHubToken(); err != nil {
						return fmt.Errorf("GitHub token required but not available: %w", err)
					} else {
						config.Env[key] = token
					}
				} else {
					// Automatically try to get GitHub token for GitHub-related environment variables
					if key == "GITHUB_PERSONAL_ACCESS_TOKEN" || key == "GITHUB_TOKEN" || key == "GH_TOKEN" {
						if actualValue := os.Getenv(key); actualValue == "" {
							// Try to automatically get the GitHub token
							if token, err := parser.GetGitHubToken(); err == nil {
								config.Env[key] = token
							} else {
								return fmt.Errorf("GitHub token required for '%s' but not available: %w", key, err)
							}
						}
					} else {
						// For backward compatibility: check if environment variable with this name exists
						// This preserves the original behavior for existing tests
						if actualValue := os.Getenv(key); actualValue == "" {
							return fmt.Errorf("environment variable '%s' not set", key)
						}
					}
				}
			}
		}
		return nil
	}

	// Check availability of required secrets
	mcpValidationLog.Printf("Checking availability of %d required secrets", len(requiredSecrets))
	secretsStatus := checkSecretsAvailability(requiredSecrets, useActionsSecrets)

	// Separate secrets by availability
	var availableSecrets []SecretInfo
	var missingSecrets []SecretInfo

	for _, secret := range secretsStatus {
		if secret.Available {
			availableSecrets = append(availableSecrets, secret)
		} else {
			missingSecrets = append(missingSecrets, secret)
		}
	}

	// Display information about secrets
	if verbose {
		if len(availableSecrets) > 0 {
			fmt.Fprintln(os.Stderr, console.FormatInfoMessage(fmt.Sprintf("Found %d available secret(s):", len(availableSecrets))))
			for _, secret := range availableSecrets {
				source := "environment"
				if secret.Source == "actions" {
					source = "GitHub Actions"
				}
				fmt.Fprintln(os.Stderr, console.FormatSuccessMessage(fmt.Sprintf("  ✓ %s (from %s)", secret.Name, source)))
			}
		}
	}

	// Warn about missing secrets
	if len(missingSecrets) > 0 {
		mcpValidationLog.Printf("Found %d missing secrets", len(missingSecrets))
		fmt.Fprintln(os.Stderr, console.FormatWarningMessage(fmt.Sprintf("⚠️  %d required secret(s) not found:", len(missingSecrets))))
		for _, secret := range missingSecrets {
			fmt.Fprintln(os.Stderr, console.FormatWarningMessage(fmt.Sprintf("  ✗ %s", secret.Name)))
		}
	}

	mcpValidationLog.Printf("Secret validation completed: available=%d, missing=%d", len(availableSecrets), len(missingSecrets))
	return nil
}

// checkGitHubActionsPermissions checks if the GitHub token has permissions to read workflow runs
// Returns true if the token has the required permissions, false otherwise
func checkGitHubActionsPermissions() bool {
	mcpValidationLog.Print("Checking GitHub token permissions for Actions workflow runs")

	// Create a context with timeout for the API call
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to list workflow runs using the GitHub API
	// We use a minimal query to test permissions without fetching too much data
	cmd := workflow.ExecGHContext(ctx, "api", "repos/{owner}/{repo}/actions/runs", "--jq", ".total_count", "-X", "GET")
	output, err := cmd.CombinedOutput()

	if err != nil {
		mcpValidationLog.Printf("GitHub Actions permissions check failed: %v", err)
		mcpValidationLog.Printf("Output: %s", string(output))

		// Check for specific permission errors
		outputStr := string(output)
		if strings.Contains(outputStr, "403") || strings.Contains(outputStr, "Forbidden") {
			mcpValidationLog.Print("Token does not have 'actions:read' permission (403 Forbidden)")
			fmt.Fprintln(os.Stderr, console.FormatWarningMessage("⚠️  GitHub token does not have 'actions:read' permission"))
			fmt.Fprintln(os.Stderr, console.FormatInfoMessage("The 'logs' and 'audit' tools will not be available"))
			return false
		}
		if strings.Contains(outputStr, "404") || strings.Contains(outputStr, "Not Found") {
			mcpValidationLog.Print("Repository not found or token lacks access (404 Not Found)")
			fmt.Fprintln(os.Stderr, console.FormatWarningMessage("⚠️  Repository not found or token lacks repository access"))
			fmt.Fprintln(os.Stderr, console.FormatInfoMessage("The 'logs' and 'audit' tools will not be available"))
			return false
		}
		if strings.Contains(outputStr, "401") || strings.Contains(outputStr, "Unauthorized") {
			mcpValidationLog.Print("Token authentication failed (401 Unauthorized)")
			fmt.Fprintln(os.Stderr, console.FormatWarningMessage("⚠️  GitHub token authentication failed"))
			fmt.Fprintln(os.Stderr, console.FormatInfoMessage("The 'logs' and 'audit' tools will not be available"))
			return false
		}

		// For other errors, log but don't assume permission failure
		mcpValidationLog.Printf("Unable to verify GitHub Actions permissions: %v", err)
		fmt.Fprintln(os.Stderr, console.FormatWarningMessage(fmt.Sprintf("⚠️  Could not verify GitHub Actions permissions: %v", err)))
		fmt.Fprintln(os.Stderr, console.FormatInfoMessage("The 'logs' and 'audit' tools may not be available"))
		return false
	}

	// Successfully queried workflow runs - token has required permissions
	mcpValidationLog.Printf("GitHub Actions permissions check succeeded. Output: %s", strings.TrimSpace(string(output)))
	fmt.Fprintln(os.Stderr, console.FormatSuccessMessage("✅ GitHub token has 'actions:read' permission"))
	return true
}

// validateMCPServerConfiguration validates that the CLI is properly configured
// by running the status command as a test
func validateMCPServerConfiguration(cmdPath string) error {
	mcpValidationLog.Printf("Validating MCP server configuration: cmdPath=%s", cmdPath)

	// Try to run the status command to verify CLI is working
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var cmd *exec.Cmd
	if cmdPath != "" {
		mcpValidationLog.Printf("Using custom command path: %s", cmdPath)
		// Use custom command path
		cmd = exec.CommandContext(ctx, cmdPath, "status")
	} else {
		mcpValidationLog.Print("Using default gh aw command with proper token handling")
		// Use default gh aw command with proper token handling
		cmd = workflow.ExecGHContext(ctx, "aw", "status")
	}
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check for common error cases
		if ctx.Err() == context.DeadlineExceeded {
			mcpValidationLog.Print("Status command timed out")
			errMsg := "status command timed out - this may indicate a configuration issue"
			fmt.Fprintln(os.Stderr, console.FormatErrorMessage(errMsg))
			return fmt.Errorf("status command timed out - this may indicate a configuration issue")
		}

		mcpValidationLog.Printf("Status command failed: %v", err)

		// If the command failed, provide helpful error message
		if cmdPath != "" {
			errMsg := fmt.Sprintf("failed to run status command with custom command '%s': %v\nOutput: %s\n\nPlease ensure:\n  - The command path is correct and executable\n  - You are in a git repository with .github/workflows directory", cmdPath, err, string(output))
			fmt.Fprintln(os.Stderr, console.FormatErrorMessage(errMsg))
			return fmt.Errorf("failed to run status command with custom command '%s': %w\nOutput: %s\n\nPlease ensure:\n  - The command path is correct and executable\n  - You are in a git repository with .github/workflows directory", cmdPath, err, string(output))
		}
		errMsg := fmt.Sprintf("failed to run status command: %v\nOutput: %s\n\nPlease ensure:\n  - gh CLI is installed and in PATH\n  - gh aw extension is installed (run: gh extension install github/gh-aw)\n  - You are in a git repository with .github/workflows directory", err, string(output))
		fmt.Fprintln(os.Stderr, console.FormatErrorMessage(errMsg))
		return fmt.Errorf("failed to run status command: %w\nOutput: %s\n\nPlease ensure:\n  - gh CLI is installed and in PATH\n  - gh aw extension is installed (run: gh extension install github/gh-aw)\n  - You are in a git repository with .github/workflows directory", err, string(output))
	}

	// Status command succeeded - configuration is valid
	mcpValidationLog.Print("MCP server configuration validated successfully")
	fmt.Fprintln(os.Stderr, console.FormatSuccessMessage("✅ Configuration validated successfully"))
	return nil
}
