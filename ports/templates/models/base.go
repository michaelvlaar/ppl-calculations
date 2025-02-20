package models

type Base struct {
	Step string
	CSRF string
}

type Step string

const (
	StepOverview Step = "overview"
	StepWeight   Step = "weight"
	StepFuel     Step = "fuel"
	StepStats    Step = "stats"
	StepExport   Step = "export"
	StepView     Step = "view"
)
