import { describe, it, expect } from "vitest";
import { sanitizeForShell } from "./sanitize_shell.cjs";

describe("sanitizeForShell", () => {
  it("should pass through non-string values unchanged", () => {
    expect(sanitizeForShell(123)).toBe(123);
    expect(sanitizeForShell(null)).toBe(null);
    expect(sanitizeForShell(undefined)).toBe(undefined);
    expect(sanitizeForShell(true)).toBe(true);
  });

  it("should escape dollar signs, backticks, and backslashes", () => {
    expect(sanitizeForShell("test$var")).toBe("test\\$var");
    expect(sanitizeForShell("test`cmd`")).toBe("test\\`cmd\\`");
    expect(sanitizeForShell("test\\path")).toBe("test\\\\path");
  });

  it("should remove dangerous shell metacharacters", () => {
    expect(sanitizeForShell("test;cmd")).toBe("testcmd");
    expect(sanitizeForShell("test&&cmd")).toBe("testcmd");
    expect(sanitizeForShell("test|cmd")).toBe("testcmd");
    expect(sanitizeForShell("test<file")).toBe("testfile");
    expect(sanitizeForShell("test>file")).toBe("testfile");
    expect(sanitizeForShell("test(cmd)")).toBe("testcmd");
  });

  it("should limit output length to 1000 characters", () => {
    const longString = "a".repeat(1500);
    const result = sanitizeForShell(longString);
    expect(result.length).toBe(1000);
  });

  it("should handle combinations of dangerous characters", () => {
    expect(sanitizeForShell("test$var;cmd|other")).toBe("test\\$varcmdother");
    expect(sanitizeForShell("test`cmd`&&other")).toBe("test\\`cmd\\`other");
  });
});
