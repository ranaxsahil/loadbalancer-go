package proxy

import (
	"errors"
	"net/http/httputil"
	"net/url"
)

type Proxies struct {
	pr map[string]*httputil.ReverseProxy
}

func NewProxy() *Proxies {
	return &Proxies{make(map[string]*httputil.ReverseProxy)}
}

func (p *Proxies) IntiProxy(servers []string) {
	for _, server := range servers {
		parseUrl, err := url.Parse(server)
		if err != nil {
			continue
		}
		proxy := httputil.NewSingleHostReverseProxy(parseUrl)
		p.pr[server] = proxy
	}
}

func (p *Proxies) GetProxy(server string) (*httputil.ReverseProxy, error) {
	proxy, ok := p.pr[server]
	if !ok { // i could try to make a connection here to the server and then if it still can not make it return error
		return nil, errors.New("can not get proxy")
	}
	return proxy, nil
}
