package result

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data,omitempty"`
}

const (
	Success = 0

	UserTokenExpired = 1401
	UserTokenInvalid = 1402
)
