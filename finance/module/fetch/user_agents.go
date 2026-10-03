package fetch

import "math/rand"

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Linux; U; Linux i562 x86_64) Gecko/20130401 Firefox/74.7",
	"Mozilla/5.0 (Linux; U; Linux x86_64; en-US) AppleWebKit/603.18 (KHTML, like Gecko) Chrome/51.0.1529.328 Safari/536",
	"Mozilla/5.0 (Windows; U; Windows NT 6.1; WOW64) AppleWebKit/603.23 (KHTML, like Gecko) Chrome/49.0.3091.126 Safari/535",
	"Mozilla/5.0 (U; Linux x86_64; en-US) AppleWebKit/603.8 (KHTML, like Gecko) Chrome/54.0.2931.195 Safari/535",
	"Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 10.4; Win64; x64; en-US Trident/5.0)",
	"Mozilla/5.0 (Windows; Windows NT 10.0;; en-US) AppleWebKit/533.45 (KHTML, like Gecko) Chrome/54.0.3756.341 Safari/600.7 Edge/13.96533",
	"Mozilla/5.0 (Macintosh; U; Intel Mac OS X 10_11_9; en-US) Gecko/20100101 Firefox/62.2",
	"Mozilla/5.0 (Linux; Linux x86_64; en-US) AppleWebKit/533.29 (KHTML, like Gecko) Chrome/50.0.3123.288 Safari/537",
}

func randomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}
