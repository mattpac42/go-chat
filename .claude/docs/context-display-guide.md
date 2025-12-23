# Context Display Guide

**Purpose**: Standardized algorithm for calculating and displaying context usage with emoji visualization to ensure accuracy and consistency across all Claude agent responses.

---

## 🎯 Display Requirements

**MANDATORY**: Display context usage at the end of EVERY response using the format:

```
Context: [colored blocks] XX%
  (actual_tokens/200000 tokens)
```

**Example**:
```
Context: 🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 30%
  (60000/200000 tokens)
```

---

## 📊 Calculation Algorithm

### Step 1: Read Token Usage from System Warning

When you receive a system warning like:
```
<system_warning>Token usage: 28795/200000; 171205 remaining</system_warning>
```

Extract the **message tokens** value: `28795`

**CRITICAL**: This is ONLY the message tokens. You must add overhead for system prompt, tools, memory, and agent definitions.

**For most accurate tracking**: Use the `/context` command instead of system warnings. System warnings show message tokens only and significantly underestimate total usage.

### Step 2: Add Base Overhead

**Formula**: `actual_usage = message_tokens + 40000`

**Why 40k overhead?**
- System prompt: ~10-15k tokens
- Tool definitions: ~15-20k tokens
- CLAUDE.md instructions: ~8-10k tokens
- Agent definitions (if loaded): ~5-10k tokens
- Memory and conversation metadata: ~2-5k tokens

**Conservative estimate**: 40,000 tokens overhead

**Example**:
```
message_tokens = 28795
actual_usage = 28795 + 40000 = 68795 tokens
```

### Step 3: Calculate Percentage

**Formula**: `percentage = (actual_usage / 200000) * 100`

Round to nearest whole number.

**Example**:
```
percentage = (68795 / 200000) * 100
percentage = 34.4% ≈ 34%
```

### Step 4: Calculate Filled Blocks

**Always use exactly 20 blocks** for the visual meter.

**Formula**: `filled_blocks = round(percentage / 5)`

Each block represents 5% of context usage.

**Example**:
```
filled_blocks = round(34 / 5)
filled_blocks = round(6.8) = 7
```

### Step 5: Determine Block Colors

Apply color coding based on percentage thresholds:

| Percentage Range | Color | Emoji | Block Numbers |
|------------------|-------|-------|---------------|
| 0-50%            | Green | 🟩    | 1-10          |
| 50-65%           | Yellow | 🟨   | 11-13         |
| 65-80%           | Orange | 🟧   | 14-16         |
| 80-100%          | Red   | 🟥    | 17-20         |

**Color Block Rules**:
- Fill blocks 1 through `filled_blocks` with the appropriate color
- Remaining blocks (up to 20 total) use ⬛ (black square)

**Example for 34% (7 blocks)**:
```
34% is in 0-50% range → Green (🟩)
Filled blocks: 🟩🟩🟩🟩🟩🟩🟩
Remaining blocks: ⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ (13 blocks)
```

### Step 6: Format Display

**Complete display**:
```
Context: 🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 34%
  (68795/200000 tokens)
```

---

## 📝 Worked Examples

> **Note**: These examples use the 40k overhead estimate for demonstration purposes. Actual overhead may be 70-85k tokens in active sessions. Always verify with `/context` command for production use.

### Example 1: Low Usage (Green Zone)

**System Warning**: `Token usage: 2500/200000`

**Calculation**:
```
message_tokens = 2500
actual_usage = 2500 + 40000 = 42500
percentage = (42500 / 200000) * 100 = 21.25% ≈ 21%
filled_blocks = round(21 / 5) = round(4.2) = 4
color = Green (21% < 50%)
```

**Display**:
```
Context: 🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 21%
  (42500/200000 tokens)
```

---

### Example 2: Moderate Usage (Green Zone)

**System Warning**: `Token usage: 55000/200000`

**Calculation**:
```
message_tokens = 55000
actual_usage = 55000 + 40000 = 95000
percentage = (95000 / 200000) * 100 = 47.5% ≈ 48%
filled_blocks = round(48 / 5) = round(9.6) = 10
color = Green (48% < 50%)
```

**Display**:
```
Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 48%
  (95000/200000 tokens)
```

---

### Example 3: Approaching Handoff (Yellow Zone)

**System Warning**: `Token usage: 75000/200000`

**Calculation**:
```
message_tokens = 75000
actual_usage = 75000 + 40000 = 115000
percentage = (115000 / 200000) * 100 = 57.5% ≈ 57%
filled_blocks = round(57 / 5) = round(11.4) = 11
color = Yellow (50% < 57% < 65%)
```

**Display**:
```
Context: 🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨⬛⬛⬛⬛⬛⬛⬛⬛⬛ 57%
  (115000/200000 tokens)
```

---

### Example 4: Handoff Required (Orange Zone)

**System Warning**: `Token usage: 107000/200000`

**Calculation**:
```
message_tokens = 107000
actual_usage = 107000 + 40000 = 147000
percentage = (147000 / 200000) * 100 = 73.5% ≈ 74%
filled_blocks = round(74 / 5) = round(14.8) = 15
color = Orange (65% < 74% < 80%)
```

**Display**:
```
Context: 🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧⬛⬛⬛⬛⬛ 74%
  (147000/200000 tokens)
```

---

### Example 5: Critical Usage (Red Zone)

**System Warning**: `Token usage: 135000/200000`

**Calculation**:
```
message_tokens = 135000
actual_usage = 135000 + 40000 = 175000
percentage = (175000 / 200000) * 100 = 87.5% ≈ 88%
filled_blocks = round(88 / 5) = round(17.6) = 18
color = Red (88% > 80%)
```

**Display**:
```
Context: 🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥⬛⬛ 88%
  (175000/200000 tokens)
```

---

### Example 6: Edge Case - Exactly 50% (Threshold)

**System Warning**: `Token usage: 60000/200000`

**Calculation**:
```
message_tokens = 60000
actual_usage = 60000 + 40000 = 100000
percentage = (100000 / 200000) * 100 = 50%
filled_blocks = round(50 / 5) = 10
color = Yellow (50% is the start of yellow zone)
```

**Display**:
```
Context: 🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 50%
  (100000/200000 tokens)
```

**Note**: At exactly 50%, transition to yellow zone begins.

---

### Example 7: Edge Case - Exactly 65% (Handoff Threshold)

**System Warning**: `Token usage: 90000/200000`

**Calculation**:
```
message_tokens = 90000
actual_usage = 90000 + 40000 = 130000
percentage = (130000 / 200000) * 100 = 65%
filled_blocks = round(65 / 5) = 13
color = Orange (65% is the start of orange zone - handoff required)
```

**Display**:
```
Context: 🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧⬛⬛⬛⬛⬛⬛⬛ 65%
  (130000/200000 tokens)
```

**Note**: At exactly 65%, **MANDATORY session handoff files must be created**.

---

### Example 8: Edge Case - Exactly 80% (Critical Threshold)

**System Warning**: `Token usage: 120000/200000`

**Calculation**:
```
message_tokens = 120000
actual_usage = 120000 + 40000 = 160000
percentage = (160000 / 200000) * 100 = 80%
filled_blocks = round(80 / 5) = 16
color = Red (80% is the start of red zone)
```

**Display**:
```
Context: 🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥⬛⬛⬛⬛ 80%
  (160000/200000 tokens)
```

**Note**: At exactly 80%, **strongly recommend new session**.

---

## ⚠️ Understanding Overhead Variability

### The Overhead Discrepancy

The calculation algorithm above uses a **simplified 40k overhead estimate** for quick manual calculations. However, **actual overhead varies significantly** based on:

1. **Visible overhead** (~32-35k):
   - System prompt: 2-3k
   - Tool definitions: 17-20k
   - Memory files (CLAUDE.md): 8-10k
   - Agent definitions: 3-6k
   - MCP tools: 1-2k

2. **Hidden overhead** (~40-50k):
   - Conversation message formatting (XML structures, timestamps)
   - Tool call/response formatting
   - System reminders and metadata
   - Context management structures
   - Agent invocation overhead

**Total actual overhead: 70-85k tokens** in active development sessions.

### Why Use 40k Estimate?

The 40k estimate is intentionally conservative and used for:
- **Quick manual calculations** when system warning is visible
- **Approximate planning** of context usage
- **Simplified mental model** for developers

However, this estimate will **underreport actual usage by 30-45k tokens**.

### Source of Truth: `/context` Command

**ALWAYS use the `/context` command** for accurate token usage:

```bash
> /context
Total: 108k/200k tokens (54%)
```

This shows **actual total usage** including all hidden overhead.

### Recommended Approach

**For accurate tracking**:
1. Run `/context` command regularly to see actual usage
2. Use the percentage shown by `/context` as source of truth
3. Display that percentage in your context visualization

**For manual estimation** (when `/context` unavailable):
1. Use 40k overhead estimate from system warnings
2. Understand this will underestimate by ~40k tokens
3. Add safety buffer when approaching limits

### Updated Display Requirements

**Best practice**: Match `/context` output exactly:

```
Context: 🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨⬛⬛⬛⬛⬛⬛⬛⬛⬛ 54%
  (108000/200000 tokens)
```

This reflects the actual 108k usage from `/context`, not the 70k you'd calculate from 30.7k messages + 40k overhead.

---

## ⚠️ Common Mistakes to Avoid

### Mistake 1: Using Message Tokens Directly
❌ **WRONG**:
```
System warning shows: 28795 tokens
Display: Context: ... 14% (28795/200000 tokens)
```

✅ **CORRECT**:
```
System warning shows: 28795 tokens
Add overhead: 28795 + 40000 = 68795
Display: Context: ... 34% (68795/200000 tokens)
```

**Why**: Message tokens don't include system prompt, tools, or memory overhead.

---

### Mistake 2: Mismatched Block Count to Percentage
❌ **WRONG**:
```
Percentage: 39%
Display: Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛ 39%
(16 blocks shown = 80%, not 39%)
```

✅ **CORRECT**:
```
Percentage: 39%
Filled blocks: round(39 / 5) = round(7.8) = 8
Display: Context: 🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 39%
(8 blocks shown = 40%, close to 39%)
```

**Why**: Each block = 5%, so 39% should show ~8 blocks, not 16.

---

### Mistake 3: Wrong Color for Percentage
❌ **WRONG**:
```
Percentage: 71%
Display: Context: 🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛ 71%
(Green shown, but 71% is in orange zone)
```

✅ **CORRECT**:
```
Percentage: 71%
Color: Orange (65% < 71% < 80%)
Display: Context: 🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧⬛⬛⬛⬛⬛⬛ 71%
```

**Why**: 71% exceeds 65% threshold, requiring orange blocks.

---

### Mistake 4: Forgetting to Round Percentage
❌ **WRONG**:
```
Calculation: (63795 / 200000) * 100 = 31.8975%
Display: Context: ... 31.8975%
```

✅ **CORRECT**:
```
Calculation: (63795 / 200000) * 100 = 31.8975% ≈ 32%
Display: Context: ... 32%
```

**Why**: Display should show whole number percentages for clarity.

---

### Mistake 5: Not Using Exactly 20 Blocks
❌ **WRONG**:
```
Display: Context: 🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛ 30%
(Only 12 blocks shown)
```

✅ **CORRECT**:
```
Display: Context: 🟩🟩🟩🟩🟩🟩⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 30%
(Exactly 20 blocks total)
```

**Why**: Meter should always show full 20-block scale for consistency.

---

## 📐 Quick Reference Table

| Token Range | Actual Usage | Percentage | Color | Filled Blocks | Action Required |
|-------------|--------------|------------|-------|---------------|-----------------|
| 0-60k       | 40k-100k     | 0-50%      | 🟩    | 0-10          | Normal operation |
| 60k-90k     | 100k-130k    | 50-65%     | 🟨    | 11-13         | Warning - approaching handoff |
| 90k-120k    | 130k-160k    | 65-80%     | 🟧    | 14-16         | **MANDATORY: Create handoff files** |
| 120k-160k   | 160k-200k    | 80-100%    | 🟥    | 17-20         | **Strongly recommend new session** |

**Note**: "Token Range" = message tokens from system warning, "Actual Usage" = message tokens + 40k overhead

---

## 🎯 Context Threshold Actions

### 50% Threshold (100k tokens)
**Display**:
```
Context: 🟨🟨🟨🟨🟨🟨🟨🟨🟨🟨⬛⬛⬛⬛⬛⬛⬛⬛⬛⬛ 50%
  (100000/200000 tokens)

⚠️ WARNING: Approaching handoff threshold. Session handoff files will be created automatically at 65%.
```

**Action**: Display warning message about approaching handoff.

---

### 65% Threshold (130k tokens)
**Display**:
```
Context: 🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧🟧⬛⬛⬛⬛⬛⬛⬛ 65%
  (130000/200000 tokens)

🚨 HANDOFF REQUIRED: Creating session handoff files now.
```

**MANDATORY Actions**:
1. Create `[NUMBER]-SESSION.md` using `session-summary-template.md`
2. Create/overwrite `HANDOFF-SESSION.md` using `handoff-session-template.md`
3. Include all active tasks, agent history, and session context
4. Continue working in current session

---

### 80% Threshold (160k tokens)
**Display**:
```
Context: 🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥🟥⬛⬛⬛⬛ 80%
  (160000/200000 tokens)

🔴 CRITICAL: Context usage at 80%. Strongly recommend starting new session.
   Use HANDOFF-SESSION.md to restore context in new session.
```

**Recommended Action**: Inform user to start new session with handoff file.

---

## ✅ Validation Checklist

Before displaying context usage, verify:

- [ ] Read message tokens from system warning
- [ ] Added 40,000 token overhead to get actual usage
- [ ] Calculated percentage: `(actual_usage / 200000) * 100`
- [ ] Rounded percentage to nearest whole number
- [ ] Calculated filled blocks: `round(percentage / 5)`
- [ ] Applied correct color based on percentage threshold
- [ ] Used exactly 20 blocks total in visual meter
- [ ] Remaining blocks filled with ⬛ (black square)
- [ ] Included token count: `(actual_tokens/200000 tokens)`
- [ ] Block count matches percentage (within ±5%)

**Quick validation**: `filled_blocks * 5` should be within ±5% of displayed percentage.

**Example**:
```
Displayed: 32%
Filled blocks: 6
Validation: 6 * 5 = 30% (within 2% of 32%) ✅
```

---

## 🔧 Implementation Template

Copy this template for consistent implementation:

```python
# Read from system warning
message_tokens = <value_from_system_warning>

# Step 1: Add overhead
actual_usage = message_tokens + 40000

# Step 2: Calculate percentage
percentage = round((actual_usage / 200000) * 100)

# Step 3: Calculate filled blocks
filled_blocks = round(percentage / 5)

# Step 4: Determine color
if percentage < 50:
    color = "🟩"  # Green
elif percentage < 65:
    color = "🟨"  # Yellow
elif percentage < 80:
    color = "🟧"  # Orange
else:
    color = "🟥"  # Red

# Step 5: Build display
filled = color * filled_blocks
empty = "⬛" * (20 - filled_blocks)
meter = filled + empty

# Step 6: Format output
display = f"Context: {meter} {percentage}%\n  ({actual_usage}/200000 tokens)"
```

---

## 📚 Related Documentation

- **CLAUDE.md**: Context & Session Management section
- **Session Handoff Protocol**: `.claude/docs/session-handoff-guide.md` (if exists)
- **Agent History Guide**: `.claude/docs/agent-history-guide.md`

---

**Last Updated**: 2025-11-15
**Version**: 1.2.0
