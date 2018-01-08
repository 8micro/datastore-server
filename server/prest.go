package server

import jwt "github.com/dgrijalva/jwt-go"

import (
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"net/http"
	"time"
)

//PrestRequest is exported
type PrestRequest struct {
	PrestURL string
	Headers  http.Header
}

//NewPrestRequest is exported
func NewPrestRequest(prestConfig *PrestConfig) (*PrestRequest, error) {

	prestURL, err := accessPrestURL(prestConfig)
	if err != nil {
		return nil, err
	}

	prestToken, err := signingPrestToken(prestConfig)
	if err != nil {
		return nil, err
	}

	headers := requestHeaders(prestToken)
	return &PrestRequest{
		PrestURL: prestURL,
		Headers:  headers,
	}, nil
}

func requestHeaders(token string) http.Header {

	tick := time.Now().UnixNano() / int64(time.Millisecond)
	timeStamp := fmt.Sprintf("DATASTORESERVER#%d", tick)
	return map[string][]string{
		"Content-Type":    []string{`application/json;charset=utf-8`},
		"Authorization":   []string{`Bearer ` + token},
		"X-8MircoService": []string{base64.StdEncoding.EncodeToString([]byte(timeStamp))},
	}
}

func accessPrestURL(prestConfig *PrestConfig) (string, error) {

	host := hashHost(prestConfig.Hosts)
	if host == "" {
		return "", fmt.Errorf("prest host invalid.")
	}
	prestURL := "http://" + host + "/" + prestConfig.DataBase + "/" + prestConfig.Schema
	return prestURL, nil
}

func hashHost(hosts []string) string {

	size := len(hosts)
	if size == 0 {
		return ""
	}

	index := (uint32)(0)
	if size > 1 {
		index = crc32.ChecksumIEEE([]byte(time.Now().String())) % (uint32)(size)
	}
	return hosts[index]
}

func signingPrestToken(prestConfig *PrestConfig) (string, error) {

	expired, err := time.ParseDuration(prestConfig.TokenExpired)
	if err != nil {
		expired, _ = time.ParseDuration("120s")
	}

	key := []byte(prestConfig.Secret)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": time.Now().Add(expired).Unix(),
		"nbf": time.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC).Unix(),
	})

	signature, err := token.SignedString(key)
	if err != nil {
		return "", fmt.Errorf("prest access signing token failure.")
	}
	return signature, nil
}
