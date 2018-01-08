package server

import "github.com/8micro/datastore-server/api/request"
import "github.com/8micro/datastore-server/types"
import "github.com/8micro/gounits/httpx"
import "github.com/8micro/gounits/rand"

import (
	"fmt"
	"net"
	"net/http"
	"time"
)

//PrestConfig is exported
type PrestConfig struct {
	DataBase     string
	Schema       string
	Hosts        []string
	Secret       string
	TokenExpired string
}

//DataServer is exported
type DataServer struct {
	Key         string
	PrestConfig *PrestConfig
	client      *httpx.HttpClient
}

//NewDataServer is exported
func NewDataServer(key string, prestConfig *PrestConfig) *DataServer {

	dataServer := &DataServer{
		Key:         key,
		PrestConfig: prestConfig,
		client: httpx.NewClient().
			SetTransport(&http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 60 * time.Second,
				}).DialContext,
				DisableKeepAlives:     false,
				MaxIdleConns:          25,
				MaxIdleConnsPerHost:   25,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   http.DefaultTransport.(*http.Transport).TLSHandshakeTimeout,
				ExpectContinueTimeout: http.DefaultTransport.(*http.Transport).ExpectContinueTimeout,
			}),
	}
	return dataServer
}

//CreateUserResource is exported
func (server *DataServer) CreateUserResource(r *request.CreateUserResourceRequest) (string, error) {

	//验证AuthHeader
	//检查用户ID是否有效存在
	createUserResource := &types.CreateUserResource{
		UserId:       r.UserId,
		FileUniqueId: rand.UUID(false),
		FileName:     r.FileName,
		Duration:     r.Duration,
		Rate:         r.Rate,
		Resolution:   r.Resolution,
		VerifyCode:   r.VerifyCode,
		UploadAt:     r.UploadAt,
		ExpiredAt:    0,
		State:        0,
	}

	prestRequest, err := NewPrestRequest(server.PrestConfig)
	if err != nil {
		return "", err
	}

	respData, err := server.client.PostJSON(prestRequest.PrestURL+"/resources", nil, createUserResource, prestRequest.Headers)
	if err != nil {
		return "", err
	}

	defer respData.Close()
	statusCode := respData.StatusCode()
	if statusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("HTTP prest CreateUserResource failure %d.", statusCode)
	}
	return createUserResource.FileUniqueId, nil
}
