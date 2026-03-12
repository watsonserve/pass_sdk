package pass_sdk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func passportRPC(app, secret, scope, method, ct string, reqBody []byte) ([]byte, error) {
	_url := fmt.Sprintf("%s/api/%s.json", PASSPORT_ORIGIN, strings.ReplaceAll(scope, "_", "-"))
	req, err := http.NewRequest(method, _url, bytes.NewBuffer(reqBody))
	if nil != err {
		return nil, err
	}

	req.Header.Set("Content-Type", ct)
	req.SetBasicAuth(app, secret)
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

func LoadByCode(clientId, secret, code, scope string) ([]byte, error) {
	reqBody := fmt.Sprintf("grant_type=authorization_code&client_id=%s&code=%s", clientId, code)
	return passportRPC(clientId, secret, scope, "POST", "application/x-www-form-urlencoded", []byte(reqBody))
}

func LoadByToken(clientId, secret, token, scope string) ([]byte, error) {
	reqBody := fmt.Sprintf("grant_type=access_token&client_id=%s&access_token=%s", clientId, token)
	return passportRPC(clientId, secret, scope, "POST", "application/x-www-form-urlencoded", []byte(reqBody))
}
