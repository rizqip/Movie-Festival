package models

import (
	"encoding/json"
	"movie-festival/structs"
	"movie-festival/system"
	"os"
	"strconv"
	"time"
)

func CreateMovie(req map[string]interface{}) (map[string]interface{}, error) {
	var responseMovie map[string]interface{}
	delete(req,"ArtisName")
	delete(req,"GenreName")
	dataReq := system.CreatedTimeNow(req)

	tx := ModelsDb.Table("Movies").Create(dataReq).Limit(1).Order("MovieId Desc").Find(&responseMovie)

	responseMovie["CreatedAt"] = responseMovie["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")
	responseMovie["UpdatedAt"] = responseMovie["UpdatedAt"].(time.Time).Format("2006-01-02 15:04:05")

	return responseMovie, tx.Error
}

func UpdateMovie(req map[string]interface{}) (map[string]interface{}, error) {
	var responseMovie map[string]interface{}
	delete(req,"ArtisName")
	delete(req,"GenreName")
	if req["Url"] == "" {
		delete(req,"Url")
	}
	dataReq := system.UpdatedTimeNow(req)

	tx := ModelsDb.Table("Movies").Where("MovieId = ?", req["MovieId"]).Updates(dataReq).Find(&responseMovie)

	responseMovie["CreatedAt"] = responseMovie["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")
	responseMovie["UpdatedAt"] = responseMovie["UpdatedAt"].(time.Time).Format("2006-01-02 15:04:05")

	return responseMovie, tx.Error
}

func ListMovie(req structs.Filter)(map[string]interface{}, error){
	var ListMovie []map[string]interface{}

	tx := ModelsDb.Select("MovieId,Title,Description")
	tx.Table("Movies")

	if req.Search != ""{
		tx.Where("Title LIKE ? OR Description LIKE ?", "%"+req.Search+"%","%"+req.Search+"%")
	}

	var count int64
	tx.Count(&count)

	if req.Limit > 0{
		tx.Limit(req.Limit)
	}
	if req.Offset > 0{
		tx.Offset(req.Offset)
	}

	tx.Find(&ListMovie)

	resp := map[string]interface{}{
		"list": ListMovie,
		"total_items": count,
	}

	return resp, nil
}

func DetailMovie(MovieId int64)(map[string]interface{}, error ){
	/* Get Detail Movie */
	var DetailMovie map[string]interface{}
	ModelsDb.Table("Movies").Where("MovieId = ?", MovieId).Find(&DetailMovie)
	/* Get Genre Name */
	var GenreId []int
	json.Unmarshal([]byte(DetailMovie["GenreId"].(string)), &GenreId)
	var arrGenre []map[string]interface{}
	ModelsDb.Table("Genres").Where("GenreId IN ?", GenreId).Find(&arrGenre)

	/* Get Artist Name */
	var ArtistId []int
	json.Unmarshal([]byte(DetailMovie["ArtistId"].(string)), &ArtistId)
	var arrArtist []string
	ModelsDb.Select("ArtistName").Table("Artists").Where("ArtistId IN ?", ArtistId).Find(&arrArtist)

	/* Add New Response */
	DetailMovie["GenreName"] = arrGenre
	DetailMovie["ArtistName"] = arrArtist


	DetailMovie["CreatedAt"] = DetailMovie["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")
	DetailMovie["UpdatedAt"] = DetailMovie["UpdatedAt"].(time.Time).Format("2006-01-02 15:04:05")


	MovieIdStr := strconv.Itoa(int(MovieId))
	DetailMovie["Url"] = os.Getenv("BASE_URL") + "/general/viewMovieRecord/" + MovieIdStr

	return DetailMovie, nil
}

func MostViewedMovie()([]map[string]interface{}, error){
	var arrMovies []map[string]interface{}
	tx := ModelsDb.Table("Movies").Select("Title, ViewCount").Order("ViewCount DESC").Find(&arrMovies)

	return arrMovies, tx.Error
}