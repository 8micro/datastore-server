package api

import "github.com/8micro/datastore-server/api/response"
import "github.com/8micro/datastore-server/server"

import (
	"net/http"
)

func postUserResource(c *Context) error {

	resp := &response.ResponseImpl{}
	request := ResolveCreateUserResourceRequest(c)
	if request == nil {
		resp.SetContent(response.ErrRequestResolveInvaild.Error())
		return c.JSON(http.StatusBadRequest, resp)
	}

	dataServer := c.Get("DataServer").(*server.DataServer)
	FileUniqueId, err := dataServer.CreateUserResource(request)
	if err != nil {
		resp.SetContent(err.Error())
		return c.JSON(http.StatusInternalServerError, resp)
	}

	respData := &response.CreateUserResourceResponse{
		UserId:       request.UserId,
		FileUniqueId: FileUniqueId,
	}
	resp.SetContent(response.ErrRequestSuccessed.Error())
	resp.SetData(respData)
	return c.JSON(http.StatusOK, resp)
}
