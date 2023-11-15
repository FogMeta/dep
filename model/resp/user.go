package resp

type UserResp struct {
	UID    int    `json:"uid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Wallet string `json:"wallet"`
	Type   int    `json:"type"`
	Token  string `json:"token"`
}

type UserInfoResp struct {
	UID    int    `json:"uid"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Wallet string `json:"wallet"`
}
