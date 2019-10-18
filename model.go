package main

import (
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"net/url"
)

type Config struct {
	Proto        string `json:"proto"`
	Host         string `json:"host"`
	Port         string `json:"port"`
	ClientID     string `json:"client-id"`
	ClientSecret string `json:"client-secret"`
	Scope        string `json:"scope"`
	RedirectUrl  string `json:"redirect-url"`
}

type ConnectConfig struct {
	Protocol string
	Host     string
	Port     string
}

type AuthCodeRequest struct {
	ResponseType string
	State        string
	ClientID     string
	Scope        string
	LoginHint    string
	RedirectUrl  string
}

type TokenRequest struct {
	GrantType    string
	Code         string
	ClientID     string
	ClientSecret string
	RedirectUri  string
}

func NewConnectionConfig(config *Config) *ConnectConfig {
	return &ConnectConfig{
		Protocol: config.Proto,
		Host:     config.Host,
		Port:     config.Port,
	}
}

func NewAuthRequest(config *Config) *AuthCodeRequest {
	return &AuthCodeRequest{
		ResponseType: "code",
		State:        "123e4567-e89b-12d3-a456-426655440000",
		ClientID:     config.ClientID,
		Scope:        config.Scope,
		LoginHint:    "123e4567-e89b-12d3-a456-426655440000",
		RedirectUrl:  config.RedirectUrl,
	}
}

func NewTokenRequest(config *Config, code string) *TokenRequest {
	return &TokenRequest{
		GrantType:    "authorization_code",
		Code:         code,
		ClientID:     config.ClientID,
		ClientSecret: sha384Base64Encoding(config.ClientSecret),
		RedirectUri:  config.RedirectUrl,
	}

}

func (conn *ConnectConfig) ParseUrl() string {
	return conn.Protocol + "://" + conn.Host + ":" + conn.Port
}

func (authReq *AuthCodeRequest) ParseUrlParameters(req *http.Request) {
	q := req.URL.Query()
	q.Add("response_type", authReq.ResponseType)
	q.Add("state", authReq.State)
	q.Add("client_id", authReq.ClientID)
	q.Add("scope", authReq.Scope)
	q.Add("login_hint", authReq.LoginHint)
	q.Add("redirect_uri", authReq.RedirectUrl)
	req.URL.RawQuery = q.Encode()
}

func (tokenReq *TokenRequest) ParsePostData() *url.Values {

	data := url.Values{}
	data.Set("grant_type", tokenReq.GrantType)
	data.Set("code", tokenReq.Code)
	data.Set("client_id", tokenReq.ClientID)
	data.Set("client_secret", tokenReq.ClientSecret)
	data.Set("redirect_uri", tokenReq.RedirectUri)
	return &data
}

func sha384Base64Encoding(s string) string {
	h := sha512.New384()
	h.Write([]byte(s))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))

}
