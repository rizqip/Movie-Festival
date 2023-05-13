package models

import (
	"encoding/json"
	"movie-festival/structs"
	"movie-festival/system"
	"time"
)

func Register(req map[string]interface{}) (map[string]interface{}, error) {
	var responseUser map[string]interface{}
	dataReq := system.CreatedTimeNow(req)

	tx := ModelsDb.Table("Users").Create(dataReq).Limit(1).Order("UserId Desc").Find(&responseUser)

	responseUser["CreatedAt"] = responseUser["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")
	responseUser["UpdatedAt"] = responseUser["UpdatedAt"].(time.Time).Format("2006-01-02 15:04:05")

	return responseUser, tx.Error
}

func CheckUser(email string)(alreadyRegis bool){
	var User map[string]interface{}
	ModelsDb.Table("Users").Where("Email = ?", email).Find(&User)

	alreadyRegis = false
	if User != nil {
		alreadyRegis = true
	}
	
	return alreadyRegis
}

func GetUser(req structs.User)(map[string]interface{}, error){
	var DataUser map[string]interface{}
	tx := ModelsDb.Table("Users")

	if req.UserId > 0{
		tx.Where("UserId = ?", req.UserId)
	}
	if req.Email != "" {
		tx.Where("Email = ?", req.Email)
	}
	
	tx.Find(&DataUser)

	return DataUser, nil
}

func UpdateStatus(UserId, Status int64) error{
	tx := ModelsDb.Table("Users").Where("UserId = ?", UserId).Update("Status", Status)

	return tx.Error
}

func ViewMovieRecord(req map[string]interface{})(string, error){
	/* Insert Record View */
	req["CreatedAt"] = time.Now().Format("2006-01-02 15:04:05")
	tx := ModelsDb.Table("UserRecords").Create(req)

	/* Update View Count Movie */
	var DetailMovie map[string]interface{}
	ModelsDb.Table("Movies").Where("MovieId = ?", req["MovieId"]).Find(&DetailMovie)
	ModelsDb.Table("Movies").Where("MovieId = ?", req["MovieId"]).Updates(map[string]interface{}{"ViewCount": DetailMovie["ViewCount"].(int64)+1})

	/* Get Genre Name */
	var GenreId []int
	json.Unmarshal([]byte(DetailMovie["GenreId"].(string)), &GenreId)
	var arrGenre []map[string]interface{}
	ModelsDb.Table("Genres").Where("GenreId IN ?", GenreId).Find(&arrGenre)
	
	for _, val := range arrGenre{
		ModelsDb.Table("Genres").Where("GenreId = ?", val["GenreId"]).Update("ViewCount", val["ViewCount"].(int64)+1)
	}

	var url string
	ModelsDb.Select("Url").Table("Movies").Where("MovieId = ?",req["MovieId"]).Find(&url)

	return url, tx.Error
}

func RecordViewedMovie(UserId int64)([]map[string]interface{}, error){
	var listMovie []map[string]interface{}
	tx := ModelsDb.Select("UsrRec.CreatedAt, UsrRec.MovieId, Mov.Title, UsrRec.UserId")
	tx.Table("UserRecords UsrRec")
	tx.Joins("Join Movies Mov On UsrRec.MovieId = Mov.MovieId")
	tx.Where("UsrRec.UserId = ?", UserId)
	tx.Find(&listMovie)

	for _, val := range listMovie{
		val["CreatedAt"] = val["CreatedAt"].(time.Time).Format("2006-01-02 15:04:05")
	}
	
	return listMovie, nil
}