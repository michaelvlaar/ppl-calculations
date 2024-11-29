package models

type Export struct {
	Base
}

func ExportFromExportSheet(csrf string) interface{} {
	template := Export{
		Base: Base{
			Step: string(StepExport),
			CSRF: csrf,
		},
	}

	return template
}
