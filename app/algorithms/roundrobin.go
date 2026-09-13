package algorithms

import (
	"errors"
	"sync"
	"sync/atomic"
)

type robinServer struct {
	server  string
	isAlive atomic.Bool
}

type RoundRobin struct {
	servers      []*robinServer
	findServer   map[string]*robinServer
	aliveServers atomic.Int32
	count        int
	mutex        sync.Mutex
}

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (r *RoundRobin) AddServers(servers []string) (string, error) {
	if len(servers) == 0 {
		return "", errors.New("No input for the servers")
	}
	for _, server := range servers {
		currServer := robinServer{server: server}
		currServer.isAlive.Store(true)

		// saving in the servers array
		r.servers = append(r.servers, &currServer)

		// saving in the map
		r.findServer[server] = &currServer
		r.aliveServers.Add(1)
	}
	return "OK", nil
}

func (r *RoundRobin) GetServer() (string, error) {
	if len(r.servers) == 0 {
		return "", errors.New("No Server in the list and trying to get one [503]")
	}
	if r.aliveServers.Load() == 0 {
		return "", errors.New("NO Alive Server")
	}

	r.mutex.Lock()
	var server string
	for {
		if r.servers[r.count].isAlive.Load() == true {
			server = r.servers[r.count].server
			r.count = (r.count + 1) % len(r.servers)
			break
		}
		r.count = (r.count + 1) % len(r.servers)
	}

	// server := r.servers[r.count]
	// r.count = (r.count + 1) % len(r.servers) // for avoiding overflows
	r.mutex.Unlock()

	return server, nil
}

func (r *RoundRobin) Reset() string {

	r.mutex.Lock()
	r.count = 0
	r.mutex.Unlock()

	return "OK"
}

// TODO -- Make sure the aliveServers never go below 0 --

func (r *RoundRobin) Dead(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return errors.New("NO Server in the list")
	}

	foundServer.isAlive.Store(false)
	r.aliveServers.Add(-1)
	return nil
}

func (r *RoundRobin) Alive(server string) error {
	foundServer, ok := r.findServer[server]
	if !ok {
		return errors.New("NO Server in the list")
	}

	foundServer.isAlive.Store(true)
	r.aliveServers.Add(1)
	return nil
}
