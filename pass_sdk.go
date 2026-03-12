// 提供通用接入passport方法
//
// AUTHOR: JamesWatson (c) 2019 watsonserve.com
package pass_sdk

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
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
		SrvInfo: *srvInfo,
		bao:     bao,
	}
	route.Use(am.pageFilter)
	route.Set(srvInfo.AuthPathname, am.auth)
	return nil
}

/**
 * @param {string} authAddr 授权路径，例如：/auth
 * @param {*url.URL} raw 当前路径，将被转换为：%2Fpathname%3Fsearch%23hash
 * @return string /auth?r=%2Fpathname%3Fsearch%23hash
 */
func (am *authMgr) getAuthAddr(salt, stamp, scope, rd string) string {
	redirect := url.URL{
		Scheme: am.Scheme,
		Host:   am.Host,
		Path:   am.AuthPathname,
	}
	q := redirect.Query()
	if "" != rd {
		q.Set("s", salt)
		q.Set("t", stamp)
		q.Set("c", scope)
		q.Set("rd", rd)
		redirect.RawQuery = q.Encode()
	}
	return redirect.String()
}

func (am *authMgr) auth(res http.ResponseWriter, req *http.Request) {
	// 校验来源
	if "GET" != req.Method || !chkReferer(req, PASSPORT_ORIGIN) {
		am.bao.Error(res, req, http.StatusMethodNotAllowed, "")
		return
	}

	query := req.URL.Query()
	authCode := query.Get("code")
	state := query.Get("state")
	rd := query.Get("rd")
	s := query.Get("s")
	t := query.Get("t")
	if "user_info" != query.Get("c") {
		am.bao.Error(res, req, http.StatusBadRequest, "Scope Not Allowed")
		return
	}

	// check state
	if stat, err := am.genState(s + t + "user_info" + rd); nil != err || stat != state {
		msg := ""
		if nil != err {
			msg = err.Error()
		}
		am.bao.Error(res, req, http.StatusBadRequest, msg)
		return
	}

	req.URL.Scheme = am.Scheme
	req.URL.Host = am.Host
	strJson, err := LoadByCode(am.AppId, am.Secret, authCode, "user-info")
	if nil != err {
		am.bao.Error(res, req, http.StatusBadRequest, err.Error())
		return
	}

	usr := &UserData{}
	err = json.Unmarshal(strJson, usr)
	if nil != err {
		am.bao.Error(res, req, http.StatusBadRequest, err.Error())
		return
	}
	am.bao.User(res, req, usr, rd)
}

func (am *authMgr) genState(txt string) (string, error) {
	key, err := base32.StdEncoding.DecodeString(am.Secret)
	if nil != err {
		return "", err
	}

	encoder := hmac.New(sha1.New, key)
	_, err = encoder.Write([]byte(txt))
	if nil != err {
		return "", err
	}

	hash := encoder.Sum(nil)
	// get offset
	offset := int(hash[len(hash)-1] & 0xf)
	n := (uint(hash[offset]&0x7f) << 24) | (uint(hash[offset+1]) << 16) | (uint(hash[offset+2]) << 8) | (uint(hash[offset+3]&0xff) << 0)

	format := fmt.Sprintf("%%0%dd", 16)
	code := fmt.Sprintf(format, n)
	rawLen := len(code)
	if 16 < rawLen {
		code = code[rawLen-16:]
	}
	return code, nil
}

func (am *authMgr) GetPassportUrl(uri *url.URL, scope string) (string, error) {
	// 随机字符串
	salt := goutils.RandomString(16)
	stamp := fmt.Sprintf("%d", goutils.Now())
	state, err := am.genState(salt + stamp + scope + cutUri(uri))
	if nil != err {
		return "", err
	}

	// passport成功后回跳地址
	redirect := am.getAuthAddr(salt, stamp, scope, cutUri(uri))
	// 组织参数
	passParams := url.Values{
		"response_type": []string{"code"},
		"client_id":     []string{am.AppId},
		"redirect_uri":  []string{redirect},
		"scope":         []string{scope},
		// auth server return this msg without any changed.
		"state": []string{state},
	}

	// jump
	return fmt.Sprintf("%s/?%s", PASSPORT_ORIGIN, passParams.Encode()), nil
}

func (am *authMgr) pageFilter(rsp http.ResponseWriter, req *http.Request) bool {
	// 授权接口地址 || 已登录
	pass := req.URL.Path == am.AuthPathname || am.bao.IsCheckedIn(rsp, req)

	// 未登录 jump
	if !pass {
		u, err := am.GetPassportUrl(req.URL, "user_info")
		if nil != err {
			am.bao.Error(rsp, req, http.StatusForbidden, "")
			return false
		}
		am.bao.Error(rsp, req, http.StatusForbidden, u)
	}

	return pass
}
