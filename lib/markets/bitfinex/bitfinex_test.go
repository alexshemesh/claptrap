package bitfinexClient

import (
	"testing"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/logs"
)

func Test_NewBitfinexClient( t *testing.T){
	settings := vault.NewVaultTestKit()
	settings.SetValue("bitfinex/apikey", "apikey")
	settings.SetValue("bitfinex/apisecret", "secret")
	log := logs.NewLogger("bitfinex test")
	bitfinexClient := NewBitfinexClient(*log, settings)
	if bitfinexClient == nil {
		t.Errorf("Unable to create bitfinex client")
	}
}