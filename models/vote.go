package models

import (
	"movie-festival/structs"
)

func Vote(req map[string]interface{})(bool,error){
	var alreadyVote bool = true
	var UserVote map[string]interface{}
	delete(req, "Type")

	tx := ModelsDb.Table("MovieVotes").Where("UserId = ? && MovieId = ?", req["UserId"], req["MovieId"]).Find(&UserVote)

	if UserVote == nil {
		tx = ModelsDb.Table("MovieVotes").Create(req)
		alreadyVote = false
	}

	return alreadyVote,tx.Error
}

func UnVote(req structs.Vote)error{
	tx := ModelsDb.Table("MovieVotes").Where("UserId = ? && MovieId = ?", req.UserId, req.MovieId).Delete(req)

	return tx.Error
}

func MostVotedMovie(req structs.Filter)(map[string]interface{},error){
	var ListMovieVoted []map[string]interface{}
	var ListUserId []map[string]interface{}
	tx := ModelsDb.Select("MVote.MovieId, Mov.Title, Count(MVote.MovieId) AS VotedCount")
	tx.Table("MovieVotes MVote")
	tx.Joins("Join Movies Mov On MVote.MovieId = Mov.MovieId")

	if req.Limit > 0 {
		tx.Limit(req.Limit)
	}
	if req.Offset > 0 {
		tx.Offset(req.Offset)
	}
	tx.Group("MVote.MovieId")
	
	var count int64
	tx.Count(&count)

	tx.Order("VotedCount DESC")
	tx.Find(&ListMovieVoted)

	tx = ModelsDb.Select("MovieVotes.MovieId, MovieVotes.UserId, Users.Name")
	tx.Table("MovieVotes")
	tx.Joins("Join Users On MovieVotes.UserId = Users.UserId")
	tx.Find(&ListUserId)

	for _, val := range ListMovieVoted{
		NameUser := make([]string,0)
		for _, user := range ListUserId{
			if val["MovieId"] == user["MovieId"]{
				NameUser = append(NameUser, user["Name"].(string))
			}
		}
		val["ListUser"] = NameUser
	}
	resp := map[string]interface{}{
		"list" : ListMovieVoted,
		"total_items" : count,
	}

	return resp, nil
}