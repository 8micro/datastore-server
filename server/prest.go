package server

import "github.com/8micro/gounits/rand"
import jwt "github.com/dgrijalva/jwt-go"

import (
	"encoding/base64"
	"fmt"
	"hash/crc32"
	"net/http"
	"time"
)

//Jwt is exported
type JwtConfig struct {
	Secret  string
	Expired string
}

//PrestConfig is exported
type PrestConfig struct {
	DataBase     string
	Schema       string
	Hosts        []string
	Jwt 		 *JwtConfig
}

//PrestRequest is exported
type PrestRequest struct {
	PrestURL string
	Headers  http.Header
}

//NewPrestRequest is exported
func NewPrestRequest(prestConfig *PrestConfig) (*PrestRequest, error) {

	var (
		err error
		prestURL string
		prestSignature string
	)

	prestURL, err = accessPrestURL(prestConfig)
	if err != nil {
		return nil, err
	}

	if prestConfig.Jwt != nil {
		prestSignature, err = signingPrestToken(prestConfig.Jwt)
		if err != nil {
			return nil, err
		}
	}

	headers := requestHeaders(prestSignature)
	return &PrestRequest{
		PrestURL: prestURL,
		Headers:  headers,
	}, nil
}

func requestHeaders(signature string) http.Header {

	tick := time.Now().UnixNano() / int64(time.Millisecond)
	timeStamp := fmt.Sprintf("DATASTORE-SERVER#%d", tick)
	headers := map[string][]string{
		"Content-Type": []string{`application/json;charset=utf-8`} ,
		"X-8MircoService": []string{base64.StdEncoding.EncodeToString([]byte(timeStamp))},
	}

	if signature != "" {
		headers["Authorization"] = []string{`Bearer ` + signature}
	}
	return headers
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
		key := rand.UUID(true)
		index = crc32.ChecksumIEEE([]byte(key)) % (uint32)(size)
	}
	return hosts[index]
}

func signingPrestToken(jwtConfig *JwtConfig) (string, error) {

	expired, err := time.ParseDuration(jwtConfig.Expired)
	if err != nil {
		expired, _ = time.ParseDuration("120s")
	}

	key := []byte(jwtConfig.Secret)
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
