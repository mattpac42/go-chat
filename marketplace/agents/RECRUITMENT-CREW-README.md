# Recruitment Agent Crew Documentation

## Overview

This specialized crew of 7 recruitment agents provides a complete, automated resume screening and candidate evaluation system designed to process 100s of resumes efficiently with consistent, objective scoring.

**Key Features**:
- Priority weighting: Cultural Fit (40%) > Technical Skills (35%) > Experience (25%)
- Comprehensive 0-100 scoring scale with detailed breakdowns
- Multi-format resume support (PDF, DOCX, Markdown)
- Fully automated pipeline with seamless agent handoffs
- Industry-agnostic design adaptable to any market
- Partial scoring with penalties for incomplete data
- Intelligent inference for missing information

---

## Agent Crew Architecture

### Workflow Pipeline

```
┌─────────────────────────────────────────────────────────────────────┐
│                         RECRUITMENT PIPELINE                         │
└─────────────────────────────────────────────────────────────────────┘

1. Job Requirements Analyst
   └─> Parses job description
   └─> Creates weighted evaluation framework (40/35/25)
   └─> Generates keyword lists for parsing

2. Resume Parser
   └─> Extracts structured data from resumes (PDF/DOCX/MD)
   └─> Normalizes candidate profiles
   └─> Flags missing information

3. Cultural Fit Analyst (40% weight)
   └─> Scores career stability, leadership, collaboration
   └─> Provides 0-40 point cultural fit assessment

4. Skills Matcher (35% weight)
   └─> Matches technical skills against requirements
   └─> Identifies gaps and transferable skills
   └─> Provides 0-35 point technical skills score

5. Experience Evaluator (25% weight)
   └─> Assesses relevant experience depth
   └─> Analyzes career progression quality
   └─> Provides 0-25 point experience score

6. Ranking Coordinator
   └─> Aggregates all scores into composite (0-100)
   └─> Classifies candidates: Strong/Potential/Weak Match
   └─> Generates final rankings

7. Report Generator
   └─> Creates individual candidate reports
   └─> Produces executive summaries
   └─> Builds comparative dashboards
```

---

## Agent Details

### 1. Job Requirements Analyst
**File**: `recruitment-job-requirements-analyst.md`
**Color**: Teal (#14b8a6 / #0d9488)

**Purpose**: Parse job descriptions and create weighted evaluation frameworks

**Key Outputs**:
- Structured requirements extraction (Cultural Fit, Technical Skills, Experience)
- 0-100 scoring rubric with weighted breakdowns (40/35/25)
- Cultural fit indicators and keywords
- Technical skills requirements with proficiency levels
- Experience criteria with relevance definitions
- Keyword lists for resume parser optimization

**When to Use**: At the start of every recruitment process to establish evaluation criteria

---

### 2. Resume Parser
**File**: `recruitment-resume-parser.md`
**Color**: Emerald (#10b981 / #059669)

**Purpose**: Extract structured data from resumes in multiple formats

**Key Outputs**:
- Structured candidate profiles (JSON + Markdown)
- Normalized data (skills, experience, education, contact info)
- Data quality flags for missing or incomplete information
- Partial scoring penalties applied per framework
- Extraction confidence scores

**When to Use**: After requirements analysis, to process all candidate resumes

**Supported Formats**: PDF, DOCX, Markdown

---

### 3. Cultural Fit Analyst
**File**: `recruitment-cultural-fit-analyst.md`
**Color**: Cyan (#06b6d4 / #0891b2)

**Purpose**: Evaluate candidate alignment with company culture (highest priority - 40%)

**Key Outputs**:
- Cultural fit scores (0-40 points) with sub-scores:
  - Career Stability (0-12 points)
  - Leadership Indicators (0-10 points)
  - Collaboration & Teamwork (0-10 points)
  - Communication Style (0-8 points)
- Supporting evidence with resume excerpts
- Cultural fit rankings within candidate pool

**When to Use**: After resume parsing, as first evaluation dimension (highest weight)

---

### 4. Skills Matcher
**File**: `recruitment-skills-matcher.md`
**Color**: Purple (#8b5cf6 / #7c3aed)

**Purpose**: Match candidate technical skills against job requirements (35%)

**Key Outputs**:
- Technical skills scores (0-35 points) with breakdowns:
  - Core required skills match (0-25 points)
  - Preferred skills match (0-10 points)
- Skill match percentages with proficiency adjustments
- Equivalent and transferable skills recognition
- Comprehensive skill gap analysis
- Skills-based candidate rankings

**When to Use**: After resume parsing, as second evaluation dimension

---

### 5. Experience Evaluator
**File**: `recruitment-experience-evaluator.md`
**Color**: Amber (#f59e0b / #d97706)

**Purpose**: Assess relevant work experience quality and depth (25%)

**Key Outputs**:
- Experience scores (0-25 points) with sub-scores:
  - Base experience (0-15 points)
  - Career progression quality (0-6 points)
  - Industry/domain relevance (0-4 points)
- Career trajectory analysis
- Relevant experience calculations with recency weighting
- Experience-based candidate rankings

**When to Use**: After resume parsing, as third evaluation dimension

---

### 6. Ranking Coordinator
**File**: `recruitment-ranking-coordinator.md`
**Color**: Pink (#ec4899 / #db2777)

**Purpose**: Synthesize all evaluations and create final candidate rankings

**Key Outputs**:
- Weighted composite scores (0-100) = Culture (40) + Skills (35) + Experience (25)
- Candidate tier classifications:
  - Strong Match (80-100): Ready for interview
  - Potential Match (60-79): Consider with focused evaluation
  - Weak Match (0-59): Pass recommendation
- Overall and tier-specific rankings
- Comparative analysis and unbalanced score detection
- Top-N shortlists with justifications

**When to Use**: After all evaluation agents complete scoring, to produce final rankings

---

### 7. Report Generator
**File**: `recruitment-report-generator.md`
**Color**: Indigo (#6366f1 / #4f46e5)

**Purpose**: Create comprehensive evaluation reports and summaries

**Key Outputs**:
- Individual candidate assessment reports with complete score breakdowns
- Executive summaries (1-2 pages) with top candidate highlights
- Comparative dashboards showing side-by-side candidate analysis
- Tier-specific candidate lists with within-tier rankings
- Top-N shortlist reports with interview preparation guidance
- Both JSON (machine-readable) and Markdown (human-readable) formats

**When to Use**: After ranking coordinator completes, to produce hiring decision support documentation

---

## Usage Workflows

### Complete Recruitment Process

**Step 1: Initialize Evaluation Framework**
```
User: "Analyze this job description and create an evaluation framework"
→ Invoke: recruitment-job-requirements-analyst
→ Output: Weighted scoring rubric, keyword lists, evaluation criteria
```

**Step 2: Process Candidate Resumes**
```
User: "Parse these 150 resumes and extract structured candidate data"
→ Invoke: recruitment-resume-parser
→ Output: Structured candidate profiles (JSON + Markdown)
```

**Step 3: Evaluate Cultural Fit (40% weight)**
```
User: "Score cultural fit for all candidates"
→ Invoke: recruitment-cultural-fit-analyst
→ Output: 0-40 point cultural fit scores with detailed breakdowns
```

**Step 4: Match Technical Skills (35% weight)**
```
User: "Evaluate technical skills match for all candidates"
→ Invoke: recruitment-skills-matcher
→ Output: 0-35 point technical skills scores with gap analysis
```

**Step 5: Assess Experience (25% weight)**
```
User: "Score experience relevance for all candidates"
→ Invoke: recruitment-experience-evaluator
→ Output: 0-25 point experience scores with progression analysis
```

**Step 6: Generate Final Rankings**
```
User: "Aggregate all scores and rank candidates"
→ Invoke: recruitment-ranking-coordinator
→ Output: Composite scores (0-100), tier classifications, rankings
```

**Step 7: Create Reports**
```
User: "Generate executive summary and top 10 candidate reports"
→ Invoke: recruitment-report-generator
→ Output: Executive summary, individual reports, comparative dashboard
```

### Quick Workflows

**Quick Screening (Top Candidates Only)**
```
1. Job Requirements Analyst → Framework creation
2. Resume Parser → Extract all candidate data
3. Ranking Coordinator → Use simplified scoring for quick tiers
4. Report Generator → Executive summary with top 10
```

**Focused Re-Evaluation (Single Dimension)**
```
Example: Re-score technical skills after framework update
1. Job Requirements Analyst → Update technical skills criteria
2. Skills Matcher → Re-score all candidates
3. Ranking Coordinator → Recalculate composite with new skills scores
4. Report Generator → Updated reports//////////
```

**Comparative Analysis (Finalists)**
```
Example: Deep dive on top 5 candidates
1. Report Generator → Create comparative dashboard for finalists
2. Include dimension-by-dimension breakdowns
3. Highlight relative strengths and interview focus areas
```

---

## Scoring Framework Details

### Weighting Distribution
- **Cultural Fit**: 40 points (40%)
  - Career Stability: 0-12 points
  - Leadership Indicators: 0-10 points
  - Collaboration & Teamwork: 0-10 points
  - Communication Style: 0-8 points

- **Technical Skills**: 35 points (35%)
  - Core Required Skills: 0-25 points
  - Preferred Skills: 0-10 points

- **Experience**: 25 points (25%)
  - Base Experience: 0-15 points
  - Career Progression: 0-6 points
  - Industry/Domain Relevance: 0-4 points

**Total**: 100 points

### Tier Classifications
- **Strong Match (80-100)**: Excellent candidates, ready for interview
- **Potential Match (60-79)**: Good candidates, may need focused evaluation
- **Weak Match (0-59)**: Significant gaps, pass recommendation

---

## Edge Case Handling

### Missing Data
- **Strategy**: Partial scoring with penalties + intelligent inference
- **Resume Parser**: Flags missing fields, applies penalties per framework
- **Evaluation Agents**: Score available data, document confidence levels
- **Report Generator**: Highlights data quality issues in reports

### Non-Standard Formats
- **Strategy**: Graceful degradation with maximum extraction effort
- **Resume Parser**: Handles unusual layouts, extracts what's possible
- **Quality Flags**: Documents extraction challenges for manual review
- **Partial Credit**: Scores what can be extracted, penalizes gaps

### Equivalent/Transferable Skills
- **Strategy**: Intelligent recognition with partial credit
- **Skills Matcher**: Uses equivalency mappings (React ≈ Angular, AWS ≈ Azure)
- **Credit Calculations**:
  - Direct equivalent: 100% credit
  - Close equivalent: 80-90% credit
  - Transferable: 50-70% credit

### Incomplete Resumes
- **Strategy**: Score what's available, flag for human review
- **Penalties**: Applied consistently per framework rules
- **Inference**: Attempt reasonable inference from context
- **Flagging**: Clearly mark candidates requiring manual evaluation

---

## Best Practices

### For Optimal Results

1. **Clear Job Descriptions**: Provide detailed job descriptions to Job Requirements Analyst
   - Include explicit technical requirements
   - Describe company culture and values
   - Specify experience requirements clearly

2. **Batch Processing**: Process all resumes together for consistent evaluation
   - Ensures fair comparison across entire pool
   - Optimizes agent efficiency
   - Maintains consistent scoring standards

3. **Framework Validation**: Review evaluation framework before mass scoring
   - Confirm weighting feels appropriate (40/35/25 is default)
   - Validate cultural fit indicators match company values
   - Ensure technical requirements are comprehensive

4. **Iterative Refinement**: Use insights from first pass to refine
   - Adjust framework if initial results don't match expectations
   - Re-score candidates after framework updates
   - Learn from successful hire patterns over time

5. **Human Oversight**: Use agents for screening, humans for final decisions
   - Agents provide objective data and recommendations
   - Humans make final hiring calls considering intangibles
   - Review edge cases and boundary candidates manually

### Common Pitfalls to Avoid

1. **Skipping Requirements Analysis**: Don't jump straight to parsing without framework
2. **Ignoring Data Quality Flags**: Review flagged candidates manually
3. **Over-relying on Automation**: Agents screen, humans decide
4. **Inconsistent Batching**: Process comparable candidates together
5. **Ignoring Cultural Fit**: It's 40% for a reason - don't deprioritize

---

## Customization Options

### Adjusting Weights
Default: Cultural Fit (40%) + Technical Skills (35%) + Experience (25%)

To adjust:
1. Modify framework in Job Requirements Analyst
2. Ensure weights still sum to 100%
3. Update Ranking Coordinator to use new weights
4. Re-score all candidates with updated framework

### Adding Custom Criteria
1. Extend relevant evaluation agent (Cultural Fit, Skills, Experience)
2. Add sub-criteria to scoring breakdown
3. Update framework in Job Requirements Analyst
4. Maintain consistent 0-100 total scale

### Industry-Specific Adaptations
- Healthcare: Add credential verification, license checking
- Government: Add security clearance evaluation
- Startups: Weight adaptability and breadth over depth
- Enterprise: Weight stability and process experience higher

---

## Agent Color Scheme

Visual identification for agent library:

| Agent | Light Mode | Dark Mode |
|-------|-----------|-----------|
| Job Requirements Analyst | #14b8a6 (Teal) | #0d9488 |
| Resume Parser | #10b981 (Emerald) | #059669 |
| Cultural Fit Analyst | #06b6d4 (Cyan) | #0891b2 |
| Skills Matcher | #8b5cf6 (Purple) | #7c3aed |
| Experience Evaluator | #f59e0b (Amber) | #d97706 |
| Ranking Coordinator | #ec4899 (Pink) | #db2777 |
| Report Generator | #6366f1 (Indigo) | #4f46e5 |

Colors chosen to visually distinguish recruitment agents from other domains in the agent library.

---

## Integration Notes

### Data Flow Format
- **Agent-to-Agent**: JSON structured data for automation
- **Agent-to-Human**: Markdown formatted reports for readability
- **Hybrid Output**: Both formats provided simultaneously

### Session History
- Each agent creates session history file in `.claude/context/agent-history/`
- Naming pattern: `[YYYYMMDD-HHMMSS]-[agent-name]-[SEQUENCE].md`
- Complete audit trail of all recruitment decisions

### Context Management
- Agents designed to stay under 65% context window usage
- Batch processing optimized for 100s of candidates
- Efficient data structures minimize token consumption

---

## Support and Troubleshooting

### Common Issues

**Issue**: Resumes not parsing correctly
- **Solution**: Check file format (PDF/DOCX/MD), validate file integrity, review Resume Parser output flags

**Issue**: Cultural fit scores seem inconsistent
- **Solution**: Review company values definition in framework, ensure resume quality is adequate, check for missing soft skill indicators

**Issue**: Too many/too few candidates in Strong Match tier
- **Solution**: Review framework criteria stringency, adjust sub-criteria weights if needed, validate scoring thresholds

**Issue**: Skills matcher missing obvious matches
- **Solution**: Check keyword lists, add equivalent skill mappings, verify skill normalization in Resume Parser

### Getting Help

For agent improvements or bug reports:
1. Review agent session history files for debugging insights
2. Consult `.claude/docs/agent-history-guide.md` for agent workflow details
3. Use `agent-editor` agent to refine agent instructions if needed

---

## Version History

- **v1.0 (2025-11-12)**: Initial recruitment crew creation
  - 7 specialized agents covering complete recruitment pipeline
  - Fully automated workflow with agent handoffs
  - 40/35/25 weighted scoring (Culture/Skills/Experience)
  - Multi-format resume support (PDF/DOCX/MD)
  - High-volume processing (100s of candidates)
  - Industry-agnostic design

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────┐
│           RECRUITMENT CREW QUICK REFERENCE               │
├─────────────────────────────────────────────────────────┤
│ 1. Job Requirements Analyst → Create framework          │
│ 2. Resume Parser → Extract candidate data               │
│ 3. Cultural Fit Analyst → Score culture (0-40)          │
│ 4. Skills Matcher → Score skills (0-35)                 │
│ 5. Experience Evaluator → Score experience (0-25)       │
│ 6. Ranking Coordinator → Generate rankings (0-100)      │
│ 7. Report Generator → Create reports & summaries        │
├─────────────────────────────────────────────────────────┤
│ Scoring: Culture 40% + Skills 35% + Experience 25%      │
│ Tiers: Strong (80-100), Potential (60-79), Weak (0-59)  │
│ Formats: PDF, DOCX, Markdown resumes supported          │
│ Scale: Optimized for 100s of candidates                 │
└─────────────────────────────────────────────────────────┘
```

---

*Documentation Generated: 2025-11-12*
*Agent Crew Version: 1.0*
*For questions or improvements, consult the agent-creator or agent-editor agents*
