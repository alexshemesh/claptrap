package kraken

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/logs"
)


func Test_NewKrakenClient(t *testing.T){
	settings := *vault.NewVaultTestKit()

	log := *logs.NewLogger("kraken test")

	val :=  NewKrakenClient( log, settings )
	if ( val == nil){
		t.Error("Cannot create NewKrakenClient")
	}
}