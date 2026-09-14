package algorithms

import (
	"errors"
	"sync"
)

var (
	ErrEmptyInput         = errors.New("NO Servers in the List")
	ErrNoAliveServer      = errors.New("NO Server is Alive")
	ErrNoAvailabeServer   = errors.New("No Servers Available")
	ErrNoServerFoundDead  = errors.New("NO Server found for Dead")
	ErrNoServerFoundAlive = errors.New("No Server found for Alive")
)

type robinServer struct {
	server  string
	isAlive bool
}

type RoundRobin struct {
	servers      []*robinServer
	findServer   map[string]*robinServer
	aliveServers int
	count        int
	mutex        sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{findServer: make(map[string]*robinServer)}
}

func (r *RoundRobin) AddServers(servers []string) error {
	if len(servers) == 0 {
		return ErrEmptyInput
	}
	for _, server := range servers {
		if _, ok := r.findServer[server]; ok {
			continue
		}
		currServer := robinServer{server: server, isAlive: true}

		// saving in the servers array
		r.servers = append(r.servers, &currServer)

		// saving in the map
		r.findServer[server] = &currServer
		r.aliveServers += 1
	}
	return nil
}

func (r *RoundRobin) GetServer() (string, error) {
	if len(r.servers) == 0 {
		return "", ErrNoAvailabeServer
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if r.aliveServers == 0 {
		return "", ErrNoAliveServer
	}
	var server string
	start := r.count
	for {
		if r.servers[r.count].isAlive {
			server = r.servers[r.count].server
			r.count = (r.count + 1) % len(r.servers)
			break
		}
		r.count = (r.count + 1) % len(r.servers)
		if r.count == start {
			return "", ErrNoAliveServer
		}
	}

	return server, nil
}

func (r *RoundRobin) Reset() {
	r.mutex.Lock() // make every server alive and aliveserver count to len(robinServers) and count to zero
	defer r.mutex.Unlock()

	for key := range r.servers {
		r.servers[key].isAlive = true
	}
	r.aliveServers = len(r.servers)
	r.count = 0
}

// TODO -- Make sure the aliveServers never go below 0 --

func (r *RoundRobin) Dead(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return ErrNoServerFoundDead
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if foundServer.isAlive { // if alive then
		foundServer.isAlive = false
		r.aliveServers -= 1
	}
	return nil
}

func (r *RoundRobin) Alive(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return ErrNoServerFoundAlive
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !foundServer.isAlive { // if not alive then
		foundServer.isAlive = true
		r.aliveServers += 1
	}

	return nil
}
