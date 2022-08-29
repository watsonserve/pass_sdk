package pass_sdk

import (
    "github.com/watsonserve/goengine"
    "net/http"
    "net/url"
)

func cutUri(raw *url.URL) string {
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
func getAuthAddr(scheme string, host string, authAddr string, raw *url.URL) string {
    curAddr := cutUri(raw)
    redirect := url.URL {
        Scheme: scheme,
        Host: host,
        Path: authAddr,
    }
    q := redirect.Query()
    q.Set("rd", curAddr)
    redirect.RawQuery = q.Encode()
    return redirect.String()
}

// 从session检出用户id
func WhoIsUser(session *goengine.Session) map[string]interface{} {
    // 检出数据
    user := session.Get(SESSION_USER_KEY)
    if nil == user {
        return nil
    }
    return user.(map[string]interface{})
}

// 检查ref
func chkReferer(req *http.Request, selfDomain string) *url.URL {
	referer := req.Header.Get("referer")
	if "" == referer {
			return nil
	}
	refUri, err := url.Parse(referer)
	if nil != err {
			return nil
	}
	
	refHost := refUri.Scheme + "://" + refUri.Host
	if selfDomain != refHost {
			return nil
	}
	return refUri
}

func passportRPC(ticket string, app string, secret string) (*UserData, error) {
    // @TODO
    return nil, nil
}
