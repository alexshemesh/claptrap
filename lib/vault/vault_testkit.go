package vault

import (
	"fmt"
	"github.com/alexshemesh/claptrap/lib/contracts"
)

type VaultTestKit struct {
	values map[string]string
	errors map[string]error
	loginRes bool
	tokenToSet string
	userName string
}


func NewVaultTestKit()(retVal *VaultTestKit){
	retVal = &VaultTestKit{}
	retVal.values = make(map[string]string)
	retVal.errors = make(map[string]error)
	retVal.loginRes = true
	retVal.tokenToSet = "token"
	retVal.userName = "defaultuser"
	return retVal
}

func (this VaultTestKit)SetValue(path string, value string) (err error){
	this.values[ path] = value
	return err
}

func (this VaultTestKit)GetValue(path string) (retVal string, err error){
	if this.values[ path] != "" {
		retVal = this.values[path]
		err = this.errors[path]
	}else if this.errors[path] != nil {
		err = this.errors[path]
	} else {
		err = fmt.Errorf("404 Not Found")
	}

	return retVal, err
}

func (this VaultTestKit)SetErrorEmulation(path string, value  error ) {
	this.errors[ path] = value
}


func (this VaultTestKit)LogIn( userName string, password string ) ( retObjSettings contracts.Settings,retObjAuth contracts.Auth ,err error){
	if !this.loginRes {
		err = fmt.Errorf("Authentication Failed")
	}else{
		retObjAuth = VaultTestKit{userName:userName, tokenToSet: this.tokenToSet, loginRes:this.loginRes}
	}
	return this , retObjAuth, err

}

func (this VaultTestKit)LogInCached(userName string)(retObjSettings contracts.Settings,retObjAuth contracts.Auth ,err error){
	if !this.loginRes {
		err = fmt.Errorf("Authentication Failed")
	}
	return this , this, err

}


func (this VaultTestKit)SaveCachedToken(userName string , token string) (err error){

	return err
}


func (this VaultTestKit)GetUserName()(retVal string){
	return this.userName
}

func (this VaultTestKit)GetUserToken()(retVal string){
	return this.tokenToSet
}


func (this VaultTestKit)RenewToken()(err error){
	return err
}



func (this *VaultTestKit) SetLogiResult(newRes bool){
	this.loginRes = newRes
}

func (this *VaultTestKit) SetTokenToReturn(newRes string){
	this.tokenToSet = newRes
}