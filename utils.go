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
