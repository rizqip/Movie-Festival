package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"movie-festival/helpers"
	"movie-festival/models"
	"movie-festival/structs"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func Register(e echo.Context) error {
	log.Println("Starting process Register User")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Binding Request Payload to Struct */
	var dataReq map[string]interface{}
	var req structs.User
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	AlreadyRegis := models.CheckUser(req.Email)
	if AlreadyRegis{
		response.Data = req
		response.Message = "You Account Already Register"
		response.StatusCode = http.StatusBadRequest

		return e.JSON(http.StatusBadRequest, response)
	}

	/* convert Struct to Map */
	jsonreq, _ := json.Marshal(req)
	Jsonerr := json.Unmarshal([]byte(jsonreq), &dataReq)
	if Jsonerr != nil {
		fmt.Println("error:", Jsonerr)
	}

	/* Proses Insert Movie */
	rsp, tx := models.Register(dataReq)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Create User Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func Login(e echo.Context) error {
	log.Println("Starting process Login User")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Binding Request Payload to Struct */
	var req structs.User
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

	/* Proses Insert Movie */
	DataUser, tx := models.GetUser(req)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	if DataUser == nil {
		response.Data = nil
		response.Message = "Akun Anda Tidak Ditemukan"
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	if DataUser["Password"] != req.Password{
		response.Data = nil
		response.Message = "Password Anda Tidak Sesuai"
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	/* Generate JWT Token */
	token := helpers.GenerateTokenAuth(DataUser["Name"].(string), DataUser["Email"].(string), DataUser["UserId"].(int64), DataUser["Type"].(int64))

	models.UpdateStatus(DataUser["UserId"].(int64), 1)

	response.Data = token
	response.Message = "Login Success"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func Logout(e echo.Context) error {
	log.Println("Starting process Logout User")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Binding Request Payload to Struct */
	var req structs.User
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

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

	var UserId int64
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
	}
	
	models.UpdateStatus(UserId, 0)

	response.Data = nil
	response.Message = "Log Out Success"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func Vote(e echo.Context) error {
	log.Println("Starting process Vote Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Binding Request Payload to Struct */
	var req structs.Vote
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}

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

	var UserId int64
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
	}
	req.UserId = UserId

	/* convert Struct to Map */
	var dataReq map[string]interface{}
	jsonreq, _ := json.Marshal(req)
	Jsonerr := json.Unmarshal([]byte(jsonreq), &dataReq)
	if Jsonerr != nil {
		fmt.Println("error:", Jsonerr)
	}

	if req.Type == 1{
		alreadyVote, tx := models.Vote(dataReq)
		if alreadyVote{
			response.Data = nil
			response.Message = "You Already Vote This Movie"
			response.StatusCode = http.StatusBadRequest
		
			return e.JSON(http.StatusBadRequest, response)
		}
		if tx != nil{
			response.Data = ""
			response.Message = tx.Error()
			response.StatusCode = http.StatusBadRequest
		
			return e.JSON(http.StatusBadRequest, response)
		}
	
		response.Data = nil
		response.Message = "Vote Successfully"
		response.StatusCode = http.StatusOK
	
		return e.JSON(http.StatusOK, response)
	}
	
	tx := models.UnVote(req)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = nil
	response.Message = "Unvote Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func MostVoted(e echo.Context) error {
	log.Println("Starting process Get List Most Voted Movie")

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
	rsp, tx := models.MostVotedMovie(req)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest

		return e.JSON(http.StatusBadRequest, response)
	}
	response.Data = rsp
	response.Message = "Get Most Voted Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func ViewMovieRecord(e echo.Context) error {
	log.Println("Starting process Record View Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	MovieId, _ := strconv.Atoi(e.Param("MovieId"))
	/* Binding Request Payload to Struct */
	var dataReq map[string]interface{}
	var req structs.Vote
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}
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

	var UserId int64
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
	}
	req.UserId = UserId
	req.MovieId = int64(MovieId)
	
	/* convert Struct to Map */
	jsonreq, _ := json.Marshal(req)
	Jsonerr := json.Unmarshal([]byte(jsonreq), &dataReq)
	if Jsonerr != nil {
		fmt.Println("error:", Jsonerr)
	}

	/* Proses Insert Movie */
	url, tx := models.ViewMovieRecord(dataReq)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = url
	response.Message = "Add Record Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}

func RecordViewedMovie(e echo.Context) error {
	log.Println("Starting process Record Viewed Movie")

	/* Declare Variable Response */
	response := new(helpers.JSONResponse)

	/* Binding Request Payload to Struct */
	var req structs.Vote
	if err := e.Bind(&req); err != nil {
		response = &helpers.JSONResponse{
			StatusCode: http.StatusBadRequest,
			Message:    "Bad request",
			Data:       nil,
		}
		return e.JSON(response.StatusCode, response)
	}
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

	var UserId int64
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
	}
	req.UserId = UserId
	
	/* Proses Insert Movie */
	rsp,tx := models.RecordViewedMovie(UserId)
	if tx != nil{
		response.Data = ""
		response.Message = tx.Error()
		response.StatusCode = http.StatusBadRequest
	
		return e.JSON(http.StatusBadRequest, response)
	}

	response.Data = rsp
	response.Message = "Get Record Successfully"
	response.StatusCode = http.StatusOK

	return e.JSON(http.StatusOK, response)
}