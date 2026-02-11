// @ts-check
/// <reference types="@actions/github-script" />

/**
 * Sanitizes a value for safe use in shell commands by removing/escaping shell metacharacters
 * @param {any} value - The value to sanitize
 * @returns {any} - Sanitized value (strings are processed, other types returned as-is)
 */
function sanitizeForShell(value) {
  if (typeof value !== "string") return value;

  // Remove shell metacharacters
  return value
    .replace(/[$`\\]/g, "\\$&") // Escape $, `, \
    .replace(/[();&|<>]/g, "") // Remove dangerous chars
    .substring(0, 1000); // Limit length
}

module.exports = { sanitizeForShell };
