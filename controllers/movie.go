package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie-festival/helpers"
	"movie-festival/models"
	"movie-festival/structs"
	"movie-festival/system"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
)

func CreateMovie(e echo.Context) error {
	log.Println("Starting process Create Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Decode JWT */
	token := e.Request().Header.Get("Authorization")
	if token == "" || len(token) <= 7 {
		log.Println("error token is empty")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	/* Upload File */
	file, _ := system.UploadFile(e, "movie", "./")
	if file == nil {
		response.Data = nil
		response.Message = "File Is Empty or Upload New File"
		response.StatusCode = http.StatusBadRequest
		return e.JSON(response.StatusCode, response)
	}

	/* Binding Request Payload to Struct */
	var dataReq map[string]interface{}
	var req structs.Movie
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	var UserId, Type int64
	if data, isSuccess := helpers.DecodeTokenJwt(token[7:]); !isSuccess {
		log.Println("error decode token:")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	} else {
		UserId = int64(data["UserId"].(float64))
		Type = int64(data["UserTypes"].(float64))
	}
	req.CreatedBy = UserId

	if Type != 2{
		response.Data = nil
		response.Message = "Maaf Anda Bukan Admin"
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	/* Get Directory Url Movie */
	mydir, err := os.Getwd()
    if err != nil {
        fmt.Println(err)
    }
	url := mydir+"/"+file.(string)
	req.Url = url

	/* Get Id Artis */
	var artisName []string
	json.Unmarshal([]byte(req.ArtisName), &artisName)
	var ArtisId []int64
	for _, artis := range artisName {
		ArtisId = append(ArtisId, models.CheckArtis(artis))
	}
	jsonArtisId, _ := json.Marshal(ArtisId)
	req.ArtistId = string(jsonArtisId)

	/* Get Id Genre */
	var GenreName []string
	json.Unmarshal([]byte(req.GenreName), &GenreName)
	var GenreId []int64
	for _, genre := range GenreName {
		GenreId = append(GenreId, models.CheckGenre(genre))
	}
	jsonGenreIdId, _ := json.Marshal(GenreId)
	req.GenreId	 = string(jsonGenreIdId)

	/* convert Struct to Map */
	jsonreq, _ := json.Marshal(req)
	Jsonerr := json.Unmarshal([]byte(jsonreq), &dataReq)
	if Jsonerr != nil {
		fmt.Println("error:", Jsonerr)
	}

	/* Proses Insert Movie */
	rsp, tx := models.CreateMovie(dataReq)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Create Movie Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func UpdateMovie(e echo.Context) error {
	log.Println("Starting process Update Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Decode JWT */
	token := e.Request().Header.Get("Authorization")
	if token == "" || len(token) <= 7 {
		log.Println("error token is empty")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	var Type int64
	if data, isSuccess := helpers.DecodeTokenJwt(token[7:]); !isSuccess {
		log.Println("error decode token:")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	} else {
		Type = int64(data["UserTypes"].(float64))
	}

	if Type != 2{
		response.Data = nil
		response.Message = "Maaf Anda Bukan Admin"
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	/* Binding Request Payload to Struct */
	var dataReq map[string]interface{}
	var req structs.Movie
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	/* Upload File */
	file, _ := system.UploadFile(e, "movie", "./")
	if file != nil {
		/* Get Directory Url Movie */
		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		url := mydir+"/"+file.(string)
		req.Url = url
	}

	/* Get MovieId */
	MoviesId, _ := strconv.Atoi(e.Param("MovieId"))
	req.MoviesId = int64(MoviesId)

	/* Get Id Artis */
	if req.ArtisName != ""{
		var artisName []string
		json.Unmarshal([]byte(req.ArtisName), &artisName)
		var ArtisId []int64
		for _, artis := range artisName {
			ArtisId = append(ArtisId, models.CheckArtis(artis))
		}
		jsonArtisId, _ := json.Marshal(ArtisId)
		req.ArtistId = string(jsonArtisId)
	}

	/* Get Id Genre */
	if req.GenreName != "" {
		var GenreName []string
		json.Unmarshal([]byte(req.GenreName), &GenreName)
		var GenreId []int64
		for _, genre := range GenreName {
			GenreId = append(GenreId, models.CheckGenre(genre))
		}
		jsonGenreIdId, _ := json.Marshal(GenreId)
		req.GenreId	 = string(jsonGenreIdId)
	}

	/* convert Struct to Map */
	jsonreq, _ := json.Marshal(req)
	Jsonerr := json.Unmarshal([]byte(jsonreq), &dataReq)
	if Jsonerr != nil {
		fmt.Println("error:", Jsonerr)
	}

	/* Proses Update Movie */
	rsp, tx := models.UpdateMovie(dataReq)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Update Movie Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func ListMovie(e echo.Context) error {
	log.Println("Starting process Get List Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Decode JWT */
	token := e.Request().Header.Get("Authorization")
	if token == "" || len(token) <= 7 {
		log.Println("error token is empty")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	/* Binding Request Payload to Struct */
	var req structs.Filter
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	/* Proses List Movie */
	rsp, tx := models.ListMovie(req)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Get List Movie Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func DetailMovie(e echo.Context) error {
	log.Println("Starting process Get Detail Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Get MovieId */
	MoviesId, _ := strconv.Atoi(e.Param("MovieId"))

	/* Proses Update Movie */
	rsp, tx := models.DetailMovie(int64(MoviesId))
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Get Detail Movie Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func MostViewed(e echo.Context) error {
	log.Println("Starting process Get Most Viewed")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Get MovieId */
	Type, _ := strconv.Atoi(e.QueryParam("Type"))

	/* Decode JWT */
	token := e.Request().Header.Get("Authorization")
	if token == "" || len(token) <= 7 {
		log.Println("error token is empty")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	var UserType int64
	if data, isSuccess := helpers.DecodeTokenJwt(token[7:]); !isSuccess {
		log.Println("error decode token:")
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Token is not valid",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	} else {
		UserType = int64(data["UserTypes"].(float64))
	}

	if UserType != 2{
		response.Data = nil
		response.Message = "Maaf Anda Bukan Admin"
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	/* Proses Update Movie */
	var rsp []map[string]interface{}
	var tx error
	var msg string
	if Type == 1{
		rsp, tx = models.MostViewedMovie()
		msg = "Success Get Most Movie Viewed"
	}else{
		rsp, tx = models.MostViewedGenre()
		msg = "Success Get Most Genre Viewed"
	}
	
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = msg
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}