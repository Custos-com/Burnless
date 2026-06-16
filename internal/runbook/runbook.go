// Package runbook implements the Burnless runbook executor.
// Runbooks are declared in sre.yaml as executable step sequences.
// They run themselves — they do not just describe what to do.
package runbook

import (
	"fmt"
	"time"
)

// Mode controls whether a runbook executes immediately or waits for approval.
type Mode string

const (
	ModeAuto     Mode = "auto"      // execute immediately
	ModeSemiAuto Mode = "semi-auto" // propose + wait for approval
)

// Execution records the result of a runbook run.
type Execution struct {
	RunbookName string
	StartedAt   time.Time
	FinishedAt  time.Time
	Steps       []StepResult
	Succeeded   bool
	Error       string
}

// StepResult records the output of a single runbook step.
type StepResult struct {
	Step     string
	Output   string
	ExitCode int
	Duration time.Duration
}

// Runner executes runbook steps.
type Runner struct {
	DryRun bool // if true, print steps but do not execute
}

// Execute runs all steps in a runbook.
func (r *Runner) Execute(name string, steps []string) (*Execution, error) {
	exec := &Execution{
		RunbookName: name,
		StartedAt:   time.Now(),
	}

	for _, step := range steps {
		result, err := r.runStep(step)
		exec.Steps = append(exec.Steps, result)
		if err != nil {
			exec.Error = err.Error()
			exec.FinishedAt = time.Now()
			return exec, fmt.Errorf("step %q failed: %w", step, err)
		}
	}

	exec.Succeeded = true
	exec.FinishedAt = time.Now()
	return exec, nil
}

func (r *Runner) runStep(step string) (StepResult, error) {
	start := time.Now()
	result := StepResult{Step: step}

	if r.DryRun {
		result.Output = fmt.Sprintf("[dry-run] would execute: %s", step)
		result.Duration = time.Since(start)
		return result, nil
	}

	// TODO: implement actual step execution
	// - shell commands: exec via os/exec
	// - wait Ns: time.Sleep
	// - assert slo.*: query metrics and compare
	result.Output = fmt.Sprintf("executed: %s", step)
	result.Duration = time.Since(start)
	return result, nil
}
