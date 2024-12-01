package models

import "ppl-calculations/app/queries"

type DownloadRedirect struct {
	Base
	Reference string
}

func DownloadRedirectFromStatsSheet(csrf string, reference string, statsSheet queries.StatsSheetResponse) interface{} {
	template := DownloadRedirect{
		Base: Base{
			CSRF: csrf,
		},
		Reference: reference,
	}
	return template
}
