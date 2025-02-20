package models

import "ppl-calculations/app/queries"

type Export struct {
	Base
	Name string
}

func ExportFromExportSheet(sheet queries.ExportSheetResponse) Export {
	template := Export{
		Base: Base{
			Step: string(StepExport),
		},
	}

	if sheet.Name != nil {
		template.Name = sheet.Name.String()
	}

	return template
}
