package api

import "github.com/8micro/datastore-server/api/request"

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

//ResolveCreateUserResourceRequest is exported
func ResolveCreateUserResourceRequest(c *Context) *request.CreateUserResourceRequest {

	buf, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return nil
	}

	request := &request.CreateUserResourceRequest{}
	if err := json.NewDecoder(bytes.NewReader(buf)).Decode(request); err != nil {
		return nil
	}
	return request
}

/*


//ResolveJobsAllocDataRequest is exported
func ResolveJobsAllocDataRequest(c *Context) string {

	vars := mux.Vars(c.request)
	runtime := strings.TrimSpace(vars["runtime"])
	if len(runtime) == 0 {
		return ""
	}
	return runtime
}

//ResolveServersRequest is exported
func ResolveServersRequest(c *Context) string {

	vars := mux.Vars(c.request)
	runtime := strings.TrimSpace(vars["runtime"])
	if len(runtime) == 0 {
		return ""
	}
	return runtime
}

//ResolveJobActionRequest is exported
func ResolveJobActionRequest(c *Context) *JobActionRequest {

	buf, err := ioutil.ReadAll(c.request.Body)
	if err != nil {
		return nil
	}

	request := &JobActionRequest{}
	if err := json.NewDecoder(bytes.NewReader(buf)).Decode(request); err != nil {
		return nil
	}

	request.Context = c
	return request
}

//ResolveMessageRequest is exported
func ResolveMessageRequest(c *Context) *MessageRequest {

	buf, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return nil
	}

	msgHeader := &models.MsgHeader{}
	reader := bytes.NewReader(buf)
	if err := json.NewDecoder(reader).Decode(msgHeader); err != nil {
		return nil
	}

	if _, err := reader.Seek(0, 0); err != nil {
		return nil
	}

	return &MessageRequest{
		Header:  msgHeader,
		Reader:  reader,
		Context: c,
	}
}
*/
