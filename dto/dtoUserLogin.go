package dto

type ReqLogin struct {
	Phone int `json:"phone" form:"phone" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
}