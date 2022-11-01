package pass_sdk

import (
	"net/http"
	"net/url"

	"github.com/watsonserve/goengine"
)

func Goback(bao BizAO, res http.ResponseWriter, req *http.Request) {
	// 解析重定向地址
	rd := req.URL.Query().Get("rd")
	rd, _ = url.QueryUnescape(rd)
	// 重定向地址出局
	if '/' != rd[0] {
		bao.Error(res, req, 503, "Redirect Out")
		return
	}

	res.Header().Set("Location", rd)
	res.WriteHeader(302)
}

func DefaultError(res http.ResponseWriter, code int, msg string) {
	header := res.Header()
	header.Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)
	ret := &stdJSON{
		Status: false,
		Msg:    msg,
	}
	res.Write(ret.Chars())
}

// 从session检出用户id
func DefaultGet(res http.ResponseWriter, req *http.Request) map[string]interface{} {
	var session *goengine.Session
	ctx := req.Context()
	session = ctx.Value("session").(*goengine.Session)
	// 检出数据
	user := session.Get(SESSION_USER_KEY)
	if nil == user {
		return nil
	}
	return user.(map[string]interface{})
}

// 从session检出用户id
func DefaultSave(res http.ResponseWriter, req *http.Request, userData *UserData) error {
	ctx := req.Context()
	session := ctx.Value("session").(*goengine.Session)
	session.Set(SESSION_USER_KEY, userData)
	session.Save(res, 0)
	return nil
}

func DefaultScope(bao BizAO, res http.ResponseWriter, req *http.Request, token *Token_t) {
	// 获取用户信息
	userData, err := GetUserData(bao.App(), bao.Secret(), token.AccessToken)
	if nil == err {
		err = bao.Save(res, req, userData)
	}
	if nil != err {
		bao.Error(res, req, 503, err.Error())
		return
	}
	Goback(bao, res, req)
}
