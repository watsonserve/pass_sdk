package pass_sdk

import (
	"bytes"
	"crypto/sha512"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/watsonserve/otp"
)

func passportRPC(app, secret, scope, method, ct string, cookies []*http.Cookie, reqBody []byte) ([]byte, error) {
	_url := fmt.Sprintf("%s/api/%s.json", PASSPORT_ORIGIN, strings.ReplaceAll(scope, "_", "-"))
	var reqBuf *bytes.Buffer = nil
	if nil != reqBody {
		reqBuf = bytes.NewBuffer(reqBody)
	}
	req, err := http.NewRequest(method, _url, reqBuf)
	if nil != err {
		return nil, err
	}
	if "" != ct {
		req.Header.Set("Content-Type", ct)
	}
	if nil != cookies && 0 < len(cookies) {
		for _, ck := range cookies {
			req.AddCookie(ck)
		}
	}
	code, err := otp.GenTotp(sha512.New, secret)
	if nil != err {
		return nil, err
	}
	req.SetBasicAuth(app, code)
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	var resp *http.Response
	resp, err = cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	strBody, err := io.ReadAll(resp.Body)
	jsonBody := &stdJsonResp{}
	if err == nil {
		err = json.Unmarshal(strBody, jsonBody)
	}
	if err != nil {
		return nil, err
	}
	if !jsonBody.Status {
		return nil, errors.New(jsonBody.Msg)
	}
	return json.Marshal(jsonBody.Data)
}

func LoadByCode(clientId, secret, code, scope string) ([]byte, error) {
	reqBody := fmt.Sprintf("grant_type=authorization_code&client_id=%s&code=%s", clientId, code)
	return passportRPC(clientId, secret, scope, http.MethodPost, "application/x-www-form-urlencoded", nil, []byte(reqBody))
}

func LoadByToken(clientId, secret, token, scope string) ([]byte, error) {
	reqBody := fmt.Sprintf("grant_type=access_token&client_id=%s&access_token=%s", clientId, token)
	return passportRPC(clientId, secret, scope, http.MethodPost, "application/x-www-form-urlencoded", nil, []byte(reqBody))
}

func LoadByCookie(clientId, secret string, cookies []*http.Cookie) ([]byte, error) {
	return passportRPC(clientId, secret, "open-user", http.MethodGet, "", cookies, nil)
}
