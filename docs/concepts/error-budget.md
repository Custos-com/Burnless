# Error budgets and burn rate

An error budget is the allowed amount of unreliability in your SLO.

If your availability SLO is 99.9% over 30 days, your error budget is:

  0.1% × 30 days × 24 hours × 60 minutes = 43.2 minutes of allowed downtime

**Burn rate** measures how fast you are spending that budget.

| Burn rate | Meaning | Budget exhausted in |
|---|---|---|
| 1x | On track | 30 days (at SLO window end) |
| 6x | Warning | 5 days |
| 14.4x | Critical | 2 days |

Burnless uses multiwindow burn rate alerts (1hr + 6hr) — the approach
recommended in the Google SRE Workbook — to catch both fast and slow burns
without excessive false positives.
