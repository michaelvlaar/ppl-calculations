package models

import "ppl-calculations/app/queries"

type Export struct {
	Base
	Name string
}

func ExportFromExportSheet(csrf string, sheet queries.ExportSheetResponse) interface{} {
	template := Export{
		Base: Base{
			Step: string(StepExport),
			CSRF: csrf,
		},
	}

	if sheet.Name != nil {
		template.Name = sheet.Name.String()
	}

	return template
}
