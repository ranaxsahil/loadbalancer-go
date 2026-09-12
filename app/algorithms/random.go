package algorithms

import (
	"errors"
	"math/rand"
	"sync/atomic"
)

type randomServer struct {
	server string
	conn   atomic.Int32
}

type Random struct {
	servers []randomServer
}

func (r *Random) AddServers(servers []byte) (string, error) {
	return "", nil
}

func (r *Random) GetServer() (string, error) {
	if len(r.servers) == 0 {
		return "", errors.New("NO server in the list")
	}

	ser1 := rand.Intn(len(r.servers))
	ser2 := rand.Intn(len(r.servers))

	for ser1 == ser2 {
		ser2 = rand.Intn(len(r.servers))
	}

	first := r.servers[ser1].conn.Load()
	sec := r.servers[ser2].conn.Load()

	// first the min conn then smaller lexographically
	if first < sec {
		r.servers[ser1].conn.Add(1)
	} else if first > sec {
		r.servers[ser2].conn.Add(1)
	} else if r.servers[ser1].server < r.servers[ser2].server {
		r.servers[ser1].conn.Add(1)
	} else {
		r.servers[ser2].conn.Add(1)
	}
	return "OK", nil
}

func (r *Random) Reset() string {
	for i := range r.servers {
		r.servers[i].conn.Store(0)
	}
	return "OK"
}

func (r *randomServer) Done() string {
	if r.conn.Load() <= 0 {
		return "OK"
	}
	r.conn.Add(-1)
	return "OK"
}
