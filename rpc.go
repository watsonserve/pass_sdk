package pass_sdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func passportRPC(app, secret string, uri *url.URL, method, ct string, reqBody []byte) ([]byte, error) {
	to := PASSPORT_ORIGIN + cutUri(uri)
	req, err := http.NewRequest(method, to, bytes.NewBuffer(reqBody))

	if nil != err {
		return nil, err
	}
	req.Header.Set("Content-Type", ct)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(app+":"+secret)))
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:       10,
			IdleConnTimeout:    30 * time.Second,
			DisableCompression: true,
		},
	}
	var resp *http.Response
	resp, err = cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func loadToken(clientId, secret, code, redirect string) (*Token_t, error) {
	reqBody := fmt.Sprintf("grant_type=authorization_code&client_id=%s&code=%s&redirect_uri=%s", clientId, code, redirect)
	body, err := passportRPC(clientId, secret, &url.URL{}, "POST", "application/x-www-form-urlencoded", []byte(reqBody))
	if err != nil {
		return nil, err
	}
	respData := &Token_t{}
	json.Unmarshal(body, respData)
	return respData, err
}

func GetUserData(app, secret, accessToken string) (*UserData, error) {
	reqBody := fmt.Sprintf("grant_type=access_token&client_id=%s&access_token=%s", app, accessToken)
	body, err := passportRPC(app, secret, &url.URL{Path: "/api/user_info"}, "POST", "application/x-www-form-urlencoded", []byte(reqBody))
	if err != nil {
		return nil, err
	}
	usrInfo := &UserData{
		OpenId: accessToken,
	}
	json.Unmarshal(body, usrInfo)
	return usrInfo, nil
}
