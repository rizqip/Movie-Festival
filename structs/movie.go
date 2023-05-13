package structs

type Movie struct {
	MoviesId int64 `json:"MovieId,omitempty" form:"MovieId" query:"MovieId" gorm:"column:MovieId"`
	Title string `json:"Title,omitempty" form:"Title" query:"Title" gorm:"column:Title"`
	Description string `json:"Description,omitempty" form:"Description" query:"Description" gorm:"column:Description"`
	Duration string `json:"Duration,omitempty" form:"Duration" query:"Duration" gorm:"column:Duration"`
	ArtisName string `json:"ArtisName,omitempty" form:"ArtisName" query:"ArtisName"`
	ArtistId string `json:"ArtistId,omitempty" form:"ArtistId" query:"ArtistId" gorm:"column:ArtistId"`
	GenreName string `json:"GenreName,omitempty" form:"GenreName" query:"GenreName"`
	GenreId string `json:"GenreId,omitempty" form:"GenreId" query:"GenreId" gorm:"column:GenreId"`
	Url string `json:"Url" form:"Url,omitempty" query:"Url" gorm:"column:Url"`
	CreatedAt string `json:"CreatedAt,omitempty" form:"CreatedAt" query:"CreatedAt" gorm:"column:CreatedAt"`
	CreatedBy int64 `json:"CreatedBy,omitempty" form:"CreatedBy" query:"CreatedBy" gorm:"column:CreatedBy"`
	UpdatedAt string `json:"UpdatedAt,omitempty" form:"UpdatedAt" query:"UpdatedAt" gorm:"column:UpdatedAt"`
}