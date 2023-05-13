package structs

type User struct {
	UserId int64 `json:"UserId" form:"UserId" query:"UserId" gorm:"column:UserId"`
	Name string `json:"Name" form:"Name" query:"Name" gorm:"column:Name"`
	Email string `json:"Email" form:"Email" query:"Email" gorm:"column:Email"`
	Password string `json:"Password" form:"Password" query:"Password" gorm:"column:Password"`
	Type int64 `json:"Type" form:"Type" query:"Type" gorm:"column:Type"`
	Status int64 `json:"Status" form:"Status" query:"Status" gorm:"column:Status"`
	CreatedAt string `json:"CreatedAt,omitempty" form:"CreatedAt" query:"CreatedAt" gorm:"column:CreatedAt"`
	UpdatedAt string `json:"UpdatedAt,omitempty" form:"UpdatedAt" query:"UpdatedAt" gorm:"column:UpdatedAt"`
}

type Vote struct {
	UserId int64 `json:"UserId" form:"UserId" query:"UserId" gorm:"column:UserId"`
	MovieId int64 `json:"MovieId,omitempty" form:"MovieId" query:"MovieId" gorm:"column:MovieId"`
	Type int64 `json:"Type,omitempty" form:"Type" query:"Type" gorm:"column:Type"`
}