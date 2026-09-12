package algorithms

import (
	"errors"
	"sync"
)

type leastServers struct {
	server string
	conn   int
}

type LeastConnections struct {
	servers []leastServers
	instant map[string]int
	sync.Mutex
}

func (l *LeastConnections) AddServers(servers []byte) (string, error) {
	return "", nil
}

func (l *LeastConnections) GetServer() (string, error) {
	if len(l.servers) == 0 {
		return "", errors.New("NO server in the servers List")
	}
	return "", nil

}

func (l *LeastConnections) Reset() string {
	l.Lock()
	for i := range l.servers {
		l.servers[i].conn = 0
	}
	l.Unlock()
	return "OK"
}

func (l *LeastConnections) Done(server string) string {
	l.Lock()
	// for i := range l.servers {
	// 	if server == l.servers[i].server {
	// 		l.servers[i].conn--
	// 		break
	// 	}
	// }
	found := l.instant[server]
	l.servers[found].conn--
	l.Unlock()
	return "OK"
}
