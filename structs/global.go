package structs

type Filter struct {
	Limit int `json:"limit" query:"limit" formbody:"limit" xml:"limit"`
	Offset int `json:"offset" query:"offset" formbody:"offset" xml:"offset"`
	Search string `json:"search" query:"search" formbody:"search" xml:"search"`
}