package pass_sdk

import (
    "encoding/json"
    "errors"
    "net/http"
)

type stdJSON struct {
    Status    bool      `json:"status"`
    Msg       string    `json:"msg"`
}

func (this *stdJSON) Chars() []byte {
    foo, _ := json.Marshal(this)
    return foo
}


type defaultBizAO struct {
    db map[string]*UserData
}

func (this *defaultBizAO) Get(id string) (*UserData, error) {
    var err error = nil
    ret := this.db[id]
    if nil == ret {
        err = errors.New("")
    }
    return ret, err
}

func (this *defaultBizAO) Save(id string, userData *UserData) error {
    if "" == id || nil == userData {
        return errors.New("")
    }
    this.db[id] = userData
    return nil
}

func (this *defaultBizAO) Error(res http.ResponseWriter, code int, msg string) {
    header := res.Header()
    header.Set("Content-Type", "application/json; charset=utf-8");
    res.WriteHeader(code)
    ret := &stdJSON {
        Status: false,
        Msg:    msg,
    }
    res.Write(ret.Chars())
}
