# GitLab CI/CD Configuration Guide

## Overview

This document defines critical constraints and best practices for GitLab CI/CD configurations. These rules prevent common YAML parsing errors and ensure pipeline reliability.

**CRITICAL**: When writing `.gitlab-ci.yml` files or any GitLab CI/CD configurations, you MUST follow these constraints to avoid validation errors and pipeline failures.

---

## 1. Echo Statement Formatting

### The Problem
GitLab's YAML parser interprets colons (`:`) as key-value separators, causing validation errors when used in echo statements.

### The Rule
- ❌ **NEVER** use colons (`:`) to separate labels from variables in echo statements
- ✅ **ALWAYS** use double dashes (`--`) instead

### Examples

```yaml
# ❌ WRONG - GitLab parses colons as YAML syntax
- echo "Job: ${CI_JOB_NAME}"
- echo "Branch: ${CI_COMMIT_REF_NAME}"

# ✅ CORRECT - Use double dashes instead
- echo "Job--${CI_JOB_NAME}"
- echo "Branch--${CI_COMMIT_REF_NAME}"
```

---

## 2. Multi-line Shell Scripts

### The Rule
- Use YAML literal block scalar (`|`) for multi-line if/else statements
- Never use inline YAML syntax for shell conditionals

### Example

```yaml
script:
  - some_command || EXIT_CODE=$?
  - |
    if [ -n "${EXIT_CODE}" ]; then
      echo "Error occurred"
      exit ${EXIT_CODE}
    fi
```

**Why**: The pipe (`|`) character tells YAML to treat the following indented block as a literal multi-line string, preserving shell script formatting and preventing YAML parsing conflicts.

---

## 3. YAML Anchors for before_script

### The Rule
- Anchor must reference array directly, not a hash with `before_script` key
- Use `before_script: *anchor_name` not `<<: *anchor_name`

### Example

```yaml
# ✅ CORRECT
.common_before_script: &common_before_script
  - echo "Setup step 1"
  - echo "Setup step 2"

job:
  before_script: *common_before_script
```

**Why**: The anchor references the array of commands directly. Using `<<: *anchor_name` is for merging hash structures, not for referencing arrays.

---

## Why This Matters

**GitLab YAML Parser Behavior**: GitLab's YAML parser interprets colons as key-value separators, causing validation errors when they appear in strings without proper escaping. Using `--` instead of `:` prevents parsing issues while maintaining readability in pipeline logs.

**Common Symptoms of Violations**:
- Pipeline validation fails with "mapping values are not allowed here"
- Unexpected YAML structure errors
- before_script or script sections not executing as expected

---

## Enforcement

When delegating CI/CD work to agents or implementing pipelines:

1. **Pre-validation**: Review all `.gitlab-ci.yml` changes for these constraints
2. **Reject violations**: Any pipeline configuration using colons in echo statements must be corrected
3. **Test locally**: Use `gitlab-ci-lint` or GitLab's CI Lint tool before committing
4. **Include in agent briefings**: Reference this guide when delegating CI/CD tasks

---

## Related Documentation

- GitLab CI/CD YAML Syntax: https://docs.gitlab.com/ee/ci/yaml/
- GitLab CI/CD Variables: https://docs.gitlab.com/ee/ci/variables/
- YAML Anchors and Aliases: https://docs.gitlab.com/ee/ci/yaml/yaml_optimization.html

---

## Quick Reference Checklist

Before committing `.gitlab-ci.yml` changes:

- [ ] No colons in echo statements (use `--` instead)
- [ ] Multi-line shell scripts use `|` literal block scalar
- [ ] YAML anchors for before_script reference arrays directly
- [ ] Pipeline passes local validation (gitlab-ci-lint)
- [ ] All string interpolations properly formatted
