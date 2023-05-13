package system

import (
	"bytes"
	"crypto"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	b64 "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nfnt/resize"
)

/* Encryption */
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func EncryptBase64(data string)string{
	encrypt := b64.StdEncoding.EncodeToString([]byte(data))
	return encrypt
}

func DecodeBase64(data string)string{
	decode, _ := b64.StdEncoding.DecodeString(data)

	return string(decode)
}

func HashSha256(data string)string{
	hash := sha256.New()
	hash.Write([]byte(data))
    bs := hash.Sum(nil)

    hashing := fmt.Sprintf("%x\n", bs)
	return hashing
}

func VerifySignature(PublicKey, Data, Signature string)bool{
	block, _ := pem.Decode([]byte(PublicKey))
    if block == nil {
        fmt.Println("Invalid PEM Block")
        return false
    }

    key, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        fmt.Println(err)
        return false
    }

    pubKey := key.(*rsa.PublicKey)

    signature, err := base64.StdEncoding.DecodeString(Signature)
    if err != nil {
        fmt.Println(err)
        return false
    }

	sha256 := sha256.New()
	sha256.Write([]byte(Data))
	hash := sha256.Sum(nil)

    err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature)
    if err != nil {
        fmt.Println(err)
        return false
    }

    fmt.Println("Successfully verified message with signature and public key")
    return true
}

/* Image File */
func getExtention(base64 string) string {
	png := strings.Contains(base64, "image/png")
	jpeg := strings.Contains(base64, "image/jpeg")
	gif := strings.Contains(base64, "image/gif")

	ex := ".jpg"

	if png {
		ex = ".png"
	} else if jpeg {
		ex = ".jpeg"
	} else if gif {
		ex = ".gif"
	}

	return ex
}

func Compress(nmfile string, height string, width string) error {
	file, err := os.Open(nmfile)
	if err != nil {
		log.Fatal(err)
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	hg, err := strconv.Atoi(height)
	uhg := uint(hg)

	wd, err := strconv.Atoi(width)
	uwd := uint(wd)

	// (width, height, input file, kernel sampling)
	m := resize.Resize(uwd, uhg, img, resize.Lanczos3)

	out, err := os.Create(nmfile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	fmt.Println("Ini Masuk Compress")
	return jpeg.Encode(out, m, nil)
}

func CheckFormValidation(r map[string]interface{}) (res bool) {
	var v = 0
	for _, val := range r {
		if val == nil || val == "" {
			v += 1
		}
	}
	return v > 0
}

func UploadBase64ToImg(param string, path string, name string) (data interface{}, err error) {
	// Read form fields
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)

		if err != nil {
			fmt.Println("tidak bisa create")
		}
	}

	if param != "" {
		b64data := param[strings.IndexByte(param, ',')+1:]
		dec, err := b64.StdEncoding.DecodeString(b64data)
		if err != nil {
			return data, err
		}

		time := time.Now().Unix()
		//fmt.Println(nm, "alex")s
		nmfile := name + strconv.Itoa(int(time)) + getExtention(param)
		f, err := os.Create(path + nmfile)
		if err != nil {
			return data, err
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			return data, err
		}
		if err := f.Sync(); err != nil {
			return data, err
		}
		fmt.Println(nmfile)
		Compress(path+nmfile, "0", "700")
		fmt.Println("Berhasil Compress")
		return nmfile, err
	}
	fmt.Println("alex2")
	return data, err
}

func IsImageFile(typefile string) (valid bool) {
	return strings.Contains(typefile, "image")
}

func IntToCharStr(i int) string {
	return string('A' - 1 + i)
}

func UploadFile(c echo.Context, param string, path string) (data interface{}, err error) {
	// Read form fields
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0777)

		if err != nil {
			fmt.Println("tidak bisa create")
		}
	}
	file, err := c.FormFile(param)
	if err != nil {
		return nil, err
	}
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(path + file.Filename)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return nil, err
	}

	return file.Filename, err
}

/* Proses Request */
func RequestGetParams(e echo.Context) (res map[string]interface{}) {
	// e.MultipartForm()
	data, _ := e.MultipartForm()
	fromParam := make(map[string]interface{})
	fromJson := make(map[string]interface{})
	r := e.Request()
	if err := r.ParseForm(); err != nil {
	}
	jsn, _ := json.Marshal(r.Form)
	if err := json.Unmarshal(jsn, &fromParam); err != nil {
	}
	if err := json.NewDecoder(e.Request().Body).Decode(&fromJson); err != nil {
	}

	if len(fromJson) > 0 {
		res = fromJson
	} else if len(fromParam) > 0 {
		var data = make(map[string]interface{})
		for k, val := range fromParam {
			var afterDec []interface{}
			v, _ := json.Marshal(val)
			if err := json.Unmarshal(v, &afterDec); err != nil {
			}
			data[k] = afterDec[0]
		}
		res = data
	}

	if res == nil {
		res = map[string]interface{}{}
	}
	if data != nil && len(data.File) > 0 {
		for k, vals := range data.File {
			res[k] = vals
		}
	}

	return res
}

func ReqTrimSpace(req map[string]interface{}) (res map[string]interface{}) {
	// for k, val := range req {
	// 	if reflect.TypeOf(val).String() == "string" {
	// 		req[k] = strings.TrimSpace(val.(string))
	// 		fmt.Println(req[k], "----- alex")
	// 	}
	// 	// valstr := val.(string)
	// 	// val = strings.TrimSpace(valstr)
	// 	// req[k] = val
	// }
	return req
}

func MultiReqTrimSpace(req []interface{}) (res []interface{}) {
	// for _, val := range req {
	// 	data_req := val.(map[string]interface{})
	// 	for k, v := range data_req {
	// 		if reflect.TypeOf(v).String() == "string" {
	// 			data_req[k] = strings.TrimSpace(v.(string))
	// 		}
	// 	}
	// }
	return req
}

func HandleNil(req map[string]interface{}) (rsp map[string]interface{}) {
	for key, val := range req {
		if val == nil || val == "" {
			delete(req, key)
		}
	}
	return req
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

/* CURL / HTTP Request */
func CurlHeader(method string, link string, params string, header map[string]string, formbody interface{}) (res map[string]interface{}, err error) {
	// param := url.Values{}
	// for k, v := range formbody{
	// 	param.Add(k, v.(string))
	// }
	// var payload = bytes.NewBufferString(param.Encode())
	var json_data, _ = json.Marshal(formbody)
	var payload = bytes.NewBuffer(json_data)

	if params != "" {
		link = link + params
	}

	req, err := http.NewRequest(method, link, payload)

	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if body != nil {
		json.Unmarshal(body, &res)
	}
	return res, nil
}

func CurlBasicAuthJson(method, link string, body map[string]interface{}) (Status int64, rsp map[string]interface{}) {
	var json_data, _ = json.Marshal(body)
	var payload = bytes.NewBuffer(json_data)

	username := os.Getenv("USERNAME_VALIDATION")
	password := os.Getenv("PASSWORD_VALIDATION")

	req, err := http.NewRequest(method, link, payload)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(username, password)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	var response map[string]interface{}
	json.Unmarshal([]byte(resBody), &response)
	if response == nil {
		Status = 400
	} else {
		Status = int64(response["Status"].(float64))
	}

	if Status == 200 {
		rsp = response["Data"].(map[string]interface{})
	}

	return
}

/* Convert Variable */
func ToBulan(tipe string, lng string, i int, month string) (bulan map[string]string) {
	if tipe == "string" {
		if lng == "id" {
			switch month {
			case "Januari":
				bulan = map[string]string{
					"id": "Januari",
					"en": "January",
				}
				return bulan
			case "Februari":
				bulan = map[string]string{
					"id": "Februari",
					"en": "February",
				}
				return bulan
			case "Maret":
				bulan = map[string]string{
					"id": "Maret",
					"en": "March",
				}
				return bulan
			case "April":
				bulan = map[string]string{
					"id": "April",
					"en": "April",
				}
				return bulan
			case "Mei":
				bulan = map[string]string{
					"id": "Mei",
					"en": "May",
				}
				return bulan
			case "Juni":
				bulan = map[string]string{
					"id": "Juni",
					"en": "June",
				}
				return bulan
			case "Juli":
				bulan = map[string]string{
					"id": "Juli",
					"en": "July",
				}
				return bulan
			case "Agustus":
				bulan = map[string]string{
					"id": "Agustus",
					"en": "August",
				}
				return bulan
			case "September":
				bulan = map[string]string{
					"id": "September",
					"en": "September",
				}
				return bulan
			case "Oktober":
				bulan = map[string]string{
					"id": "Oktober",
					"en": "October",
				}
				return bulan
			case "November":
				bulan = map[string]string{
					"id": "November",
					"en": "November",
				}
				return bulan
			case "Desember":
				bulan = map[string]string{
					"id": "Desember",
					"en": "December",
				}
				return bulan

			default:
				bulan = map[string]string{
					"id": "Bulan Tidak Ditemukan",
					"en": "Month not found",
				}
				return bulan
			}
		} else if lng == "en" {

		}
	} else if tipe == "number" {
		switch i {
		case 1:
			bulan = map[string]string{
				"id": "Januari",
				"en": "January",
			}
			return bulan
		case 2:
			bulan = map[string]string{
				"id": "Februari",
				"en": "February",
			}
			return bulan
		case 3:
			bulan = map[string]string{
				"id": "Maret",
				"en": "March",
			}
			return bulan
		case 4:

			bulan = map[string]string{
				"id": "April",
				"en": "April",
			}
			return bulan
		case 5:
			bulan = map[string]string{
				"id": "Mei",
				"en": "May",
			}
			return bulan
		case 6:
			bulan = map[string]string{
				"id": "Juni",
				"en": "June",
			}
			return bulan
		case 7:
			bulan = map[string]string{
				"id": "Juli",
				"en": "July",
			}
			return bulan
		case 8:
			bulan = map[string]string{
				"id": "Agustus",
				"en": "August",
			}
			return bulan
		case 9:

			bulan = map[string]string{
				"id": "September",
				"en": "September",
			}
			return bulan
		case 10:
			bulan = map[string]string{
				"id": "Oktober",
				"en": "October",
			}
			return bulan
		case 11:
			bulan = map[string]string{
				"id": "November",
				"en": "November",
			}
			return bulan
		case 12:
			bulan = map[string]string{
				"id": "Desember",
				"en": "December",
			}
			return bulan

		default:
			bulan = map[string]string{
				"id": "Bulan Tidak Ditemukan",
				"en": "Month not found",
			}
			return bulan
		}
	}
	return bulan
}

func ToTerbilang(num int) string {
	var s string
	satuan := [12]string{"", "satu", "dua", "tiga", "empat", "lima", "enam", "tujuh", "delapan", "sembilan", "sepuluh", "sebelas"}
	if num < 12 {
		s = satuan[num]
	} else if num < 20 {
		s = fmt.Sprintf("%s belas", ToTerbilang(num-10))
	} else if num < 100 {
		s = fmt.Sprintf("%s puluh %s", ToTerbilang(num/10), ToTerbilang(num%10))
	} else if num < 200 { // ratus
		s = fmt.Sprintf("seratus %s", ToTerbilang(num-100))
	} else if num < 1000 {
		s = fmt.Sprintf("%s ratus %s", ToTerbilang(num/100), ToTerbilang(num%100))
	} else if num < 2000 { // ribu
		s = fmt.Sprintf("seribu %s", ToTerbilang(num-1000))
	} else if num < 1000000 {
		s = fmt.Sprintf("%s ribu %s", ToTerbilang(num/1000), ToTerbilang(num%1000))
	} else if num < 2000000 { // juta
		s = fmt.Sprintf("satu juta %s", ToTerbilang(num-1000000))
	} else if num < 1000000000 {
		s = fmt.Sprintf("%s juta %s", ToTerbilang(num/1000000), ToTerbilang(num%1000000))
	} else if num < 2000000000 { // milyar
		s = fmt.Sprintf("satu milyar %s", ToTerbilang(num-1000000000))
	} else if num < 1000000000000 {
		s = fmt.Sprintf("%s milyar %s", ToTerbilang(num/1000000000), ToTerbilang(num%1000000000))
	} else if num < 2000000000000 { // triliun
		s = fmt.Sprintf("satu triliun %s", ToTerbilang(num-1000000000000))
	} else if num < 1000000000000000 {
		s = fmt.Sprintf("%s triliun %s", ToTerbilang(num/1000000000000), ToTerbilang(num%1000000000000))
	}
	return strings.TrimSpace(s)
}

func pow(i int, p int) int {
	return int(math.Pow(float64(i), float64(p)))
}

func OrdinalNumber(n int) string {
	to19 := []string{"first", "second", "third", "fourth", "fiveth", "sixth", "seventh", "eighth", "ninth", "tenth", "eleventh", "twelveth",
		"thirteenth", "fourteenth", "fifteenth", "sixteenth", "seventeenth", "eighteenth", "nineteenth"}

	tens := []string{"twentieth", "thirtieth", "fortieth", "fiftieth", "sixtieth", "seventieth", "eightieth", "ninetieth"}
	if n == 0 {
		return ""
	}
	if n < 20 {
		return to19[n-1]
	}
	if n < 100 {
		return tens[n/10-2] + " " + OrdinalNumber(n%10)
	}
	if n < 1000 {
		return to19[n/100-1] + " hundred " + OrdinalNumber(n%100)
	}

	for idx, w := range []string{"thousand", "million", "billion"} {
		p := idx + 1
		if n < pow(1000, (p+1)) {
			return OrdinalNumber(n/pow(1000, p)) + " " + w + " " + OrdinalNumber(n%pow(1000, p))
		}
	}

	return "error"
}

func RangeDate(start, end time.Time) func() time.Time {
	y, m, d := start.Date()
	start = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
	y, m, d = end.Date()
	end = time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	return func() time.Time {
		if start.After(end) {
			return time.Time{}
		}
		date := start
		start = start.AddDate(0, 0, 1)
		return date
	}
}

func CreatedTimeNow(req map[string]interface{}) (rsp map[string]interface{}) {
	req["CreatedAt"] = time.Now().Format("2006-01-02 15:04:05")
	req["UpdatedAt"] = time.Now().Format("2006-01-02 15:04:05")
	//HandleNil(req)
	return req
}

func UpdatedTimeNow(req map[string]interface{}) (rsp map[string]interface{}) {
	req["UpdatedAt"] = time.Now().Format("2006-01-02 15:04:05")
	//HandleNil(req)
	return req
}