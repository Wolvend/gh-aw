// @ts-check
/// <reference types="@actions/github-script" />

/**
 * @typedef {import('./types/handler-factory').HandlerFactoryFunction} HandlerFactoryFunction
 */

const { createCloseEntityHandler } = require("./close_entity_helpers.cjs");

/**
 * Get issue details using REST API
 * @param {any} github - GitHub REST API instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} issueNumber - Issue number
 * @returns {Promise<{number: number, title: string, labels: Array<{name: string}>, html_url: string, state: string}>} Issue details
 */
async function getIssueDetails(github, owner, repo, issueNumber) {
  const { data: issue } = await github.rest.issues.get({
    owner,
    repo,
    issue_number: issueNumber,
  });

  if (!issue) {
    throw new Error(`Issue #${issueNumber} not found in ${owner}/${repo}`);
  }

  return issue;
}

/**
 * Add comment to a GitHub Issue using REST API
 * @param {any} github - GitHub REST API instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} issueNumber - Issue number
 * @param {string} message - Comment body
 * @returns {Promise<{id: number, html_url: string}>} Comment details
 */
async function addIssueComment(github, owner, repo, issueNumber, message) {
  const { data: comment } = await github.rest.issues.createComment({
    owner,
    repo,
    issue_number: issueNumber,
    body: message,
  });

  return comment;
}

/**
 * Close a GitHub Issue using REST API
 * @param {any} github - GitHub REST API instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} issueNumber - Issue number
 * @returns {Promise<{number: number, html_url: string, title: string}>} Issue details
 */
async function closeIssue(github, owner, repo, issueNumber) {
  const { data: issue } = await github.rest.issues.update({
    owner,
    repo,
    issue_number: issueNumber,
    state: "closed",
  });

  return issue;
}

/**
 * Main handler factory for close_issue
 * Returns a message handler function that processes individual close_issue messages
 * @type {HandlerFactoryFunction}
 */
async function main(config = {}) {
  // Use the generic factory with issue-specific APIs
  return createCloseEntityHandler(
    "issue",
    "issue_number",
    "issue",
    {
      getDetails: getIssueDetails,
      addComment: addIssueComment,
      closeEntity: closeIssue,
    },
    config
  );
}

module.exports = { main };
