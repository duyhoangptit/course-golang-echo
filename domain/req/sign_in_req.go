package req

type SignInReq struct {
	Email    string `json:"email" validate:"email,required"`
	Password string `json:"password" validate:"pwd,required"`
}
