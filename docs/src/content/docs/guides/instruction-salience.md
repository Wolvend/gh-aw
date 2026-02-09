---
title: Instruction Salience
description: Learn how to identify and work around current AI model limitations that cause some instructions to be followed less reliably than others.
sidebar:
  order: 6
---

## The Ideal vs. Current Reality

**The Ideal State**: Ideally, AI agents should consider ALL instructions equally and completely, following every instruction with the same level of attention and compliance. Users should not need to emphasize, position, or format instructions specially.

**Current Reality**: In practice, current AI models exhibit varying compliance rates based on instruction characteristics such as position, formatting, and complexity. This is a model limitation, not intended behavior.

**This Guide's Purpose**: This guide helps you identify where instructions may be at risk of being skipped due to current model limitations, and provides workarounds to improve compliance until model behavior improves.

---

**Instruction salience** measures how noticeable, memorable, and "dominant" a particular instruction is to the AI model at the moment it decides what to do next. This is an observed characteristic of current model behavior, not a design goal.

This guide explains how to analyze instruction salience and apply workarounds using the [Instruction Salience Analyzer](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) workflow.

## Understanding Current Model Limitations

Current AI models, when processing workflows with multiple possible actions, exhibit these observable behaviors:

- **Simplicity bias**: "Create 4 issues" is more likely to be followed than "Create 4 issues, update Projects, set fields, post status update"
- **Recency bias**: Instructions near the end of prompts are followed more reliably than early instructions
- **Emphasis sensitivity**: Bold, emphasized instructions with visual markers are followed more reliably

These are limitations we work around, not features to embrace. As models improve, these workarounds should become less necessary.

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

**Result**: Due to current model limitations, the agent often skips the permission check and proceeds directly to creating issues. This is undesirable behavior we work around.

## The Salience Scoring Formula

The [Instruction Salience Analyzer](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) uses a quantitative formula to estimate the risk that an instruction will be skipped based on observed model behavior patterns:

```
Combined Score = (Position × 0.4) + (Emphasis × 0.3) + (Semantic × 0.3)
```

**Important**: This formula measures the likelihood that current models will follow an instruction, not how important the instruction should be. Low scores indicate instructions at risk of being skipped that may need workarounds.

### Position Score (0-10) - Weight: 40%

Due to recency effects in current models, instructions at the end of prompts are followed more reliably:

```
Position Score = 10 × (1 - line_number / total_lines)
```

- **End of prompt**: Score ≈ 9-10 (most reliably followed)
- **Middle**: Score ≈ 4-6
- **Beginning**: Score ≈ 0-2 (least reliably followed)

**Key observation**: End-of-prompt instructions have ~4x higher compliance rates in current models.

### Emphasis Score (0-10) - Weight: 30%

Visual markers and formatting improve compliance in current models:

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

Different sections of the prompt show varying compliance rates in current models:

| Section Type | Points | Examples |
|--------------|--------|----------|
| Runtime context | +5 | `${{ github.event.* }}` variables |
| Main workflow body | +6 | Primary instructions |
| Template prompts | +4 | Concrete examples |
| Custom instructions | +4 | AGENTS.md sections |
| Imported agents | +3 | `imports:` files |
| Tool configurations | +3 | `tools:` settings |

## The Compliance Hierarchy (Current Model Limitations)

Based on empirical analysis of current AI model behavior in gh-aw, instructions show varying compliance rates across five tiers. These tiers reflect observed model limitations, not intended behavior:

### Tier 1: Highest Compliance (9.0-10.0)

**Characteristics**: End-of-prompt, strong emphasis, runtime context

**Examples**:
- Final checklists before completion
- Runtime feedback from previous steps
- Structured output requirements
- Gated instructions ("Do not proceed until...")

**Observed compliance**: >95% - These instructions are usually followed

### Tier 2: Good Compliance (7.0-8.9)

**Characteristics**: Main workflow body, concrete templates, strong formatting

**Examples**:
- Primary task instructions with examples
- Copy-pastable templates
- Emphasized requirements (bold + emoji)
- Step-by-step procedures

**Observed compliance**: 80-95% - These instructions are usually followed

### Tier 3: Moderate Compliance (5.0-6.9)

**Characteristics**: Imported instructions, explicit MUST/REQUIRED language

**Examples**:
- Imported agent instructions
- Security requirements with MUST language
- Configuration guidelines
- Best practices sections

**Observed compliance**: 60-80% - These instructions are sometimes skipped

### Tier 4: Low Compliance (3.0-4.9)

**Characteristics**: Middle sections, weak emphasis, optional suggestions

**Examples**:
- General guidelines in AGENTS.md middle sections
- Tool configuration notes
- Suggested (not required) actions
- Background context

**Observed compliance**: 30-60% - These instructions are often skipped

### Tier 5: Very Low Compliance (0.0-2.9)

**Characteristics**: Early prompt sections, no emphasis, buried instructions

**Examples**:
- Early AGENTS.md sections
- Optional suggestions ("you should consider...")
- Informational text
- Multiple alternative approaches

**Observed compliance**: <30% - These instructions are rarely followed

**Important Reminder**: All instructions SHOULD be followed equally. These tiers describe current model limitations we work around, not desired behavior.

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

## Workarounds for Current Model Limitations

The following techniques help work around observed compliance issues in current models. As models improve, these workarounds should become less necessary.

**Remember**: The goal is not to make instructions "high salience" - the goal is to ensure ALL instructions are followed. These techniques are temporary workarounds for current model limitations.

### 1. Position Optimization

**Workaround**: Move at-risk instructions closer to where they'll be executed:

❌ **Problematic**: Critical instruction far from action it guards
```markdown
Line 50: Before proceeding, validate user permissions.
...
Line 350: Now create the issues...
```

✅ **Workaround**: Place instruction at decision point
```markdown
Line 340: ## Pre-Creation Validation
🚨 **STOP**: Validate user permissions before creating issues.

Line 350: Now create the issues...
```

### 2. Emphasis Enhancement

**Workaround**: Add visual markers to improve compliance:

❌ **Problematic**: Plain text often skipped
```markdown
You must validate all inputs before processing.
```

✅ **Workaround**: Multiple emphasis markers
```markdown
## ⚠️ INPUT VALIDATION REQUIRED

**🚨 CRITICAL**: You **MUST** validate ALL inputs:

- ✅ Check data types
- ✅ Sanitize user input
- ✅ Verify ranges

❌ **NEVER** process unvalidated input.
```

### 3. Gating Language

**Workaround**: Force execution order with explicit gates:

❌ **Problematic**: Suggestion easily ignored
```markdown
You should verify the token has correct permissions.
```

✅ **Workaround**: Gating with clear checkpoint
```markdown
## Checkpoint 1: Token Validation

**Do not proceed to Step 2 until token validation completes:**

1. Verify token scopes include `repo:write`
2. Test token with read-only operation
3. Confirm success before proceeding

**Verification complete?** ✅ Proceed to Step 2.
```

### 4. Success Criteria Checklists

**Workaround**: Define completion requirements explicitly:

❌ **Problematic**: Vague criteria easily skipped
```markdown
After creating the PR, make sure everything is set up correctly.
```

✅ **Workaround**: Explicit checklist
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

**Workaround**: Provide copy-pastable examples:

❌ **Problematic**: Abstract description
```markdown
Create a detailed issue with all the necessary information.
```

✅ **Workaround**: Concrete template
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

**Workaround**: Use consistent terminology:

❌ **Problematic**: Multiple terms cause confusion
```markdown
Line 100: Update the project board...
Line 200: Modify the project...
Line 300: Edit the Projects v2 item...
```

✅ **Workaround**: Consistent terminology
```markdown
Line 100: Update the project (using `update_project` tool)...
Line 200: Update the project status field...
Line 300: Update the project assignee...
```

## Common Compliance Issues (Model Limitations)

These patterns describe common ways instructions get skipped due to current model limitations:

### 1. The Buried Security Check

**Issue**: Critical security checks appear early with no emphasis and are often skipped

**Workaround**: Move to decision point with strong emphasis and gating

### 2. The Optional Requirement

**Issue**: Required actions phrased as suggestions ("you should consider...") are often ignored

**Workaround**: Use explicit language (MUST, REQUIRED, DO NOT PROCEED)

### 3. The Competing Simple Path

**Issue**: Simple action (create issue) appears before complex required steps and gets executed first

**Workaround**: Use checklists to enforce order, put simple actions last

### 4. The Ambiguous Function Name

**Issue**: Multiple ways to do the same thing (create-issue vs create_issue vs createIssue) cause confusion

**Workaround**: Pick one convention and use consistently

### 5. The Long Preamble

**Issue**: 200 lines of background before the actual task means early instructions are rarely followed

**Workaround**: Move context to appendix, start with the task

## Monitoring Compliance

After applying workarounds, monitor workflow runs to verify effectiveness:

1. **Check workflow outputs**: Are instructions being followed now?
2. **Review failure patterns**: Which instructions are still being skipped?
3. **Iterate on problem areas**: Apply additional workarounds where needed
4. **Re-run analysis**: Verify improvements with the analyzer

## Best Practices

1. **Run analysis early**: Identify at-risk instructions during workflow design, not after failures
2. **Focus on critical paths**: Prioritize workarounds for security, data integrity, user experience
3. **Test incrementally**: Apply 2-3 workarounds, test, then continue
4. **Target Tier 2+ (7.0+)**: For critical instructions, aim for scores that predict >80% compliance
5. **Balance readability**: Don't over-emphasize everything (emphasis inflation reduces effectiveness)
6. **Validate with real runs**: Scores predict behavior but aren't perfect - verify with actual workflow runs
7. **Remember the goal**: These are workarounds for current limitations, not permanent solutions

## Further Reading

- [Editing Workflows](/gh-aw/guides/editing-workflows/) - Understanding workflow structure
- [Instruction Salience Analyzer workflow](https://github.com/github/gh-aw/blob/main/.github/workflows/instruction-salience-analyzer.md) - Source code
- [Creating Workflows](/gh-aw/setup/creating-workflows/) - Workflow creation guide

## Quick Reference Card

| Score Range | Tier | Observed Compliance | Workaround Needed |
|-------------|------|-------------------|---------------|
| 9.0-10.0 | Tier 1 | >95% | ✅ Usually followed |
| 7.0-8.9 | Tier 2 | 80-95% | ✅ Acceptable |
| 5.0-6.9 | Tier 3 | 60-80% | ⚠️ Consider workarounds |
| 3.0-4.9 | Tier 4 | 30-60% | ⚠️ Workarounds recommended |
| 0.0-2.9 | Tier 5 | <30% | 🚨 Workarounds required |

**Target for critical instructions**: Score 7.0+ (predicts >80% compliance in current models)

**Remember**: These scores reflect current model limitations. The ideal state is 100% compliance for ALL instructions regardless of position, emphasis, or formatting.
