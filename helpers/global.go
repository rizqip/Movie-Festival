package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func OnlyInt(str string) int {

	re, _ := regexp.Compile(`^[0-9]+`)

	str = re.FindString(str)

	result, _ := strconv.Atoi(str)

	return result
}

func CreateHMACSHA256(params interface{}) (result string, err error) {
	/* create json string from parameters */
	jsonParam, err := json.Marshal(params)
	if err != nil {
		return
	}
	log.Println(string(jsonParam))
	/* create HMAC from secret key 3f0ee534c3d86cb16f4413fe1a76a12f */
	key := []byte("3f0ee534c3d86cb16f4413fe1a76a12f")
	message := string(jsonParam)

	sig := hmac.New(sha256.New, key)
	sig.Write([]byte(message))

	result = hex.EncodeToString(sig.Sum(nil))
	return
}

func MarshalUnmarshal(param interface{}, result interface{}) error {
	paramByte, err := json.Marshal(param)
	if err != nil {
		log.Println("Error marshal", err.Error())
		return err
	}

	err = json.Unmarshal(paramByte, &result)
	if err != nil {
		log.Println("Error unmarshal", err.Error())
		return err
	}

	return nil
}

func Log(title string, data interface{}) {
	header := strings.Repeat("=", 10) + " " + title + " " + strings.Repeat("=", 10)
	log.Println()
	log.Println(header)
	if data != nil {
		footer := strings.Repeat("=", 10) + strings.Repeat("=", len(title)) + strings.Repeat("=", 10)
		if reflect.ValueOf(data).Kind() == reflect.Map || reflect.ValueOf(data).Kind() == reflect.Struct {
			dataByte, err := json.MarshalIndent(data, "", "\t")
			if err != nil {
				fmt.Println("Error logger:", err)
			}
			log.Println(string(dataByte))
			log.Println(footer)
		}
	}
	log.Println()
}
