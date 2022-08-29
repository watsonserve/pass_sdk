package pass_sdk

import (
    "net/http"
)

const (
    PASSPORT_ORIGIN = "https://passport.watsonserve.com"
    SESSION_USER_KEY = "user"
)

type SrvInfo struct {
    AuthPathname string
    AppId string
    Scheme string
    Host string
    Secret string
}

// 业务访问对象
//
// Get: 从业务自己的数据库读取用户信息
//
// Save: 向业务自己的数据库保存用户信息
//
// Error: 输出错误页面
type BizAO interface {
    Get(id string) (*UserData, error)
    Save(id string, userData *UserData) error
    Error(res http.ResponseWriter, code int, msg string)
}

type UserData struct {
    UserId string
    Name   string
    Avatar string
}

type userService struct {
    bao BizAO
    app string
    secret string
}

type authMgr struct {
    userService
    authAddr string
    scheme string
    host string
}
