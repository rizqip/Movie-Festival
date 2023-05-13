package models

func CheckGenre(GenreName string) (GenresId int64){
	var GenreId int64
	ModelsDb.Select("GenreId").Table("Genres").Where("GenreName = ?", GenreName).Find(&GenreId)

	if GenreId == 0 {
		req := map[string]interface{}{
			"GenreName":GenreName,
		}
		ModelsDb.Table("Genres").Create(req).Select("GenreId").Table("Genres").Where("GenreName = ?", GenreName).Find(&GenreId)
	}

	return GenreId
}

func MostViewedGenre()([]map[string]interface{}, error){
	var arrGenres []map[string]interface{}
	tx := ModelsDb.Table("Genres").Select("GenreName, ViewCount").Order("ViewCount DESC").Find(&arrGenres)

	return arrGenres, tx.Error
}