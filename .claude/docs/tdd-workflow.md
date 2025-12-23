# Test-Driven Development (TDD) Workflow - Comprehensive Reference

This document serves as the complete, authoritative reference for Test-Driven Development (TDD) practices in The Garden Claude Agent System. All agents performing code implementation, testing, or bug fixes MUST follow these protocols.

**NO EXCEPTIONS** - TDD is MANDATORY for ALL code implementation tasks unless explicitly exempted by the main agent.

---

## TDD Requirements Template

This template MUST be included in ALL agent briefings for code implementation, testing, or bug fix tasks.

### Template for Agent Briefings

```markdown
## TDD Requirements (MANDATORY)

You MUST follow Test-Driven Development workflow:

### Step 1: Write Failing Test First
- **Test File**: `tests/[unit|integration]/test_[feature_name].py`
- **Run Command**: `pytest tests/[path]/test_[feature_name].py -v -s`
- **Expected Result**: Test FAILS (reproducing the bug or showing feature not implemented)
- **Report Back**: Paste test failure output showing the exact error

### Step 2: Implement Minimal Code
- **File to Modify**: `src/[module]/[file].py`
- **Approach**: Write smallest code change to make test pass
- **No Gold Plating**: Only implement what's needed for test to pass

### Step 3: Verify Unit Tests Pass
- **Run Command**: `pytest tests/[path]/test_[feature_name].py -v --cov=src/[module]`
- **Expected Result**: All tests PASS ✅
- **Coverage Target**: >80% for modified code
- **Report Back**: Paste test success output with coverage percentage

### Step 4: Run Integration/E2E Test
- **Command**: [Specify exact API call, demo script, or CLI command]
- **Expected Behavior**: [Describe specific success criteria]
- **Verification**: [What to check in response/output]
- **Report Back**: Paste actual command output showing success

### Step 5: Commit Only If All Tests Pass
- **Pre-Commit Check**: All unit tests ✅, integration test ✅
- **Commit Command**: `git add tests/ src/ && git commit -m "[descriptive message]"`
- **Report Back**: Commit hash and verification that tests are included

### Step 6: Report TDD Completion

Provide summary:
- [ ] Test file created: [path]
- [ ] Initial test run: FAIL (expected)
- [ ] Code implemented: [files modified]
- [ ] Unit tests: PASS ✅ (coverage: X%)
- [ ] Integration test: PASS ✅
- [ ] Committed: [hash]

**DO NOT REPORT TASK COMPLETE UNTIL ALL 6 STEPS DONE WITH EVIDENCE**
```

---

## TDD Workflow (Required Sequence)

Every specialized agent performing code implementation MUST follow this exact sequence. This is the core TDD cycle that ensures quality, maintainability, and correctness.

### 1. Write Failing Test First

**Purpose**: Establish clear success criteria before writing any implementation code.

**Requirements**:
- Create test file BEFORE any implementation code
- Test must reproduce the requirement or bug
- Run test, verify it fails with clear error message
- Report test failure output as evidence

**Why This Matters**:
- Proves the test can detect the issue
- Defines expected behavior upfront
- Prevents false positives (tests that always pass)
- Documents requirements as executable specifications

**Example Output to Report**:
```
FAILED tests/unit/test_authentication.py::test_login_with_valid_credentials - AssertionError: Expected user object, got None
```

### 2. Implement Minimal Code

**Purpose**: Write only the code necessary to make the test pass.

**Requirements**:
- Write smallest amount of code to make test pass
- No premature optimization
- No "extra features" beyond requirement
- Focus on making the test green

**Why This Matters**:
- Prevents over-engineering
- Keeps codebase simple and maintainable
- Reduces scope creep
- Makes code easier to understand and debug

**Anti-Pattern to Avoid**:
```python
# ❌ WRONG - Adding features not required by test
def login(username, password):
    user = authenticate(username, password)
    log_audit_trail(user)  # Not required by current test
    send_notification(user)  # Not required by current test
    return user

# ✅ CORRECT - Only what's needed for test
def login(username, password):
    return authenticate(username, password)
```

### 3. Verify Test Passes

**Purpose**: Confirm the implementation satisfies the requirements.

**Requirements**:
- Run test suite, confirm all tests pass
- Fix any failures before proceeding
- No "it should work" claims without test evidence
- Report test success output as evidence

**Why This Matters**:
- Provides immediate feedback on correctness
- Catches integration issues with existing code
- Establishes baseline for refactoring
- Creates confidence in the implementation

**Example Output to Report**:
```
tests/unit/test_authentication.py::test_login_with_valid_credentials PASSED [100%]

---------- coverage: platform darwin, python 3.11.5 -----------
Name                     Stmts   Miss  Cover
--------------------------------------------
src/auth/login.py           15      2    87%
```

### 4. Refactor If Needed

**Purpose**: Improve code quality while maintaining correctness.

**Requirements**:
- Clean up code while keeping tests green
- Improve readability, remove duplication
- Re-run tests after each refactor
- Ensure tests still pass after refactoring

**Why This Matters**:
- Maintains code quality over time
- Makes future changes easier
- Reduces technical debt
- Improves team velocity

**Refactoring Safety Net**:
- Tests must pass BEFORE refactoring starts
- Run tests after EACH refactoring step
- If tests fail, revert the refactoring
- Only commit when tests are green again

### 5. Run Integration/E2E Tests

**Purpose**: Verify the feature works in the actual system context.

**Requirements**:
- Verify feature works in actual system
- Test via API call, CLI command, or demo script
- Confirm user-facing behavior is correct
- Report integration test success

**Why This Matters**:
- Unit tests verify components in isolation
- Integration tests verify components work together
- Catches configuration issues, environment problems, deployment issues
- Validates the actual user experience

**Example Output to Report**:
```bash
$ curl -X POST http://localhost:8000/api/login \
  -d '{"username":"test","password":"pass123"}' \
  -H "Content-Type: application/json"

{"status":"success","user":{"id":1,"username":"test"},"token":"eyJ0eXAiOiJ..."}

✅ Integration test PASSED - User authenticated successfully
```

### 6. Commit Only If All Pass

**Purpose**: Ensure repository always contains working, tested code.

**Requirements**:
- Unit tests: PASS ✅
- Integration tests: PASS ✅
- E2E verification: PASS ✅
- Only then execute: `git commit`

**Why This Matters**:
- Maintains clean git history
- Prevents broken commits
- Enables safe rollbacks
- Supports continuous integration

**Commit Message Format**:
```bash
git add tests/ src/ && git commit -m "feat: add user authentication with password validation

- Implemented login() function with credential validation
- Added unit tests with 87% coverage
- Verified via integration test (API endpoint working)

Tests: tests/unit/test_authentication.py
Coverage: 87% (13/15 lines)"
```

---

## Enforcement Rules

### Main Agent Responsibilities

The main Claude agent orchestrating work MUST:

1. **Include TDD workflow in every agent briefing for code tasks**
   - Use the TDD Requirements Template above
   - Specify exact test file paths and commands
   - Define clear success criteria for each step

2. **Require agents to report test results before claiming success**
   - Demand test failure output from Step 1
   - Demand test success output from Step 3
   - Demand integration test output from Step 5

3. **Reject agent deliverables that skip TDD steps**
   - If agent reports "task complete" without test evidence → REJECT
   - If agent skips Step 1 (failing test first) → REJECT
   - If agent commits without all tests passing → REJECT

4. **Never accept "this should work" without test evidence**
   - "I implemented the feature" → INSUFFICIENT
   - "I ran the tests and they passed" → INSUFFICIENT
   - "Here's the test output showing PASS" → ACCEPTABLE

5. **Verify test files exist and are committed with code changes**
   - Check git commit includes both `src/` and `tests/` files
   - Verify test files follow naming conventions
   - Confirm test coverage meets 80% threshold

### Specialized Agent Responsibilities

Agents performing code implementation MUST:

1. **Report each TDD step completion with evidence**
   - Step 1: Paste test failure output
   - Step 3: Paste test success output with coverage
   - Step 5: Paste integration test command and output
   - Step 6: Provide commit hash and file list

2. **Include test file paths in deliverables**
   - Absolute paths to test files created
   - Absolute paths to implementation files modified
   - Test command used with full arguments

3. **Show test pass/fail status before claiming task complete**
   - Use pytest's verbose output (`-v` flag)
   - Include coverage report (`--cov` flag)
   - Show actual vs expected values in failures

4. **Create regression tests for bugs before fixing them**
   - Write test that reproduces the bug
   - Verify test fails (proves bug exists)
   - Fix the bug
   - Verify test passes (proves bug is fixed)

5. **Document test coverage percentage**
   - Report coverage for modified code only
   - Target: >80% line coverage
   - Identify any untested code paths
   - Explain why certain lines are uncovered (if applicable)

### Violation Consequences

When TDD requirements are not followed:

1. **Agent work is rejected and must be redone with proper TDD**
   - No partial credit for "mostly working" code
   - Agent must restart from Step 1
   - Main agent documents the violation

2. **Task is not marked complete until tests exist and pass**
   - Task remains "in_progress" status
   - Cannot move to next task
   - Blocks dependent work

3. **No commits are made until full TDD cycle completes**
   - All 6 steps must complete successfully
   - Test evidence must be provided
   - Integration test must pass

4. **Main agent must restart the task with TDD requirements**
   - Re-brief the agent with explicit TDD template
   - Monitor progress more closely
   - Escalate if repeated violations occur

---

## TDD Exemptions (Rare)

**Default assumption: TDD is REQUIRED unless explicitly exempted by main agent**

Only these scenarios can skip TDD (must be explicitly stated in task briefing):

### 1. Documentation-Only Changes
- **Examples**: README updates, markdown files, code comments
- **Rationale**: No executable logic to test
- **Requirement**: Must be purely documentation, no code changes

### 2. Configuration Changes with No Logic
- **Examples**: JSON files, YAML files, environment variable files
- **Rationale**: Configuration is validated at runtime, not via unit tests
- **Requirement**: No conditional logic, only data values

### 3. Experimental Spike/Prototype
- **Examples**: Research tasks, proof-of-concept implementations
- **Rationale**: Code is temporary and will be discarded
- **Requirement**: Must be marked as temporary AND include follow-up task for proper TDD implementation

### 4. Emergency Hotfix
- **Examples**: Production outage, critical security fix
- **Rationale**: Speed is critical, tests added after fix is deployed
- **Requirement**: Must include follow-up task to add tests within 24 hours

**Exemption Request Format**:
```
🔔 TDD EXEMPTION REQUESTED

**Task**: [Describe the task]
**Exemption Reason**: [One of the 4 scenarios above]
**Justification**: [Explain why TDD is not appropriate]
**Follow-up Plan**: [If applicable, describe when tests will be added]

Requesting explicit approval to proceed without TDD.
```

---

## Test Quality Standards

### Unit Tests

Unit tests verify individual functions, methods, or classes in isolation.

**Requirements**:
- **Minimum 80% code coverage for new code**
- **Test edge cases**: null, empty strings, invalid inputs, boundary values
- **Test error conditions**: exceptions, error returns, validation failures
- **Use descriptive test names**: `test_feature_scenario_expectedResult`
- **Use fixtures for test data setup**: DRY principle, consistent test data

**Example Quality Unit Test**:
```python
import pytest
from src.auth.login import authenticate

@pytest.fixture
def valid_user():
    return {"username": "testuser", "password": "pass123"}

@pytest.fixture
def invalid_user():
    return {"username": "baduser", "password": "wrongpass"}

def test_authenticate_with_valid_credentials_returns_user_object(valid_user):
    """Test that valid credentials return authenticated user object."""
    result = authenticate(valid_user["username"], valid_user["password"])
    assert result is not None
    assert result["username"] == valid_user["username"]
    assert "id" in result

def test_authenticate_with_invalid_password_returns_none(valid_user):
    """Test that invalid password returns None."""
    result = authenticate(valid_user["username"], "wrongpassword")
    assert result is None

def test_authenticate_with_nonexistent_user_returns_none():
    """Test that nonexistent user returns None."""
    result = authenticate("nonexistent", "anypassword")
    assert result is None

def test_authenticate_with_empty_username_raises_value_error():
    """Test that empty username raises ValueError."""
    with pytest.raises(ValueError, match="Username cannot be empty"):
        authenticate("", "password")

def test_authenticate_with_none_username_raises_type_error():
    """Test that None username raises TypeError."""
    with pytest.raises(TypeError, match="Username must be a string"):
        authenticate(None, "password")
```

**Coverage Target**:
- New code: >80% line coverage
- Modified code: >80% line coverage
- Critical paths: 100% coverage (authentication, payments, security)
- Error handling: 100% coverage (all exceptions tested)

### Integration Tests

Integration tests verify that multiple components work together correctly.

**Requirements**:
- **Test complete feature workflows end-to-end**
- **Use realistic data matching production scenarios**
- **Verify API responses, database state, file outputs**
- **Test error handling and rollback behavior**

**Example Quality Integration Test**:
```python
import pytest
import requests
from tests.integration.fixtures import test_database, test_server

def test_user_login_workflow_end_to_end(test_database, test_server):
    """Test complete user login workflow from API request to database."""
    # Step 1: Register new user
    register_response = requests.post(
        f"{test_server}/api/register",
        json={"username": "newuser", "password": "secure123", "email": "new@test.com"}
    )
    assert register_response.status_code == 201
    user_id = register_response.json()["user_id"]

    # Step 2: Verify user exists in database
    user = test_database.query("SELECT * FROM users WHERE id = ?", (user_id,))
    assert user is not None
    assert user["username"] == "newuser"
    assert user["email"] == "new@test.com"

    # Step 3: Login with valid credentials
    login_response = requests.post(
        f"{test_server}/api/login",
        json={"username": "newuser", "password": "secure123"}
    )
    assert login_response.status_code == 200
    assert "token" in login_response.json()
    token = login_response.json()["token"]

    # Step 4: Use token to access protected endpoint
    profile_response = requests.get(
        f"{test_server}/api/profile",
        headers={"Authorization": f"Bearer {token}"}
    )
    assert profile_response.status_code == 200
    assert profile_response.json()["username"] == "newuser"

    # Step 5: Verify login event logged in database
    login_event = test_database.query(
        "SELECT * FROM audit_log WHERE user_id = ? AND event = 'login'",
        (user_id,)
    )
    assert login_event is not None
```

**Integration Test Scope**:
- API endpoints (HTTP requests/responses)
- Database operations (CRUD, transactions)
- External service integration (mocked or test environments)
- File system operations (read/write/delete)
- Message queues (pub/sub, events)

### Regression Tests

Regression tests ensure that bugs, once fixed, never recur.

**Requirements**:
- **Every bug fix MUST include regression test**
- **Test reproduces the original bug before fix**
- **Test passes after fix is applied**
- **Test prevents bug from recurring**

**Regression Test Workflow**:
1. User reports bug: "Login fails when password contains special characters"
2. Write test that reproduces bug (test FAILS initially)
3. Fix the bug in implementation
4. Verify test now PASSES
5. Commit both test and fix together

**Example Regression Test**:
```python
def test_login_with_special_characters_in_password_succeeds():
    """
    Regression test for bug #123: Login failed with special chars in password.

    Bug Report: Users with passwords containing '@#$%' could not login.
    Root Cause: Password validation regex rejected special characters.
    Fix: Updated regex pattern to allow all printable ASCII characters.
    """
    # This test would have FAILED before the fix
    user = create_test_user(username="testuser", password="P@ssw0rd!#$%")
    result = authenticate("testuser", "P@ssw0rd!#$%")

    # After fix, this test PASSES
    assert result is not None
    assert result["username"] == "testuser"
```

**Regression Test Documentation**:
- Reference bug report number in docstring
- Describe original bug behavior
- Explain root cause
- Document the fix applied
- Ensure test would fail without the fix

---

## TDD Success Metrics

Agents must track and report these metrics in session history files:

### 1. Test Coverage Percentage
- **Target**: >80% for modified code
- **Measurement**: Use `pytest --cov=src/module --cov-report=term-missing`
- **Report Format**: "Unit test coverage: 87% (13/15 lines)"

### 2. Number of Tests Written
- **Count**: Total test functions created
- **Breakdown**: Unit tests vs integration tests vs regression tests
- **Report Format**: "Tests created: 8 unit, 2 integration, 1 regression"

### 3. Test Execution Time
- **Measurement**: Time to run full test suite
- **Target**: <5 seconds for unit tests, <30 seconds for integration tests
- **Report Format**: "Test execution time: unit=2.3s, integration=18.7s"

### 4. Integration Test Pass/Fail Status
- **Status**: All integration tests must pass before commit
- **Report Format**: "Integration tests: 3/3 PASS ✅"

### 5. TDD Workflow Adherence
- **Checklist**: All 6 steps completed with evidence
- **Report Format**:
  ```
  TDD Workflow Adherence:
  ✅ Step 1: Failing test written (evidence provided)
  ✅ Step 2: Minimal implementation completed
  ✅ Step 3: Unit tests pass (coverage: 87%)
  ✅ Step 4: Code refactored (no test failures)
  ✅ Step 5: Integration test pass
  ✅ Step 6: Committed (hash: a1b2c3d)

  TDD Adherence: 100% (6/6 steps completed)
  ```

---

## Test Requirements (Mandatory Gate)

Every task involving code changes MUST satisfy these test requirements before marking complete.

This is a **quality gate** - no exceptions.

### Checklist for Task Completion

- [ ] **Unit tests exist** with >80% coverage of modified code
- [ ] **Unit tests pass** - run `pytest tests/unit/... -v` successfully
- [ ] **Integration test exists** demonstrating feature works end-to-end
- [ ] **Integration test passes** - verified via API call, demo script, or CLI
- [ ] **Test output documented** in agent session history or commit message
- [ ] **Test files committed** with code changes in same commit

**NO EXCEPTIONS** - Code without passing tests is incomplete and task cannot be marked done.

### Quality Gate Enforcement

**Main Agent Verification**:
- Check that test files exist before accepting agent deliverables
- Verify test output shows PASS status
- Confirm coverage meets 80% threshold
- Review integration test results

**Agent Proof Requirements**:
- Provide test execution output (not just "tests passed")
- Include coverage report with percentage
- Show integration test command and output
- List test files created/modified

**Task Status Management**:
- Tasks remain "in_progress" until all test requirements met
- Cannot move task to "completed" without test evidence
- Cannot start dependent tasks until tests pass
- Cannot merge pull request until all tests pass

**Commit Rejection Criteria**:
- Commits without test files are rejected
- Commits with failing tests are rejected
- Commits without coverage report are rejected
- Commits without integration test evidence are rejected

---

## Summary

Test-Driven Development is not optional in The Garden Claude Agent System. It is a **mandatory discipline** that ensures:

1. **Correctness**: Code works as specified (tests prove it)
2. **Maintainability**: Code can be safely refactored (tests catch regressions)
3. **Documentation**: Tests serve as executable specifications
4. **Confidence**: Changes don't break existing functionality
5. **Quality**: 80%+ coverage means most code paths are verified

**The TDD Workflow**:
1. Write failing test first
2. Implement minimal code
3. Verify test passes
4. Refactor if needed
5. Run integration/E2E test
6. Commit only if all pass

**Remember**: Tests are not a burden - they are a **safety net** that enables rapid, confident development.

**For Agents**: Follow the TDD workflow exactly. Report test evidence at each step. Never claim completion without passing tests.

**For Main Agent**: Include TDD requirements in every code task briefing. Reject deliverables without test evidence. Enforce the quality gate.

**NO EXCEPTIONS** - Test-Driven Development is how we build reliable, maintainable software in The Garden.
