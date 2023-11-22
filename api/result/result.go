package result

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

const (
	Success = 0

	InternalError = 1000
	DBError       = 1001
	RedisError    = 1002
	InvalidPara   = 1101
	URLInvalid    = 1102

	UserEmailRegistered  = 1201
	UserEmailCodeInvalid = 1202
	UserEmailSendFailed  = 1203
	UserTokenInvalid     = 1204
	UserTokenExpired     = 1205
	UserNotFound         = 1206

	UserEmailPasswordIncorrect = 1301

	SpaceURLInvalid     = 2001
	SpaceWalletNotMatch = 2002
)
