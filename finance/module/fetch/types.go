package fetch

import (
	"net/url"
	"sort"
	"strings"
)

type Map map[string]string

func (m Map) Set(key string, val string) {
	m[key] = val
}

func (m Map) Get(key string) string {
	ret, ok := m[key]
	if !ok {
		return ""
	}
	return ret
}

func (m Map) ToQueryString() string {
	values := url.Values{}
	for k, v := range m {
		values.Set(k, v)
	}
	return values.Encode()
}

// ToQueryStringAWSStandard
// it's ToQueryString but for AWS request
func (m Map) ToQueryStringAWSStandard() string {
	var a []string
	for k, v := range m {
		vEscaped := url.QueryEscape(v)
		vEscaped = strings.ReplaceAll(vEscaped, "+", "%20")
		vEscaped = strings.ReplaceAll(vEscaped, "*", "%2A")
		vEscaped = strings.ReplaceAll(vEscaped, "%7E", "~")
		a = append(a, url.QueryEscape(k)+"="+vEscaped)
	}
	sort.Strings(a)
	str := ""
	for i, s := range a {
		if i > 0 {
			str += "&"
		}
		str += s
	}
	return str
}

type Options struct {
	Method   string
	Url      string
	Query    Map
	Headers  Map
	ProxyURL string
	Body     []byte
}

type Response struct {
	Status  int
	Headers Map
	Body    []byte
}

func (r *Response) IsStatusError() bool {
	return r.Status >= 400
}
