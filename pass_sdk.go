// 提供通用接入passport方法
//
// AUTHOR: JamesWatson (c) 2019 watsonserve.com
package pass_sdk

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/watsonserve/goengine"
	"github.com/watsonserve/goutils"
)

// 绑定pass_sdk 授权管理器
//
// srvInfo.app(aka client_id): came from authorize server
// bao: 业务访问对象, 如果bao = nil, 使用默认bao
func BindAuthMgr(srvInfo *SrvInfo, bao BizAO, route *goengine.HttpRoute) error {
	if nil == bao {
		return errors.New("bao is required")
	}
	am := &authMgr{
		bao:      bao,
		app:      srvInfo.AppId,
		secret:   srvInfo.Secret,
		authAddr: srvInfo.AuthPathname,
		scheme:   srvInfo.Scheme,
		host:     srvInfo.Host,
	}
	route.Use(am.pageFilter)
	route.Set(srvInfo.AuthPathname, am.auth)
	return nil
}

func (am *authMgr) auth(res http.ResponseWriter, req *http.Request) {
	// 校验来源
	if "GET" != req.Method || !chkReferer(req, PASSPORT_ORIGIN) {
		am.bao.Error(res, req, 405, "Method Not Allowed")
		return
	}

	query := req.URL.Query()
	authCode := query.Get("code")
	rd := query.Get("rd")
	// state := query.Get("state")

	// if goutils.MD5(token+redirect+stamp) != state {

	// }

	redirect := getAuthAddr(am.scheme, am.host, am.authAddr, rd)
	tokenResp, err := loadToken(am.app, am.secret, authCode, redirect)
	if nil != err {
		am.bao.Error(res, req, 400, err.Error())
		return
	}

	am.bao.Scope(res, req, tokenResp)
}

func (am *authMgr) GetPassportUrl(uri *url.URL, scope string) string {
	// 随机字符串
	salt := goutils.RandomString(16)
	token := goutils.MD5(am.app + salt + am.secret)
	stamp := fmt.Sprintf("%d", goutils.Now())
	// passport成功后回跳地址
	redirect := getAuthAddr(am.scheme, am.host, am.authAddr, cutUri(uri))
	// 组织参数
	passParams := url.Values{}
	passParams.Set("response_type", "code")
	passParams.Set("client_id", am.app)
	passParams.Set("redirect_uri", redirect)
	passParams.Set("scope", scope)
	// auth server return this msg without any changed.
	passParams.Set("state", goutils.MD5(token+redirect+stamp))
	// jump
	return fmt.Sprintf("%s/?%s", PASSPORT_ORIGIN, passParams.Encode())
}

func (am *authMgr) pageFilter(res http.ResponseWriter, req *http.Request) bool {
	// 授权接口地址
	if req.URL.Path == am.authAddr {
		return true
	}

	// 已登录
	if nil != am.bao.Get(res, req) {
		return true
	}

	// 未登录
	res.Header().Set("Location", am.GetPassportUrl(req.URL, "user_info"))
	// jump
	res.WriteHeader(302)
	return false
}
