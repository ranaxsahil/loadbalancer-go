package health

import (
	"errors"
	"testing"

	"github.com/ranaxsahil/loadbalancer-go/app/loadbalancer"
)

type testAlgo struct {
	loadbalancer.ServerV1
	deadServers  map[string]bool
	aliveServers map[string]bool
}

func (ta *testAlgo) Dead(server string) {
	ta.deadServers[server] = true
	delete(ta.aliveServers, server)
}

func (ta *testAlgo) Alive(server string) {
	ta.aliveServers[server] = true
	delete(ta.deadServers, server)
}
func TestChecks(t *testing.T) {

	pingCounts := make(map[string]int) // will change the state of the servers to get different test results

	testCases := []struct {
		name      string
		in        []string // server
		runs      int
		wantErr   error
		wantDead  []string
		wantAlive []string
		ping      func(url string) error
	}{
		{
			name:      "Always Healthy (Never dies)",
			in:        []string{"NodeA"},
			runs:      10,
			wantDead:  []string{},
			wantAlive: []string{"NodeA"}, // Never died, so Alive() is never explicitly called
			ping: func(url string) error {
				return nil
			},
		},
		{
			name:      "Straight to Dead (3 consecutive fails)",
			in:        []string{"NodeB"},
			runs:      5,
			wantDead:  []string{"NodeB"},
			wantAlive: []string{},
			ping: func(url string) error {
				return errors.New("timeout") // Always fails
			},
		},
		{
			name:      "Recovers (3 fails -> Dead, then 2 successes -> Alive)",
			in:        []string{"NodeC"},
			runs:      5,
			wantDead:  []string{},
			wantAlive: []string{"NodeC"},
			ping: func(url string) error {
				pingCounts[url]++
				if pingCounts[url] <= 3 {
					return errors.New("fail") // Runs 1, 2, 3 fail
				}
				return nil // Runs 4, 5 succeed
			},
		},
		{
			name:      "Resets Fail Counter (2 fails, 1 success, 2 fails = NOT dead)",
			in:        []string{"NodeD"},
			runs:      5,
			wantDead:  []string{}, // Never hits 3 in a row
			wantAlive: []string{"NodeD"},
			ping: func(url string) error {
				pingCounts[url]++
				if pingCounts[url] == 3 {
					return nil // Success on 3rd run resets the fail counter
				}
				return errors.New("fail") // Fails on 1, 2, 4, 5
			},
		},
		{
			name:      "Resets Success Counter (Dies, 1 success, 1 fail = Stays dead)",
			in:        []string{"NodeE"},
			runs:      5,
			wantDead:  []string{"NodeE"}, // Still dead at the end
			wantAlive: []string{},
			ping: func(url string) error {
				pingCounts[url]++
				if pingCounts[url] <= 3 {
					return errors.New("fail") // Runs 1, 2, 3 fail -> Node dies
				}
				if pingCounts[url] == 4 {
					return nil // Run 4 succeeds (1 consecutive success)
				}
				return errors.New("fail") // Run 5 fails -> resets success counter!
			},
		},
		{
			name:      "Flapping (Dies, Recovers, Dies again)",
			in:        []string{"NodeF"},
			runs:      8,
			wantDead:  []string{"NodeF"}, // Final state should be dead
			wantAlive: []string{},
			ping: func(url string) error {
				pingCounts[url]++
				c := pingCounts[url]
				// 1,2,3 Fail (Dies)
				// 4,5 Succeed (Alive)
				// 6,7,8 Fail (Dies again)
				if c == 4 || c == 5 {
					return nil
				}
				return errors.New("fail")
			},
		},
	}
	for val := range testCases {
		t.Run(testCases[val].name, func(t *testing.T) {
			pingCounts = make(map[string]int)

			algo := testAlgo{
				deadServers:  make(map[string]bool),
				aliveServers: make(map[string]bool),
			}

			for _, server := range testCases[val].in {
				algo.Alive(server)
			}

			checks := NewHealthChecker(testCases[val].ping)

			checks.AddServers(testCases[val].in)

			// health checks
			for i := 0; i < testCases[val].runs; i++ {
				checks.Check(&algo, 1)
			}

			// check if the expected servers are dead or not
			for _, dead := range testCases[val].wantDead {
				if _, found := algo.deadServers[dead]; !found {
					t.Errorf("Server Not found Dead; Got NULL, Want %s", dead)
				}
			}

			// check if the expected servers are alive or not
			for _, alive := range testCases[val].wantAlive {
				if _, found := algo.aliveServers[alive]; !found {
					t.Errorf("Server Not found Alive; Got NULL, Want %s", alive)
				}
			}

			// Check if the dead/alive server in the test is in the maps
			if len(algo.aliveServers) != len(testCases[val].wantAlive) { // alive
				t.Errorf("Wanted Alive Server to be %d, Got %d", len(testCases[val].wantAlive), len(algo.aliveServers))
			}
			if len(algo.deadServers) != len(testCases[val].wantDead) { // dead
				t.Errorf("Wanted Dead Server to be %d, Got %d", len(testCases[val].wantDead), len(algo.deadServers))
			}

		})
	}
}
