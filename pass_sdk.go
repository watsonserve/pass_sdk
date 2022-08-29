// 提供通用接入passport方法
//
// AUTHOR: JamesWatson (c) 2019 watsonserve.com
package pass_sdk

import (
    "fmt"
    "github.com/watsonserve/goengine"
    "github.com/watsonserve/goutils"
    "net/http"
    "net/url"
)

// 绑定pass_sdk 授权管理器
//
// boa: 业务访问对象, 如果bao = nil, 使用默认bao
func BindAuthMgr(srvInfo *SrvInfo, bao BizAO, route *goengine.HttpRoute) {
    if nil == bao {
        bao = &defaultBizAO {
            db: make(map[string]*UserData),
        }
    }
    ret := &authMgr {
        userService: userService {
            bao: bao,
            app: srvInfo.AppId,
            secret: srvInfo.Secret,
        },
        authAddr: srvInfo.AuthPathname,
        scheme: srvInfo.Scheme,
        host: srvInfo.Host,
    }
    route.Use(ret.pageFilter)
    route.Set(srvInfo.AuthPathname, ret.auth)
}

func (this *authMgr) auth(res http.ResponseWriter, session *goengine.Session, req *http.Request) {
    var refUri *url.URL
    refUri = nil
    header := res.Header()
    if "GET" == req.Method {
        refUri = chkReferer(req, PASSPORT_ORIGIN)
    }

    code := 405
    msg := "Method Not Allowed"

    for nil != refUri {
        code = 400
        ticket := req.FormValue("ticket")
        redirect := req.FormValue("rd")
        if "" == ticket {
            msg = "Bad Request"
            break
        }
        // 解析重定向地址
        redirect, _ = url.QueryUnescape(redirect)
        // 重定向地址出局
        if '/' != redirect[0] {
            msg = "Redirect Out"
            break
        }

        // 获取用户信息
        var userData *UserData
        userData, err := this.getUserInfo(ticket)
        if nil != err {
            code = 404
            msg = err.Error()
            break
        }

        // 在session中存储user
        session.Set(SESSION_USER_KEY, userData)
        err = session.Save(0)
        if nil != err {
            code = 503
            msg = err.Error()
            break
        }

        header.Set("Location", redirect)
        res.WriteHeader(302)
        return
    }
    this.bao.Error(res, code, msg)
}

/**
 * 请求授权信息
 */
func (this *authMgr) reqAuth(salt string) string {
    return goutils.MD5(this.app + this.secret + salt)
}

func (this *authMgr) pageFilter(res http.ResponseWriter, session *goengine.Session, req *http.Request) bool {
    // 授权接口地址
    if req.URL.Path == this.authAddr {
        return true
    }
    // 已登录
    userMap := WhoIsUser(session)
    if nil != userMap {
        return true
    }
    // 未登录
    header := res.Header()
    // 随机字符串
    salt := goutils.RandomString(16)
    token := this.reqAuth(salt)
    stamp := fmt.Sprintf("%d", goutils.Now())
    // passport成功后回跳地址
    redirect := getAuthAddr(this.scheme, this.host, this.authAddr, req.URL)
    // 组织参数
    app := this.app
    passParams := url.Values{}
    passParams.Set("app", app)
    passParams.Set("token", token)
    passParams.Set("redirect", redirect)
    passParams.Set("salt", salt)
    passParams.Set("stamp", stamp)
    passParams.Set("signal", goutils.MD5(app + token + salt + redirect + stamp))
    // 跳转
    header.Set("Location", fmt.Sprintf("%s/?%s", PASSPORT_ORIGIN, passParams.Encode()))
    res.WriteHeader(302)
    return false
}
