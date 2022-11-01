package pass_sdk

import (
	"encoding/json"
	"net/http"
)

const (
	PASSPORT_ORIGIN  = "https://passport.watsonserve.com"
	SESSION_USER_KEY = "user"
)

type SrvInfo struct {
	AuthPathname string
	AppId        string
	Scheme       string
	Host         string
	Secret       string
}

// 业务访问对象
//
// Get: 从业务自己的数据库读取用户信息
//
// Save: 向业务自己的数据库保存用户信息
//
// Error: 输出错误页面
type BizAO interface {
	App() string
	Secret() string
	Get(res http.ResponseWriter, req *http.Request) map[string]interface{}
	Save(res http.ResponseWriter, req *http.Request, userData *UserData) error
	Error(res http.ResponseWriter, req *http.Request, code int, msg string)
	Scope(res http.ResponseWriter, req *http.Request, tokenResp *Token_t)
}

type UserData struct {
	OpenId string
	Name   string
	Avatar string
}

type authMgr struct {
	bao      BizAO
	app      string
	secret   string
	authAddr string
	scheme   string
	host     string
}

type Token_t struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type stdJSON struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func (std *stdJSON) Chars() []byte {
	foo, _ := json.Marshal(std)
	return foo
}
