package algorithms

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"testing"
)

// test for empty input
// test for empty input and then making all the servers dead and then try to get the server
// try to do a test in a go routine for concurrency
// test for conitnuous up and down of server
// test for long term servers
// test for single server
// test for 100 servers
// some other tests

func TestUnitTests(t *testing.T) {
	// tests for adding servers

	AddServertests := []struct {
		name    string
		in      []string
		wantErr error
		wantLen int
	}{
		{ // setting the servers properly and all
			name:    "Adding Servers",
			in:      []string{"s1", "s2", "s3", "S4", "S5", "S6"},
			wantErr: nil,
			wantLen: 6,
		},
		{ // response on the empty input
			name:    "Empty List",
			in:      []string{},
			wantErr: ErrEmptyInput,
			wantLen: 0,
		},
		{
			name:    "Dublicate inputs",
			in:      []string{"s1", "s1", "s1", "s2", "s3", "s4"},
			wantErr: nil,
			wantLen: 4,
		},
	}

	for _, test := range AddServertests {
		t.Run(test.name, func(t *testing.T) {
			round := NewRoundRobin()
			err := round.AddServers(test.in)

			if err != test.wantErr {
				t.Fatalf("Got Error %s; Want %s", err, test.wantErr)
			}

			if test.wantErr != nil {
				return
			}

			if len(round.servers) != test.wantLen {
				t.Fatalf("Got Servers in list %d; want %d", len(round.servers), test.wantLen)
			}

			for server := range test.in {
				_, found := round.findServer[test.in[server]]
				if !found {
					t.Errorf("Server %s, Not found in the map", test.in[server])
				}
			}
		})

	}

	// tests for Getting server

	GetServertests := []struct {
		name    string
		in      []string
		dead    []string
		alive   []string
		wantErr error
		wantOut []string
	}{
		{
			name:    "No Servers",
			in:      []string{},
			dead:    []string{},
			alive:   []string{},
			wantErr: ErrNoAvailabeServer,
			wantOut: []string{},
		},
		{
			name: "No Alive Server",
			// Scaled up: 5 servers added, all 5 crash
			in:      []string{"s1", "s2", "s3", "s4", "s5"},
			dead:    []string{"s1", "s2", "s3", "s4", "s5"},
			alive:   []string{},
			wantErr: ErrNoAliveServer,
			wantOut: []string{},
		},
		{
			name: "One Dead Server",
			// Scaled up: 8 servers, exactly 1 goes down.
			in:      []string{"A", "B", "C", "D", "E", "F", "G", "H"},
			dead:    []string{"D"}, // D is dead
			alive:   []string{},
			wantErr: nil,
			// Pulling 16 times to ensure it cleanly wraps around multiple times skipping 'D'
			wantOut: []string{
				"A", "B", "C", "E", "F", "G", "H",
				"A", "B", "C", "E", "F", "G", "H",
				"A", "B",
			},
		},
		{
			name: "few Dead Server",
			// Scaled up: 10 servers, scattered outages
			in:      []string{"s1", "s2", "s3", "s4", "s5", "s6", "s7", "s8", "s9", "s10"},
			dead:    []string{"s2", "s3", "s7", "s9"}, // 4 servers are dead
			alive:   []string{},
			wantErr: nil,
			// Alive servers: s1, s4, s5, s6, s8, s10 (6 servers)
			// Pulling 15 times to test deep wrap-around logic
			wantOut: []string{
				"s1", "s4", "s5", "s6", "s8", "s10",
				"s1", "s4", "s5", "s6", "s8", "s10",
				"s1", "s4", "s5",
			},
		},
		{
			name: "Dead Then Alive",
			// Scaled up: 7 servers. Multiple die, but some recover before requests start.
			in:      []string{"web1", "web2", "web3", "web4", "web5", "web6", "web7"},
			dead:    []string{"web2", "web4", "web6", "web7"}, // 4 die
			alive:   []string{"web2", "web6"},                 // 2 come back online (web4 and web7 remain dead)
			wantErr: nil,
			// Alive servers: web1, web2, web3, web5, web6
			// Pulling 12 times
			wantOut: []string{
				"web1", "web2", "web3", "web5", "web6",
				"web1", "web2", "web3", "web5", "web6",
				"web1", "web2",
			},
		},
		{
			name:    "All Dead Then All Alive",
			in:      []string{"nodeA", "nodeB", "nodeC", "nodeD", "nodeE", "nodeF"},
			dead:    []string{"nodeA", "nodeB", "nodeC", "nodeD", "nodeE", "nodeF"},
			alive:   []string{"nodeA", "nodeB", "nodeC", "nodeD", "nodeE", "nodeF"},
			wantErr: nil,
			wantOut: []string{
				"nodeA", "nodeB", "nodeC", "nodeD", "nodeE", "nodeF",
				"nodeA", "nodeB", "nodeC", "nodeD", "nodeE", "nodeF",
				"nodeA", "nodeB", "nodeC",
			},
		},
	}

	for _, test := range GetServertests {
		t.Run(test.name, func(t *testing.T) {
			round := NewRoundRobin()
			round.AddServers(test.in)

			for _, server := range test.dead {
				round.Dead(server)
			}
			for _, server := range test.alive {
				round.Alive(server)
			}

			for _, out := range test.wantOut {
				server, err := round.GetServer()
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("Got Error %s; want %s", err, test.wantErr)
				}

				if err != nil {
					return
				}

				if server != out {
					t.Fatalf("Got Server %s; want %s", server, out)
				}
			}
		})
	}

	testResetTests := []struct {
		name      string
		in        []string
		wantAlive int
		makeDead  int
	}{
		{
			name:      "Empty Server",
			in:        []string{},
			wantAlive: 0,
			makeDead:  0,
		},
		{
			name:      "Some Servers",
			in:        []string{"s1", "s2", "s3", "s4"},
			wantAlive: 4,
			makeDead:  0,
		},
		{
			name:      "One Server",
			in:        []string{"standalone-node"},
			wantAlive: 1,
			makeDead:  0,
		},
		{
			name:      "Make Server Dead",
			in:        []string{"node1", "node2", "node3", "node4", "node5", "node6", "node7", "node8", "node9", "node10"},
			makeDead:  4,
			wantAlive: 10,
		},
		{
			name:      "Make All servers Dead",
			in:        []string{"web01", "web02", "web03", "web04", "web05", "web06", "web07", "web08", "web09", "web10", "web11", "web12", "web13", "web14", "web15"},
			makeDead:  15,
			wantAlive: 15,
		},
	}

	for _, test := range testResetTests {
		t.Run(test.name, func(t *testing.T) {
			round := NewRoundRobin()
			round.AddServers(test.in)

			for i := 0; i < test.makeDead; i++ {
				round.Dead(round.servers[i].server)
			}
			round.count = -5

			round.Reset()

			if round.count != 0 {
				t.Fatalf("Got Count %d, Want 0", round.count)
			}
			if round.aliveCount != len(round.servers) {
				t.Fatalf("Got Alive Servers %d, Want %d", round.aliveCount, len(round.servers))
			}

		})
	}

	// Random Alive and Dead of the servers
	t.Run("Random Alive and Dead", func(t *testing.T) {
		round := NewRoundRobin()
		testRuns := 1000
		SERVERS := 100

		// putting random server in the strings
		var servers []string
		for i := 1; i <= SERVERS; i++ {
			servers = append(servers, fmt.Sprintf("server-%d", i))
		}
		round.AddServers(servers)

		dead := make(map[string]int)

		for i := 0; i < testRuns; i++ {
			// toggle 4 random servers
			for j := 0; j < 4; j++ {
				toggleIdx := rand.Intn(len(round.servers))
				targetServer := round.servers[toggleIdx].server

				if round.servers[toggleIdx].isAlive {
					round.Dead(targetServer)
					dead[targetServer] = 1
				} else {
					round.Alive(targetServer)
					delete(dead, targetServer)
				}
			}

			// check for the errors here
			for k := 0; k < len(round.servers); k++ {
				server, err := round.GetServer()

				// all the servers are dead
				if len(dead) == len(round.servers) {
					if err == nil {
						t.Fatalf("Run %d: expected error as all servers are dead, but got server: %s", i, server)
					}
					break
				}

				// If not all servers are dead, we should NOT get an error
				if err != nil {
					t.Fatalf("Run %d, Fetch %d: unexpected error: %v", i, k, err)
				}

				// verify the returned server is actually alive
				if _, found := dead[server]; found {
					t.Fatalf("Run %d, Fetch %d: Got Server %s; but it is marked as dead!", i, k, server)
				}
			}
		}
	})
}

func TestConcurrency(t *testing.T) {
	servers := []string{"A", "B", "C", "D", "E", "F", "G"}
	rr := NewRoundRobin()
	rr.AddServers(servers)

	const Workers = 100
	const Work = 200
	var wg sync.WaitGroup

	for worker := 0; worker < Workers; worker++ {
		wg.Add(1)
		go func(workerID int, server string) {
			defer wg.Done()

			for curr := 0; curr < Work; curr++ {
				rr.GetServer()

				random := curr % 2

				if random == 0 {
					rr.Dead(server)
				} else {
					rr.Alive(server)
				}
			}

		}(worker, servers[worker%len(servers)])
	}
	wg.Wait()

	if rr.aliveCount < 0 || rr.aliveCount > 7 {
		t.Fatalf("Alive Server Wanted > 0 || < 7; Got %d", rr.aliveCount)
	}

}
