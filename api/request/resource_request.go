package request

//CreateUserResourceRequest is exported
type CreateUserResourceRequest struct {
	UserId     int    `json:"UserId"`
	FileName   string `json:"FileName"`
	Duration   int    `json:"Duration"`
	Rate       int    `json:"Rate"`
	Resolution string `json:"Resolution"`
	VerifyCode string `json:"VerifyCode"`
	UploadAt   int64  `json:"UploadAt"`
}
