package pass_sdk_test

import (
	"net/http"
	"pass_sdk"

	"github.com/watsonserve/goengine"
)

type bao struct {
	db     map[string]*pass_sdk.UserData
	appId  string
	secret string
}

func (b *bao) App() string {
	return b.appId
}

func (b *bao) Secret() string {
	return b.secret
}

func (b *bao) Get(res http.ResponseWriter, req *http.Request) map[string]interface{} {
	return pass_sdk.DefaultGet(res, req)
}

func (b *bao) Save(res http.ResponseWriter, req *http.Request, userData *pass_sdk.UserData) error {
	b.db[userData.OpenId] = userData
	return pass_sdk.DefaultSave(res, req, userData)
}

func (b *bao) Error(res http.ResponseWriter, req *http.Request, code int, msg string) {
	pass_sdk.DefaultError(res, code, msg)
}

func (b *bao) Scope(res http.ResponseWriter, req *http.Request, token *pass_sdk.Token_t) {
	pass_sdk.DefaultScope(b, res, req, token)
}

func ExampleBindAuthMgr() {
	b := &bao{
		db:     make(map[string]*pass_sdk.UserData),
		appId:  "appId",
		secret: "secret",
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
		return
	}
	engine := goengine.New(nil)

	engine.UseRouter(router)
	engine.ListenTCP(":8080")
}
