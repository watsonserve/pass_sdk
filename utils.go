package pass_sdk

import (
	"net/http"
	"net/url"
)

func cutUri(raw *url.URL) string {
	if nil == raw {
		return ""
	}
	uri := raw.Path
	if "" != raw.RawQuery {
		uri += "?" + raw.RawQuery
	}
	if "" != raw.Fragment {
		uri += "#" + raw.Fragment
	}
	return uri
}

/**
 * @param {string} authAddr 授权路径，例如：/auth
 * @param {*url.URL} raw 当前路径，将被转换为：%2Fpathname%3Fsearch%23hash
 * @return string /auth?r=%2Fpathname%3Fsearch%23hash
 */
func getAuthAddr(scheme, host, authAddr, rd string) string {
	redirect := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   authAddr,
	}
	q := redirect.Query()
	if "" != rd {
		q.Set("rd", rd)
		redirect.RawQuery = q.Encode()
	}
	return redirect.String()
}

// 检查ref
func chkReferer(req *http.Request, selfDomain string) bool {
	referer := req.Header.Get("referer")
	if "" == referer {
		return false
	}
	refUri, err := url.Parse(referer)
	if nil != err {
		return false
	}

	refHost := refUri.Scheme + "://" + refUri.Host
	if selfDomain != refHost {
		return false
	}
	return true
}
