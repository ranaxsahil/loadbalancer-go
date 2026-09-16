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
	servers    []robinServer
	findServer map[string]int // [server]position
	aliveCount int
	count      int
	mu         sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{findServer: make(map[string]int)}
}

func (r *RoundRobin) AddServers(servers []string) error {
	if len(servers) == 0 {
		return ErrEmptyInput
	}
	var pos int = -1
	for _, server := range servers {
		pos++
		if _, ok := r.findServer[server]; ok {
			continue
		}
		currServer := robinServer{server: server, isAlive: true}

		// saving in the servers array
		r.servers = append(r.servers, currServer)

		// saving in the map
		r.findServer[server] = pos
		r.aliveCount += 1
	}
	return nil
}

func (r *RoundRobin) GetServer() (string, error) {
	if len(r.servers) == 0 {
		return "", ErrNoAvailabeServer
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.aliveCount == 0 {
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
	r.mu.Lock() // make every server alive and aliveserver count to len(robinServers) and count to zero
	defer r.mu.Unlock()

	for key := range r.servers {
		r.servers[key].isAlive = true
	}
	r.aliveCount = len(r.servers)
	r.count = 0
}

// TODO -- Make sure the aliveCount never go below 0 --

func (r *RoundRobin) Dead(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return ErrNoServerFoundDead
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if r.servers[foundServer].isAlive { // if alive then
		r.servers[foundServer].isAlive = false
		r.aliveCount -= 1
	}
	return nil
}

func (r *RoundRobin) Alive(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return ErrNoServerFoundAlive
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.servers[foundServer].isAlive { // if not alive then
		r.servers[foundServer].isAlive = true
		r.aliveCount += 1
	}

	return nil
}
