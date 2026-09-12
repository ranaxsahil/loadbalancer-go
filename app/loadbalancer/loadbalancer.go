package loadbalancer

import (
	"net/http/httputil"

	"github.com/ranaxsahil/loadbalancer-go/app/algorithms"
	"github.com/ranaxsahil/loadbalancer-go/app/proxy"
)

type loadbalancer struct {
	p *proxy.Proxies
	a ServerV1
}

var loadb = loadbalancer{}

func StartServer(servers []string) { // init the proxy and then add server to the algorithm and return
	proxy := proxy.NewProxy()
	loadb.p = proxy
	proxy.IntiProxy(servers)
	algo := algorithms.NewRoundRobin()
	loadb.a = algo
	algo.AddServers(servers)
}

func MakeConnection() *httputil.ReverseProxy { // make the connection to the server and return back the request
	for {
		server, _ := loadb.a.GetServer()
		proxy, err := loadb.p.GetProxy(server)
		if err != nil {
			continue
		} else {
			return proxy
		}
	}
}

func CloseServer() { // Close the connections and then process the current requests and then terminate the process

}
