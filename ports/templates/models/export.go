package models

import "ppl-calculations/app/queries"

type Export struct {
	Name string
}

func ExportFromExportSheet(sheet queries.ExportSheetResponse) Export {
	template := Export{}

	if sheet.Name != nil {
		template.Name = sheet.Name.String()
	}

	return template
}
