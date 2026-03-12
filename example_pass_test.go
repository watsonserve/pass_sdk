package pass_sdk_test

import (
	"encoding/json"
	"net/http"

	"github.com/watsonserve/goengine"
	"github.com/watsonserve/pass_sdk"
)

func redirect(rsp http.ResponseWriter, to string) {
	rsp.Header().Set("Location", to)
	rsp.WriteHeader(http.StatusFound)
	rsp.Write(nil)
}

type bao struct {
	db     map[string]*pass_sdk.UserData
	appId  string
	secret string
	smgr   goengine.SessionManager
}

func (b *bao) IsCheckedIn(res http.ResponseWriter, req *http.Request) bool {
	si := b.smgr.LoadSession(res, req)
	// 检出数据
	uid := ""
	err := si.Load("usr", &uid)
	return nil == err && "" != uid
}

func (b *bao) Error(res http.ResponseWriter, req *http.Request, code int, explain string) {
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.WriteHeader(code)

	ret := map[string]interface{}{"status": false}
	key := "msg"
	if http.StatusForbidden == code {
		key = "location"
	}
	if "" == explain {
		explain = http.StatusText(code)
	}
	ret[key] = explain
	body, _ := json.Marshal(ret)
	res.Write(body)
}

func (b *bao) User(rsp http.ResponseWriter, req *http.Request, usr *pass_sdk.UserData, rd string) {
	sess := b.smgr.LoadSession(rsp, req)
	sess.Set("usr", usr)
	err := b.smgr.UpData(sess, 0)
	if nil != err {
		b.Error(rsp, req, http.StatusServiceUnavailable, err.Error())
		return
	}
	redirect(rsp, rd)
}

func ExampleBindAuthMgr() {
	redis := goengine.NewRedisStore("localhost", "password", 0)
	b := &bao{
		db:     make(map[string]*pass_sdk.UserData),
		appId:  "appId",
		secret: "secret",
		smgr:   goengine.InitSessionManager(redis, "sess", "cookie_prefix", "pass_sdk_test", ""),
	}
	router := goengine.InitHttpRoute()
	err := pass_sdk.BindAuthMgr(&pass_sdk.SrvInfo{
		AppId:        "appId",
		Secret:       "secret",
		AuthPathname: "/auth",
		Scheme:       "https",
		Host:         "localhost",
	}, b, router)
	if nil != err {
		panic(err)
	}

	goengine.New(router).Listen("tcp", ":8080")
}
