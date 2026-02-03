package workflow

// ========================================
// Safe Outputs - Tool Configuration Registry
// ========================================
//
// This file implements a registry pattern to eliminate repetitive conditional
// logic in config generation. Instead of 84+ if-checks, we use a data-driven
// approach with handler functions for each tool type.

// SafeOutputToolHandler is a function that generates configuration for a specific safe output tool
type SafeOutputToolHandler func(safeOutputs *SafeOutputsConfig) (config map[string]any, shouldInclude bool)

// toolRegistry maps tool names to their configuration handlers
var toolRegistry = map[string]SafeOutputToolHandler{
	"create_issue":                          handleCreateIssue,
	"create_agent_session":                  handleCreateAgentSession,
	"add_comment":                           handleAddComment,
	"create_discussion":                     handleCreateDiscussion,
	"close_discussion":                      handleCloseDiscussion,
	"close_issue":                           handleCloseIssue,
	"create_pull_request":                   handleCreatePullRequest,
	"create_pull_request_review_comment":    handleCreatePullRequestReviewComment,
	"create_code_scanning_alert":            handleCreateCodeScanningAlert,
	"autofix_code_scanning_alert":           handleAutofixCodeScanningAlert,
	"add_labels":                            handleAddLabels,
	"remove_labels":                         handleRemoveLabels,
	"add_reviewer":                          handleAddReviewer,
	"assign_milestone":                      handleAssignMilestone,
	"assign_to_agent":                       handleAssignToAgent,
	"assign_to_user":                        handleAssignToUser,
	"update_issue":                          handleUpdateIssue,
	"update_discussion":                     handleUpdateDiscussion,
	"update_pull_request":                   handleUpdatePullRequest,
	"mark_pull_request_as_ready_for_review": handleMarkPullRequestAsReadyForReview,
	"push_to_pull_request_branch":           handlePushToPullRequestBranch,
	"upload_asset":                          handleUploadAsset,
	"missing_tool":                          handleMissingTool,
	"missing_data":                          handleMissingData,
	"update_project":                        handleUpdateProject,
	"create_project_status_update":          handleCreateProjectStatusUpdate,
	"create_project":                        handleCreateProject,
	"update_release":                        handleUpdateRelease,
	"link_sub_issue":                        handleLinkSubIssue,
	"noop":                                  handleNoOp,
	"hide_comment":                          handleHideComment,
	"create_missing_tool_issue":             handleCreateMissingToolIssue,
	"create_missing_data_issue":             handleCreateMissingDataIssue,
}

// Tool configuration handler functions

func handleCreateIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateIssues == nil {
		return nil, false
	}
	config := generateMaxWithAllowedLabelsConfig(
		safeOutputs.CreateIssues.Max,
		1, // default max
		safeOutputs.CreateIssues.AllowedLabels,
	)
	// Add group flag if enabled
	if safeOutputs.CreateIssues.Group {
		config["group"] = true
	}
	// Add expires value if set (0 means explicitly disabled or not set)
	if safeOutputs.CreateIssues.Expires > 0 {
		config["expires"] = safeOutputs.CreateIssues.Expires
	}
	return config, true
}

func handleCreateAgentSession(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateAgentSessions == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.CreateAgentSessions.Max, 1), true
}

func handleAddComment(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AddComments == nil {
		return nil, false
	}
	return generateMaxWithTargetConfig(
		safeOutputs.AddComments.Max,
		1, // default max
		safeOutputs.AddComments.Target,
	), true
}

func handleCreateDiscussion(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateDiscussions == nil {
		return nil, false
	}
	config := generateMaxWithAllowedLabelsConfig(
		safeOutputs.CreateDiscussions.Max,
		1, // default max
		safeOutputs.CreateDiscussions.AllowedLabels,
	)
	// Add expires value if set (0 means explicitly disabled or not set)
	if safeOutputs.CreateDiscussions.Expires > 0 {
		config["expires"] = safeOutputs.CreateDiscussions.Expires
	}
	return config, true
}

func handleCloseDiscussion(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CloseDiscussions == nil {
		return nil, false
	}
	return generateMaxWithDiscussionFieldsConfig(
		safeOutputs.CloseDiscussions.Max,
		1, // default max
		safeOutputs.CloseDiscussions.RequiredCategory,
		safeOutputs.CloseDiscussions.RequiredLabels,
		safeOutputs.CloseDiscussions.RequiredTitlePrefix,
	), true
}

func handleCloseIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CloseIssues == nil {
		return nil, false
	}
	return generateMaxWithRequiredFieldsConfig(
		safeOutputs.CloseIssues.Max,
		1, // default max
		safeOutputs.CloseIssues.RequiredLabels,
		safeOutputs.CloseIssues.RequiredTitlePrefix,
	), true
}

func handleCreatePullRequest(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreatePullRequests == nil {
		return nil, false
	}
	return generatePullRequestConfig(
		safeOutputs.CreatePullRequests.AllowedLabels,
		safeOutputs.CreatePullRequests.AllowEmpty,
		safeOutputs.CreatePullRequests.AutoMerge,
		safeOutputs.CreatePullRequests.Expires,
	), true
}

func handleCreatePullRequestReviewComment(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreatePullRequestReviewComments == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.CreatePullRequestReviewComments.Max, 10), true
}

func handleCreateCodeScanningAlert(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateCodeScanningAlerts == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.CreateCodeScanningAlerts.Max, 0), true
}

func handleAutofixCodeScanningAlert(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AutofixCodeScanningAlert == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.AutofixCodeScanningAlert.Max, 10), true
}

func handleAddLabels(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AddLabels == nil {
		return nil, false
	}
	return generateMaxWithAllowedConfig(
		safeOutputs.AddLabels.Max,
		3, // default max
		safeOutputs.AddLabels.Allowed,
	), true
}

func handleRemoveLabels(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.RemoveLabels == nil {
		return nil, false
	}
	return generateMaxWithAllowedConfig(
		safeOutputs.RemoveLabels.Max,
		3, // default max
		safeOutputs.RemoveLabels.Allowed,
	), true
}

func handleAddReviewer(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AddReviewer == nil {
		return nil, false
	}
	return generateMaxWithReviewersConfig(
		safeOutputs.AddReviewer.Max,
		3, // default max
		safeOutputs.AddReviewer.Reviewers,
	), true
}

func handleAssignMilestone(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AssignMilestone == nil {
		return nil, false
	}
	return generateMaxWithAllowedConfig(
		safeOutputs.AssignMilestone.Max,
		1, // default max
		safeOutputs.AssignMilestone.Allowed,
	), true
}

func handleAssignToAgent(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AssignToAgent == nil {
		return nil, false
	}
	return generateAssignToAgentConfig(
		safeOutputs.AssignToAgent.Max,
		safeOutputs.AssignToAgent.DefaultAgent,
		safeOutputs.AssignToAgent.Target,
		safeOutputs.AssignToAgent.Allowed,
	), true
}

func handleAssignToUser(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.AssignToUser == nil {
		return nil, false
	}
	return generateMaxWithAllowedConfig(
		safeOutputs.AssignToUser.Max,
		1, // default max
		safeOutputs.AssignToUser.Allowed,
	), true
}

func handleUpdateIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UpdateIssues == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.UpdateIssues.Max, 1), true
}

func handleUpdateDiscussion(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UpdateDiscussions == nil {
		return nil, false
	}
	return generateMaxWithAllowedLabelsConfig(
		safeOutputs.UpdateDiscussions.Max,
		1, // default max
		safeOutputs.UpdateDiscussions.AllowedLabels,
	), true
}

func handleUpdatePullRequest(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UpdatePullRequests == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.UpdatePullRequests.Max, 1), true
}

func handleMarkPullRequestAsReadyForReview(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.MarkPullRequestAsReadyForReview == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.MarkPullRequestAsReadyForReview.Max, 10), true
}

func handlePushToPullRequestBranch(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.PushToPullRequestBranch == nil {
		return nil, false
	}
	return generateMaxWithTargetConfig(
		safeOutputs.PushToPullRequestBranch.Max,
		0, // default: unlimited
		safeOutputs.PushToPullRequestBranch.Target,
	), true
}

func handleUploadAsset(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UploadAssets == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.UploadAssets.Max, 0), true
}

func handleMissingTool(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.MissingTool == nil {
		return nil, false
	}
	config := make(map[string]any)
	// Add max if set
	if safeOutputs.MissingTool.Max > 0 {
		config["max"] = safeOutputs.MissingTool.Max
	}
	return config, true
}

func handleMissingData(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.MissingData == nil {
		return nil, false
	}
	config := make(map[string]any)
	// Add max if set
	if safeOutputs.MissingData.Max > 0 {
		config["max"] = safeOutputs.MissingData.Max
	}
	return config, true
}

func handleUpdateProject(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UpdateProjects == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.UpdateProjects.Max, 10), true
}

func handleCreateProjectStatusUpdate(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateProjectStatusUpdates == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.CreateProjectStatusUpdates.Max, 10), true
}

func handleCreateProject(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.CreateProjects == nil {
		return nil, false
	}
	config := generateMaxConfig(safeOutputs.CreateProjects.Max, 1)
	// Add target-owner if specified
	if safeOutputs.CreateProjects.TargetOwner != "" {
		config["target_owner"] = safeOutputs.CreateProjects.TargetOwner
	}
	// Add title-prefix if specified
	if safeOutputs.CreateProjects.TitlePrefix != "" {
		config["title_prefix"] = safeOutputs.CreateProjects.TitlePrefix
	}
	return config, true
}

func handleUpdateRelease(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.UpdateRelease == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.UpdateRelease.Max, 1), true
}

func handleLinkSubIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.LinkSubIssue == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.LinkSubIssue.Max, 5), true
}

func handleNoOp(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.NoOp == nil {
		return nil, false
	}
	return generateMaxConfig(safeOutputs.NoOp.Max, 1), true
}

func handleHideComment(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.HideComment == nil {
		return nil, false
	}
	return generateHideCommentConfig(
		safeOutputs.HideComment.Max,
		5, // default max
		safeOutputs.HideComment.AllowedReasons,
	), true
}

func handleCreateMissingToolIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.MissingTool == nil || !safeOutputs.MissingTool.CreateIssue {
		return nil, false
	}
	config := make(map[string]any)
	config["max"] = 1 // Only create one issue per workflow run

	if safeOutputs.MissingTool.TitlePrefix != "" {
		config["title_prefix"] = safeOutputs.MissingTool.TitlePrefix
	}

	if len(safeOutputs.MissingTool.Labels) > 0 {
		config["labels"] = safeOutputs.MissingTool.Labels
	}
	return config, true
}

func handleCreateMissingDataIssue(safeOutputs *SafeOutputsConfig) (map[string]any, bool) {
	if safeOutputs.MissingData == nil || !safeOutputs.MissingData.CreateIssue {
		return nil, false
	}
	config := make(map[string]any)
	config["max"] = 1 // Only create one issue per workflow run

	if safeOutputs.MissingData.TitlePrefix != "" {
		config["title_prefix"] = safeOutputs.MissingData.TitlePrefix
	}

	if len(safeOutputs.MissingData.Labels) > 0 {
		config["labels"] = safeOutputs.MissingData.Labels
	}
	return config, true
}
