package pass_sdk

import (
	"net/http"
)

const (
	PASSPORT_ORIGIN = "https://passport.watsonserve.com"
)

type SrvInfo struct {
	WebAuthPathname string
	CliAuthPathname string
	AppId           string
	Scheme          string
	Host            string
	Secret          string
}

type UserData struct {
	OpenId string `json:"open_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type BizAO interface {
	IsCheckedIn(res http.ResponseWriter, req *http.Request) bool
	Error(res http.ResponseWriter, req *http.Request, code int, explain string)
	User(res http.ResponseWriter, req *http.Request, usr *UserData, rd string)
}

type authMgr struct {
	SrvInfo
	bao BizAO
}

type stdJsonResp struct {
	Status bool        `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

// type Token_t struct {
// 	TokenType    string `json:"token_type"`
// 	ExpiresIn    int    `json:"expires_in"`
// 	AccessToken  string `json:"access_token"`
// 	RefreshToken string `json:"refresh_token"`
// 	Scope        string `json:"scope"`
// }
