# Permissions Builder Migration Guide

## Overview

The `PermissionsBuilder` provides a fluent API for constructing GitHub Actions permissions, replacing the previous factory explosion pattern (23 constructors) with a composable builder pattern that scales to any permission combination.

## Why the Change?

### Before: Factory Explosion

The old approach required a separate factory function for each permission combination:

```go
// 23 separate constructor functions
NewPermissionsContentsRead()
NewPermissionsContentsReadIssuesWrite()
NewPermissionsContentsReadIssuesWritePRWrite()
NewPermissionsContentsReadIssuesWritePRWriteDiscussionsWrite()
// ... 19 more variants
```

**Problems:**
- Combinatorial explosion: N permissions → 2^N possible combinations
- Not scalable: Each new combination requires a new function
- Violates DRY: Repetitive code patterns
- Hard to discover: Users must know exact function names

### After: Builder Pattern

The new builder pattern provides a composable, fluent API:

```go
// One builder, infinite combinations
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithIssues(PermissionWrite).
    WithPullRequests(PermissionWrite).
    Build()
```

**Benefits:**
- ✅ Scales to any permission combination
- ✅ Type-safe and discoverable (IDE autocomplete)
- ✅ Follows established Go patterns (similar to CompilerOption)
- ✅ Backward compatible (old functions still work)

## Migration Examples

### Single Permission

**Before:**
```go
perms := NewPermissionsContentsRead()
```

**After:**
```go
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    Build()
```

### Two Permissions

**Before:**
```go
perms := NewPermissionsContentsReadIssuesWrite()
```

**After:**
```go
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithIssues(PermissionWrite).
    Build()
```

### Three Permissions

**Before:**
```go
perms := NewPermissionsContentsReadIssuesWritePRWrite()
```

**After:**
```go
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithIssues(PermissionWrite).
    WithPullRequests(PermissionWrite).
    Build()
```

### Four Permissions

**Before:**
```go
perms := NewPermissionsContentsReadIssuesWritePRWriteDiscussionsWrite()
```

**After:**
```go
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithIssues(PermissionWrite).
    WithPullRequests(PermissionWrite).
    WithDiscussions(PermissionWrite).
    Build()
```

### Complex Combinations

**Before:**
```go
// No factory function exists for this combination
// Had to use NewPermissionsFromMap()
perms := NewPermissionsFromMap(map[PermissionScope]PermissionLevel{
    PermissionActions:      PermissionWrite,
    PermissionContents:     PermissionWrite,
    PermissionIssues:       PermissionWrite,
    PermissionPullRequests: PermissionWrite,
    PermissionDiscussions:  PermissionWrite,
    PermissionSecurityEvents: PermissionWrite,
})
```

**After:**
```go
// Clean, readable, type-safe
perms := NewPermissionsBuilder().
    WithActions(PermissionWrite).
    WithContents(PermissionWrite).
    WithIssues(PermissionWrite).
    WithPullRequests(PermissionWrite).
    WithDiscussions(PermissionWrite).
    WithSecurityEvents(PermissionWrite).
    Build()
```

## Available Builder Methods

The builder provides methods for all GitHub Actions permission scopes:

| Method | Description |
|--------|-------------|
| `WithActions(level)` | actions permission |
| `WithAttestations(level)` | attestations permission |
| `WithChecks(level)` | checks permission |
| `WithContents(level)` | contents permission |
| `WithDeployments(level)` | deployments permission |
| `WithDiscussions(level)` | discussions permission |
| `WithIdToken(level)` | id-token permission |
| `WithIssues(level)` | issues permission |
| `WithMetadata(level)` | metadata permission |
| `WithModels(level)` | models permission |
| `WithPackages(level)` | packages permission |
| `WithPages(level)` | pages permission |
| `WithPullRequests(level)` | pull-requests permission |
| `WithRepositoryProjects(level)` | repository-projects permission |
| `WithOrganizationProjects(level)` | organization-projects permission |
| `WithSecurityEvents(level)` | security-events permission |
| `WithStatuses(level)` | statuses permission |

## Permission Levels

Each scope accepts one of three levels:

- `PermissionRead` - read access
- `PermissionWrite` - write access (implies read)
- `PermissionNone` - explicitly no access

## Advanced Usage

### Dynamic Permission Building

```go
builder := NewPermissionsBuilder().WithContents(PermissionRead)

if needsIssueAccess {
    builder = builder.WithIssues(PermissionWrite)
}

if needsPRAccess {
    builder = builder.WithPullRequests(PermissionWrite)
}

perms := builder.Build()
```

### Overwriting Permissions

Setting a permission twice overwrites the previous value:

```go
perms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithContents(PermissionWrite).  // Overwrites to write
    Build()
```

### Method Chaining

All builder methods return the builder for chaining:

```go
builder := NewPermissionsBuilder()
builder.WithContents(PermissionRead)  // Returns builder
builder.WithIssues(PermissionWrite)   // Returns builder
perms := builder.Build()
```

## Backward Compatibility

**All existing factory functions remain available** and are now implemented as wrappers around the builder:

```go
// These still work (marked as deprecated):
NewPermissionsContentsRead()
NewPermissionsContentsReadIssuesWrite()
NewPermissionsContentsReadIssuesWritePRWrite()
// ... all 23 functions
```

**No breaking changes** - existing code continues to work unchanged.

## When to Migrate?

- ✅ **New code**: Always use the builder pattern
- ⚠️ **Existing code**: Optional, but recommended for:
  - Complex permission combinations
  - Code that would benefit from improved readability
  - Areas under active development
- ❌ **Don't migrate**: Stable code that works fine (no rush)

## Testing

The builder pattern maintains full backward compatibility:

```go
// Old and new produce identical results
oldPerms := NewPermissionsContentsReadIssuesWrite()
newPerms := NewPermissionsBuilder().
    WithContents(PermissionRead).
    WithIssues(PermissionWrite).
    Build()

// Renders identical YAML
assert.YAMLEq(t, oldPerms.RenderToYAML(), newPerms.RenderToYAML())
```

## Summary

| Aspect | Before | After |
|--------|--------|-------|
| **Approach** | 23 factory functions | Fluent builder |
| **Scalability** | O(2^N) functions | O(1) builder |
| **Discovery** | Know exact function names | IDE autocomplete |
| **Flexibility** | Limited to predefined combos | Any combination |
| **Backward Compat** | N/A | ✅ 100% compatible |

The builder pattern provides a more scalable, maintainable, and user-friendly API while maintaining full backward compatibility with existing code.
