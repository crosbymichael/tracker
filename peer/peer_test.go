package peer

import (
	"testing"
)

func TestSerialize(t *testing.T) {
	tests := []struct {
		p        Peer
		expected string
	}{
		{
			p: Peer{
				IP:   "10.10.10.10",
				Port: 55555,
			},
			expected: "d2:ip11:10.10.10.104:porti55555ee",
		},
		{
			p: Peer{
				ID:        "1000",
				IP:        "10.10.10.10",
				Port:      55555,
				InfoHash:  "deadbeef",
				Key:       "secret_key",
				BytesLeft: 10000,
			},
			expected: "d2:ip11:10.10.10.104:porti55555ee",
		},
	}

	for _, test := range tests {
		got, err := test.p.BTSerialize()
		if err != nil {
			t.Fatal(err)
		}
		if got != test.expected {
			t.Errorf("expected %q, got %q", test.expected, got)
		}
	}
}
