package algorithms

import (
	"sync"
	// "golang.org/x/text/message"
)

type backends struct {
	server  string
	weight  int
	current int
}

type SmoothWeightedRR struct {
	servers     []backends
	totalWeight int
	count       int
	mutex       sync.Mutex
}

func NewWeightedRoundRobin() *SmoothWeightedRR {
	return &SmoothWeightedRR{}
}

// func (s *SmoothWeightedRR) AddServers(servers []byte) (string, error) {
// 	server := strings.Split(strings.TrimSpace(string(servers)), " ")
// 	message := ""
// 	for i := 0; i < len(server); i++ {
// 		data := strings.Split(strings.TrimSpace(server[i]), ":")
// 		weight := 0
// 		if len(data) == 1 {
// 			message += fmt.Sprintf("Server %s; Weight Not Assigned; Given ONE\n", data[0])
// 		} else {
// 			w, err := strconv.Atoi(data[1])
// 			if err != nil {
// 				message += fmt.Sprintf("Server %s; Invalid Weight '%s'; Given ONE\n", data[0], data[1])
// 				w = 1
// 			}
// 			weight = w
// 		}

// 		s.servers = append(s.servers, data[0])
// 		s.weights = append(s.weights, weight)
// 		// s.current = append(s.current, weight)
// 		s.totalWeight += weight
// 	}
// 	if message == "" {
// 		return "OK", nil
// 	}
// 	return "", errors.New(message)
// }

// func (s *SmoothWeightedRR) GetServer() (string, error) {
// 	if len(s.servers) == 0 {
// 		return "", errors.New("No Server in the list and trying to get one [503]")
// 	}
// 	for i := 0; i < len(s.servers); i++ {
// 		s.current[i] += s.weights[i]
// 	}
// 	max := 0
// 	for i := 0; i < len(s.servers); i++ {
// 		if s.current[i] > max {
// 			max = i
// 		}
// 	}
// 	server := s.servers[max]
// 	s.current[max] -= s.totalWeight

// 	return server, nil
// }

// func (s *SmoothWeightedRR) Reset() string {
// 	for i := 0; i < len(s.current); i++ {
// 		s.current[i] = 0
// 	}
// 	return "OK"
// }
