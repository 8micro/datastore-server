package types

//CreateUserResource is exported
type CreateUserResource struct {
	UserId       int    `json:"UserId"`
	FileUniqueId string `json:"FileUniqueId"`
	FileName     string `json:"FileName"`
	Duration     int    `json:"Duration"`
	Rate         int    `json:"Rate"`
	Resolution   string `json:"Resolution"`
	VerifyCode   string `json:"VerifyCode"`
	UploadAt     int64  `json:"UploadAt"`
	ExpiredAt    int64  `json:"ExpiredAt"`
	State        int    `json:"State"`
}
