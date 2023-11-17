package req

type UserCreateReq struct {
	Email       string `json:"email"              binding:"required_without=WalletToken"`
	Password    string `json:"password"           binding:"required_with=Email"`
	AuthCode    string `json:"auth_code"          binding:"required_with=Email"`
	WalletToken string `json:"wallet_token"       binding:"required_without=Email"`
	Type        int    `json:"type"               binding:"oneof=1 2"`
}

type UserLoginReq struct {
	Email       string `json:"email"              binding:"required_without=WalletToken"`
	Password    string `json:"password"           binding:"required_with=Email"`
	WalletToken string `json:"wallet_token"       binding:"required_without=Email"`
	Type        int    `json:"type"               binding:"oneof=1 2"`
}

type UserResetPasswordReq struct {
	Email    string `json:"email"              binding:"required"`
	Password string `json:"password"           binding:"required"`
	AuthCode string `json:"auth_code"          binding:"required"`
}

type UserUpdatePasswordReq struct {
	Password    string `json:"password"     binding:"required"`
	OldPassword string `json:"old_password" binding:"required"`
}

type EmailReq struct {
	Email string `json:"email"       binding:"required"`
}
