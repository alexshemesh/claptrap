package quoinex

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/logs"
)

func TestNewQuoineClient(t *testing.T) {
	log := *logs.NewLogger("kraken test")
	client := NewQuoineClient(log, "aaa", "bbb")
	if client == nil {
		t.Error("Cannot create QuoineClient")
	}
}