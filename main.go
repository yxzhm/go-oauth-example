package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"reflect"
	"strings"
)

var config *Config
var connectionConfig *ConnectConfig
var authCodeRequest *AuthCodeRequest

const HydraAuthUri = "/oauth2/auth"
const HydraTokenUri = "/oauth2/token"

func init() {
	//Read config
	jsonFile, err := os.Open("./cert.json")
	if err != nil {
		Error.Fatalln("Can't open config file")
	}
	defer jsonFile.Close()
	if err = json.NewDecoder(jsonFile).Decode(&config); err != nil {
		Error.Fatalln("Can't parse config file")
	}

	//Init struct
	connectionConfig = NewConnectionConfig(config)
	authCodeRequest = NewAuthRequest(config)
}

func main() {
	req, _ := http.NewRequest("GET", connectionConfig.ParseUrl()+HydraAuthUri, nil)
	authCodeRequest.ParseUrlParameters(req)
	Info.Println(req.URL.String())

	cook, _ := cookiejar.New(nil)
	client := &http.Client{}
	client.Jar = cook
	defer client.CloseIdleConnections()

	// Get code thru http get
	_, err := client.Get(req.URL.String())

	// There is err due to the callback url is faked
	if err != nil && reflect.TypeOf(err).String() == "*url.Error" {
		e := err.(*url.Error)
		if strings.HasPrefix(e.URL, authCodeRequest.RedirectUrl+"?code=") {
			rawData, _ := url.Parse(e.URL)
			data, _ := url.ParseQuery(rawData.RawQuery)
			code := data["code"][0]
			Info.Println(code)

			// Get JWT token thru http post
			tokenReq := NewTokenRequest(config, code)
			res, _ := client.Post(connectionConfig.ParseUrl()+HydraTokenUri, "application/x-www-form-urlencoded", strings.NewReader(tokenReq.ParsePostData().Encode()))
			body, _ := ioutil.ReadAll(res.Body)
			Info.Println(string(body))

		}
	}
}
