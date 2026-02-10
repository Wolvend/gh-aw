// @ts-check
import { describe, it, expect, beforeEach, afterEach, vi } from "vitest";

describe("expiration_helpers", () => {
  let mockCore;
  let originalEnv;
  let addExpirationComment;

  beforeEach(async () => {
    // Save original environment
    originalEnv = { ...process.env };

    // Reset mocks before each test
    mockCore = {
      info: vi.fn(),
      messages: [],
    };

    // Capture logged messages
    mockCore.info = vi.fn(msg => {
      mockCore.messages.push({ level: "info", message: msg });
    });

    // Set globals
    global.core = mockCore;

    // Clear module cache and re-require to get fresh instance
    vi.resetModules();
    const module = await import("./expiration_helpers.cjs");
    addExpirationComment = module.addExpirationComment;
  });

  afterEach(() => {
    // Restore original environment
    process.env = originalEnv;
  });

  describe("addExpirationComment", () => {
    it("should add expiration comment when expires is set to valid hours", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "24";
      const bodyLines = ["Line 1", "Line 2"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(3);
      expect(bodyLines[2]).toContain("<!-- gh-aw-expires:");
      expect(bodyLines[2]).toContain(" -->");
      expect(mockCore.info).toHaveBeenCalledWith(expect.stringMatching(/Discussion will expire on .+ \(24 hours\)/));
    });

    it("should not add expiration comment when expires is not set", () => {
      const bodyLines = ["Line 1", "Line 2"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(2);
      expect(mockCore.info).not.toHaveBeenCalled();
    });

    it("should not add expiration comment when expires is invalid (NaN)", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "invalid";
      const bodyLines = ["Line 1", "Line 2"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(2);
      expect(mockCore.info).not.toHaveBeenCalled();
    });

    it("should not add expiration comment when expires is zero", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "0";
      const bodyLines = ["Line 1", "Line 2"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(2);
      expect(mockCore.info).not.toHaveBeenCalled();
    });

    it("should not add expiration comment when expires is negative", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "-5";
      const bodyLines = ["Line 1", "Line 2"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(2);
      expect(mockCore.info).not.toHaveBeenCalled();
    });

    it("should work with different environment variable names", () => {
      process.env.GH_AW_ISSUE_EXPIRES = "48";
      const bodyLines = ["Line 1"];

      addExpirationComment(bodyLines, "GH_AW_ISSUE_EXPIRES", "Issue");

      expect(bodyLines.length).toBe(2);
      expect(bodyLines[1]).toContain("<!-- gh-aw-expires:");
      expect(mockCore.info).toHaveBeenCalledWith(expect.stringMatching(/Issue will expire on .+ \(48 hours\)/));
    });

    it("should work with different entity types", () => {
      process.env.GH_AW_PR_EXPIRES = "72";
      const bodyLines = ["Line 1"];

      addExpirationComment(bodyLines, "GH_AW_PR_EXPIRES", "Pull Request");

      expect(bodyLines.length).toBe(2);
      expect(bodyLines[1]).toContain("<!-- gh-aw-expires:");
      expect(mockCore.info).toHaveBeenCalledWith(expect.stringMatching(/Pull Request will expire on .+ \(72 hours\)/));
    });

    it("should append to existing body lines", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "12";
      const bodyLines = ["Existing line 1", "Existing line 2", "Existing line 3"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(4);
      expect(bodyLines[0]).toBe("Existing line 1");
      expect(bodyLines[1]).toBe("Existing line 2");
      expect(bodyLines[2]).toBe("Existing line 3");
      expect(bodyLines[3]).toContain("<!-- gh-aw-expires:");
    });

    it("should handle decimal hours by parsing to integer", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "24.5";
      const bodyLines = ["Line 1"];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(2);
      expect(bodyLines[1]).toContain("<!-- gh-aw-expires:");
      expect(mockCore.info).toHaveBeenCalledWith(expect.stringMatching(/Discussion will expire on .+ \(24 hours\)/));
    });

    it("should calculate correct expiration date", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "1";
      const bodyLines = [];
      const beforeTime = new Date();

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      // Extract ISO date from the comment
      const comment = bodyLines[0];
      const match = comment.match(/<!-- gh-aw-expires: (.+?) -->/);
      expect(match).not.toBeNull();

      const expirationDate = new Date(match[1]);
      const afterTime = new Date();
      afterTime.setHours(afterTime.getHours() + 1);

      // Expiration should be approximately 1 hour from now (allow 1 second variance for test execution)
      const expectedMs = beforeTime.getTime() + 60 * 60 * 1000;
      const actualMs = expirationDate.getTime();
      expect(Math.abs(actualMs - expectedMs)).toBeLessThan(1000);
    });

    it("should handle empty body lines array", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "24";
      const bodyLines = [];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      expect(bodyLines.length).toBe(1);
      expect(bodyLines[0]).toContain("<!-- gh-aw-expires:");
    });

    it("should use ISO format for expiration date in comment", () => {
      process.env.GH_AW_DISCUSSION_EXPIRES = "24";
      const bodyLines = [];

      addExpirationComment(bodyLines, "GH_AW_DISCUSSION_EXPIRES", "Discussion");

      const comment = bodyLines[0];
      const match = comment.match(/<!-- gh-aw-expires: (.+?) -->/);
      expect(match).not.toBeNull();

      // Verify it's a valid ISO 8601 date
      const isoDate = match[1];
      expect(isoDate).toMatch(/^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{3}Z$/);
    });
  });
});
