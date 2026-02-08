// @ts-check
/// <reference types="@actions/github-script" />

/**
 * @typedef {import('./types/handler-factory').HandlerFactoryFunction} HandlerFactoryFunction
 */

const { createCloseEntityHandler } = require("./close_entity_helpers.cjs");

/**
 * Get discussion details using GraphQL with pagination for labels
 * @param {any} github - GitHub GraphQL instance
 * @param {string} owner - Repository owner
 * @param {string} repo - Repository name
 * @param {number} discussionNumber - Discussion number
 * @returns {Promise<{id: string, number: number, title: string, category: {name: string}, labels: Array<{name: string}>, html_url: string, state: string}>} Discussion details
 */
async function getDiscussionDetails(github, owner, repo, discussionNumber) {
  // Fetch all labels with pagination
  const allLabels = [];
  let hasNextPage = true;
  let cursor = null;
  let discussion = null;

  while (hasNextPage) {
    const query = await github.graphql(
      `
      query($owner: String!, $repo: String!, $num: Int!, $cursor: String) {
        repository(owner: $owner, name: $repo) {
          discussion(number: $num) {
            id
            title
            category {
              name
            }
            url
            labels(first: 100, after: $cursor) {
              nodes {
                name
              }
              pageInfo {
                hasNextPage
                endCursor
              }
            }
          }
        }
      }`,
      { owner, repo, num: discussionNumber, cursor }
    );

    if (!query?.repository?.discussion) {
      throw new Error(`Discussion #${discussionNumber} not found in ${owner}/${repo}`);
    }

    // Store the discussion metadata from the first query
    if (!discussion) {
      discussion = {
        id: query.repository.discussion.id,
        title: query.repository.discussion.title,
        category: query.repository.discussion.category,
        url: query.repository.discussion.url,
      };
    }

    const labels = query.repository.discussion.labels?.nodes || [];
    allLabels.push(...labels);

    hasNextPage = query.repository.discussion.labels?.pageInfo?.hasNextPage || false;
    cursor = query.repository.discussion.labels?.pageInfo?.endCursor || null;
  }

  if (!discussion) {
    throw new Error(`Failed to fetch discussion #${discussionNumber}`);
  }

  // Adapt the response to match the expected interface
  return {
    id: discussion.id,
    number: discussionNumber,
    title: discussion.title,
    category: discussion.category,
    labels: allLabels, // Already in {name: string} format
    html_url: discussion.url,
    state: "open", // Discussions don't expose state via GraphQL, assume open
  };
}

/**
 * Add comment to a GitHub Discussion using GraphQL
 * Note: Discussions use GraphQL node IDs instead of (owner, repo, number)
 * @param {any} github - GitHub GraphQL instance
 * @param {string} owner - Repository owner (unused for GraphQL)
 * @param {string} repo - Repository name (unused for GraphQL)
 * @param {number} discussionNumber - Discussion number
 * @param {string} message - Comment body
 * @returns {Promise<{id: number, html_url: string}>} Comment details
 */
async function addDiscussionComment(github, owner, repo, discussionNumber, message) {
  // First fetch the discussion to get its node ID
  const discussion = await getDiscussionDetails(github, owner, repo, discussionNumber);

  const result = await github.graphql(
    `
    mutation($dId: ID!, $body: String!) {
      addDiscussionComment(input: { discussionId: $dId, body: $body }) {
        comment { 
          id 
          url
        }
      }
    }`,
    { dId: discussion.id, body: message }
  );

  // Adapt the response to match the expected interface
  return {
    id: result.addDiscussionComment.comment.id,
    html_url: result.addDiscussionComment.comment.url,
  };
}

/**
 * Close a GitHub Discussion using GraphQL
 * Note: Discussions use GraphQL node IDs instead of (owner, repo, number)
 * @param {any} github - GitHub GraphQL instance
 * @param {string} owner - Repository owner (unused for GraphQL)
 * @param {string} repo - Repository name (unused for GraphQL)
 * @param {number} discussionNumber - Discussion number
 * @returns {Promise<{number: number, html_url: string, title: string}>} Discussion details
 */
async function closeDiscussion(github, owner, repo, discussionNumber) {
  // First fetch the discussion to get its node ID and title
  const discussion = await getDiscussionDetails(github, owner, repo, discussionNumber);

  const mutation = `
    mutation($dId: ID!) {
      closeDiscussion(input: { discussionId: $dId }) {
        discussion { 
          id
          url
        }
      }
    }`;

  const result = await github.graphql(mutation, { dId: discussion.id });

  // Adapt the response to match the expected interface
  return {
    number: discussionNumber,
    html_url: result.closeDiscussion.discussion.url,
    title: discussion.title,
  };
}

/**
 * Main handler factory for close_discussion
 * Returns a message handler function that processes individual close_discussion messages
 * @type {HandlerFactoryFunction}
 */
async function main(config = {}) {
  // Use the generic factory with discussion-specific APIs
  return createCloseEntityHandler(
    "discussion",
    "discussion_number",
    "discussion",
    {
      getDetails: getDiscussionDetails,
      addComment: addDiscussionComment,
      closeEntity: closeDiscussion,
    },
    config
  );
}

module.exports = { main };
