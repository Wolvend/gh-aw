---
title: Instruction Salience
description: Learn how to optimize your agentic workflow instructions for maximum model compliance by understanding and improving instruction salience.
sidebar:
  order: 6
---

**Instruction salience** is how noticeable, memorable, and "dominant" a particular instruction is to the AI model at the moment it decides what to do next. High-salience instructions are consistently followed; low-salience instructions tend to get skipped, especially when competing with simpler alternatives.

This guide explains how to analyze and improve instruction salience in your workflows using the [Instruction Salience Analyzer](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) workflow.

## Why Instruction Salience Matters

AI agents often have multiple plausible completion paths when processing workflows. When faced with choices, they tend to follow:

- **Simpler paths**: "Create 4 issues" is easier than "Create 4 issues, update Projects, set fields, post status update"
- **Prominent instructions**: Bold, emphasized, end-of-prompt instructions
- **Concrete templates**: Copy-pastable examples with clear structure

Low-salience instructions for complex paths tend to get skipped, even when they're critical to your workflow's success.

### Real-World Example

Consider this workflow instruction:

```markdown
Before making changes, check if the user has proper permissions.

## Step 1: Create Issues

Create 4 issues based on the analysis results...
```

The permission check (low salience):
- Appears early in a long prompt
- Has no emphasis markers
- Competes with the concrete "Create Issues" step

**Result**: The agent often skips the permission check and proceeds directly to creating issues.

## The Salience Scoring Formula

The [Instruction Salience Analyzer](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) uses a quantitative formula to score each instruction (0-10 scale):

```
Combined Score = (Position × 0.4) + (Emphasis × 0.3) + (Semantic × 0.3)
```

### Position Score (0-10) - Weight: 40%

Instructions at the end of prompts have higher salience due to recency effects:

```
Position Score = 10 × (1 - line_number / total_lines)
```

- **End of prompt**: Score ≈ 9-10 (most salient)
- **Middle**: Score ≈ 4-6
- **Beginning**: Score ≈ 0-2 (least salient)

**Key insight**: End-of-prompt instructions have 4x higher salience than beginning instructions.

### Emphasis Score (0-10) - Weight: 30%

Visual markers and formatting increase instruction noticeability:

| Marker | Points | Example |
|--------|--------|---------|
| Emoji markers | +2 per emoji (max 2) | 🚨 ⚠️ ✅ ❌ |
| Bold text | +1 | **CRITICAL** |
| All caps | +1 | MUST, REQUIRED, NEVER |
| Code blocks | +1 | \`\`\`yaml ... \`\`\` |
| XML/HTML tags | +2 | `<critical>`, `<important>` |
| Repetition | +2 | Repeated across sections |
| List formatting | +1 | - [ ] Checklist items |

### Semantic Score (0-10) - Weight: 30%

Different sections of the prompt have inherent salience levels:

| Section Type | Points | Examples |
|--------------|--------|----------|
| Runtime context | +5 | `${{ github.event.* }}` variables |
| Main workflow body | +6 | Primary instructions |
| Template prompts | +4 | Concrete examples |
| Custom instructions | +4 | AGENTS.md sections |
| Imported agents | +3 | `imports:` files |
| Tool configurations | +3 | `tools:` settings |

## The Salience Hierarchy

Based on empirical analysis of gh-aw agent behavior, instructions fall into five tiers:

### Tier 1: High Salience (9.0-10.0)

**Characteristics**: End-of-prompt, strong emphasis, runtime context

**Examples**:
- Final checklists before completion
- Runtime feedback from previous steps
- Structured output requirements
- Gated instructions ("Do not proceed until...")

**Model behavior**: Consistently followed (>95% compliance)

### Tier 2: Good Salience (7.0-8.9)

**Characteristics**: Main workflow body, concrete templates, strong formatting

**Examples**:
- Primary task instructions with examples
- Copy-pastable templates
- Emphasized requirements (bold + emoji)
- Step-by-step procedures

**Model behavior**: Usually followed (80-95% compliance)

### Tier 3: Medium Salience (5.0-6.9)

**Characteristics**: Imported instructions, explicit MUST/REQUIRED language

**Examples**:
- Imported agent instructions
- Security requirements with MUST language
- Configuration guidelines
- Best practices sections

**Model behavior**: Sometimes followed (60-80% compliance)

### Tier 4: Low Salience (3.0-4.9)

**Characteristics**: Middle sections, weak emphasis, optional suggestions

**Examples**:
- General guidelines in AGENTS.md middle sections
- Tool configuration notes
- Suggested (not required) actions
- Background context

**Model behavior**: Often skipped (30-60% compliance)

### Tier 5: Very Low Salience (0.0-2.9)

**Characteristics**: Early prompt sections, no emphasis, buried instructions

**Examples**:
- Early AGENTS.md sections
- Optional suggestions ("you should consider...")
- Informational text
- Multiple alternative approaches

**Model behavior**: Rarely followed (<30% compliance)

## Using the Instruction Salience Analyzer

### Triggering the Analysis

The analyzer can be triggered in two ways:

#### 1. Via Issue (Recommended)

Create an issue with "[Salience Analysis]" in the title:

```markdown
Title: [Salience Analysis] Review my-workflow instructions

Body:
I'd like to analyze the instruction salience for `.github/workflows/my-workflow.md`
```

#### 2. Via Manual Workflow Dispatch

Navigate to Actions → Instruction Salience Analyzer → Run workflow

- **workflow_path**: `.github/workflows/my-workflow.md`

### Understanding the Analysis Report

The analyzer creates a discussion with:

1. **Executive Summary**: Overall scores and issue counts
2. **Salience Distribution**: Instructions by tier (1-5)
3. **Detailed Scores**: Every instruction with Position/Emphasis/Semantic breakdown
4. **Critical Issues**: Low-salience instructions that should be high-salience
5. **Recommendations**: Prioritized improvements with before/after examples

### Example Analysis Output

```markdown
### Issue 1: Critical Security Check Has Low Salience

**Problem**: Security validation appears early with weak emphasis

**Current Instruction** (Line 45, Score: 2.8/10):
```
Before making changes, check if the user has proper permissions.
```

**Breakdown**:
- Position Score: 1.2/10 (line 45 of 412 = 10.9% position)
- Emphasis Score: 0.0/10 (no formatting)
- Semantic Score: 6.0/10 (main workflow body)

**Recommended Fix** (Projected Score: 8.8/10):
```markdown
## ⚠️ SECURITY VALIDATION - DO NOT PROCEED WITHOUT COMPLETING

**CRITICAL**: Before ANY repository modifications:

1. 🚨 **VERIFY USER PERMISSIONS**: Call `check_permissions` API
2. ✅ **CONFIRM**: User has `write` or `admin` access
3. ❌ **ABORT** if permissions check fails

**Do not proceed to the next section until verification succeeds.**
```

**Improvements**:
- Moved to decision point (before actions)
- Added emojis (⚠️, 🚨, ✅, ❌): +8 emphasis
- Bold and caps: +2 emphasis
- Gating language: "Do not proceed until"
- Numbered checklist
```

## Improving Instruction Salience

### 1. Position Optimization

**Move critical instructions to high-salience positions**:

❌ **Bad**: Critical instruction at line 50 of 400-line prompt
```markdown
Line 50: Before proceeding, validate user permissions.
...
Line 350: Now create the issues...
```

✅ **Good**: Critical instruction at decision point
```markdown
Line 340: ## Pre-Creation Validation
🚨 **STOP**: Validate user permissions before creating issues.

Line 350: Now create the issues...
```

### 2. Emphasis Enhancement

**Add visual markers to increase noticeability**:

❌ **Bad**: Plain text (Emphasis: 0/10)
```markdown
You must validate all inputs before processing.
```

✅ **Good**: Multiple markers (Emphasis: 10/10)
```markdown
## ⚠️ INPUT VALIDATION REQUIRED

**🚨 CRITICAL**: You **MUST** validate ALL inputs:

- ✅ Check data types
- ✅ Sanitize user input
- ✅ Verify ranges

❌ **NEVER** process unvalidated input.
```

### 3. Gating Language

**Force execution order with explicit gates**:

❌ **Bad**: Suggestion without enforcement
```markdown
You should verify the token has correct permissions.
```

✅ **Good**: Gating with clear checkpoint
```markdown
## Checkpoint 1: Token Validation

**Do not proceed to Step 2 until token validation completes:**

1. Verify token scopes include `repo:write`
2. Test token with read-only operation
3. Confirm success before proceeding

**Verification complete?** ✅ Proceed to Step 2.
```

### 4. Success Criteria Checklists

**Define completion requirements explicitly**:

❌ **Bad**: Vague completion criteria
```markdown
After creating the PR, make sure everything is set up correctly.
```

✅ **Good**: Explicit checklist
```markdown
## Success Criteria

A successful PR creation MUST complete ALL items:

- [ ] PR created with correct base branch
- [ ] Description includes issue reference (#123)
- [ ] Labels applied: `needs-review`, `feature`
- [ ] Reviewers assigned: @team/reviewers
- [ ] CI checks triggered successfully

**Incomplete runs will be flagged as failures.**
```

### 5. Concrete Templates

**Provide copy-pastable examples**:

❌ **Bad**: Abstract description
```markdown
Create a detailed issue with all the necessary information.
```

✅ **Good**: Concrete template
```markdown
Create an issue using this exact template:

\`\`\`markdown
## Summary
[One-line description]

## Details
[Detailed explanation]

## Acceptance Criteria
- [ ] Criterion 1
- [ ] Criterion 2

## Related
- Blocked by: #123
- Relates to: #456
\`\`\`

Copy this template and fill in the bracketed sections.
```

### 6. Consistency in Naming

**Use consistent terminology throughout**:

❌ **Bad**: Multiple terms for same concept
```markdown
Line 100: Update the project board...
Line 200: Modify the project...
Line 300: Edit the Projects v2 item...
```

✅ **Good**: Consistent terminology
```markdown
Line 100: Update the project (using `update_project` tool)...
Line 200: Update the project status field...
Line 300: Update the project assignee...
```

## Common Salience Anti-Patterns

### 1. The Buried Security Check

**Problem**: Critical security checks appear early with no emphasis

**Fix**: Move to decision point with strong emphasis and gating

### 2. The Optional Requirement

**Problem**: Required actions phrased as suggestions ("you should consider...")

**Fix**: Use explicit language (MUST, REQUIRED, DO NOT PROCEED)

### 3. The Competing Simple Path

**Problem**: Simple action (create issue) appears before complex required steps

**Fix**: Use checklists to enforce order, put simple actions last

### 4. The Ambiguous Function Name

**Problem**: Multiple ways to do the same thing (create-issue vs create_issue vs createIssue)

**Fix**: Pick one convention and use consistently

### 5. The Long Preamble

**Problem**: 200 lines of background before the actual task

**Fix**: Move context to appendix, start with the task

## Monitoring Compliance

After improving salience, monitor workflow runs:

1. **Check workflow outputs**: Are instructions being followed?
2. **Review failure patterns**: Which instructions are still skipped?
3. **Iterate on low performers**: Increase salience for problem areas
4. **Re-run analysis**: Verify improvements with the analyzer

## Best Practices

1. **Run analysis early**: Check salience during workflow design, not after failures
2. **Focus on critical paths**: Prioritize security, data integrity, user experience
3. **Test incrementally**: Improve 2-3 instructions, test, then continue
4. **Use the hierarchy**: Target Tier 2+ (7.0+) for all critical instructions
5. **Balance readability**: Don't over-emphasize everything (salience inflation)
6. **Validate with real runs**: Salience scores predict behavior but aren't perfect

## Further Reading

- [Editing Workflows](/gh-aw/guides/editing-workflows/) - Understanding workflow structure
- [Instruction Salience Analyzer workflow](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) - Source code
- [Creating Workflows](/gh-aw/setup/creating-workflows/) - Workflow creation guide

## Quick Reference Card

| Score Range | Tier | Typical Compliance | Action Needed |
|-------------|------|-------------------|---------------|
| 9.0-10.0 | Tier 1 | >95% | ✅ Optimal |
| 7.0-8.9 | Tier 2 | 80-95% | ✅ Good |
| 5.0-6.9 | Tier 3 | 60-80% | ⚠️ Review |
| 3.0-4.9 | Tier 4 | 30-60% | ⚠️ Improve |
| 0.0-2.9 | Tier 5 | <30% | 🚨 Critical |

**Target**: All critical instructions should score 7.0+ (Tier 2 or higher)
