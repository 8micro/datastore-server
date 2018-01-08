package response

//CreateUserResourceResponse is exported
type CreateUserResourceResponse struct {
	UserId       int    `json:"userid"`
	FileUniqueId string `json:"fileuniqueid"`
}
