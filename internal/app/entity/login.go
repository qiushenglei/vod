package entity

type LoginReq struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
