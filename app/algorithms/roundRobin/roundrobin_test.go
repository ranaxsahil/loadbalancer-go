package roundrobin

import (
	"testing"
)

func TestRoundRobin(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{input: "POOL s1 s2 s3", output: "OK"},
		{input: "PICK", output: "s1"},
		{input: "PICK", output: "s2"},
		{input: "RESET", output: "OK"},
		{input: "PICK", output: "s1"},
	}
	servers := RoundRobin{}
	for _, tt := range tests {
		output, err := Round([]byte(tt.input), &servers)
		if err != nil {
			t.Fatalf("ERROR: %s, Test Failed Input %s", err, tt.input)
		}
		if output != tt.output {
			t.Fatalf("ERROR: Wrong Output Wanted %s; Got %s", tt.output, output)
		}
	}
}
