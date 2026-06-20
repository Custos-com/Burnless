package toil

// Cost estimation for toil.
//
// Toil is measured in engineer time. Converting that time into a dollar figure
// makes the ROI of automation legible to managers and budget owners — "we spent
// $4,200 last month on work a script could do" lands harder than "70 hours".

// DefaultAnnualSalaryUSD is the average fully-loaded SRE salary used as the
// default benchmark when no salary is configured. Teams can override it to match
// their own compensation; the default works out to ~$86.54 per hour.
const DefaultAnnualSalaryUSD = 180_000.0

// WorkHoursPerYear is the number of paid working hours in a year
// (40 hours/week × 52 weeks). It converts an annual salary into an hourly rate.
const WorkHoursPerYear = 2080.0

// CostModel converts toil duration into an estimated dollar cost.
//
// AnnualSalaryUSD is configurable so teams can plug in their own compensation
// benchmark (for example, loaded from a config file). When it is zero or
// negative the default SRE salary benchmark is used.
type CostModel struct {
	AnnualSalaryUSD float64
}

// DefaultCostModel returns a CostModel using the default SRE salary benchmark.
func DefaultCostModel() CostModel {
	return CostModel{AnnualSalaryUSD: DefaultAnnualSalaryUSD}
}

// annualSalary returns the configured salary, falling back to the default
// benchmark when none has been set.
func (m CostModel) annualSalary() float64 {
	if m.AnnualSalaryUSD > 0 {
		return m.AnnualSalaryUSD
	}
	return DefaultAnnualSalaryUSD
}

// HourlyRateUSD returns the cost of one hour of engineer time.
func (m CostModel) HourlyRateUSD() float64 {
	return m.annualSalary() / WorkHoursPerYear
}

// CostForHours returns the estimated dollar cost of the given number of toil
// hours. Negative input is treated as zero.
func (m CostModel) CostForHours(hours float64) float64 {
	if hours <= 0 {
		return 0
	}
	return hours * m.HourlyRateUSD()
}

// CostForMinutes returns the estimated dollar cost of the given number of toil
// minutes — the unit Event records duration in. Negative input is treated as
// zero.
func (m CostModel) CostForMinutes(minutes int) float64 {
	if minutes <= 0 {
		return 0
	}
	return m.CostForHours(float64(minutes) / 60.0)
}

// Cost returns the estimated dollar cost of a single toil event.
func (m CostModel) Cost(e Event) float64 {
	return m.CostForMinutes(e.DurationMins)
}
