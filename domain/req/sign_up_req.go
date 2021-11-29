package req

type SignUpReq struct {
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"pwd,required"`
}
