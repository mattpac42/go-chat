# Agent Invocation Examples

## Task Tool Invocation Template

```
Task: [Specific task description]

Context:
- [Essential context item 1]
- [Essential context item 2]
- [Essential context item 3]

Constraints:
- [Constraint 1]
- [Constraint 2]

Expected Deliverables:
- [Deliverable 1]
- [Deliverable 2]

Success Criteria:
- [Criteria 1]
- [Criteria 2]

Please provide your analysis and recommendations following your specialized prompt format.
```

## Agent-Specific Examples

### Platform Engineer
```
Task: Design auto-scaling strategy for web application

Context:
- Node.js application on Kubernetes
- Current load: 1000 req/min peak
- Budget constraints: $500/month

Constraints:
- Must maintain 99.9% uptime
- Response time < 200ms

Expected Deliverables:
- HPA/VPA configuration files
- Load balancer setup
- Monitoring configuration

Success Criteria:
- Handles 5x traffic spikes
- Cost stays within budget
- Zero downtime deployment
```

### Security Engineer
```
Task: Analyze authentication system for vulnerabilities

Context:
- JWT-based authentication
- User registration and login endpoints
- Password reset functionality

Constraints:
- Must comply with SOC2 requirements
- No breaking changes to existing API

Expected Deliverables:
- Vulnerability assessment report
- Security hardening recommendations
- Compliance gap analysis

Success Criteria:
- All critical vulnerabilities identified
- Actionable remediation steps provided
- SOC2 compliance achieved
```
