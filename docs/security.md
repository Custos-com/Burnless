# Enterprise Security & Data Handling

## What data does Burnless actually touch?

Burnless is designed with a minimal data footprint:

| Data type | What it is | Where it lives |
|-----------|-----------|----------------|
| `sre.yaml` | Your SLO rules and thresholds | Your Git repo |
| Metrics | Uptime numbers like `99.85%` | Your Prometheus |
| Event logs | What Burnless did and when | Local file or your DB |

> **Key point: Burnless never touches your business data.**
> No customer records, no payments, no PII — only metric numbers.

---

## Data flow

1. Your service emits metrics (uptime, error rate, latency)
2. Prometheus scrapes and stores those metrics
3. Burnless agent reads your `sre.yaml` and queries Prometheus every 60 seconds
4. If burn rate is too high → runs runbook automatically
5. Sends notification to Slack or PagerDuty
6. Logs what it did

---

## 4 security layers for enterprise

### 1. Data isolation
Each company gets their own isolated database.
Company A can never see Company B's data. Ever.

### 2. Encryption
- **At rest** — AES-256 encryption on all stored data
- **In transit** — TLS 1.3 on all connections (same as HTTPS banking)

### 3. Access control
- SAML/SSO — login with your company Google or Okta account
- Role based access — Admin, Editor, Viewer
- API tokens with scoped permissions
- Session timeout and MFA support

### 4. Audit log
Every action is permanently recorded:
- Who changed which SLO
- When a runbook was triggered
- Who acknowledged an alert
- Non-deletable. Exportable for compliance.

---

## Deployment options

| Option | Description | Best for |
|--------|-------------|----------|
| ☁️ Kairos Cloud (SaaS) | We host it, you log in | Most companies |
| 🏢 Your VPC | Runs inside your own AWS/GCP | Data-sensitive orgs |
| 🖥️ On-premise | Fully air-gapped, no internet | Government, finance |

---

## Compliance certifications

| Certification | Status |
|--------------|--------|
| SOC 2 Type II | Planned Q4 2026 |
| ISO 27001 | Planned Q1 2027 |
| GDPR | Compliant by design |
| HIPAA | Available on request |

---

## Vulnerability reporting

Found a security issue? Please report it privately:
See [SECURITY.md](../SECURITY.md) for responsible disclosure instructions.

Do NOT open a public GitHub issue for security vulnerabilities.
