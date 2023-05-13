package models

func CheckArtis(ArtisName string) (ArtistsId int64){
	var ArtisId int64
	ModelsDb.Select("ArtistId").Table("Artists").Where("ArtistName = ?", ArtisName).Find(&ArtisId)

	if ArtisId == 0 {
		req := map[string]interface{}{
			"ArtistName":ArtisName,
		}
		ModelsDb.Table("Artists").Create(req).Select("ArtistId").Table("Artists").Where("ArtistName = ?", ArtisName).Find(&ArtisId)
	}

	return ArtisId
}