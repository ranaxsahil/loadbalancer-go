package health

import (
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/ranaxsahil/loadbalancer-go/app/loadbalancer"
)

type pingFunc func(url string) error

var ErrEmptyInput = errors.New("NO Servers in the List")

type heatlhStatus bool

const (
	DOWN   heatlhStatus = false
	ACTIVE heatlhStatus = true
)

type activeHealth struct {
	server  string
	failure int
	success int
	health  heatlhStatus
}

type ActiveServer struct {
	servers    []activeHealth
	findServer map[string]int
	ping       pingFunc
	mu         sync.Mutex
}

// make a new struct for active servers
func NewHealthChecker(function pingFunc) *ActiveServer {
	if function == nil {
		function = GetRequest
	}
	return &ActiveServer{ping: function, findServer: make(map[string]int)}
}

func (a *ActiveServer) AddServers(servers []string) error {
	if len(servers) == 0 {
		return ErrEmptyInput
	}
	for _, server := range servers {
		if _, found := a.findServer[server]; found { // removing dublicates
			continue
		}
		currServer := activeHealth{server: server, failure: 0, success: 0, health: ACTIVE}

		// saving in the servers array
		a.servers = append(a.servers, currServer)
		a.findServer[server] = 1
	}
	return nil
}

//	--- TODO ----
//
// make the code concurrent so that it can be called directly and full run will not take more than 1 sec MAX
// --- ---
func (a *ActiveServer) Check(algorithm loadbalancer.ServerV1, sec int) { // 3 drops then remove // 2 ups then add again
	const FALIURE = 3 // NO of consicutive faliures
	const SUCCESS = 2 // NO of consicutive success

	time.Sleep(time.Duration(sec) * time.Millisecond)
	for val := range a.servers {
		err := a.ping(a.servers[val].server)

		a.mu.Lock()

		if a.servers[val].health {
			if err != nil {
				a.servers[val].failure++
			} else {
				a.servers[val].failure = 0
			}
			if a.servers[val].failure == FALIURE {
				algorithm.Dead(a.servers[val].server)
				a.servers[val].health = DOWN
				a.servers[val].failure = 0
				a.servers[val].success = 0
			}
		} else {
			if err != nil {
				a.servers[val].success = 0
			} else {
				a.servers[val].success++
			}
			if a.servers[val].success == SUCCESS {
				algorithm.Alive(a.servers[val].server)
				a.servers[val].health = ACTIVE
				a.servers[val].failure = 0
				a.servers[val].success = 0

			}
		}

		a.mu.Unlock()
	}

}

func GetRequest(server string) error {
	res, err := http.Get(server)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
