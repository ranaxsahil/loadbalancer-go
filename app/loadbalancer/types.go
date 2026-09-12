package loadbalancer

// Simpler Algorithms
type ServerV1 interface {
	GetServer() (string, error)
	AddServers([]string) (string, error)
	RemoveServer(string)
	Reset()
}

// More complicated Algorithms
type ServerV2 interface {
	GetServer() (string, error)
	AddServers(string)
	RemoveServer(string)
	Reset()
	Done()
}
