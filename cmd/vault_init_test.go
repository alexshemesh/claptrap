package cmd

import (
	"testing"
	"github.com/bouk/monkey"
	"github.com/alexshemesh/claptrap/lib/vault"
	"github.com/alexshemesh/claptrap/lib/logs"
)

func Test_vault_initialize(t *testing.T){

	var unsealed bool
	monkey.Patch(vault.VaultClient.IsInitialized, func(this vault.VaultClient)(retVal bool, err error){
		return false, nil
	})
	monkey.Patch(vault.VaultClient.Initialize, func(this vault.VaultClient)(keys []string, token string, err error){
		keys = append(keys,"somekey")
		token = "token"
		return keys,token, nil
	})


	monkey.Patch(vault.VaultClient.Unseal, func(thsi vault.VaultClient,keys []string)(err error){
		unsealed = true

		return err
	})
	vaultClient := vault.NewVaultClient("someserver",*logs.NewLogger("whatever"))
	vaultClient.Initialize()
	defer func(){
		monkey.Unpatch(vault.VaultClient.IsInitialized)
		monkey.Unpatch(vault.VaultClient.Initialize)
		monkey.Unpatch(vault.VaultClient.Unseal)
	}()

	err := runInit(*logs.NewLogger("testloger"),"someulr")
	if err != nil {
		t.Error(err)
	}

	if unsealed ==false{
		t.Errorf("Cannot unseal vault")
	}

}
