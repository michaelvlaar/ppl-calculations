package models

type DownloadRedirect struct {
	Base
	Reference string
}

func DownloadRedirectFromStatsSheet(csrf string, reference string) interface{} {
	template := DownloadRedirect{
		Base: Base{
			CSRF: csrf,
		},
		Reference: reference,
	}
	return template
}
