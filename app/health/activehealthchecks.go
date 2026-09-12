package health

import (
	"net/http"

	"github.com/ranaxsahil/loadbalancer-go/app/loadbalancer"
)

type activeHealth struct {
	server string
	health int
}

type ActiveServer struct {
	servers []activeHealth
}

func (a *ActiveServer) AddServers(servers []string) {
	for val := range servers {
		newServer := activeHealth{servers[val], 0}
		a.servers = append(a.servers, newServer)
	}
}

func (a *ActiveServer) CheckServer(algorithm loadbalancer.ServerV1, time int) { // 3 drops then remove // 2 ups then add again
	for val := range a.servers {
		_, err := http.Get(a.servers[val].server)
		if err != nil {
			a.servers[val].health -= 1
		} else {
			if a.servers[val].health <= -3 {
				a.servers[val].health -= 1

			} else {
				a.servers[val].health += 1

			}
		}
		if a.servers[val].health == -5 {
			a.servers[val].health = 0
			algorithm.AddServers(a.servers[val].server) // make a add server func in the algos and also a addservers
		} else if a.servers[val].health == -3 {
			algorithm.RemoveServer(a.servers[val].server)
		}
	}

}
