# Task: [Task Name]

**TDD Workflow**: This template follows Test-Driven Development methodology. Use for code implementation tasks requiring high quality and test coverage.

---

## Task Overview

**Objective**: [Clear description of what needs to be implemented]

**Success Criteria**:
- Feature/fix works as specified
- Unit tests exist and pass (>80% coverage recommended)
- Integration test demonstrates feature works end-to-end
- All tests pass before commit

---

## TDD Workflow Checklist

**Complete in order - DO NOT skip steps:**

### Step 1: Write Failing Test ❌

- [ ] **Create test file**
  - File: `tests/[unit|integration]/test_[feature_name].py` (or appropriate path/language)
  - Purpose: [What this test validates]

- [ ] **Write test that reproduces requirement/bug**
  ```python
  # Example structure - adapt to your language/framework
  def test_[feature_name]():
      # Arrange: Set up test data
      # Act: Call the function/feature
      # Assert: Verify expected behavior
      assert expected == actual
  ```

- [ ] **Run test and verify it fails**
  - Command: `[test command, e.g., pytest tests/unit/test_feature.py -v]`
  - Expected: FAIL ❌ (test should fail before implementation)
  - Output:
    ```
    [Paste test failure output here]
    ```

**Why this step matters**: Ensures test actually validates the requirement and isn't a false positive.

---

### Step 2: Implement Minimal Code ✅

- [ ] **Write smallest code change to make test pass**
  - File(s) modified: `[path/to/implementation/file]`
  - Lines changed: [N additions, M deletions]
  - Approach: [Brief 1-2 sentence description]

- [ ] **Implementation notes**:
  ```
  [Key decisions, algorithms, or logic used]
  ```

**Why this step matters**: Prevents over-engineering and keeps code focused on actual requirements.

---

### Step 3: Verify Unit Tests Pass ✅

- [ ] **Run unit tests**
  - Command: `[test command]`
  - Expected: PASS ✅
  - Coverage: [% if available]
  - Output:
    ```
    [Paste test success output here]
    ```

- [ ] **All existing tests still pass**
  - Command: `[full test suite command]`
  - Result: [N passed, 0 failed]

**Why this step matters**: Confirms fix works and didn't break existing functionality.

---

### Step 4: Integration/E2E Verification ✅

- [ ] **Run integration or end-to-end test**
  - Test type: [API call | Demo script | CLI command | Manual test]
  - Command/Steps:
    ```bash
    [Specific command or step-by-step manual test]
    ```

- [ ] **Verify expected behavior**
  - Expected result: [What should happen]
  - Actual result: [What actually happened]
  - Status: PASS ✅
  - Output/Screenshot:
    ```
    [Paste actual output or attach screenshot]
    ```

**Why this step matters**: Proves feature works in actual system, not just in isolation.

---

### Step 5: Refactor (If Needed) 🔧

- [ ] **Code cleanup** (while keeping tests green)
  - Refactoring done: [None | Specific improvements made]
  - Tests still passing: ✅

**Why this step matters**: Maintains code quality without breaking functionality.

---

### Step 6: Commit Changes 📝

- [ ] **Pre-commit verification**
  - Unit tests: PASS ✅
  - Integration tests: PASS ✅
  - E2E verification: PASS ✅
  - Code review (if applicable): [Approved | Self-reviewed]

- [ ] **Create commit**
  - Commit hash: `[hash]`
  - Commit message:
    ```
    [Type]: [Clear description]

    - [What was changed]
    - [Why it was changed]
    - Tests: [unit|integration|e2e] passing
    ```

**Why this step matters**: Only commit working, tested code.

---

## Implementation Details

### Files Modified
- `[path/file1]` - [What changed]
- `[path/file2]` - [What changed]

### Files Created
- `tests/[type]/test_[name]` - [Test description]
- `[other new files]` - [Purpose]

### Dependencies Added/Changed
- [Package/library]: [Version] - [Why needed]

---

## Testing Summary

| Test Type | Status | Coverage | Notes |
|-----------|--------|----------|-------|
| Unit Tests | ✅ PASS | [%] | [N tests passing] |
| Integration Tests | ✅ PASS | - | [Description] |
| E2E Verification | ✅ PASS | - | [Description] |

---

## Notes & Learnings

**Challenges encountered**:
- [Challenge 1 and how it was resolved]

**Key decisions**:
- [Decision 1 and rationale]

**Future improvements**:
- [Potential enhancement 1]
- [Potential enhancement 2]

---

## Task Complete ✅

**Completion criteria met**:
- [ ] All TDD steps completed in order
- [ ] All tests passing
- [ ] Code committed with test evidence
- [ ] Documentation updated (if needed)

**Task status**: [In Progress | Blocked | Complete]
**Completed by**: [Name/Agent]
**Date**: [YYYY-MM-DD]
