package algorithms

import (
	"errors"
	"sync"
)

type RoundRobin struct {
	servers []string
	count   int
	mutex   sync.Mutex
}

// func Round(command []byte, config *RoundRobin) (string, error) {
// 	input := strings.Split(strings.TrimSpace(string(command)), " ")
// 	if len(input) < 1 {
// 		return "", errors.New("ERROR: Invalid Input in the RoundRobin algorithm (Empty Input)")
// 	}

// 	// var count int = 0
// 	switch input[0] {
// 	case "POOL":
// 		for i := 1; i < len(input); i++ {
// 			config.servers = append(config.servers, input[i])
// 		}
// 		return "OK", nil
// 	case "PICK":
// 		if len(config.servers) == 0 {
// 			// fmt.Println(servers)
// 			return "", errors.New("ERROR: NO Server in the list")
// 		}
// 		current := config.servers[config.count%len(config.servers)]
// 		config.count++
// 		return current, nil
// 	case "RESET":
// 		config.count = 0
// 		return "OK", nil
// 	default:
// 		return "", errors.New("ERROR: Invalid Input")
// 	}
// }

func NewRoundRobin() *RoundRobin {
	return &RoundRobin{}
}

func (r *RoundRobin) AddServers(servers []string) (string, error) {
	// servers := strings.Split(strings.TrimSpace(string(list)), " ")
	if len(servers) == 0 {
		return "", errors.New("No input for the servers")
	}
	r.servers = append(r.servers, servers...)
	return "OK", nil
}

func (r *RoundRobin) GetServer() (string, error) {
	if len(r.servers) == 0 {
		return "", errors.New("No Server in the list and trying to get one [503]")
	}

	r.mutex.Lock()
	server := r.servers[r.count]
	r.count = (r.count + 1) % len(r.servers) // for avoiding overflows
	r.mutex.Unlock()

	return server, nil
}

func (r *RoundRobin) Reset() string {

	r.mutex.Lock()
	r.count = 0
	r.mutex.Unlock()

	return "OK"
}
