# Tasks: PRD-001 Tool Use for File Operations

## Task Breakdown

### 1. Tool Definition Layer
- **Agent**: developer
- **Files**: `backend/internal/service/claude.go`

#### 1.1 Define Tool Structs
Create Go structs for tool definitions matching Claude API schema.

#### 1.2 Add Tool Definitions
Add `write_file`, `read_file`, `list_files` tool definitions to API request.

#### 1.3 Update System Prompt
Update system prompt to instruct Claude to use tools for file operations.

---

### 2. Response Parsing Layer
- **Agent**: developer
- **Files**: `backend/internal/service/claude.go`, `backend/internal/service/chat.go`

#### 2.1 Parse Tool Use Blocks
Handle `content_block` events with `type: tool_use` in streaming response.

#### 2.2 Extract Tool Inputs
Parse tool name and input parameters from tool_use blocks.

#### 2.3 Validate Tool Inputs
Validate required parameters before execution.

---

### 3. Tool Execution Layer
- **Agent**: developer
- **Files**: `backend/internal/service/chat.go`

#### 3.1 Implement write_file Handler
Execute file write using existing FileRepository.

#### 3.2 Implement read_file Handler
Read file content from FileRepository.

#### 3.3 Implement list_files Handler
List files in project using FileRepository.

#### 3.4 Error Handling
Capture errors and format for tool_result response.

---

### 4. Tool Result Feedback
- **Agent**: developer
- **Files**: `backend/internal/service/claude.go`

#### 4.1 Send Tool Result
After tool execution, send tool_result message back to Claude.

#### 4.2 Handle Multi-turn
Support Claude continuing after receiving tool results.

#### 4.3 Stream Final Response
Ensure final text response streams to WebSocket correctly.

---

### 5. Testing
- **Agent**: developer
- **Files**: `backend/internal/service/*_test.go`

#### 5.1 Unit Tests for Tool Parsing
Test tool_use block parsing with various inputs.

#### 5.2 Integration Tests
Test full flow: message → tool use → file created → response.

#### 5.3 Fallback Tests
Verify markdown parsing still works as fallback.

---

### 6. Documentation
- **Agent**: researcher
- **Files**: `README.md`, inline comments

#### 6.1 Update API Documentation
Document new tool-based file creation flow.

#### 6.2 Add Code Comments
Document tool handling logic inline.
