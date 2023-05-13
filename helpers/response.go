package helpers

import (
	"encoding/json"
	"log"
)

func HandleError(message string, err interface{}) {
	log.Println("========== Start Error Message ==========")
	log.Println("Message => " + message + ".")
	if err != nil {
		log.Println("Error => ", err)
	}
	log.Println("========== End Of Error Message ==========")
	log.Println()
}

type JSONResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func JSONEncode(data interface{}) string {
	jsonResult, _ := json.Marshal(data)

	return string(jsonResult)
}
