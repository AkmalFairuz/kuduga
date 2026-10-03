package fetch

import (
	"encoding/json"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpproxy"
	url2 "net/url"
)

func Get(url string, options ...*Options) (*Response, error) {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &Options{}
	}

	opt.Method = "GET"
	opt.Url = url

	return Request(opt)
}

func PostJSON(url string, data any, options ...*Options) (*Response, error) {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &Options{}
	}

	opt.Method = "POST"
	opt.Url = url

	if opt.Headers == nil {
		opt.Headers = Map{}
	}

	if opt.Headers.Get("Content-Type") == "" {
		opt.Headers.Set("Content-Type", "application/json")
	}
	if opt.Headers.Get("User-Agent") == "" {
		opt.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	}

	encoded, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	opt.Body = encoded

	return Request(opt)
}

func PostForm(url string, data map[string]string, options ...*Options) (*Response, error) {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &Options{}
	}

	opt.Method = "POST"
	opt.Url = url

	if opt.Headers == nil {
		opt.Headers = Map{}
	}

	if opt.Headers.Get("Content-Type") == "" {
		opt.Headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	form := url2.Values{}
	for k, v := range data {
		form.Set(k, v)
	}

	opt.Body = []byte(form.Encode())

	return Request(opt)
}

func Post(url string, options ...*Options) (*Response, error) {
	var opt *Options
	if len(options) > 0 {
		opt = options[0]
	} else {
		opt = &Options{}
	}

	opt.Method = "POST"
	opt.Url = url

	if opt.Headers == nil {
		opt.Headers = Map{}
	}
	if opt.Headers.Get("User-Agent") == "" {
		opt.Headers.Set("User-Agent", randomUserAgent())
	}

	return Request(opt)
}

func Request(options *Options) (*Response, error) {
	client := &fasthttp.Client{}

	if options.ProxyURL != "" {
		client.Dial = fasthttpproxy.FasthttpHTTPDialer(options.ProxyURL)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethod(options.Method)

	url := options.Url
	if options.Query != nil {
		url += "?" + options.Query.ToQueryString()
	}
	req.SetRequestURI(url)

	if options.Headers == nil {
		options.Headers = Map{}
	}
	if options.Headers.Get("User-Agent") == "" {
		options.Headers.Set("User-Agent", "User-Agent: Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0")
	}

	for key, val := range options.Headers {
		req.Header.Set(key, val)
	}
	if options.Body != nil {
		req.SetBody(options.Body)
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	respHeaders := Map{}
	for _, key := range resp.Header.PeekKeys() {
		respHeaders.Set(string(key), string(resp.Header.Peek(string(key))))
	}

	ret := &Response{
		Status:  resp.Header.StatusCode(),
		Headers: respHeaders,
		Body:    resp.Body(),
	}

	return ret, nil
}
