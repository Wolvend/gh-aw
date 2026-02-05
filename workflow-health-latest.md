# Workflow Health Dashboard - 2026-02-05T11:32:17Z

## Overview
- **Total workflows**: 147 (146 executable, 1 shared import)
- **Healthy**: 143 (97.9%)
- **Warning**: 0 (0%)
- **Critical**: 3 (2.1%)
- **Inactive**: 0 (0%)
- **Compilation coverage**: 146/146 (100% ✅)
- **Overall health score**: 75/100 (↓ -20 from 95/100)

## 🚨 STATUS: DEGRADED - ACTION REQUIRED

### Critical Issues 🚨

#### Issue #1: Missing Action Files Breaking Workflows (P1 - High)

**Affected Workflows**: 3 workflows failing
- Daily Fact About gh-aw
- Copilot PR Conversation NLP Analysis
- The Great Escapi

**Root Cause**: Missing JavaScript action files:
- `/opt/gh-aw/actions/parse_mcp_gateway_log.cjs`
- `/opt/gh-aw/actions/handle_agent_failure.cjs`

**Error Pattern**:
```
Error: Cannot find module '/opt/gh-aw/actions/parse_mcp_gateway_log.cjs'
```

**Impact**:
- Workflows using MCP gateway logging cannot complete
- Agent failure handling is broken
- 3 workflows consistently failing

**Action Taken**: Issue created for immediate investigation

**Priority**: P1 (Within 24 hours)

### Health Assessment

**Status: DEGRADED**

**Health Summary:**
- ✅ **100% compilation coverage** (146/146 workflows)
- ✅ **Zero outdated lock files** (sustained)
- ✅ **Zero missing lock files** (sustained)
- ⚠️ **3 workflow failures** in last 24 hours (NEW)
- ⚠️ **Health score**: 75/100 (↓ -20, downgraded to DEGRADED)

**Recent Activity (Last 7 Days):**
- Total runs: 27
- Success: 14 (51.9%)
- Failure: 3 (11.1%)
- Skipped: 10 (37.0%)
- Success rate: 82.4% (excluding skipped)

**Key Changes Since Last Run (2026-02-04):**
- ⚠️ 3 new failures introduced (missing action files)
- ↓ Health score decreased by -20 points (95 → 75)
- ⚠️ Success rate dropped from 100% to 82.4%
- ✅ Compilation coverage maintained at 100%

## Workflow Statistics

### Compilation Status
- **Total .md files**: 147 (146 executable + 1 shared import)
- **Total .lock.yml files**: 146 (1 shared import correctly excluded)
- **Missing lock files**: 0
- **Outdated lock files**: 0 ✅
- **Compilation success rate**: 100%

### Engine Distribution

**Workflow breakdown by AI engine:**
- **Copilot**: 69 workflows (47.3%)
- **Claude**: 29 workflows (19.9%)
- **Codex**: 9 workflows (6.2%)
- **Unknown/No engine**: 39 workflows (26.7%)

### Safe Outputs Usage

- **136 workflows** (93.2%) have safe-outputs configured
- **9 workflows** (6.2%) do not use safe-outputs
- **High adoption rate** indicates excellent security practices

### Workflow Categories

- **Regular workflows**: 146 (100%)
- **Campaign orchestrators**: 0
- **Campaign specs**: 0
- **Shared imports**: 1 (intentionally not compiled)

## Recent Workflow Activity

### Most Active Workflows (Last 7 Days)

1. **Issue Monster** - 3 runs (all success)
2. **Agentic Maintenance** - 2 runs (all success)
3. **Daily Workflow Updater** - 1 run (success)
4. **Daily Code Metrics** - 1 run (success)
5. **Typist** - 1 run (success)

### Failed Workflows (Last 24 Hours)

1. **Daily Fact About gh-aw** - 1 failure (missing action files)
2. **Copilot PR Conversation NLP Analysis** - 1 failure (agent execution)
3. **The Great Escapi** - 1 failure (agent execution)

### Conclusion Breakdown

- **success**: 14 runs (51.9%)
- **failure**: 3 runs (11.1%)
- **skipped**: 10 runs (37.0%)
- **action_required**: 0 runs (0%)

## Error Analysis

### Error Pattern: Missing Action Files

**Frequency**: 2 workflows affected directly
**Severity**: High (P1)
**First Seen**: 2026-02-05

**Error Message**:
```
Error: Cannot find module '/opt/gh-aw/actions/parse_mcp_gateway_log.cjs'
Error: Cannot find module '/opt/gh-aw/actions/handle_agent_failure.cjs'
```

**Root Cause**:
- Action files not being copied during setup
- Possible recent changes to action file structure
- Missing files in actions/setup/js/ directory

**Affected Workflows**:
1. Daily Fact About gh-aw (run 21709325824)
2. Copilot PR Conversation NLP Analysis (run 21707633912) - secondary impact
3. The Great Escapi (run 21705733653) - secondary impact

**Resolution Steps**:
1. Verify action files exist in source
2. Update actions/setup to copy all required files
3. Add validation to ensure files are present
4. Test fix on failing workflows

## Trends

- Overall health score: 75/100 (↓ -20 from last run)
- Compilation coverage: 100% (sustained)
- Recent failure rate: 11.1% (3/27 runs - ↑ from 0%)
- Safe outputs adoption: 93.2% (stable)
- Outdated lock files: 0 (sustained)

**Health Trend**: ↓ **DEGRADED** (95/100 → 75/100)

## Actions Taken This Run

- ✅ Analyzed 146 executable workflows
- ✅ Verified 100% compilation coverage
- ⚠️ Detected 3 new failures (missing action files)
- ✅ Created issue for missing action files (P1)
- ✅ Analyzed error logs for root cause
- ✅ Updated health score: 75/100 (↓ -20)
- ✅ Documented error patterns and resolution steps

## Recommendations

### High Priority (P1 - Within 24 hours)
1. **Fix missing action files** (Issue created)
   - Verify parse_mcp_gateway_log.cjs exists
   - Verify handle_agent_failure.cjs exists
   - Update actions/setup to copy all files
   - Test fix on 3 failing workflows

### Medium Priority (P2 - This week)
1. Add validation to ensure all action files are present
2. Improve error messages when action files are missing
3. Add health checks for critical action file availability
4. Document action file dependencies

### Low Priority (P3 - Next sprint)
1. Add automated testing for action file availability
2. Monitor safe outputs adoption for remaining 6.8% of workflows
3. Continue tracking workflow run success rates

## System Status Summary

### ⚠️ Degraded Health - Action Required

**Infrastructure Health:**
- Compilation: 100% ✅
- Execution: 82.4% success rate ⚠️ (down from 100%)
- Safe outputs: 93.2% adoption ✅
- Lock files: 100% up-to-date ✅

**Quality Metrics:**
- 3 failures in last 24 hours ⚠️ (new)
- Zero timeout issues ✅
- Zero permission errors ✅
- Zero missing lock files ✅
- Zero outdated lock files ✅

**Operational Status:**
- **NEW ISSUE**: Missing action files (P1)
- Health score at 75/100 (degraded) ⚠️
- 3 workflows need immediate attention ⚠️
- 143 workflows operating normally ✅

## Resolution Plan

### Immediate Actions (Next 24 Hours)

1. **Investigate missing files** (ETA: 2 hours)
   - Check if files exist in repo
   - Review recent commits
   - Identify when files were removed/renamed

2. **Fix actions/setup** (ETA: 2 hours)
   - Update file copy logic
   - Add validation checks
   - Test on failing workflows

3. **Verify fix** (ETA: 1 hour)
   - Manually trigger 3 failing workflows
   - Confirm successful execution
   - Monitor for new failures

### Expected Outcome

**Target Health Score**: 95/100 (restored to excellent)  
**Target Success Rate**: 100% (no failures)  
**Target Timeline**: Within 24 hours  

---
> **Last updated**: 2026-02-05T11:32:17Z  
> **Next check**: 2026-02-06 (daily schedule)  
> **Health Trend**: ↓ Degraded (95/100 → 75/100)  
> **Status**: 🚨 **DEGRADED - ACTION REQUIRED**  
> **Priority Action**: Fix missing action files (P1)
