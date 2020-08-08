package datamodels

//LoginReq ask for args
type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//result of response
type LoginResult struct {
	Token          string `json:"token"`
	Nickname       string `json:"nickname"`
	ProfilePicture string `json:"profilePicture"`
}
