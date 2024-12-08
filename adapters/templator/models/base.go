package models

type Base struct {
	Step string
	CSRF string
}

type Step string

const (
	StepWeight Step = "weight"
	StepFuel   Step = "fuel"
	StepStats  Step = "stats"
	StepExport Step = "export"
)
