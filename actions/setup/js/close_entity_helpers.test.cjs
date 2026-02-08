import { describe, it, expect, beforeEach, vi } from "vitest";
const mockCore = { debug: vi.fn(), info: vi.fn(), warning: vi.fn(), error: vi.fn(), setFailed: vi.fn(), setOutput: vi.fn(), summary: { addRaw: vi.fn().mockReturnThis(), write: vi.fn().mockResolvedValue() } },
  mockContext = { eventName: "issues", runId: 12345, repo: { owner: "testowner", repo: "testrepo" }, payload: { issue: { number: 42 }, pull_request: { number: 100 }, repository: { html_url: "https://github.com/testowner/testrepo" } } },
  mockGithub = {};
((global.core = mockCore), (global.context = mockContext), (global.github = mockGithub));
const { checkLabelFilter, checkTitlePrefixFilter, parseEntityConfig, resolveEntityNumber, escapeMarkdownTitle, createCloseEntityHandler, ISSUE_CONFIG, PULL_REQUEST_CONFIG } = require("./close_entity_helpers.cjs");
describe("close_entity_helpers", () => {
  (beforeEach(() => {
    (vi.clearAllMocks(),
      delete process.env.GH_AW_CLOSE_ISSUE_REQUIRED_LABELS,
      delete process.env.GH_AW_CLOSE_ISSUE_REQUIRED_TITLE_PREFIX,
      delete process.env.GH_AW_CLOSE_ISSUE_TARGET,
      delete process.env.GH_AW_CLOSE_PR_REQUIRED_LABELS,
      delete process.env.GH_AW_CLOSE_PR_REQUIRED_TITLE_PREFIX,
      delete process.env.GH_AW_CLOSE_PR_TARGET,
      (global.context.eventName = "issues"),
      (global.context.payload.issue = { number: 42 }),
      (global.context.payload.pull_request = { number: 100 }));
  }),
    describe("checkLabelFilter", () => {
      (it("should return true when no required labels specified", () => {
        expect(checkLabelFilter([{ name: "bug" }], [])).toBe(!0);
      }),
        it("should return true when entity has one of the required labels", () => {
          expect(checkLabelFilter([{ name: "bug" }, { name: "enhancement" }], ["bug", "wontfix"])).toBe(!0);
        }),
        it("should return false when entity has none of the required labels", () => {
          expect(checkLabelFilter([{ name: "bug" }], ["enhancement", "wontfix"])).toBe(!1);
        }),
        it("should return false when entity has no labels and required labels specified", () => {
          expect(checkLabelFilter([], ["bug"])).toBe(!1);
        }));
    }),
    describe("checkTitlePrefixFilter", () => {
      (it("should return true when no required prefix specified", () => {
        expect(checkTitlePrefixFilter("Some Title", "")).toBe(!0);
      }),
        it("should return true when title starts with required prefix", () => {
          expect(checkTitlePrefixFilter("[bug] Fix something", "[bug]")).toBe(!0);
        }),
        it("should return false when title does not start with required prefix", () => {
          expect(checkTitlePrefixFilter("Fix something", "[bug]")).toBe(!1);
        }),
        it("should be case-sensitive", () => {
          expect(checkTitlePrefixFilter("[BUG] Fix something", "[bug]")).toBe(!1);
        }));
    }),
    describe("parseEntityConfig", () => {
      (it("should return defaults when no environment variables set", () => {
        const config = parseEntityConfig("GH_AW_CLOSE_ISSUE");
        (expect(config.requiredLabels).toEqual([]), expect(config.requiredTitlePrefix).toBe(""), expect(config.target).toBe("triggering"));
      }),
        it("should parse required labels from environment", () => {
          process.env.GH_AW_CLOSE_ISSUE_REQUIRED_LABELS = "bug, enhancement, stale";
          const config = parseEntityConfig("GH_AW_CLOSE_ISSUE");
          expect(config.requiredLabels).toEqual(["bug", "enhancement", "stale"]);
        }),
        it("should parse required title prefix from environment", () => {
          process.env.GH_AW_CLOSE_ISSUE_REQUIRED_TITLE_PREFIX = "[refactor]";
          const config = parseEntityConfig("GH_AW_CLOSE_ISSUE");
          expect(config.requiredTitlePrefix).toBe("[refactor]");
        }),
        it("should parse target from environment", () => {
          process.env.GH_AW_CLOSE_ISSUE_TARGET = "*";
          const config = parseEntityConfig("GH_AW_CLOSE_ISSUE");
          expect(config.target).toBe("*");
        }),
        it("should work with PR environment variable prefix", () => {
          ((process.env.GH_AW_CLOSE_PR_REQUIRED_LABELS = "ready-to-close"), (process.env.GH_AW_CLOSE_PR_TARGET = "123"));
          const config = parseEntityConfig("GH_AW_CLOSE_PR");
          (expect(config.requiredLabels).toEqual(["ready-to-close"]), expect(config.target).toBe("123"));
        }));
    }),
    describe("resolveEntityNumber", () => {
      (describe("with target '*'", () => {
        (it("should resolve from item number field", () => {
          const result = resolveEntityNumber(ISSUE_CONFIG, "*", { issue_number: 50 }, !0);
          (expect(result.success).toBe(!0), expect(result.number).toBe(50));
        }),
          it("should handle string number field", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "*", { issue_number: "75" }, !0);
            (expect(result.success).toBe(!0), expect(result.number).toBe(75));
          }),
          it("should fail when number field is missing", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "*", {}, !0);
            (expect(result.success).toBe(!1), expect(result.message).toContain("no issue_number specified"));
          }),
          it("should fail when number field is invalid", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "*", { issue_number: "abc" }, !0);
            (expect(result.success).toBe(!1), expect(result.message).toContain("Invalid issue number specified"));
          }),
          it("should fail when number is zero or negative", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "*", { issue_number: -5 }, !0);
            (expect(result.success).toBe(!1), expect(result.message).toContain("Invalid issue number specified"));
          }),
          it("should fail when number is zero (falsy)", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "*", { issue_number: 0 }, !0);
            (expect(result.success).toBe(!1), expect(result.message).toContain("no issue_number specified"));
          }));
      }),
        describe("with explicit target number", () => {
          (it("should resolve from target configuration", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "123", {}, !0);
            (expect(result.success).toBe(!0), expect(result.number).toBe(123));
          }),
            it("should fail when target is not a valid number", () => {
              const result = resolveEntityNumber(ISSUE_CONFIG, "invalid", {}, !0);
              (expect(result.success).toBe(!1), expect(result.message).toContain("Invalid issue number in target configuration"));
            }));
        }),
        describe("with target 'triggering'", () => {
          (it("should resolve from context in issue event", () => {
            const result = resolveEntityNumber(ISSUE_CONFIG, "triggering", {}, !0);
            (expect(result.success).toBe(!0), expect(result.number).toBe(42));
          }),
            it("should fail when not in entity context", () => {
              const result = resolveEntityNumber(ISSUE_CONFIG, "triggering", {}, !1);
              (expect(result.success).toBe(!1), expect(result.message).toContain("Not in issue context"));
            }),
            it("should fail when context payload has no number", () => {
              global.context.payload.issue = {};
              const result = resolveEntityNumber(ISSUE_CONFIG, "triggering", {}, !0);
              (expect(result.success).toBe(!1), expect(result.message).toContain("no issue found in payload"));
            }));
        }),
        describe("for pull requests", () => {
          (beforeEach(() => {
            global.context.eventName = "pull_request";
          }),
            it("should resolve PR number from item with target '*'", () => {
              const result = resolveEntityNumber(PULL_REQUEST_CONFIG, "*", { pull_request_number: 200 }, !0);
              (expect(result.success).toBe(!0), expect(result.number).toBe(200));
            }),
            it("should resolve PR number from triggering context", () => {
              const result = resolveEntityNumber(PULL_REQUEST_CONFIG, "triggering", {}, !0);
              (expect(result.success).toBe(!0), expect(result.number).toBe(100));
            }));
        }));
    }),
    describe("escapeMarkdownTitle", () => {
      (it("should escape square brackets", () => {
        expect(escapeMarkdownTitle("[feature] Add new thing")).toBe("\\[feature\\] Add new thing");
      }),
        it("should escape parentheses", () => {
          expect(escapeMarkdownTitle("Fix bug (urgent)")).toBe("Fix bug \\(urgent\\)");
        }),
        it("should escape all markdown special characters", () => {
          expect(escapeMarkdownTitle("[test] (foo) [bar]")).toBe("\\[test\\] \\(foo\\) \\[bar\\]");
        }),
        it("should not modify titles without special characters", () => {
          expect(escapeMarkdownTitle("Simple title")).toBe("Simple title");
        }));
    }),
    describe("ISSUE_CONFIG", () => {
      (it("should have correct entity type", () => {
        expect(ISSUE_CONFIG.entityType).toBe("issue");
      }),
        it("should have correct item type", () => {
          expect(ISSUE_CONFIG.itemType).toBe("close_issue");
        }),
        it("should have correct item type display", () => {
          expect(ISSUE_CONFIG.itemTypeDisplay).toBe("close-issue");
        }),
        it("should have correct context events", () => {
          (expect(ISSUE_CONFIG.contextEvents).toContain("issues"), expect(ISSUE_CONFIG.contextEvents).toContain("issue_comment"));
        }),
        it("should have correct URL path", () => {
          expect(ISSUE_CONFIG.urlPath).toBe("issues");
        }));
    }),
    describe("PULL_REQUEST_CONFIG", () => {
      (it("should have correct entity type", () => {
        expect(PULL_REQUEST_CONFIG.entityType).toBe("pull_request");
      }),
        it("should have correct item type", () => {
          expect(PULL_REQUEST_CONFIG.itemType).toBe("close_pull_request");
        }),
        it("should have correct item type display", () => {
          expect(PULL_REQUEST_CONFIG.itemTypeDisplay).toBe("close-pull-request");
        }),
        it("should have correct context events", () => {
          (expect(PULL_REQUEST_CONFIG.contextEvents).toContain("pull_request"), expect(PULL_REQUEST_CONFIG.contextEvents).toContain("pull_request_review_comment"));
        }),
        it("should have correct URL path", () => {
          expect(PULL_REQUEST_CONFIG.urlPath).toBe("pull");
        }));
    }),
    describe("createCloseEntityHandler", () => {
      (it("should create a handler function that respects max count", async () => {
        const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [], html_url: "http://test", state: "open" });
        const mockAddComment = vi.fn().mockResolvedValue({ id: 1, html_url: "http://comment" });
        const mockCloseEntity = vi.fn().mockResolvedValue({ number: 1, html_url: "http://test", title: "Test" });
        const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, { max: 1 });
        const result1 = await handler({ issue_number: 1 }, {});
        (expect(result1.success).toBe(true), expect(mockCloseEntity).toHaveBeenCalled());
        const result2 = await handler({ issue_number: 2 }, {});
        (expect(result2.success).toBe(false), expect(result2.error).toContain("Max count"), expect(mockCloseEntity).toHaveBeenCalledTimes(1));
      }),
        it("should handle missing entity number field", async () => {
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: vi.fn(), addComment: vi.fn(), closeEntity: vi.fn() }, {});
          (global.context.payload.issue = undefined);
          const result = await handler({}, {});
          (expect(result.success).toBe(false), expect(result.error).toContain("No issue number available"));
        }),
        it("should handle invalid entity number", async () => {
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: vi.fn(), addComment: vi.fn(), closeEntity: vi.fn() }, {});
          const result = await handler({ issue_number: "invalid" }, {});
          (expect(result.success).toBe(false), expect(result.error).toContain("Invalid issue number"));
        }),
        it("should handle already closed entity", async () => {
          const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [], html_url: "http://test", state: "closed" });
          const mockAddComment = vi.fn();
          const mockCloseEntity = vi.fn();
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, {});
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(true), expect(result.alreadyClosed).toBe(true), expect(mockCloseEntity).not.toHaveBeenCalled());
        }),
        it("should validate required labels", async () => {
          const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [{ name: "bug" }], html_url: "http://test", state: "open" });
          const mockAddComment = vi.fn();
          const mockCloseEntity = vi.fn();
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, { required_labels: ["enhancement"] });
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(false), expect(result.error).toContain("Missing required labels"), expect(mockCloseEntity).not.toHaveBeenCalled());
        }),
        it("should validate required title prefix", async () => {
          const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [], html_url: "http://test", state: "open" });
          const mockAddComment = vi.fn();
          const mockCloseEntity = vi.fn();
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, { required_title_prefix: "[bug]" });
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(false), expect(result.error).toContain(`Title doesn't start with "[bug]"`), expect(mockCloseEntity).not.toHaveBeenCalled());
        }),
        it("should add comment when configured", async () => {
          const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [], html_url: "http://test", state: "open" });
          const mockAddComment = vi.fn().mockResolvedValue({ id: 1, html_url: "http://comment" });
          const mockCloseEntity = vi.fn().mockResolvedValue({ number: 1, html_url: "http://test", title: "Test" });
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, { comment: "Closing this issue" });
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(true), expect(mockAddComment).toHaveBeenCalledWith(expect.any(Object), "testowner", "testrepo", 1, "Closing this issue"), expect(mockCloseEntity).toHaveBeenCalled());
        }),
        it("should close entity successfully", async () => {
          const mockGetDetails = vi.fn().mockResolvedValue({ number: 1, title: "Test", labels: [], html_url: "http://test", state: "open" });
          const mockAddComment = vi.fn().mockResolvedValue({ id: 1, html_url: "http://comment" });
          const mockCloseEntity = vi.fn().mockResolvedValue({ number: 1, html_url: "http://closed", title: "Test Closed" });
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, {});
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(true), expect(result.number).toBe(1), expect(result.url).toBe("http://closed"), expect(result.title).toBe("Test Closed"), expect(mockCloseEntity).toHaveBeenCalledWith(expect.any(Object), "testowner", "testrepo", 1));
        }),
        it("should handle errors gracefully", async () => {
          const mockGetDetails = vi.fn().mockRejectedValue(new Error("API Error"));
          const mockAddComment = vi.fn();
          const mockCloseEntity = vi.fn();
          const handler = createCloseEntityHandler("issue", "issue_number", "issue", { getDetails: mockGetDetails, addComment: mockAddComment, closeEntity: mockCloseEntity }, {});
          const result = await handler({ issue_number: 1 }, {});
          (expect(result.success).toBe(false), expect(result.error).toContain("API Error"));
        }));
    }));
});
